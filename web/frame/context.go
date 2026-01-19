package frame

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

type Context struct {
	Writer     http.ResponseWriter
	Req        *http.Request
	Method     string
	Path       string
	Params     []string
	StatusCode int

	api string
}

func NewContext(w http.ResponseWriter, r *http.Request) *Context {
	return &Context{
		Writer: w,
		Req:    r,
		Method: r.Method,
		Path:   r.URL.Path,
	}
}

func (c *Context) SetHeader(key, value string) {
	c.Writer.Header().Set(key, value)
}

func (c *Context) SetStatus(code int) {
	c.StatusCode = code
	c.Writer.WriteHeader(code)
}

func (c *Context) UpdateParams(params *[]string) {
	c.Params = *params
}

func (c *Context) Printf(format string, a ...any) error {
	_, err := fmt.Fprintf(c.Writer, format, a...)
	if err != nil {
		return err
	}
	return nil
}

func (c *Context) HTML(code int, html string) {
	c.SetHeader("Context-Type", "text/html")
	c.Writer.Write([]byte(html))
}

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
