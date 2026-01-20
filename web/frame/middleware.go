package frame

import "fmt"

// test if the wrap before and after works
func SimpleMiddleware(ctx *Context) {
	fmt.Println("before handling")
	ctx.Next()
	fmt.Println("after handling")

}
