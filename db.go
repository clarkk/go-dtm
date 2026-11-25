package dtm

import (
	"database/sql"
	"github.com/clarkk/go-api"
	"github.com/clarkk/go-dbd/sqlc"
)

type DB_query struct {
	Env_context
	schema_base
	table_extension		string
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
	joins.apply(q)
	return q
}

func (s DB_query) Select_query_id_explicit(id uint64, fields []string, joins Explicit_joins) *sqlc.Select_query {
	q := sqlc.Select_id(s.Resource(), id).
		Select(fields)
	joins.apply(q)
	return q
}

func (s DB_query) Select_query(fields []string, joins Joins) *sqlc.Select_query {
	return s.Select_query_table(s.Resource(), fields, joins)
}

func (s DB_query) Select_query_explicit(fields []string, joins Explicit_joins) *sqlc.Select_query {
	return s.Select_query_table_explicit(s.Resource(), fields, joins)
}

func (s DB_query) Select_query_table(table string, fields []string, joins Joins) *sqlc.Select_query {
	q := sqlc.Select(table).
		Select(fields).
		Optimize_joins()
	joins.apply(q)
	return q
}

func (s DB_query) Select_query_table_explicit(table string, fields []string, joins Explicit_joins) *sqlc.Select_query {
	q := sqlc.Select(table).
		Select(fields)
	joins.apply(q)
	return q
}

func (s DB_query) Exec(q sqlc.SQL) (sql.Result, error){
	return s.tx.Exec(q)
}

func (s DB_query) Query_row(q sqlc.SQL, scan []any) (empty bool, err error){
	return s.tx.Query_row(q, scan)
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

func (s DB_query) Insert(fields sqlc.Map) (id uint64, err error){
	return s.tx.Insert(
		sqlc.Insert(s.Resource()).
		Fields(fields))
}

func (s DB_query) Insert_update_duplicate(fields sqlc.Map, update_fields []string) (id uint64, err error){
	return s.tx.Insert(
		sqlc.Insert(s.Resource()).
		Fields(fields).
		Update_duplicate(update_fields))
}

func (s DB_query) Insert_extension(value any, fields sqlc.Map) error {
	table, key := s.get_extension()
	fields[key] = value
	return s.tx.Insert_no_return(
		sqlc.Insert(table).
		Fields(fields))
}

func (s DB_query) Update_id(id uint64, fields sqlc.Map, where *sqlc.Where_clause) error {
	return s.tx.Update(
		sqlc.Update_id(s.Resource(), id).Fields(fields).
		Where(where))
}

func (s DB_query) Update(fields sqlc.Map, where *sqlc.Where_clause) error {
	return s.tx.Update(
		sqlc.Update(s.Resource()).Fields(fields).
		Where(where))
}

func (s DB_query) Update_extension(value any, fields sqlc.Map, where *sqlc.Where_clause) error {
	table, key := s.get_extension()
	return s.tx.Update(
		sqlc.Update(table).Fields(fields).
		Where(where.Eq(key, value)))
}

func (s DB_query) Delete_id(id uint64, where *sqlc.Where_clause) (empty bool, err error){
	return s.tx.Delete(
		sqlc.Delete_id(s.Resource(), id).
		Where(where))
}

func (s DB_query) Delete(where *sqlc.Where_clause) (empty bool, err error){
	return s.tx.Delete(
		sqlc.Delete(s.Resource()).
		Where(where))
}

func (s DB_query) Match_etag(id uint64, etag_header uint32, where *sqlc.Where_clause) (found bool, match bool, err error){
	q := s.Select_query_id(id, []string{"etag"}, nil).
		Where(where)
	
	var etag uint32
	empty, err := s.tx.Query_row(q, []any{&etag})
	if empty {
		err = nil
		return
	}
	if err != nil {
		return
	}
	
	found = true
	match = etag_header == etag
	return
}

func (s DB_query) Constraints(wheres []*sqlc.Where_clause, joins Joins) (exists bool, err error){
	for _, where := range wheres {
		exists, err = s.Exists(where, joins)
		if exists || err != nil {
			return
		}
	}
	return
}

func (s DB_query) Exists_id(id uint64, where *sqlc.Where_clause, joins Joins) (exists bool, err error){
	q := sqlc.Select_id(s.Resource(), id).
		Optimize_joins()
	joins.apply(q)
	return s.existence(q, where)
}

func (s DB_query) Exists_id_lock_for_update(id uint64, where *sqlc.Where_clause, joins Joins) (exists bool, err error){
	q := sqlc.Select_id(s.Resource(), id).
		Lock_for_update().
		Optimize_joins()
	joins.apply(q)
	return s.existence(q, where)
}

func (s DB_query) Exists(where *sqlc.Where_clause, joins Joins) (exists bool, err error){
	q := sqlc.Select(s.Resource()).
		Optimize_joins()
	joins.apply(q)
	return s.existence(q, where)
}

func (s DB_query) Get_id_exists(where *sqlc.Where_clause) (id uint64, err error){
	return s.get_id_exists(where, false)
}

func (s DB_query) Get_id_exists_lock_for_update(where *sqlc.Where_clause) (id uint64, err error){
	return s.get_id_exists(where, true)
}

func (s DB_query) Get_id_exists_omit(where *sqlc.Where_clause, omit_id uint64) (id uint64, err error){
	if omit_id != 0 {
		where.Not_eq("id", omit_id)
	}
	return s.Get_id_exists(where)
}

func (s DB_query) Update_limit_count(limit *api.Limit, where *sqlc.Where_clause) (out_of_range bool, err error){
	count, err := s.count(where)
	if err != nil {
		return false, err
	}
	return !limit.Count(count), nil
}

func (s DB_query) Get_ids(where *sqlc.Where_clause, lock_for_update bool) ([]uint64, error){
	query := s.Select_query([]string{"id"}, nil).
		Where(where)
	
	if lock_for_update {
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

func (s DB_query) get_id_exists(where *sqlc.Where_clause, lock_for_update bool) (id uint64, err error){
	q := s.Select_query([]string{"id"}, nil).
		Where(where).
		Limit(0, 1)
	
	if lock_for_update {
		q.Lock_for_update()
	}
	
	empty, err := s.tx.Query_row(q, []any{&id})
	if empty {
		err = nil
	}
	return
}

func (s DB_query) existence(q *sqlc.Select_query, where *sqlc.Where_clause) (exist bool, err error){
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

func (s DB_query) get_extension() (table string, key string){
	tbl := s.Resource()
	return tbl+"_"+s.table_extension, tbl+"_id"
}