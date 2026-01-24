package session

import (
	"orm/clause"
	"reflect"
)

func (s *Session) Insert(values ...any) (int64, error) {
	recordValues := make([]any, 0)
	for _, val := range values {
		table := s.Model(val).Schema()
		s.clause.Set(clause.INSERT, table.Name, table.FieldsName())
		recordValues = append(recordValues, table.RecordValues(val))
	}

	s.clause.Set(clause.VALUES, recordValues...)
	sql, vars := s.clause.Build(clause.INSERT, clause.VALUES)
	res, err := s.Raw(sql, vars...).Exec()
	if err != nil {
		return 0, err
	}
	return res.RowsAffected()
}
func (s *Session) Find(values any) error {
	destSlice := reflect.Indirect(reflect.ValueOf(values))
	destType := destSlice.Type().Elem()
	table := s.Model(reflect.New(destType).Elem().Interface()).Schema()

	s.clause.Set(clause.SELECT, table.Name, table.FieldsName())
	sql, vars := s.clause.Build(clause.SELECT, clause.WHERE, clause.ORDERBY, clause.LIMIT)
	rows, err := s.Raw(sql, vars...).QueryRows()
	if err != nil {
		return err
	}

	for rows.Next() {
		dest := reflect.New(destType).Elem()
		var values []any
		for _, name := range table.FieldsName() {
			values = append(values, dest.FieldByName(name).Addr().Interface())
		}
		if err := rows.Scan(values...); err != nil {
			return err
		}
		destSlice.Set(reflect.Append(destSlice, dest))
	}
	return rows.Close()
}
