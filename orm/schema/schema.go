package schema

import (
	"go/ast"
	"orm/dialect"
	"reflect"
)

type Field struct {
	Name string
	Type string
	Tag  string
}

type Schema struct {
	Name       string
	Model      any
	fieldsName []string
	fields     []*Field
	fieldsMap  map[string]*Field
}

func (s *Schema) GetField(name string) *Field {
	return s.fieldsMap[name]
}

func (s *Schema) Len() int {
	return len(s.fields)
}

func Parse(dialect dialect.Dialect, model any) *Schema {
	modelType := reflect.Indirect(reflect.ValueOf(model)).Type()
	s := &Schema{
		Model:      model,
		Name:       modelType.Name(),
		fieldsName: make([]string, 0),
		fields:     make([]*Field, 0),
		fieldsMap:  map[string]*Field{},
	}

	for i := range modelType.NumField() {
		p := modelType.Field(i)
		if !p.Anonymous && ast.IsExported(p.Name) {
			field := &Field{
				Name: p.Name,
				Type: dialect.TypeOf(reflect.Indirect(reflect.New(p.Type))),
			}
			if v, ok := p.Tag.Lookup("orm"); ok {
				field.Tag = v
			}
			s.fields = append(s.fields, field)
			s.fieldsName = append(s.fieldsName, p.Name)
			s.fieldsMap[p.Name] = field
		}
	}
	return s
}

func (s *Schema) Fields() []*Field {
	return s.fields
}

func (s *Schema) FieldsName() []string {
	return s.fieldsName
}

func (s *Schema) RecordValues(dest any) []any {
	values := reflect.Indirect(reflect.ValueOf(dest))
	fieldValues := make([]any, 0)
	for _, field := range s.fields {
		fieldValues = append(fieldValues, values.FieldByName(field.Name).Interface())
	}
	return fieldValues
}
