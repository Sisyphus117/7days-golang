package orm

import (
	"database/sql"
	"fmt"
	"orm/dialect"
	"orm/log"
	"orm/session"
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
