package dtm

import (
	"database/sql"
	"github.com/clarkk/go-api"
	"github.com/clarkk/go-dbd/sqlc"
)

type DB_query struct {
	Env_context
	schema_base
	table_extension string
}

func Int64_ptr(v *int64) int64 {
	if v == nil {
		return 0
	}
	return *v
}

func (s DB_query) Select_query_id(id uint64, fields []string, joins Joins) *sqlc.Select_query {
	q := sqlc.Select_id(s.Resource(), id).
		Select(fields).
		Optimize_joins()
	apply_joins(q, joins)
	return q
}

func (s DB_query) Select_query_id_raw(id uint64, fields []string, joins Raw_joins) *sqlc.Select_query {
	q := sqlc.Select_id(s.Resource(), id).
		Select(fields)
	apply_raw_joins(q, joins)
	return q
}

func (s DB_query) Select_query(fields []string, joins Joins) *sqlc.Select_query {
	return s.Select_query_table(s.Resource(), fields, joins)
}

func (s DB_query) Select_query_raw(fields []string, joins Raw_joins) *sqlc.Select_query {
	return s.Select_query_table_raw(s.Resource(), fields, joins)
}

func (s DB_query) Select_query_table(table string, fields []string, joins Joins) *sqlc.Select_query {
	q := sqlc.Select(table).
		Select(fields).
		Optimize_joins()
	apply_joins(q, joins)
	return q
}

func (s DB_query) Select_query_table_raw(table string, fields []string, joins Raw_joins) *sqlc.Select_query {
	q := sqlc.Select(table).
		Select(fields)
	apply_raw_joins(q, joins)
	return q
}

func (s DB_query) Exec(q sqlc.SQL) (sql.Result, error){
	return s.tx.Exec(q)
}

func (s DB_query) Query_row(q sqlc.SQL, dest []any) (bool, error){
	return s.tx.Query_row(q, dest)
}

func (s DB_query) Fetch(q sqlc.SQL, scan func(rows *sql.Rows) error) error {
	rows, err := s.tx.Query(q)
	if err != nil {
		return err
	}
	defer rows.Close()
	for rows.Next() {
		if err := scan(rows); err != nil {
			return err
		}
	}
	return rows.Err()
}

func (s DB_query) Insert(fields sqlc.Map) (uint64, error){
	return s.tx.Insert(
		sqlc.Insert(s.Resource()).
		Fields(fields))
}

func (s DB_query) Insert_update_duplicate(fields sqlc.Map, update_fields []string) (uint64, error){
	return s.tx.Insert(
		sqlc.Insert(s.Resource()).
		Fields(fields).
		Update_duplicate(update_fields))
}

func (s DB_query) Insert_extension(value any, fields sqlc.Map) error {
	table, key := s.get_extension(s.table_extension)
	fields[key] = value
	return s.tx.Insert_no_return(
		sqlc.Insert(table).
		Fields(fields))
}

func (s DB_query) Update(id uint64, fields sqlc.Map, where *sqlc.Where_clause) error {
	return s.tx.Update(
		sqlc.Update_id(s.Resource(), id).Fields(fields).
		Where(where))
}

func (s DB_query) Update_where(fields sqlc.Map, where *sqlc.Where_clause) error {
	return s.tx.Update(
		sqlc.Update(s.Resource()).Fields(fields).
		Where(where))
}

func (s DB_query) Update_extension(value any, fields sqlc.Map, where *sqlc.Where_clause) error {
	table, key := s.get_extension(s.table_extension)
	return s.tx.Update(
		sqlc.Update(table).Fields(fields).
		Where(where.Eq(key, value)))
}

func (s DB_query) Delete(id uint64, where *sqlc.Where_clause) (bool, error){
	return s.tx.Delete(
		sqlc.Delete_id(s.Resource(), id).
		Where(where))
}

func (s DB_query) Delete_where(where *sqlc.Where_clause) (bool, error){
	return s.tx.Delete(
		sqlc.Delete(s.Resource()).
		Where(where))
}

func (s DB_query) Match_etag(id uint64, etag_header uint32, where *sqlc.Where_clause) (bool, bool, error){
	q := s.Select_query_id(id, []string{"etag"}, nil).
		Where(where)
	
	var etag uint32
	empty, err := s.tx.Query_row(q, []any{&etag})
	if empty {
		return false, false, nil
	}
	if err != nil {
		return false, false, err
	}
	return true, etag_header == etag, nil
}

func (s DB_query) Constraints(wheres []*sqlc.Where_clause, joins Joins) (bool, error){
	for _, where := range wheres {
		used, err := s.Exists(where, joins)
		if used || err != nil {
			return used, err
		}
	}
	return false, nil
}

func (s DB_query) Exists_id(id uint64, where *sqlc.Where_clause, joins Joins) (bool, error){
	q := sqlc.Select_id(s.Resource(), id).
		Optimize_joins()
	apply_joins(q, joins)
	return s.existence(q, where)
}

func (s DB_query) Exists_id_read_lock(id uint64, where *sqlc.Where_clause, joins Joins) (bool, error){
	q := sqlc.Select_id(s.Resource(), id).
		Lock_for_update().
		Optimize_joins()
	apply_joins(q, joins)
	return s.existence(q, where)
}

func (s DB_query) Exists(where *sqlc.Where_clause, joins Joins) (bool, error){
	q := sqlc.Select(s.Resource()).
		Optimize_joins()
	apply_joins(q, joins)
	return s.existence(q, where)
}

func (s DB_query) Get_id_exists(where *sqlc.Where_clause) (id uint64, err error){
	q := s.Select_query([]string{"id"}, nil).
		Where(where).
		Limit(0, 1)
	
	var empty bool
	empty, err = s.tx.Query_row(q, []any{&id})
	if empty {
		err = nil
	}
	return
}

func (s DB_query) Get_id_exists_read_lock(where *sqlc.Where_clause) (id uint64, err error){
	q := s.Select_query([]string{"id"}, nil).
		Where(where).
		Limit(0, 1).
		Lock_for_update()
	
	var empty bool
	empty, err = s.tx.Query_row(q, []any{&id})
	if empty {
		err = nil
	}
	return
}

func (s DB_query) Get_id_exists_omit(where *sqlc.Where_clause, omit_id uint64) (id uint64, err error){
	if omit_id != 0 {
		where.Not_eq("id", omit_id)
	}
	return s.Get_id_exists(where)
}

func (s DB_query) Update_limit_count(limit *api.Limit, where *sqlc.Where_clause) (bool, error){
	count, err := s.count(where)
	if err != nil {
		return false, err
	}
	return !limit.Count(count), nil
}

func (s DB_query) Get_ids(where *sqlc.Where_clause, read_lock bool) ([]uint64, error){
	query := s.Select_query([]string{"id"}, nil).
		Where(where)
	
	if read_lock {
		query.Lock_for_update()
	}
	
	entries := []uint64{}
	err := s.Fetch(query, func(rows *sql.Rows) error {
		var id uint64
		if err := rows.Scan(&id); err != nil {
			return err
		}
		entries = append(entries, id)
		return nil
	})
	return entries, err
}

func (s DB_query) Has_entries(where *sqlc.Where_clause, left_join Joins) (bool, error){
	q := s.Select_query([]string{"id"}, left_join).
		Where(where).
		Limit(0, 1)
	
	var id uint64
	empty, err := s.tx.Query_row(q, []any{&id})
	if empty {
		return false, nil
	}
	if err != nil {
		return false, err
	}
	return true, nil
}

func (s DB_query) existence(q *sqlc.Select_query, where *sqlc.Where_clause) (bool, error){
	q.Select([]string{"id"}).
		Where(where).
		Limit(0, 1)
		
	var id uint64
	empty, err := s.tx.Query_row(q, []any{&id})
	if empty {
		return false, nil
	}
	if err != nil {
		return false, err
	}
	return true, nil
}

func (s DB_query) count(where *sqlc.Where_clause) (count uint32, err error){
	q := sqlc.Select(s.Resource()).
		Select([]string{"count|id"})
	if where != nil {
		q.Where(where)
	}
	_, err = s.tx.Query_row(q, []any{&count})
	return
}

func (s DB_query) get_extension(extension string) (string, string){
	tbl := s.Resource()
	return tbl+"_"+extension, tbl+"_id"
}

func apply_joins(query *sqlc.Select_query, joins Joins){
	for t, join := range joins {
		if join.Fixed_condition {
			switch join.Type {
			case Join_left:
				query.Left_join_fixed(join.Table, t, join.Field, join.Field_foreign, join.Fixed_field, join.Fixed_value)
			case Join_inner:
				query.Inner_join_fixed(join.Table, t, join.Field, join.Field_foreign, join.Fixed_field, join.Fixed_value)
			}
		} else {
			switch join.Type {
			case Join_left:
				query.Left_join(join.Table, t, join.Field, join.Field_foreign)
			case Join_inner:
				query.Inner_join(join.Table, t, join.Field, join.Field_foreign)
			case Join_cross:
				query.Cross_join(join.Table, t)
			}
		}
	}
}

func apply_raw_joins(query *sqlc.Select_query, joins Raw_joins){
	for _, join := range joins {
		switch join.Type {
		case Join_left:
			query.Left_join_multi(join.Table, join.Alias, join.On)
		case Join_inner:
			query.Inner_join_multi(join.Table, join.Alias, join.On)
		case Join_cross:
			query.Cross_join(join.Table, join.Alias)
		}
	}
}