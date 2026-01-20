package frame

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

// context of a request
type Context struct {
	Writer     http.ResponseWriter
	Req        *http.Request
	Method     string
	Path       string
	StatusCode int
	//params parsed from path
	Params map[string]string
	//middlewares wrapped and handler, in call order
	handlers []HandlerFunc
	index    int
	api      string
	engine   *Engine
}

func NewContext(w http.ResponseWriter, r *http.Request, e *Engine) *Context {
	return &Context{
		Writer: w,
		Req:    r,
		Method: r.Method,
		Path:   r.URL.Path,
		engine: e,
	}
}

func (c *Context) SetHeader(key, value string) {
	c.Writer.Header().Set(key, value)
}

// set the statuscode
func (c *Context) SetStatus(code int) {
	c.StatusCode = code
	c.Writer.WriteHeader(code)
}

// get a param
func (c *Context) Param(key string) string {
	return c.Params[key]
}

// replace all params
func (c *Context) UpdateParams(params *map[string]string) {
	c.Params = *params
}

func (c *Context) Printf(format string, a ...any) error {
	_, err := fmt.Fprintf(c.Writer, format, a...)
	if err != nil {
		return err
	}
	return nil
}

// html context
func (c *Context) HTML(code int, name string, data any) {
	c.SetHeader("Context-Type", "text/html")
	if err := c.engine.templates.ExecuteTemplate(c.Writer, name, data); err != nil {
		c.SetStatus(http.StatusInternalServerError)
		return
	}
}

// json context
func (c *Context) JSON(code int, obj any) {
	c.SetHeader("Context-Type", "text/json")

	var buf bytes.Buffer
	encoder := json.NewEncoder(&buf)
	err := encoder.Encode(obj)
	if err != nil {
		c.SetHeader("Context-Type", "text/plain")
		c.SetStatus(http.StatusInternalServerError)
		fmt.Fprintf(c.Writer, "json encoding error: %v", err)
		return
	}
	c.SetStatus(code)
	c.Writer.Write(buf.Bytes())
}

// call the next middleware or handler on the chain, nonblock
func (c *Context) Next() {
	c.index++
	for c.index < len(c.handlers) {
		c.handlers[c.index](c)
		c.index++
	}
}

// set handlers and set the call index to -1
func (c *Context) SetHandlers(handlers []HandlerFunc) {
	c.handlers = handlers
	c.index = -1
}
