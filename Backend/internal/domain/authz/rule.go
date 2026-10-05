package authz

import (
	"bytes"
	"encoding/json"
	"fmt"
)

const (
	ActionManage = "manage"

	ActionRead   = "read"
	ActionCreate = "create"
	ActionUpdate = "update"
	ActionDelete = "delete"
	ActionExport = "export"
)

const SubjectAll = "all"

type Strings []string

func (s Strings) MarshalJSON() ([]byte, error) {
	if len(s) == 1 {
		return json.Marshal(s[0])
	}
	if len(s) == 0 {
		return []byte(`[]`), nil
	}
	return json.Marshal([]string(s))
}

func (s *Strings) UnmarshalJSON(data []byte) error {
	data = bytes.TrimSpace(data)
	if len(data) == 0 || bytes.Equal(data, []byte("null")) {
		*s = nil
		return nil
	}

	if data[0] == '"' {
		var single string
		if err := json.Unmarshal(data, &single); err != nil {
			return err
		}
		*s = Strings{single}
		return nil
	}

	var many []string
	if err := json.Unmarshal(data, &many); err != nil {
		return err
	}
	*s = Strings(many)
	return nil
}

func (s Strings) Has(v string) bool {
	for _, item := range s {
		if item == v {
			return true
		}
	}
	return false
}

type Rule struct {
	Action     Strings      `json:"action"`
	Subject    Strings      `json:"subject"`
	Conditions Conditions   `json:"conditions,omitempty"`
	Fields     []string     `json:"fields,omitempty"`
	Inverted   bool         `json:"inverted,omitempty"`
	Reason     string       `json:"reason,omitempty"`
}

func Allow(action, subject string) Rule {
	return Rule{Action: Strings{action}, Subject: Strings{subject}}
}

func AllowMany(actions []string, subject string) Rule {
	return Rule{Action: Strings(actions), Subject: Strings{subject}}
}

func Deny(action, subject, reason string) Rule {
	return Rule{
		Action:   Strings{action},
		Subject:  Strings{subject},
		Inverted: true,
		Reason:   reason,
	}
}

func (r Rule) Where(c Conditions) Rule {
	r.Conditions = c
	return r
}

func (r Rule) OnFields(fields ...string) Rule {
	r.Fields = fields
	return r
}

func (r Rule) Because(reason string) Rule {
	r.Reason = reason
	return r
}

func (r Rule) matchesAction(action string) bool {
	return r.Action.Has(ActionManage) || r.Action.Has(action)
}

func (r Rule) matchesSubject(subject string) bool {
	return r.Subject.Has(SubjectAll) || r.Subject.Has(subject)
}

func (r Rule) matchesField(field string) bool {
	if len(r.Fields) == 0 {
		return true
	}
	if field == "" {
		return true
	}
	for _, f := range r.Fields {
		if f == field {
			return true
		}
	}
	return false
}

func (r Rule) Validate() error {
	if len(r.Action) == 0 {
		return fmt.Errorf("authz: قانون بدون action معتبر نیست")
	}
	if len(r.Subject) == 0 {
		return fmt.Errorf("authz: قانون بدون subject معتبر نیست")
	}
	for _, a := range r.Action {
		if a == "" {
			return fmt.Errorf("authz: نام عمل خالی است")
		}
	}
	for _, s := range r.Subject {
		if s == "" {
			return fmt.Errorf("authz: نام موضوع خالی است")
		}
	}
	return nil
}