package dtm

import (
	"strconv"
	"net/http"
	"github.com/clarkk/go-api"
	"github.com/clarkk/go-dbd"
	"github.com/clarkk/go-util/env"
)

type (
	Model_table struct {
		model_base
		update_method 	bool
		limit_max		uint8
		children		bool
	}
	
	Env_context struct {
		tx 			*dbd.Tx
		env 		*env.Environment
	}
	
	Names 			map[string]string
	
	model_base struct {
		Env_context
		schema
		table_extension string
		names 			Names
	}
	
	action_id_func 	func(id uint64) (bool, error)
)

//	Check if two pointers are equal (have the same value or nil)
func Equal_ptr[T comparable](a, b *T) bool {
	if a == b {
		return true
	}
	if a == nil || b == nil {
		return false
	}
	return *a == *b
}

func (m *Model_table) Set_update_method(){
	m.update_method = true
}

func (m *Model_table) Update_method() bool {
	return m.update_method
}

func (m *Model_table) Children(a *api.Request) bool {
	if !m.children {
		a.Errorf(http.StatusBadRequest, "Unsupported query param: children")
		return false
	}
	return true
}

func (m *Model_table) Update_limit_max(limit *api.Limit) error {
	if limit == nil {
		return NewError(http.StatusBadRequest, m.Env_lang_error("REQUEST_LIMIT", nil))
	}
	limit.Limit_max(m.limit_max)
	return nil
}

func (m *Model_table) Update_exists(fn action_id_func, id uint64) error {
	exists, err := fn(id)
	if err != nil {
		return err
	}
	if !exists {
		return ErrNotFound
	}
	return nil
}

func (m *Model_table) Delete_exists(fn action_id_func, id uint64) error {
	empty, err := fn(id)
	if empty {
		return ErrNotFound
	}
	return err
}

func (m *Model_table) Parse_etag(s string) uint32 {
	i, err := strconv.ParseUint(s, 10, 32)
	if err != nil {
		return 0
	}
	return uint32(i)
}

func (e Env_context) Tx() *dbd.Tx {
	return e.tx
}

func (e Env_context) Env() *env.Environment {
	return e.env
}

func (m *model_base) Base() *model_base {
	return m
}

func (m *model_base) Public_access(a *api.Request) bool {
	if !m.Public() || !m.env_access() {
		a.Error(http.StatusForbidden, nil)
		return false
	}
	return true
}

func (m *model_base) Source_resource() string {
	var table = m.Resource()
	if m.table_extension != "" {
		table += "_"+m.table_extension
	}
	return table
}

func (m *model_base) Env_lang_field(field string) string {
	return m.env.Lang_string(m.names[field], nil)
}

func (m *model_base) Env_lang_error(key string, replace map[string]any) error {
	return m.env.Lang_error(key, replace)
}

func (m *model_base) env_access() bool {
	env_access := m.Env_access()
	if env_access == nil {
		return false
	}
	env_data := m.env.Data()
	for _, key := range env_access {
		if val, ok := env_data[key].(uint64); !ok || val == 0 {
			return false
		}
	}
	return true
}