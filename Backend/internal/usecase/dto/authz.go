package dto

import (
	"github.com/MohammadrezaNadirkhanloo/internal/domain/authz"
	"github.com/MohammadrezaNadirkhanloo/internal/domain/model"
)

type MeResponse struct {
	User  MeUser       `json:"user"`
	Rules []authz.Rule `json:"rules"`
}

type MeUser struct {
	ID       int64    `json:"id"`
	Username string   `json:"username"`
	Roles    []string `json:"roles"`
}

func NewMeResponse(user model.User, rules []authz.Rule) MeResponse {
	if rules == nil {
		rules = []authz.Rule{}
	}
	return MeResponse{
		User: MeUser{
			ID:       user.ID,
			Username: user.Username,
			Roles:    user.RoleNames(),
		},
		Rules: rules,
	}
}

func TodoResource(m model.Todo) authz.Resource {
	return authz.Resource{
		"id":        m.ID,
		"userId":    m.UserID,
		// "completed": m.Completed,
	}
}

var TodoFieldColumns = map[string]string{
	"title":     "title",
	"completed": "completed",
}