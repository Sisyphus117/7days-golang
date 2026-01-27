package clause

import "strings"

type Type int

const (
	INSERT Type = iota
	UPDATE
	DELETE
	VALUES
	SELECT
	WHERE
	LIMIT
	ORDERBY
	COUNT
)

type Clause struct {
	sql    map[Type]string
	sqlVar map[Type][]any
}

func NewClause() *Clause {
	return &Clause{
		sql:    make(map[Type]string),
		sqlVar: make(map[Type][]any),
	}
}

func (c *Clause) Set(ty Type, vars ...any) {
	sql, sqlVars := generators[ty](vars...)
	c.sql[ty] = sql
	c.sqlVar[ty] = sqlVars
}

func (c *Clause) Build(orders ...Type) (string, []any) {
	var sql strings.Builder
	var vars []any
	for i, ty := range orders {
		sql.WriteString(c.sql[ty])
		vars = append(vars, c.sqlVar[ty]...)
		if i != len(orders)-1 {
			sql.WriteByte(' ')
		}
	}

	return sql.String(), vars
}

func Clear() *Clause {
	return NewClause()
}
