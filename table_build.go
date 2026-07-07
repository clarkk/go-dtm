package dtm

type Build_schema struct {
	Resource	string
	Enums		Enums
	Public		bool
	Env_access	[]string
}

func (b Build_schema) Build() schema {
	return schema{
		schema_base: schema_base{
			resource:	b.Resource,
			enums:		b.Enums,
		},
		public:			b.Public,
		env_access:		b.Env_access,
	}
}