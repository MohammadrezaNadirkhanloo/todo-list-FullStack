package authz

type Ability struct {
	rules []Rule
}

func NewAbility(rules []Rule) *Ability {
	return &Ability{rules: rules}
}

func (a *Ability) Rules() []Rule {
	if a == nil || a.rules == nil {
		return []Rule{}
	}
	return a.rules
}

type Decision struct {
	Allowed bool
	Reason  string
	Fields  []string
}

func (a *Ability) Can(action, subject string, resource Resource) Decision {
	return a.can(action, subject, "", resource)
}

func (a *Ability) CanField(action, subject, field string, resource Resource) Decision {
	return a.can(action, subject, field, resource)
}

func (a *Ability) can(action, subject, field string, resource Resource) Decision {
	if a == nil {
		return Decision{Allowed: false}
	}

	for i := len(a.rules) - 1; i >= 0; i-- {
		r := a.rules[i]

		if !r.matchesAction(action) || !r.matchesSubject(subject) {
			continue
		}
		if !r.matchesField(field) {
			continue
		}

		if resource == nil {
			if r.Inverted && len(r.Conditions) > 0 {
				continue
			}
			return decisionFrom(r)
		}

		if len(r.Conditions) > 0 && !r.Conditions.Match(resource) {
			continue
		}
		return decisionFrom(r)
	}

	return Decision{Allowed: false}
}

func decisionFrom(r Rule) Decision {
	return Decision{
		Allowed: !r.Inverted,
		Reason:  r.Reason,
		Fields:  r.Fields,
	}
}

func (a *Ability) PermittedFields(action, subject string, resource Resource) []string {
	if a == nil {
		return []string{}
	}

	allowed := make(map[string]bool)
	unrestricted := false

	for _, r := range a.rules {
		if r.Inverted || !r.matchesAction(action) || !r.matchesSubject(subject) {
			continue
		}
		if resource != nil && len(r.Conditions) > 0 && !r.Conditions.Match(resource) {
			continue
		}
		if len(r.Fields) == 0 {
			unrestricted = true
			continue
		}
		for _, f := range r.Fields {
			allowed[f] = true
		}
	}

	for _, r := range a.rules {
		if !r.Inverted || !r.matchesAction(action) || !r.matchesSubject(subject) {
			continue
		}
		if len(r.Conditions) > 0 {
			if resource == nil || !r.Conditions.Match(resource) {
				continue
			}
		}
		if len(r.Fields) == 0 {
			return []string{}
		}
		unrestricted = false
		for _, f := range r.Fields {
			delete(allowed, f)
		}
	}

	if unrestricted {
		return nil
	}

	out := make([]string, 0, len(allowed))
	for f := range allowed {
		out = append(out, f)
	}
	return out
}

func (a *Ability) CheckFields(action, subject string, resource Resource, fields []string) (string, bool) {
	for _, f := range fields {
		if !a.CanField(action, subject, f, resource).Allowed {
			return f, false
		}
	}
	return "", true
}

type Scope struct {
	Allowed      bool
	Unrestricted bool
	Allow        []Conditions
	Deny         []Conditions
	Reason       string
}

func (a *Ability) ScopeFor(action, subject string) Scope {
	if a == nil {
		return Scope{}
	}

	var s Scope

	for _, r := range a.rules {
		if !r.matchesAction(action) || !r.matchesSubject(subject) {
			continue
		}

		switch {
		case r.Inverted && len(r.Conditions) == 0:
			s = Scope{Reason: r.Reason}

		case r.Inverted:
			s.Deny = append(s.Deny, r.Conditions)

		case len(r.Conditions) == 0:
			s = Scope{Allowed: true, Unrestricted: true}

		default:
			s.Allowed = true
			s.Allow = append(s.Allow, r.Conditions)
		}
	}

	return s
}

func Unrestricted() Scope {
	return Scope{Allowed: true, Unrestricted: true}
}