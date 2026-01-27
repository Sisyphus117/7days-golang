package orm

import (
	"database/sql"
	"fmt"
	"orm/dialect"
	"orm/log"
	"orm/session"
	"strings"
)

type Engine struct {
	db      *sql.DB
	dialect dialect.Dialect
}

func NewEngine(driver, url string) (*Engine, error) {
	db, err := sql.Open(driver, url)
	if err != nil {
		log.Error(err)
		return nil, err
	}

	if err := db.Ping(); err != nil {
		log.Error(err)
		return nil, err
	}

	dial, ok := dialect.Get(driver)
	if !ok {
		log.Errorf("cannot find dialect of %s", driver)
		return nil, fmt.Errorf("cannot find dialect of %s", driver)
	}

	e := &Engine{db: db, dialect: dial}
	log.Info("Connect to database successfully")
	return e, nil
}

func (e *Engine) Close() {
	if err := e.db.Close(); err != nil {
		log.Error(err)
	}
	log.Info("database close successfully")
}

func (e *Engine) NewSession() *session.Session {
	return session.NewSession(e.db, e.dialect)
}

type TxFunc func(s *session.Session) (any, error)

func (e *Engine) Transaction(fn TxFunc) (res any, err error) {
	s := e.NewSession()
	if err = s.Begin(); err != nil {
		return
	}

	defer func() {
		if p := recover(); p != nil {
			s.RollBack()
			panic(p)
		} else if err != nil {
			s.RollBack()
		} else {
			err = s.Commit()
		}
	}()
	return fn(s)
}

func difference(a, b []string) []string {
	mp := make(map[string]struct{})
	fields := make([]string, 0)
	for _, field := range b {
		mp[field] = struct{}{}
	}
	for _, field := range a {
		if _, has := mp[field]; !has {
			fields = append(fields, field)
		}
	}
	return fields
}

func (e *Engine) Migrate(value any) error {
	_, err := e.Transaction(func(s *session.Session) (any, error) {
		if !s.Model(value).HasTable() {
			log.Infof("table %s doesn't exist", s.Schema().Name)
			return nil, s.CreateTable()
		}
		table := s.Schema()
		rows, _ := s.Raw(fmt.Sprintf("SELECT * FROM %s LIMIT 1", table.Name)).QueryRows()
		cols, _ := rows.Columns()
		addCols := difference(table.FieldNames(), cols)
		delCols := difference(cols, table.FieldNames())
		log.Infof("added cols:%v, deleted cols:%v", addCols, delCols)

		for _, col := range addCols {
			f := table.GetField(col)
			sql := fmt.Sprintf("ALTER TABLE %s ADD COLUMN %s %s;", table.Name, f.Name, f.Type)
			if _, err := s.Raw(sql).Exec(); err != nil {
				return nil, err
			}
		}

		if len(delCols) == 0 {
			return nil, nil
		}

		tmpName := "tmp_" + table.Name
		fields := strings.Join(table.FieldNames(), ", ")
		if _, err := s.Raw(fmt.Sprintf("CREATE TABLE %s AS SELECT %s FROM %s;", tmpName, fields, table.Name)).Exec(); err != nil {
			return nil, err
		}
		if _, err := s.Raw(fmt.Sprintf("DROP TABLE %s;", table.Name)).Exec(); err != nil {
			return nil, err
		}
		if _, err := s.Raw(fmt.Sprintf("ALTER TABLE %s RENAME TO %s;", tmpName, table.Name)).Exec(); err != nil {
			return nil, err
		}
		return nil, nil
	})
	return err
}
