package main

import (
	"fmt"
	"log"
	"net/http"
	"text/template"
	"time"
	"web/frame"
)

func FormatAsDate(t time.Time) string {
	year, month, day := t.Date()
	return fmt.Sprintf("%d-%02d-%02d", year, month, day)
}

func main() {
	e := frame.NewEngine()
	e.Static("/static", "./static")
	e.SetFuncMap(template.FuncMap{"FormatAsDate": FormatAsDate})
	e.LoadHTMLGlob("./template/*")

	e.Get("/", func(ctx *frame.Context) {
		ctx.HTML(http.StatusOK, "css.tmpl", nil)
	})
	e.Get("/date", func(ctx *frame.Context) {
		ctx.HTML(http.StatusOK, "custom_func.tmpl", map[string]any{
			"title": "date",
			"now":   time.Now(),
		})
	})

	type student struct {
		Name string
		Age  int
	}
	e.Get("/students", func(ctx *frame.Context) {
		ctx.HTML(http.StatusOK, "arr.tmpl", map[string]any{
			"title":  "students",
			"stuArr": []*student{{"ming", 12}, {"yuan", 15}},
		})
	})

	g1 := e.Group("/v1")
	g1.AppendMiddlewares(frame.SimpleMiddleware)
	g1.Get("/text/:name", func(ctx *frame.Context) {
		ctx.Printf("hello, world %s\n", ctx.Param("name"))
	})

	err := e.Run(":4000")
	if err != nil {
		log.Fatal(err)
	}
}
