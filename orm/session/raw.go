package session

import (
	"database/sql"
	"orm/clause"
	"orm/dialect"
	"orm/log"
	"orm/schema"
	"strings"
)

type Session struct {
	db      *sql.DB
	tx      *sql.Tx
	schema  *schema.Schema
	dialect dialect.Dialect
	clause  *clause.Clause
	sql     strings.Builder
	val     []any
}

type CommonDB interface {
	Exec(query string, args ...any) (sql.Result, error)
	QueryRow(query string, args ...any) *sql.Row
	Query(query string, args ...any) (*sql.Rows, error)
}

// lift from runtime check to compile check
var _ CommonDB = (*sql.DB)(nil)
var _ CommonDB = (*sql.Tx)(nil)

func NewSession(db *sql.DB, d dialect.Dialect) *Session {
	return &Session{db: db, dialect: d}
}

func (s *Session) Clear() {
	s.sql.Reset()
	s.val = nil
	s.clause = clause.Clear()
}

func (s *Session) DB() CommonDB {
	if s.tx != nil {
		return s.tx
	}
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
	res, err := s.DB().Exec(s.sql.String(), s.val...)
	if err != nil {
		log.Error(err)
	}
	return res, err
}

func (s *Session) QueryRow() *sql.Row {
	defer s.Clear()
	log.Info(s.sql.String(), s.val)
	row := s.DB().QueryRow(s.sql.String(), s.val...)
	return row
}

func (s *Session) QueryRows() (*sql.Rows, error) {
	defer s.Clear()
	log.Info(s.sql.String(), s.val)
	rows, err := s.DB().Query(s.sql.String(), s.val...)
	if err != nil {
		log.Error(err)
	}
	return rows, err
}
