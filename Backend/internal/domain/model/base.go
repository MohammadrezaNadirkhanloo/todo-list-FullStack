package model

import (
	"database/sql"
	"time"

	"gorm.io/gorm"
)

type BaseModel struct {
	ID int64 `gorm:"primaryKey"`
	CreatedAt time.Time      `gorm:"not null"`
	UpdatedAt sql.NullTime   `gorm:""`
	DeletedAt gorm.DeletedAt `gorm:"index"`
	CreatedBy sql.NullInt64  `gorm:""`
	UpdatedBy sql.NullInt64  `gorm:""`
	DeletedBy sql.NullInt64  `gorm:""`
}
