package dtm

import "net/http"

type (
	Constraints struct{
		constraints []constraint
	}
	
	env_lang_error interface {
		Env_lang_error(key string, replace map[string]any) error
	}
	
	constraint struct {
		fn		action_id_func
		id		uint64
		error	string
	}
)

func NewConstraints() *Constraints {
	return &Constraints{}
}

func (c *Constraints) Add(fn action_id_func, id uint64, error string) *Constraints {
	c.constraints = append(c.constraints, constraint{fn, id, error})
	return c
}

func (c *Constraints) Exec(m env_lang_error) error {
	for _, constraint := range c.constraints {
		used, err := constraint.fn(constraint.id)
		if err != nil {
			return err
		}
		if used {
			return NewError(http.StatusBadRequest, m.Env_lang_error(constraint.error, nil))
		}
	}
	return nil
}