package session

import (
	"fmt"
	"orm/clause"
	"reflect"
)

func (s *Session) Insert(value ...any) (int64, error) {
	recordValues := make([]any, 0)
	for _, val := range value {
		table := s.Model(val).Schema()
		s.CallMethod("BeforeInsert", val)
		s.clause.Set(clause.INSERT, table.Name, table.FieldNames())
		recordValues = append(recordValues, table.RecordValues(val))
	}

	s.clause.Set(clause.VALUES, recordValues...)
	sql, vars := s.clause.Build(clause.INSERT, clause.VALUES)
	res, err := s.Raw(sql, vars...).Exec()
	if err != nil {
		return 0, err
	}
	for _, val := range vars {
		s.CallMethod("AfterInsert", val)
	}
	return res.RowsAffected()
}
func (s *Session) Find(value any) error {
	destSlice := reflect.Indirect(reflect.ValueOf(value))
	destType := destSlice.Type().Elem()
	table := s.Model(reflect.New(destType).Elem().Interface()).Schema()
	s.CallMethod("BeforeQuery", nil)

	s.clause.Set(clause.SELECT, table.Name, table.FieldNames())
	sql, vars := s.clause.Build(clause.SELECT, clause.WHERE, clause.ORDERBY, clause.LIMIT)
	rows, err := s.Raw(sql, vars...).QueryRows()
	if err != nil {
		return err
	}

	for rows.Next() {
		dest := reflect.New(destType).Elem()
		var values []any
		for _, name := range table.FieldNames() {
			values = append(values, dest.FieldByName(name).Addr().Interface())
		}
		if err := rows.Scan(values...); err != nil {
			return err
		}

		s.CallMethod("AfterQuery", dest.Addr().Interface())
		destSlice.Set(reflect.Append(destSlice, dest))
	}
	return rows.Close()
}

func (s *Session) Update(value ...any) (int64, error) {
	mp, ok := value[0].(map[string]any)
	if !ok {
		mp = make(map[string]any)
		n := len(value) / 2
		for i := range n {
			mp[value[i].(string)] = value[i+1]
		}
	}
	s.CallMethod("BeforeUpdate", mp)
	s.clause.Set(clause.UPDATE, s.Schema().Name, mp)
	sql, vals := s.clause.Build(clause.UPDATE, clause.WHERE)
	res, err := s.Raw(sql, vals...).Exec()
	if err != nil {
		return 0, err
	}
	for _, val := range vals {
		s.CallMethod("AfterUpdate", val)
	}
	return res.RowsAffected()
}
func (s *Session) Delete() (int64, error) {
	s.CallMethod("BeforeDelete", nil)
	s.clause.Set(clause.DELETE, s.Schema().Name)
	sql, vals := s.clause.Build(clause.DELETE, clause.WHERE)
	res, err := s.Raw(sql, vals...).Exec()
	if err != nil {
		return 0, err
	}
	for _, val := range vals {
		s.CallMethod("AfterDelete", val)
	}
	return res.RowsAffected()
}
func (s *Session) Count() (int64, error) {
	s.clause.Set(clause.COUNT, s.Schema().Name)
	sql, vals := s.clause.Build(clause.COUNT, clause.WHERE)
	row := s.Raw(sql, vals...).QueryRow()
	var res int64
	if err := row.Scan(&res); err != nil {
		return 0, err
	}
	return res, nil
}

func (s *Session) Where(desc string, args ...any) *Session {
	var vars []any
	vars = append(vars, desc)
	vars = append(vars, args...)
	s.clause.Set(clause.WHERE, vars...)
	return s
}

func (s *Session) Limit(limit int) *Session {
	s.clause.Set(clause.LIMIT, limit)
	return s
}

func (s *Session) OrderBy(order string) *Session {
	s.clause.Set(clause.ORDERBY, order)
	return s
}

func (s *Session) First(value any) error {
	dest := reflect.Indirect(reflect.ValueOf(value))
	destSlice := reflect.New(reflect.SliceOf(dest.Type())).Elem()

	if err := s.Limit(1).Find(destSlice.Addr().Interface()); err != nil {
		return err
	}

	if destSlice.Len() == 0 {
		return fmt.Errorf("NOT FOUND")
	}
	dest.Set(destSlice.Index(0))
	return nil
}
