package session

import (
	"fmt"
	"orm/log"
	"orm/schema"
	"reflect"
	"strings"
)

func (s *Session) Model(val any) *Session {
	if s.schema == nil || reflect.TypeOf(val) != reflect.TypeOf(s.schema) {
		s.schema = schema.Parse(s.dialect, val)
	}
	return s
}

func (s *Session) Schema() *schema.Schema {
	if s.schema == nil {
		log.Error("Model is not set")
	}
	return s.schema
}

func (s *Session) CreateTable() error {
	table := s.Schema()
	columns := make([]string, 0)
	for _, col := range table.Fields() {
		columns = append(columns, fmt.Sprintf("%s %s %s", col.Name, col.Type, col.Tag))
	}

	dest := strings.Join(columns, ",")

	_, err := s.Raw(fmt.Sprintf("create table %s (%s);", table.Name, dest)).Exec()
	return err

}

func (s *Session) DropTable() error {
	schema := s.Schema()
	_, err := s.Raw(fmt.Sprintf("drop table if exists  %s ;", schema.Name)).Exec()
	return err
}

func (s *Session) HasTable() bool {
	sql, value := s.dialect.ExistTableSql(s.schema.Name)
	row := s.Raw(sql, value...).QueryRow()
	var tmp string
	_ = row.Scan(&tmp)
	return tmp == s.schema.Name
}
