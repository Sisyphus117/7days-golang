package frame

import (
	"net/http"
	"path"
	"text/template"
)

// router group works for middlewares attached to certain prefix
type RouterGroup struct {
	prefix      string
	middlewares []HandlerFunc
	parent      *RouterGroup
	engine      *Engine
}

// net server engine , CORE
type Engine struct {
	router    *Router
	templates *template.Template
	funcMap   template.FuncMap
	*RouterGroup
}

func NewEngine() *Engine {
	e := &Engine{
		RouterGroup: &RouterGroup{},
		router:      NewRouter(),
	}
	e.RouterGroup.engine = e
	e.RouterGroup.middlewares = []HandlerFunc{Recovery}
	return e
}

func NewGroup(name string, parent *RouterGroup) *RouterGroup {
	middlewaresCopy := make([]HandlerFunc, len(parent.middlewares))
	copy(middlewaresCopy, parent.middlewares)
	return &RouterGroup{
		prefix:      name,
		middlewares: middlewaresCopy,
		parent:      parent,
		engine:      parent.engine,
	}
}

// append middlewares to a router group
func (r *RouterGroup) AppendMiddlewares(middlewares ...HandlerFunc) {
	r.middlewares = append(r.middlewares, middlewares...)
}

// create a router group
func (r *RouterGroup) Group(name string) *RouterGroup {
	child := NewGroup(r.prefix+name, r)
	return child
}

// create a router on method GET
func (r *RouterGroup) Get(api string, handler HandlerFunc) {
	r.engine.router.addRoute("GET", r.prefix+api, append(r.middlewares, handler))
}

func (r *RouterGroup) Post(api string, handler HandlerFunc) {
	r.engine.router.addRoute("POST", r.prefix+api, append(r.middlewares, handler))
}

func (r *RouterGroup) Patch(api string, handler HandlerFunc) {
	r.engine.router.addRoute("PATCH", r.prefix+api, append(r.middlewares, handler))
}

func (r *RouterGroup) Delete(api string, handler HandlerFunc) {
	r.engine.router.addRoute("DELETE", r.prefix+api, append(r.middlewares, handler))
}

// run the net server
func (e *Engine) Run(port string) error {
	err := http.ListenAndServe(port, e)
	if err != nil {
		return err
	}
	return nil
}

// serveHTTP
func (e *Engine) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx := NewContext(w, r, e)
	e.router.handle(ctx)
}

func (r *RouterGroup) Static(relative, dir string) {
	handler := r.createStaticHandler(relative, http.Dir(dir))
	api := path.Join(relative, "/*filepath")
	r.Get(api, handler)
}

func (r *RouterGroup) createStaticHandler(relative string, fs http.FileSystem) HandlerFunc {
	path := r.prefix + relative
	fileServer := http.StripPrefix(path, http.FileServer(fs))
	return func(ctx *Context) {
		file := ctx.Param("filepath")
		if _, err := fs.Open(file); err != nil {
			ctx.SetStatus(http.StatusNotFound)
			return
		}
		fileServer.ServeHTTP(ctx.Writer, ctx.Req)
	}
}

func (e *Engine) SetFuncMap(funcMap template.FuncMap) {
	e.funcMap = funcMap
}

func (e *Engine) LoadHTMLGlob(path string) {
	e.templates = template.Must(template.New("").Funcs(e.funcMap).ParseGlob(path))
}
