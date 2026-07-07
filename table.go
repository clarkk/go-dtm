package dtm

import (
	"maps"
	"slices"
)

type (
	schema struct {
		schema_base
		public 		bool
		env_access 	[]string
	}
	
	schema_base struct {
		resource	string
		enums		Enums
	}
)

func (m schema_base) Resource() string {
	return m.resource
}

func (m schema_base) Enum(field string) Enum {
	enum, ok := m.enums[field]
	if !ok {
		panic("dtm: enum field not found in "+m.resource+": "+field)
	}
	return enum
}

func (m schema_base) Enum_values(field string) []string {
	enum, ok := m.enums[field]
	if !ok {
		return nil
	}
	return enum.Values()
}

func (m schema_base) Enums() Enums {
	return maps.Clone(m.enums)
}

func (m schema) Base() schema_base {
	return m.schema_base
}

func (m schema) Public() bool {
	return m.public
}

func (m schema) Env_access() []string {
	return slices.Clone(m.env_access)
}