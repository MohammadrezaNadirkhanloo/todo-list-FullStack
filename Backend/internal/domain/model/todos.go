package model

import (
	"time"

	"github.com/lib/pq"
)

const (
	PriorityLow    = "low"
	PriorityMedium = "medium"
	PriorityHigh   = "high"
)

type Todo struct {
	BaseModel

	UserID   int64          `gorm:"not null"                          `
	Title    string         `gorm:"size:200;not null"                 `
	Done     bool           `gorm:"not null;default:false"            `
	DueDate  *time.Time     `gorm:"type:timestamptz"                  `
	Category *string        `gorm:"size:50"                           `
	Priority string         `gorm:"size:10;not null;default:medium"   `
	Tags     pq.StringArray `gorm:"type:text[];not null;default:'{}'" `
}
