package dtm

import (
	"github.com/clarkk/go-dbd"
	"github.com/clarkk/go-util/env"
)

type Build_model_table struct {
	Schema			schema
	Tx 				*dbd.Tx
	Env 			*env.Environment
	Table_extension	string
	Endpoint		string
	Names 			Names
	Limit_max		uint8
	Children		bool
}

func (b Build_model_table) Build() Model_table {
	return Model_table{
		model_base: model_base{
			Env_context: 	Env_context{
				tx:			b.Tx,
				env:		b.Env,
			},
			schema:				b.Schema,
			table_extension:	b.Table_extension,
			endpoint:			b.Endpoint,
			names:				b.Names,
		},
		limit_max:				b.Limit_max,
		children:				b.Children,
	}
}