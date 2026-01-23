package session

import (
	"database/sql"
	"orm/dialect"
	"orm/log"
	"orm/schema"
	"strings"
)

type Session struct {
	db      *sql.DB
	schema  *schema.Schema
	dialect dialect.Dialect
	sql     strings.Builder
	val     []any
}

func NewSession(db *sql.DB, d dialect.Dialect) *Session {
	return &Session{db: db, dialect: d}
}

func (s *Session) Clear() {
	s.sql.Reset()
	s.val = nil
}

func (s *Session) DB() *sql.DB {
	return s.db
}

func (s *Session) Raw(query string, values ...any) *Session {
	s.Clear()
	s.sql.WriteString(query)
	s.val = values
	return s
}

func (s *Session) Exec() (sql.Result, error) {
	defer s.Clear()
	log.Info(s.sql.String(), s.val)
	res, err := s.db.Exec(s.sql.String(), s.val...)
	if err != nil {
		log.Error(err)
	}
	return res, err
}

func (s *Session) QueryRow() *sql.Row {
	defer s.Clear()
	log.Info(s.sql.String(), s.val)
	row := s.db.QueryRow(s.sql.String(), s.val...)
	return row
}

func (s *Session) QueryRows() (*sql.Rows, error) {
	defer s.Clear()
	log.Info(s.sql.String(), s.val)
	rows, err := s.db.Query(s.sql.String(), s.val...)
	if err != nil {
		log.Error(err)
	}
	return rows, err
}
