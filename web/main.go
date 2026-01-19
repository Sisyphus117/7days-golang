package main

import (
	"log"
	"web/frame"
)

func main() {
	e := frame.NewEngine()
	e.Get("/", func(ctx *frame.Context) error {
		ctx.Printf("hello, world\n")
		return nil
	})
	err := e.Run(":4000")
	if err != nil {
		log.Fatal(err)
	}
}
