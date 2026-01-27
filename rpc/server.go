package rpc

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net"
	"reflect"
	"rpc/codec"
	"sync"
)

const MagicNumber = 0x3bef5c

type Option struct {
	MagicNumber int
	CodecType   codec.Type
}

var defaultOption = &Option{
	MagicNumber: MagicNumber,
	CodecType:   codec.Gob,
}

type Server struct {
}

func NewServer() *Server {
	return &Server{}
}

var defaultServer = NewServer()

var invalidRequest = struct{}{}

func (s *Server) serveConn(conn io.ReadWriteCloser) {
	defer func() {
		conn.Close()
	}()

	var opt Option
	if err := json.NewDecoder(conn).Decode(&opt); err != nil {
		log.Println("rpc server: options error: ", err)
		return
	}
	if opt.MagicNumber != MagicNumber {
		log.Println("rpc servedr: invalid magic number")
		return
	}
	fn := codec.NewCodecFuncMap[opt.CodecType]
	if fn == nil {
		log.Printf("rpc server: invalid codec type %s", opt.CodecType)
		return
	}
	s.serveCodec(fn(conn))
}

func (s *Server) serveCodec(cc codec.Codec) {
	sending := new(sync.Mutex)
	wg := new(sync.WaitGroup)
	for {
		req, err := s.readRequest(cc)
		if err != nil {
			if req == nil {
				break
			}
			req.h.Error = err.Error()
			s.sendResponse(cc, req.h, invalidRequest, sending)
			continue
		}
		wg.Add(1)
		go s.heandleRequest(cc, req, sending, wg)
	}
	wg.Wait()
	cc.Close()
}

type request struct {
	h            *codec.Header
	argv, replyv reflect.Value
}

func (s *Server) readRequestHeader(cc codec.Codec) (*codec.Header, error) {
	var h codec.Header
	if err := cc.ReadHeader(&h); err != nil {
		if err != io.EOF && err != io.ErrUnexpectedEOF {
			log.Println("rpc server: read header error: ", err)
		}
		return nil, err
	}
	return &h, nil
}

func (s *Server) readRequest(cc codec.Codec) (*request, error) {
	h, err := s.readRequestHeader(cc)
	if err != nil {
		return nil, err
	}
	req := &request{h: h}
	//TODO: now we don't know the type of request argv
	req.argv = reflect.New(reflect.TypeOf(""))
	if err = cc.ReadBody(req.argv.Interface()); err != nil {
		log.Println("rpc server: read argv err: ", err)
	}
	return req, nil
}

func (s *Server) sendResponse(cc codec.Codec, h *codec.Header, body any, sending *sync.Mutex) {
	sending.Lock()
	defer sending.Unlock()
	if err := cc.Write(h, body); err != nil {
		log.Println("rpc server: write response error: ", err)
	}
}

func (s *Server) heandleRequest(cc codec.Codec, req *request, sending *sync.Mutex, wg *sync.WaitGroup) {
	//TODO should call reigistered rpc methods to get the right replyv
	defer wg.Done()
	log.Println(req.h, req.argv.Elem())
	req.replyv = reflect.ValueOf(fmt.Sprintf("geerpc resp %d", req.h.Sequence))
	s.sendResponse(cc, req.h, req.replyv.Interface(), sending)
}

func (s *Server) Accept(lis net.Listener) {
	for {
		conn, err := lis.Accept()
		if err != nil {
			log.Println("rpc server: accept error: ", err)
			return
		}
		go s.serveConn(conn)
	}

}

func Accept(lis net.Listener) {
	defaultServer.Accept(lis)
}
