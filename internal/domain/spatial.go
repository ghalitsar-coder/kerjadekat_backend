package domain

import (
	"context"
	"database/sql/driver"
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"math"
	"strings"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// NullPoint maps PostGIS GEOGRAPHY(POINT, 4326) for reads/writes through GORM.
// Longitude/Latitude order follows EPSG:4326 (X = lng, Y = lat).
type NullPoint struct {
	Lat   float64
	Lng   float64
	Valid bool
}

// Scan implements sql.Scanner for geography values returned as WKT/EWKB/text.
func (p *NullPoint) Scan(value interface{}) error {
	if value == nil {
		*p = NullPoint{}
		return nil
	}

	switch v := value.(type) {
	case string:
		s := strings.TrimSpace(v)
		if b, ok := decodeHexEWKB(s); ok {
			if lat, lng, ok := parseEWKBPointXY(b); ok {
				p.Lat, p.Lng, p.Valid = lat, lng, true
				return nil
			}
		}
		return p.scanWKT([]byte(s))
	case []byte:
		if len(v) > 0 && v[0] != 0x01 && v[0] != 0x00 {
			if v[0] == 'P' || v[0] == 'p' {
				return p.scanWKT(v)
			}
			if b, ok := decodeHexEWKB(string(v)); ok {
				if lat, lng, ok := parseEWKBPointXY(b); ok {
					p.Lat, p.Lng, p.Valid = lat, lng, true
					return nil
				}
			}
		}
		lat, lng, ok := parseEWKBPointXY(v)
		if ok {
			p.Lat, p.Lng, p.Valid = lat, lng, true
			return nil
		}
		return p.scanWKT(v)
	default:
		return fmt.Errorf("domain.NullPoint: unsupported scan type %T", value)
	}
}

func (p *NullPoint) scanWKT(b []byte) error {
	s := strings.TrimSpace(string(b))
	if s == "" {
		*p = NullPoint{}
		return nil
	}
	// Accept "POINT(lng lat)" or "SRID=4326;POINT(lng lat)" from ST_AsText / some drivers.
	if idx := strings.Index(strings.ToUpper(s), "POINT"); idx >= 0 {
		s = s[idx:]
	}
	var lng, lat float64
	n, err := fmt.Sscanf(s, "POINT(%f %f)", &lng, &lat)
	if err != nil || n != 2 {
		return fmt.Errorf("domain.NullPoint: parse WKT %q: %w", string(b), err)
	}
	p.Lng, p.Lat, p.Valid = lng, lat, true
	return nil
}

// Value implements driver.Valuer. Prefer GormValue for inserts/updates.
func (p NullPoint) Value() (driver.Value, error) {
	if !p.Valid {
		return nil, nil
	}
	return fmt.Sprintf("SRID=4326;POINT(%f %f)", p.Lng, p.Lat), nil
}

// GormDataType reports the PostgreSQL column type for migrations.
func (NullPoint) GormDataType() string {
	return "geography(POINT,4326)"
}

// GormValue writes geography using PostGIS so casts match the column type.
func (p NullPoint) GormValue(_ context.Context, _ *gorm.DB) clause.Expr {
	if !p.Valid {
		return clause.Expr{SQL: "NULL"}
	}
	wkt := fmt.Sprintf("SRID=4326;POINT(%f %f)", p.Lng, p.Lat)
	return clause.Expr{
		SQL:  "ST_GeographyFromText(?)",
		Vars: []interface{}{wkt},
	}
}

const (
	wkbZ    uint32 = 0x80000000
	wkbM    uint32 = 0x40000000
	wkbSRID uint32 = 0x20000000
)

// parseEWKBPointXY parses a PostGIS EWKB/ISO WKB Point (2D) little-endian payload.
func parseEWKBPointXY(b []byte) (lat, lng float64, ok bool) {
	if len(b) < 5+8+8 {
		return 0, 0, false
	}
	endian := b[0]
	if endian != 1 {
		return 0, 0, false
	}
	typ := binary.LittleEndian.Uint32(b[1:5])
	baseType := typ &^ (wkbZ | wkbM | wkbSRID)
	if baseType != 1 { // Point
		return 0, 0, false
	}
	off := uint32(5)
	if typ&wkbSRID != 0 {
		if uint32(len(b)) < off+4+8+8 {
			return 0, 0, false
		}
		off += 4 // SRID value present but not needed for lat/lng extraction here
	}
	if uint32(len(b)) < off+8+8 {
		return 0, 0, false
	}
	x := math.Float64frombits(binary.LittleEndian.Uint64(b[off : off+8]))
	y := math.Float64frombits(binary.LittleEndian.Uint64(b[off+8 : off+16]))
	return y, x, true
}

func decodeHexEWKB(s string) ([]byte, bool) {
	if s == "" || len(s)%2 != 0 {
		return nil, false
	}
	for i := 0; i < len(s); i++ {
		c := s[i]
		if (c < '0' || c > '9') && (c < 'a' || c > 'f') && (c < 'A' || c > 'F') {
			return nil, false
		}
	}
	b, err := hex.DecodeString(s)
	if err != nil || len(b) == 0 {
		return nil, false
	}
	return b, true
}
