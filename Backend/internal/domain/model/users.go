package model

const (
	RoleAdmin = "admin"
	RoleUser  = "user"
)

type User struct {
	BaseModel
	Username     string `gorm:"size:64;not null"`
	PasswordHash string `gorm:"type:text;not null"`
	Enabled      bool   `gorm:"not null;default:true"`
	Roles        []Role `gorm:"many2many:user_roles;"`
}

// func (u User) FullName() string {
// 	switch {
// 	case u.FirstName == "" && u.LastName == "":
// 		return u.Username
// 	case u.FirstName == "":
// 		return u.LastName
// 	case u.LastName == "":
// 		return u.FirstName
// 	default:
// 		return u.FirstName + " " + u.LastName
// 	}
// }

func (u User) RoleNames() []string {
	names := make([]string, 0, len(u.Roles))
	for _, r := range u.Roles {
		names = append(names, r.Name)
	}
	return names
}

type Role struct {
	BaseModel

	Name        string `gorm:"size:32;not null;uniqueIndex" json:"name"`
	Description string `gorm:"size:255" json:"description"`

	Users []User `gorm:"many2many:user_roles;" json:"-"`
}