package entities

import (
	"time"

	"github.com/google/uuid"
)

type Subscriptions struct {
	ID          uuid.UUID  `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`
	ServiceName string     `gorm:"size:100;not null" json:"service_name"`
	Price       float64    `gorm:"type:numeric(10,2);not null" json:"price"`
	UserID      uuid.UUID  `gorm:"type:uuid;not null;index" json:"user_id"`
	StartDate   time.Time  `gorm:"not null" json:"start_date"`
	EndDate     *time.Time `gorm:"index" json:"end_date,omitempty"`
	CreatedAt   time.Time  `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt   time.Time  `gorm:"autoUpdateTime" json:"updated_at"`
}
