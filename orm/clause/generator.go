package clause

import (
	"bytes"
	"fmt"
	"strings"
)

type generator func(...any) (string, []any)

var generators map[Type]generator

func init() {
	generators = make(map[Type]generator)
	generators[INSERT] = _insert
	generators[VALUES] = _values
	generators[SELECT] = _select
	generators[WHERE] = _where
	generators[LIMIT] = _limit
	generators[ORDERBY] = _orderby
}

func _insert(val ...any) (string, []any) {
	// INSERT INTO $tableName ($fields)
	table := val[0]
	items := strings.Join(val[1].([]string), ",")
	query := fmt.Sprintf("INSERT INTO %s (%s)", table, items)
	return query, []any{}
}

func genTamp(l int) string {
	vars := make([]string, l)
	for i := range l {
		vars[i] = "?"
	}
	return fmt.Sprintf("(%s)", strings.Join(vars, ","))
}

func _values(val ...any) (string, []any) {
	// VALUES ($v1), ($v2), ...
	var buf bytes.Buffer
	tamp := genTamp(len(val[0].([]any)))
	buf.WriteString("VALUES")
	vars := make([]any, 0)
	for i := range val {
		el := val[i].([]any)
		buf.WriteString(tamp)
		if i != len(val)-1 {
			buf.WriteByte(',')
		}
		vars = append(vars, el...)
	}
	return buf.String(), vars
}

func _select(val ...any) (string, []any) {
	// SELECT $fields FROM $tableName
	fields := strings.Join(val[1].([]string), ",")
	table := val[0]
	query := fmt.Sprintf("SELECT %s FROM %s", fields, table)
	return query, []any{}
}

func _where(val ...any) (string, []any) {
	// LIMIT $num
	query := fmt.Sprintf("WHERE %s", val[0])
	return query, val[1:]
}

func _limit(val ...any) (string, []any) {
	// WHERE $desc
	return "LIMIT ?", val
}

func _orderby(val ...any) (string, []any) {
	query := fmt.Sprintf("ORDER BY %s", val[0])
	return query, []any{}

}
