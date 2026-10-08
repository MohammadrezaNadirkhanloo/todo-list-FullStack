package authz

const (
	SubjectTodo     = "Todo"
	SubjectSettings = "Settings"
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
		Allow(ActionManage, SubjectTodo),
	},

	RoleGuest: {},
}

func RulesForRoles(roles []string) []Rule {
	out := make([]Rule, 0, len(roles)*4)
	for _, role := range roles {
		out = append(out, RoleRules[role]...)
	}
	return out
}
