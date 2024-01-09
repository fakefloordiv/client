package http1

import (
	"github.com/indigo-web/client/http"
	"github.com/indigo-web/client/http/method"
	"github.com/indigo-web/client/http/proto"
	"github.com/indigo-web/client/internal/protocol"
	"github.com/indigo-web/client/internal/tcp"
	"io"
	"os"
)

var _ protocol.Serializer = new(Serializer)

type Serializer struct {
	client tcp.Client
	buff   []byte
}

func NewSerializer(client tcp.Client, buff []byte) *Serializer {
	return &Serializer{
		client: client,
		buff:   buff,
	}
}

func (s *Serializer) Send(request *http.Request) error {
	s.method(request.Method)
	s.sp()
	s.path(request.Path)
	s.sp()
	s.proto(request.Proto)
	s.crlf()

	for _, pair := range request.Headers.Unwrap() {
		s.header(pair.Key, pair.Value)
		s.crlf()
	}

	s.crlf()

	if request.File != nil {
		return s.file(request.File)
	}

	s.buff = append(s.buff, request.Body...)

	return s.client.Write(s.buff)
}

func (s *Serializer) file(fd *os.File) error {
	// TODO: implement chunked streaming for files with size>N, where N tends to be more than 1mb
	content, err := io.ReadAll(fd)
	if err != nil {
		return err
	}

	s.buff = append(s.buff, content...)

	return s.client.Write(s.buff)
}

func (s *Serializer) method(m method.Method) {
	s.buff = append(s.buff, m...)
}

func (s *Serializer) sp() {
	s.buff = append(s.buff, ' ')
}

func (s *Serializer) path(path string) {
	s.buff = append(s.buff, path...)
}

func (s *Serializer) proto(proto proto.Protocol) {
	s.buff = append(s.buff, proto...)
}

func (s *Serializer) crlf() {
	s.buff = append(s.buff, '\r', '\n')
}

func (s *Serializer) header(key, value string) {
	s.buff = append(s.buff, key...)
	s.buff = append(s.buff, ':', ' ')
	s.buff = append(s.buff, value...)
}
