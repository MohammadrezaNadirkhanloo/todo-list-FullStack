package model

const (
	RoleAdmin = "admin"
	RoleUser  = "user"
)

type UserModel struct {
	BaseModel

	Username     string `gorm:"size:64;not null"   json:"username"`
	PasswordHash string `gorm:"type:text;not null" json:"-"`
}