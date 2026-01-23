package schema

import (
	"orm/dialect"
	"testing"
)

type User struct {
	Name string `orm:"PRIMARY KEY"`
	Age  int
}

var TestDial, _ = dialect.Get("sqlite3")

func TestParse(t *testing.T) {
	schema := Parse(TestDial, &User{})
	if schema.Name != "User" || schema.Len() != 2 {
		t.Fatal("failed to parse User struct")
	}
	if schema.GetField("Name").Tag != "PRIMARY KEY" {
		t.Fatal("failed to parse primary key")
	}
}
