package main

import (
	lru "cache/rlu"
	"cache/wrapped"
	"fmt"
)

func main() {
	r := lru.NewLRU(30)
	r.Add("ab", wrapped.Int(2))
	r.Add("ccc", wrapped.String("asgb"))
	r.Add("asmm", wrapped.Slice[int]([]int{1, 3, 5, 6}))

	fmt.Println(r.Get("ab"))
	fmt.Println(r.Get("asmm"))
}
