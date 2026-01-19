package frame

import (
	"net/http"
)

type RouterGroup struct {
	prefix      string
	middlewares []HandlerFunc
	parent      *RouterGroup
	engine      *Engine
}

type Engine struct {
	groups map[string]*RouterGroup
	router *Router
	*RouterGroup
}

func NewEngine() *Engine {
	e := &Engine{
		RouterGroup: &RouterGroup{},
		groups:      make(map[string]*RouterGroup),
		router:      NewRouter(),
	}
	e.RouterGroup.engine = e
	return e
}

func NewGroup(name string, parent *RouterGroup) *RouterGroup {
	return &RouterGroup{
		prefix:      name,
		middlewares: make([]HandlerFunc, 0),
		parent:      parent,
		engine:      parent.engine,
	}
}

func (r *RouterGroup) Group(name string) *RouterGroup {
	child := NewGroup(r.prefix+name, r)
	e := r.engine
	e.groups[child.prefix] = child
	return child
}

func (r *RouterGroup) Get(api string, handler HandlerFunc) {
	r.engine.router.addRoute("GET", r.prefix+api, handler)
}

func (r *RouterGroup) Post(api string, handler HandlerFunc) {
	r.engine.router.addRoute("POST", r.prefix+api, handler)
}

func (r *RouterGroup) Patch(api string, handler HandlerFunc) {
	r.engine.router.addRoute("PATCH", r.prefix+api, handler)
}

func (r *RouterGroup) Delete(api string, handler HandlerFunc) {
	r.engine.router.addRoute("DELETE", r.prefix+api, handler)
}

func (e *Engine) Run(port string) error {
	err := http.ListenAndServe(port, e)
	if err != nil {
		return err
	}
	return nil
}

func (e *Engine) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx := NewContext(w, r)
	e.router.handle(ctx)
}
