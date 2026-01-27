package session

import (
	"orm/log"
	"reflect"
)

func (s *Session) CallMethod(method string, value any) {
	fn := reflect.ValueOf(s.Schema().Model).MethodByName(method)
	if value != nil {
		fn = reflect.ValueOf(value).MethodByName(method)
	}
	if fn.IsValid() {
		param := []reflect.Value{reflect.ValueOf(s)}
		if v := fn.Call(param); len(v) > 0 {
			if err, ok := v[0].Interface().(error); ok {
				log.Error(err)
			}
		}
	}
}
