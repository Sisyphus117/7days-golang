package codec

import (
	"bufio"
	"encoding/gob"
	"io"
	"log"
)

type GobCodec struct {
	conn    io.ReadWriteCloser
	buf     *bufio.Writer
	decoder *gob.Decoder
	encoder *gob.Encoder
}

func NewGobCodec(conn io.ReadWriteCloser) Codec {
	return &GobCodec{
		conn:    conn,
		buf:     bufio.NewWriter(conn),
		decoder: gob.NewDecoder(conn),
		encoder: gob.NewEncoder(conn),
	}
}

func (gc *GobCodec) Close() error {
	return gc.conn.Close()
}

func (gc *GobCodec) ReadHeader(header *Header) error {
	return gc.decoder.Decode(header)

}

func (gc *GobCodec) ReadBody(body any) error {
	return gc.decoder.Decode(body)
}

func (gc *GobCodec) Write(header *Header, body any) (err error) {
	defer func() {
		gc.buf.Flush()
		if err != nil {
			gc.conn.Close()
		}
	}()
	if err = gc.encoder.Encode(header); err != nil {
		log.Println("encode header failed")
		return
	}
	if err = gc.encoder.Encode(body); err != nil {
		log.Println("encode body failed")
		return
	}
	return
}
