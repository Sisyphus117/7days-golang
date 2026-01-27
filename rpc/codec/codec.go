package codec

import "io"

type Header struct {
	ServiceMethod string
	Sequence      int64
	Error         string
}

type Codec interface {
	io.Closer
	ReadHeader(*Header) error
	ReadBody(any) error
	Write(*Header, any) error
}

type Type string

type NewCodecFunc func(io.ReadWriteCloser) Codec

const (
	Gob  Type = "application/gob"
	Json Type = "application/json"
)

var NewCodecFuncMap map[Type]NewCodecFunc

func init() {
	NewCodecFuncMap = make(map[Type]NewCodecFunc)
	NewCodecFuncMap[Gob] = NewGobCodec
}
