package frame

import (
	"fmt"
)

type HandlerFunc func(ctx *Context) error

type Router map[string]HandlerFunc

func NewRouter() Router {
	return make(Router)
}

func (r Router) addRoute(method, api string, handler HandlerFunc) {
	key := fmt.Sprintf("%s-%s", api, method)
	r[key] = handler
}

func (r Router) handle(ctx *Context) error {
	key := fmt.Sprintf("%s-%s", ctx.Path, ctx.Method)
	if handler, has := r[key]; has {
		return handler(ctx)
	} else {
		return fmt.Errorf("not found")
	}

}
