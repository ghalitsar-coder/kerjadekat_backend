package domain

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/datatypes"
)

// Kelurahan maps kelurahans (master locality with optional centroid).
type Kelurahan struct {
	ID        int       `gorm:"primaryKey;autoIncrement"`
	Name      string    `gorm:"type:varchar(100);not null"`
	Kecamatan *string   `gorm:"type:varchar(100)"`
	Kota      *string   `gorm:"type:varchar(100)"`
	Centroid *NullPoint `gorm:"type:geography(POINT,4326);index:,type:gist"`
}

// User maps users (all roles). Schema has no soft-delete column.
type User struct {
	ID            uuid.UUID  `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	PhoneNumber   string     `gorm:"type:varchar(15);uniqueIndex:idx_users_phone;not null"`
	FullName      string     `gorm:"type:varchar(100);not null"`
	Role          string     `gorm:"type:varchar(20);index:idx_users_role;not null"`
	NikHash       *string    `gorm:"type:varchar(64);uniqueIndex"`
	ProfilePhoto  *string    `gorm:"type:varchar(500)"`
	KtpPhotoRef   *string    `gorm:"type:varchar(500)"`
	Status        string     `gorm:"type:varchar(20);not null;default:pending;index:idx_users_status"`
	VerifiedBy    *uuid.UUID `gorm:"type:uuid"`
	VerifiedAt    *time.Time `gorm:"type:timestamptz"`
	RtRw          *string    `gorm:"type:varchar(20);column:rt_rw"`
	KelurahanID   *int       `gorm:"index:idx_users_kelurahan"`
	CreatedAt     time.Time  `gorm:"type:timestamptz;not null;autoCreateTime"`
	UpdatedAt     time.Time  `gorm:"type:timestamptz;not null;autoUpdateTime"`

	Kelurahan *Kelurahan `gorm:"foreignKey:KelurahanID"`
}

// WorkerProfile maps worker_profiles.
type WorkerProfile struct {
	ID                uuid.UUID  `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	UserID            uuid.UUID  `gorm:"type:uuid;uniqueIndex:idx_worker_profiles_user;not null"`
	Bio               *string    `gorm:"type:text"`
	BaseRate          *float64   `gorm:"type:numeric(10,2)"`
	Availability      string     `gorm:"type:varchar(20);not null;default:offline;index:idx_worker_profiles_availability"`
	LastLocation      *NullPoint `gorm:"type:geography(POINT,4326);index:,type:gist"`
	LocationUpdatedAt *time.Time `gorm:"type:timestamptz"`
	RatingAvg         float64    `gorm:"type:numeric(3,2);not null;default:0"`
	RatingCount       int        `gorm:"type:integer;not null;default:0"`
	TotalJobsDone     int        `gorm:"type:integer;not null;default:0"`
	CreditScore       int        `gorm:"type:integer;not null;default:0"`
	CreatedAt         time.Time  `gorm:"type:timestamptz;not null;autoCreateTime"`

	User   User          `gorm:"foreignKey:UserID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	Skills []WorkerSkill `gorm:"foreignKey:WorkerID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
}

// SkillCategory maps skill_categories.
type SkillCategory struct {
	ID          int     `gorm:"primaryKey;autoIncrement"`
	Name        string  `gorm:"type:varchar(100);uniqueIndex;not null"`
	IconURL     *string `gorm:"type:varchar(500)"`
	Description *string `gorm:"type:text"`
}

// WorkerSkill maps worker_skills (worker_id -> worker_profiles.id).
type WorkerSkill struct {
	ID       uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	WorkerID uuid.UUID `gorm:"type:uuid;not null;uniqueIndex:idx_worker_skills_pair;index:idx_worker_skills_worker"`
	SkillID  int       `gorm:"type:integer;not null;uniqueIndex:idx_worker_skills_pair;index:idx_worker_skills_skill"`
	Level    string    `gorm:"type:varchar(20);not null;default:beginner"`

	Worker WorkerProfile `gorm:"foreignKey:WorkerID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	Skill  SkillCategory `gorm:"foreignKey:SkillID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
}

// AgentTerritory maps agent_territories composite key.
type AgentTerritory struct {
	AgentID     uuid.UUID `gorm:"type:uuid;primaryKey"`
	KelurahanID int       `gorm:"type:integer;primaryKey"`
	AssignedAt  time.Time `gorm:"type:timestamptz;not null;autoCreateTime"`

	Agent     User      `gorm:"foreignKey:AgentID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	Kelurahan Kelurahan `gorm:"foreignKey:KelurahanID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
}

// Order maps orders.
type Order struct {
	ID                 uuid.UUID  `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	ConsumerID         uuid.UUID  `gorm:"type:uuid;not null;index:idx_orders_consumer"`
	WorkerID           *uuid.UUID `gorm:"type:uuid;index:idx_orders_worker"`
	SkillID            int        `gorm:"type:integer;not null;index:idx_orders_skill"`
	Status             string     `gorm:"type:varchar(30);not null;index:idx_orders_status"`
	Description        *string    `gorm:"type:text"`
	ConsumerLocation   NullPoint  `gorm:"type:geography(POINT,4326);not null;index:,type:gist"`
	ConsumerAddress    *string    `gorm:"type:varchar(255)"`
	AgreedRate         *float64   `gorm:"type:numeric(10,2)"`
	PlatformFee        float64    `gorm:"type:numeric(10,2);not null;default:2000"`
	PaymentMethodFee   *string    `gorm:"type:varchar(50)"`
	XenditInvoiceID    *string    `gorm:"type:varchar(100)"`
	FeeAuthID          *string    `gorm:"type:varchar(100)"`
	PaymentStatus      string     `gorm:"type:varchar(20);not null;default:pending"`
	ScheduledAt        *time.Time `gorm:"type:timestamptz"`
	StartedAt          *time.Time `gorm:"type:timestamptz"`
	CompletedAt        *time.Time `gorm:"type:timestamptz"`
	CancelledReason    *string    `gorm:"type:text"`
	CreatedAt          time.Time  `gorm:"type:timestamptz;not null;autoCreateTime;index:idx_orders_created"`

	Consumer User           `gorm:"foreignKey:ConsumerID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	Worker   *User          `gorm:"foreignKey:WorkerID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL"`
	Skill    SkillCategory  `gorm:"foreignKey:SkillID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	Logs     []OrderStatusLog `gorm:"foreignKey:OrderID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	Rating   *OrderRating     `gorm:"foreignKey:OrderID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
}

// OrderStatusLog maps order_status_logs.
type OrderStatusLog struct {
	ID         uuid.UUID  `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	OrderID    uuid.UUID  `gorm:"type:uuid;not null;index:idx_order_status_logs_order"`
	FromStatus *string    `gorm:"type:varchar(30)"`
	ToStatus   string     `gorm:"type:varchar(30);not null"`
	ChangedBy  *uuid.UUID `gorm:"type:uuid"`
	ChangeTime time.Time  `gorm:"type:timestamptz;not null;autoCreateTime;index:idx_order_status_logs_time"`
	Note       *string    `gorm:"type:text"`

	Order Order `gorm:"foreignKey:OrderID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	Actor *User `gorm:"foreignKey:ChangedBy;references:ID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL"`
}

// OrderRating maps order_ratings.
type OrderRating struct {
	ID        uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	OrderID   uuid.UUID `gorm:"type:uuid;uniqueIndex;not null"`
	GivenBy   uuid.UUID `gorm:"type:uuid;not null;index:idx_ratings_consumer"`
	GivenTo   uuid.UUID `gorm:"type:uuid;not null;index:idx_ratings_worker"`
	Score     int16     `gorm:"type:smallint;not null"`
	Comment   *string   `gorm:"type:text"`
	CreatedAt time.Time `gorm:"type:timestamptz;not null;autoCreateTime"`

	Order       Order `gorm:"foreignKey:OrderID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	GivenByUser User  `gorm:"foreignKey:GivenBy;references:ID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	GivenToUser User  `gorm:"foreignKey:GivenTo;references:ID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
}

// IncomeRecord maps income_records.
type IncomeRecord struct {
	ID         uuid.UUID  `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	WorkerID   uuid.UUID  `gorm:"type:uuid;not null;index:idx_income_worker"`
	OrderID    *uuid.UUID `gorm:"type:uuid"`
	Amount     float64    `gorm:"type:numeric(10,2);not null"`
	Source     string     `gorm:"type:varchar(20);not null"`
	Verified   bool       `gorm:"not null;default:false;index:idx_income_verified"`
	RecordedAt time.Time  `gorm:"type:timestamptz;not null;autoCreateTime;index:idx_income_recorded"`

	Worker User   `gorm:"foreignKey:WorkerID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	Order  *Order `gorm:"foreignKey:OrderID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL"`
}

// AuditLog maps audit_logs.
type AuditLog struct {
	ID         uuid.UUID      `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	ActorID    *uuid.UUID     `gorm:"type:uuid;index:idx_audit_actor"`
	Action     string         `gorm:"type:varchar(100);not null;index:idx_audit_action"`
	Resource   string         `gorm:"type:varchar(100);not null"`
	ResourceID *string        `gorm:"type:varchar(100)"`
	IPAddress  *string        `gorm:"type:inet"`
	Details    datatypes.JSON `gorm:"type:jsonb"`
	CreatedAt  time.Time      `gorm:"type:timestamptz;not null;autoCreateTime;index:idx_audit_created"`

	Actor *User `gorm:"foreignKey:ActorID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL"`
}
