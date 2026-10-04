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

	UserID   int64          `gorm:"not null"                          json:"userId"`
	Title    string         `gorm:"size:200;not null"                 json:"title"`
	Done     bool           `gorm:"not null;default:false"            json:"done"`
	DueDate  *time.Time     `gorm:"type:timestamptz"                  json:"dueDate"`
	Category *string        `gorm:"size:50"                           json:"category"`
	Priority string         `gorm:"size:10;not null;default:medium"   json:"priority"`
	Tags     pq.StringArray `gorm:"type:text[];not null;default:'{}'" json:"tags"`
}
