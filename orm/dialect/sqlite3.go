package dialect

import (
	"fmt"
	"reflect"
	"time"
)

type sqlite3 struct{}

func init() {
	Set("sqlite3", &sqlite3{})
}

var _ Dialect = (*sqlite3)(nil)

func (s *sqlite3) TypeOf(val reflect.Value) string {
	switch val.Kind() {
	case reflect.Int, reflect.Int16, reflect.Int32, reflect.Int8, reflect.Uint, reflect.Uint16, reflect.Uint32, reflect.Uint8:
		return "integer"
	case reflect.Int64, reflect.Uint64:
		return "bigint"
	case reflect.Float32, reflect.Float64:
		return "real"
	case reflect.Bool:
		return "bool"
	case reflect.Array, reflect.Slice:
		return "blob"
	case reflect.String:
		return "text"
	case reflect.Struct:
		if _, ok := val.Interface().(time.Time); ok {
			return "datetime"
		}
	}
	panic(fmt.Sprintf("unknown type%s (%s)", val.Type().Name(), val.Type().Kind()))
}

func (s *sqlite3) ExistTableSql(name string) (string, []any) {
	return "select name from sqlite_master where type='table' and name=?", []any{name}
}
