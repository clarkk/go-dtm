package dtm

import "github.com/clarkk/go-dbd/sqlc"

const (
	Join_left Join_type = iota
	Join_inner
	Join_cross
)

type (
	Query struct {
		Select 				[]string
		Select_jsons		[]Select_json
		Joins 				Joins
	}
	
	Raw_query struct {
		Select 				[]string
		Select_jsons		[]Select_json
		Joins 				Raw_joins
	}
	
	Select_json struct {
		Table				string
		Select_field		string
		Query				Raw_query
		Skip				func() bool
	}
	
	Join_type				int
	
	Joins map[string]Join
	Join struct {
		Type				Join_type
		Table				string
		Field 				string
		Field_foreign		string
		
		Fixed_condition		bool
		Fixed_field			string
		Fixed_value			any
	}
	
	Raw_joins []Raw_join
	Raw_join struct {
		Join_multi
		Alias				string
	}
	
	Join_multi struct {
		Type				Join_type
		Table				string
		On					sqlc.Join_conditions
	}
)

func (j Joins) apply(q *sqlc.Select_query){
	for alias, join := range j {
		if join.Fixed_condition {
			switch join.Type {
			case Join_left:
				q.Left_join_fixed(join.Table, alias, join.Field, join.Field_foreign, join.Fixed_field, join.Fixed_value)
			case Join_inner:
				q.Inner_join_fixed(join.Table, alias, join.Field, join.Field_foreign, join.Fixed_field, join.Fixed_value)
			}
		} else {
			switch join.Type {
			case Join_left:
				q.Left_join(join.Table, alias, join.Field, join.Field_foreign)
			case Join_inner:
				q.Inner_join(join.Table, alias, join.Field, join.Field_foreign)
			case Join_cross:
				q.Cross_join(join.Table, alias)
			}
		}
	}
}

func (j Raw_joins) apply(q *sqlc.Select_query){
	for _, join := range j {
		switch join.Type {
		case Join_left:
			q.Left_join_multi(join.Table, join.Alias, join.On)
		case Join_inner:
			q.Inner_join_multi(join.Table, join.Alias, join.On)
		case Join_cross:
			q.Cross_join(join.Table, join.Alias)
		}
	}
}