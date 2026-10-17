package model

import (
	"database/sql"
	"time"

	"github.com/MohammadrezaNadirkhanloo/pkg/appctx"
	"gorm.io/gorm"
)

type BaseModel struct {
	ID        int64          `gorm:"primaryKey"`
	CreatedAt time.Time      `gorm:"not null"`
	UpdatedAt sql.NullTime   `gorm:""`
	DeletedAt gorm.DeletedAt `gorm:"index"`
	CreatedBy sql.NullInt64  `gorm:""`
	UpdatedBy sql.NullInt64  `gorm:""`
	DeletedBy sql.NullInt64  `gorm:""`
}

func (m *BaseModel) BeforeCreate(tx *gorm.DB) error {
	m.CreatedAt = time.Now().UTC()
	m.CreatedBy = actorFrom(tx)
	return nil
}

func (m *BaseModel) BeforeUpdate(tx *gorm.DB) error {
	m.UpdatedAt = sql.NullTime{Time: time.Now().UTC(), Valid: true}
	m.UpdatedBy = actorFrom(tx)
	return nil
}

func actorFrom(tx *gorm.DB) sql.NullInt64 {
	if tx.Statement == nil || tx.Statement.Context == nil {
		return sql.NullInt64{}
	}
	id, ok := appctx.UserID(tx.Statement.Context)
	if !ok {
		return sql.NullInt64{}
	}
	return sql.NullInt64{Int64: id, Valid: true}
}
