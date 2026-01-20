package wrapped

import (
	"bytes"
	"fmt"
)

type String string

func (s String) Len() int {
	return len(s)
}

func (s String) String() string {
	return string(s)
}

type Int int

func (i Int) Len() int {
	return 8
}

func (i Int) String() string {
	return fmt.Sprintf("%d", int(i))
}

type Uint uint

func (u Uint) Len() int {
	return 8
}

func (u Uint) String() string {
	return fmt.Sprintf("%d", uint(u))
}

type Slice[T any] []T

func (s Slice[T]) Len() int {
	return len(s)
}

func (s Slice[T]) String() string {
	var buf bytes.Buffer
	buf.WriteByte('[')
	for i, elem := range s {
		fmt.Fprintf(&buf, "%v", elem)
		if i != len(s)-1 {
			buf.WriteByte(',')
		}
	}
	buf.WriteByte(']')
	return buf.String()
}

type ByteView struct {
	b []byte
}

func (v ByteView) Len() int {
	return len(v.b)
}

func (v ByteView) String() string {
	return string(v.b)
}

func (v ByteView) ByteSlice() []byte {
	return bytes.Clone(v.b)
}
