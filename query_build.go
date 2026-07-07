package dtm

import (
	"github.com/clarkk/go-dbd"
	"github.com/clarkk/go-util/env"
)

type Build_model_query struct {
	Resource	string
	Enums		Enums
	Public		bool
	Env_access	[]string
	Tx 			*dbd.Tx
	Env 		*env.Environment
	Names 		Names
}

func (b Build_model_query) Build() Model_query {
	return Model_query{
		model_base: model_base{
			Env_context: Env_context{
				tx:		b.Tx,
				env:	b.Env,
			},
			schema: schema{
				schema_base: schema_base{
					resource:	b.Resource,
					enums:		b.Enums,
				},
				public:			b.Public,
				env_access:		b.Env_access,
			},
			names:		b.Names,
		},
	}
}