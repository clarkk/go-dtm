package dtm

import (
	"github.com/clarkk/go-dbd"
	"github.com/clarkk/go-util/env"
)

type Build_db_query struct {
	Base 			schema_base
	Table_extension	string
	Tx 				*dbd.Tx
	Env 			*env.Environment
}

func (b Build_db_query) Build() DB_query {
	return DB_query{
		Env_context: Env_context{
			tx:				b.Tx,
			env:			b.Env,
		},
		schema_base:		b.Base,
		table_extension:	b.Table_extension,
	}
}