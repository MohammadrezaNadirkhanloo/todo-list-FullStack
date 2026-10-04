package filter

import "strings"

type Variant string

const (
	VariantText        Variant = "text"
	VariantLink        Variant = "link"
	VariantNumber      Variant = "number"
	VariantCurrency    Variant = "currency"
	VariantPercentage  Variant = "percentage"
	VariantRating      Variant = "rating"
	VariantSelect      Variant = "select"
	VariantMultiSelect Variant = "multi-select"
	VariantTags        Variant = "tags"
	VariantBoolean     Variant = "boolean"
	VariantDate        Variant = "date"
	VariantDateTime    Variant = "datetime"
	VariantTime        Variant = "time"
)

const defaultRatingMax = 5

type Field struct {
	Column     string
	Variant    Variant
	Filterable bool
	Sortable   bool
	Searchable bool
	Options    []string
	Max        int
	Array      bool
}

func (f Field) RatingMax() int {
	if f.Max <= 0 {
		return defaultRatingMax
	}
	return f.Max
}

func (f Field) AllowsOption(v string) bool {
	if len(f.Options) == 0 {
		return true
	}
	for _, o := range f.Options {
		if o == v {
			return true
		}
	}
	return false
}

type Spec struct {
	fields      map[string]Field
	names       []string
	defaultSort string
}

type FieldDef struct {
	name  string
	field Field
}

func NewSpec(defaultSort string, defs ...FieldDef) Spec {
	s := Spec{
		fields:      make(map[string]Field, len(defs)),
		names:       make([]string, 0, len(defs)),
		defaultSort: defaultSort,
	}
	if s.defaultSort == "" {
		s.defaultSort = "id desc"
	}
	for _, d := range defs {
		if d.name == "" || d.field.Column == "" {
			panic("filter: every field must have an API name and a column name")
		}
		if !validColumn(d.field.Column) {
			panic("filter: invalid column name: " + d.field.Column)
		}
		if _, dup := s.fields[d.name]; dup {
			panic("filter: duplicate field name in Spec: " + d.name)
		}
		s.fields[d.name] = d.field
		s.names = append(s.names, d.name)
	}
	return s
}

func validColumn(s string) bool {
	if s == "" || len(s) > 128 {
		return false
	}
	dots := 0
	for i, r := range s {
		switch {
		case r >= 'a' && r <= 'z', r >= 'A' && r <= 'Z', r == '_':
		case r >= '0' && r <= '9':
			if i == 0 {
				return false
			}
		case r == '.':
			dots++
			if dots > 1 || i == 0 || i == len(s)-1 {
				return false
			}
		default:
			return false
		}
	}
	return true
}

func (s Spec) Field(name string) (Field, bool) {
	f, ok := s.fields[name]
	return f, ok
}

func (s Spec) DefaultSort() string { return s.defaultSort }

func (s Spec) Names() []string { return s.names }

func (s Spec) NamesText() string { return strings.Join(s.names, ", ") }

func (s Spec) SearchableFields() []Field {
	out := make([]Field, 0, len(s.names))
	for _, n := range s.names {
		if f := s.fields[n]; f.Searchable {
			out = append(out, f)
		}
	}
	return out
}

func (s Spec) HasSearchable() bool {
	for _, f := range s.fields {
		if f.Searchable {
			return true
		}
	}
	return false
}

func def(name, column string, v Variant) FieldDef {
	return FieldDef{
		name: name,
		field: Field{
			Column:     column,
			Variant:    v,
			Filterable: true,
			Sortable:   true,
		},
	}
}

func Text(name, column string) FieldDef { return def(name, column, VariantText) }

func Link(name, column string) FieldDef { return def(name, column, VariantLink) }

func Number(name, column string) FieldDef { return def(name, column, VariantNumber) }

func Currency(name, column string) FieldDef { return def(name, column, VariantCurrency) }

func Percentage(name, column string) FieldDef { return def(name, column, VariantPercentage) }

func Rating(name, column string, max int) FieldDef {
	d := def(name, column, VariantRating)
	d.field.Max = max
	return d
}

func Select(name, column string, options ...string) FieldDef {
	d := def(name, column, VariantSelect)
	d.field.Options = options
	return d
}

func MultiSelect(name, column string, options ...string) FieldDef {
	d := def(name, column, VariantMultiSelect)
	d.field.Options = options
	return d
}

func Tags(name, column string, options ...string) FieldDef {
	d := def(name, column, VariantTags)
	d.field.Options = options
	return d
}

func Bool(name, column string) FieldDef { return def(name, column, VariantBoolean) }

func Date(name, column string) FieldDef { return def(name, column, VariantDate) }

func DateTime(name, column string) FieldDef { return def(name, column, VariantDateTime) }

func Time(name, column string) FieldDef { return def(name, column, VariantTime) }

func (d FieldDef) Search() FieldDef {
	d.field.Searchable = true
	return d
}

func (d FieldDef) NoSort() FieldDef {
	d.field.Sortable = false
	return d
}

func (d FieldDef) NoFilter() FieldDef {
	d.field.Filterable = false
	return d
}

func (d FieldDef) ArrayColumn() FieldDef {
	d.field.Array = true
	return d
}