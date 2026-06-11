package http

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	agentusecase "kerjadekat/backend/internal/agent/usecase"
	"kerjadekat/backend/internal/domain"
	"kerjadekat/backend/internal/httpx"
	"kerjadekat/backend/internal/httpapi/middleware"

	"github.com/gin-gonic/gin"
)

const maxUploadBytes = 5 << 20 // 5 MiB

type Handler struct {
	agents *agentusecase.Agents
}

func NewHandler(agents *agentusecase.Agents) *Handler {
	return &Handler{agents: agents}
}

func (h *Handler) ListTerritories(c *gin.Context) {
	cl := middleware.MustClaims(c)
	rows, err := h.agents.ListTerritories(c.Request.Context(), cl.UserID)
	if err != nil {
		httpx.WriteError(c, err)
		return
	}
	if rows == nil {
		rows = []domain.Kelurahan{}
	}
	c.JSON(http.StatusOK, gin.H{"items": rows})
}

func (h *Handler) ListWorkers(c *gin.Context) {
	cl := middleware.MustClaims(c)
	result, err := h.agents.ListWorkers(c.Request.Context(), cl.UserID)
	if err != nil {
		httpx.WriteError(c, err)
		return
	}
	if result == nil {
		result = &agentusecase.ListAgentWorkersResult{Items: []agentusecase.AgentWorkerSummary{}}
	}
	c.JSON(http.StatusOK, result)
}

func (h *Handler) RegisterWorker(c *gin.Context) {
	cl := middleware.MustClaims(c)

	if err := c.Request.ParseMultipartForm(maxUploadBytes); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid multipart form"})
		return
	}

	phone := strings.TrimSpace(c.PostForm("phone_number"))
	fullName := strings.TrimSpace(c.PostForm("full_name"))
	rtRw := strings.TrimSpace(c.PostForm("rt_rw"))
	kelurahanStr := strings.TrimSpace(c.PostForm("kelurahan_id"))

	kelurahanID, err := strconv.Atoi(kelurahanStr)
	if err != nil || kelurahanID <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid kelurahan_id"})
		return
	}

	skillIDs, err := parseSkillIDs(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var latitude, longitude float64
	if latStr := strings.TrimSpace(c.PostForm("latitude")); latStr != "" {
		latitude, _ = strconv.ParseFloat(latStr, 64)
	}
	if lngStr := strings.TrimSpace(c.PostForm("longitude")); lngStr != "" {
		longitude, _ = strconv.ParseFloat(lngStr, 64)
	}

	ktpFile, ktpHeader, err := c.Request.FormFile("ktp_photo")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ktp_photo required"})
		return
	}
	defer ktpFile.Close()

	profileFile, profileHeader, err := c.Request.FormFile("profile_photo")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "profile_photo required"})
		return
	}
	defer profileFile.Close()

	if ktpHeader.Size > maxUploadBytes || profileHeader.Size > maxUploadBytes {
		c.JSON(http.StatusBadRequest, gin.H{"error": "file too large (max 5MB)"})
		return
	}

	result, err := h.agents.RegisterWorker(c.Request.Context(), agentusecase.RegisterWorkerInput{
		AgentID:            cl.UserID,
		AgentRole:          cl.Role,
		PhoneNumber:        phone,
		FullName:           fullName,
		RtRw:               rtRw,
		KelurahanID:        kelurahanID,
		SkillIDs:           skillIDs,
		Latitude:           latitude,
		Longitude:          longitude,
		KTPPhoto:           ktpFile,
		KTPFilename:        ktpHeader.Filename,
		KTPContentType:     ktpHeader.Header.Get("Content-Type"),
		KTPSize:            ktpHeader.Size,
		ProfilePhoto:       profileFile,
		ProfileFilename:    profileHeader.Filename,
		ProfileContentType: profileHeader.Header.Get("Content-Type"),
		ProfileSize:        profileHeader.Size,
	})
	if err != nil {
		httpx.WriteError(c, err)
		return
	}

	c.JSON(http.StatusCreated, result)
}

func parseSkillIDs(c *gin.Context) ([]int, error) {
	rawParts := c.PostFormArray("skill_ids")
	if len(rawParts) == 0 {
		if single := strings.TrimSpace(c.PostForm("skill_ids")); single != "" {
			rawParts = []string{single}
		}
	}
	if len(rawParts) == 0 {
		return nil, errSimple("skill_ids required")
	}

	ids := make([]int, 0)
	for _, part := range rawParts {
		part = strings.TrimSpace(part)
		if part == "" {
			continue
		}
		if strings.HasPrefix(part, "[") {
			var arr []int
			if err := json.Unmarshal([]byte(part), &arr); err != nil {
				return nil, errSimple("skill_ids must be integers or JSON array")
			}
			ids = append(ids, arr...)
			continue
		}
		for _, tok := range strings.Split(part, ",") {
			tok = strings.TrimSpace(tok)
			if tok == "" {
				continue
			}
			id, err := strconv.Atoi(tok)
			if err != nil {
				return nil, errSimple("skill_ids must be integers")
			}
			ids = append(ids, id)
		}
	}
	if len(ids) == 0 {
		return nil, errSimple("skill_ids required")
	}
	return ids, nil
}

type simpleError string

func (e simpleError) Error() string { return string(e) }

func errSimple(msg string) error { return simpleError(msg) }
