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