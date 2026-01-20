package frame

import (
	"fmt"
	"log"
	"runtime"
	"strings"
)

func stackTrace(message string) string {
	var pcs [32]uintptr
	n := runtime.Callers(5, pcs[:])

	var str strings.Builder
	str.WriteString(message + "\nTraceback:")
	for _, pc := range pcs[:n] {
		fn := runtime.FuncForPC(pc)
		file, line := fn.FileLine(pc)
		fmt.Fprintf(&str, "\n\t%s:%d", file, line)
	}
	return str.String()
}

func Recovery(ctx *Context) {
	defer func() {
		if err := recover(); err != nil {
			ctx.Fail()
			tc := stackTrace(fmt.Sprintf("%s", err))
			log.Printf("%s", tc)
			fmt.Fprintf(ctx.Writer, "Internal serve error")
		}
	}()

	ctx.Next()
}
