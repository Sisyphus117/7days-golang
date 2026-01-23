package dialect

import "reflect"

type Dialect interface {
	TypeOf(val reflect.Value) string
	ExistTableSql(name string) (string, []any)
}

var dialectMap = make(map[string]Dialect)

func Set(name string, dialect Dialect) {
	dialectMap[name] = dialect
}

func Get(name string) (Dialect, bool) {
	dialect, has := dialectMap[name]
	return dialect, has
}
