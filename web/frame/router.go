package frame

import (
	"fmt"
	"net/http"
)

type HandlerFunc func(ctx *Context)

type Router struct {
	handlers map[string]HandlerFunc
	trie     *Trie
}

func NewRouter() *Router {
	return &Router{handlers: make(map[string]HandlerFunc), trie: NewTrie()}
}

func (r *Router) addRoute(method, api string, handler HandlerFunc) {
	err := r.trie.addKey(api)
	if err != nil {
		return
	}
	key := fmt.Sprintf("%s-%s", api, method)

	r.handlers[key] = handler
}

func (r *Router) parseRoute(ctx *Context) (string, error) {
	api, params, err := r.trie.getKey(ctx.Path)
	if err != nil {
		return "", err
	}
	ctx.UpdateParams(&params)

	return fmt.Sprintf("%s-%s", api, ctx.Method), nil
}

func (r *Router) handle(ctx *Context) {
	key, err := r.parseRoute(ctx)
	if err != nil {
		ctx.SetStatus(http.StatusNotFound)
		ctx.Printf("not found\n")
	}

	handler := r.handlers[key]
	handler(ctx)
}
