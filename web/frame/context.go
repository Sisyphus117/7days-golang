package frame

import (
	"fmt"
	"net/http"
)

type Context struct {
	Writer     http.ResponseWriter
	Req        *http.Request
	Method     string
	Path       string
	StatusCode int
}

func NewContext(w http.ResponseWriter, r *http.Request) *Context {
	return &Context{
		Writer: w,
		Req:    r,
		Method: r.Method,
		Path:   r.URL.Path,
	}
}

func (c *Context) Printf(format string, a ...any) error {
	_, err := fmt.Fprintf(c.Writer, format, a...)
	if err != nil {
		return err
	}
	return nil
}
