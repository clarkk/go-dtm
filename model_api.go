package dtm

type API_model struct {
	Model_table
}

func (a API_model) Limit() bool {
	return a.limit_max != 0
}

func (a API_model) Children() bool {
	return a.children
}

func (m Model_table) API_model() API_model {
	return API_model{m}
}