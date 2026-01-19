package frame

import (
	"log"
	"net/http"
)

type Engine struct {
	router Router
}

func NewEngine() *Engine {
	return &Engine{router: NewRouter()}
}

func (e *Engine) Get(api string, handler HandlerFunc) {
	e.router.addRoute("GET", api, handler)
}

func (e *Engine) Post(api string, handler HandlerFunc) {
	e.router.addRoute("POST", api, handler)
}

func (e *Engine) Patch(api string, handler HandlerFunc) {
	e.router.addRoute("PATCH", api, handler)
}

func (e *Engine) Delete(api string, handler HandlerFunc) {
	e.router.addRoute("DELETE", api, handler)
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
	err := e.router.handle(ctx)
	if err != nil {
		log.Fatal(err)
	}
}
