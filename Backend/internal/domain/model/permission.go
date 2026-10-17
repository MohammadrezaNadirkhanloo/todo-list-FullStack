package model

type PermissionProfile struct {
	BaseModel

	Name        string `gorm:"size:64;not null;uniqueIndex" json:"name"`
	Description string `gorm:"size:255" json:"description"`

	Enabled bool `gorm:"not null;default:true" json:"enabled"`

	Rules []PermissionRule `gorm:"foreignKey:ProfileID" json:"rules,omitempty"`
	Users []User           `gorm:"many2many:user_permission_profiles;" json:"-"`
}

type PermissionRule struct {
	BaseModel

	ProfileID int64 `gorm:"not null;index" json:"profileId"`

	Position int `gorm:"not null;default:0" json:"position"`

	Actions string `gorm:"size:255;not null" json:"actions"`

	Subjects string `gorm:"size:255;not null" json:"subjects"`

	Conditions string `gorm:"type:jsonb" json:"conditions,omitempty"`

	Fields string `gorm:"size:512" json:"fields,omitempty"`

	Inverted bool `gorm:"not null;default:false" json:"inverted"`

	Reason string `gorm:"size:255" json:"reason,omitempty"`
}
