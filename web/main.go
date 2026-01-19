package main

import (
	"log"
	"net/http"
	"web/frame"
)

func main() {
	e := frame.NewEngine()
	g1 := e.Group("/v1")
	g1.Get("/text/:name", func(ctx *frame.Context) {
		ctx.Printf("hello, world %s\n", ctx.Params[0])
	})
	e.Get("/", func(ctx *frame.Context) {
		ctx.HTML(http.StatusOK, "<h1>hello, world</h1>")
	})
	err := e.Run(":4000")
	if err != nil {
		log.Fatal(err)
	}
}
