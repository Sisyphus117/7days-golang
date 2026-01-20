package frame

import (
	"fmt"
	"net/http"
)

// handler
type HandlerFunc func(ctx *Context)

// router
type Router struct {
	handlers map[string][]HandlerFunc
	trie     *Trie
}

func NewRouter() *Router {
	return &Router{handlers: make(map[string][]HandlerFunc), trie: NewTrie()}
}

// create a route with method,api,handlers provided
func (r *Router) addRoute(method, api string, handlers []HandlerFunc) {
	api, err := r.trie.addKey(api)
	if err != nil {
		return
	}
	key := fmt.Sprintf("%s-%s", api, method)

	r.handlers[key] = handlers
}

// parse route from context,and update params to context, returns API-METHOD
func (r *Router) parseRoute(ctx *Context) (string, error) {
	api, params, err := r.trie.getKey(ctx.Path)
	if err != nil {
		return "", err
	}
	ctx.UpdateParams(&params)

	return fmt.Sprintf("%s-%s", api, ctx.Method), nil
}

// handle the request due to provided context
func (r *Router) handle(ctx *Context) {
	key, err := r.parseRoute(ctx)
	if err != nil {
		ctx.SetStatus(http.StatusNotFound)
		ctx.Printf("not found\n")
	}

	ctx.SetHandlers(r.handlers[key])
	ctx.Next()
}
