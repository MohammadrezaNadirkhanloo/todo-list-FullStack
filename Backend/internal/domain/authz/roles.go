package authz

const (
	SubjectDashboard = "Dashboard"
	SubjectCategory  = "Category"
	SubjectProduct   = "Product"
	SubjectTodo      = "Todo"
	SubjectUsers     = "Users"
	SubjectRoles     = "Roles"
	SubjectSettings  = "Settings"
)

const (
	RoleAdmin = "admin"
	RoleUser  = "user"

	RoleGuest = "guest"
)

var RoleRules = map[string][]Rule{

	RoleAdmin: {
		Allow(ActionManage, SubjectAll),
	},

	RoleUser: {
		Allow(ActionRead, SubjectDashboard),
		Allow(ActionRead, SubjectCategory),
		Allow(ActionRead, SubjectProduct),
	},

	RoleGuest: {
		Allow(ActionRead, SubjectCategory),
		Allow(ActionRead, SubjectProduct),
	},
}

func RulesForRoles(roles []string) []Rule {
	out := make([]Rule, 0, len(roles)*4)
	for _, role := range roles {
		out = append(out, RoleRules[role]...)
	}
	return out
}