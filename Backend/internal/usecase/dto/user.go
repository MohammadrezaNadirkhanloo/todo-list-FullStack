package dto

import (
	"time"

	"github.com/MohammadrezaNadirkhanloo/internal/domain/model"
)

type RegisterInput struct {
	Username string `json:"username" binding:"required,min=3,max=64"`
	Password string `json:"password" binding:"required,password"`
}

type LoginInput struct {
	Username string `json:"username" binding:"required,max=64"`
	Password string `json:"password" binding:"required,max=128"`
}

type UserOutput struct {
	ID        int64     `json:"id"`
	Username  string    `json:"username"`
	CreatedAt time.Time `json:"createdAt"`
}

func ToUserOutput(m model.UserModel) UserOutput {
	return UserOutput{
		ID:        m.ID,
		Username:  m.Username,
		CreatedAt: m.CreatedAt,
	}
}

type TokenPair struct {
	AccessToken   string    `json:"-"`
	AccessTokenID string    `json:"-"`
	RefreshToken  string    `json:"-"`
	ExpiresAt     time.Time `json:"expiresAt"`
}