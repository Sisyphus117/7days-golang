package main

import (
	"fmt"
	"orm/orm"

	_ "github.com/mattn/go-sqlite3"
)

func main() {
	engine, _ := orm.NewEngine("sqlite3", "file:orm.db?_busy_timeout=5000&_journal_mode=DELETE")
	defer engine.Close()
	s := engine.NewSession()
	_, _ = s.Raw("DROP TABLE IF EXISTS User;").Exec()
	_, _ = s.Raw("CREATE TABLE User(Name text);").Exec()
	_, _ = s.Raw("CREATE TABLE User(Name text);").Exec()
	result, _ := s.Raw("INSERT INTO User(`Name`) values (?), (?)", "Tom", "Sam").Exec()
	count, _ := result.RowsAffected()
	fmt.Printf("Exec success, %d affected\n", count)

}
