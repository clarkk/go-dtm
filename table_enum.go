package dtm

import "slices"

type (
	Enums			map[string]Enum
	Enum struct {
		values 		[]string
		unsets 		Enum_unsets
	}
	
	Enum_unsets map[string]Enum_unset
	Enum_unset struct {
		In 			bool
		Values 		[]string
	}
)

func NewEnum(values []string, unsets Enum_unsets) Enum {
	return Enum{
		values: values,
		unsets: unsets,
	}
}

func (m Enums) Enum_values(field string) []string {
	enum, ok := m[field]
	if !ok {
		return nil
	}
	return enum.Values()
}

func (m Enums) Enum_unsets(field string) Enum_unsets {
	enum, ok := m[field]
	if !ok {
		return nil
	}
	copy := make(Enum_unsets, len(enum.unsets))
	for field, rule := range enum.unsets {
		copy[field] = Enum_unset{
			In:		rule.In,
			Values:	slices.Clone(rule.Values),
		}
	}
	return copy
}

func (e Enum) Value(i int) string {
	return e.values[i]
}

func (e Enum) Values() []string {
	return slices.Clone(e.values)
}

func (e Enum) Index(value string) int {
	return slices.Index(e.values, value)
}

func (e Enum) Unset_ptr(enum_value *string, field string) bool {
	if enum_value == nil {
		return e.Unset("", field)
	}
	return e.Unset(*enum_value, field)
}

func (e Enum) Unset(enum_value, field string) bool {
	//	Check if enum value is valid
	if !slices.Contains(e.values, enum_value) {
		return true
	}
	
	rule, ok := e.unsets[field]
	if !ok {
		panic("dtm: unset rule not found for field: "+field)
	}
	
	contains := slices.Contains(rule.Values, enum_value)
	if rule.In {
		return contains
	}
	return !contains
}