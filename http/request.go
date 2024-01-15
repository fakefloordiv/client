package http

import (
	"github.com/indigo-web/client/http/headers"
	"github.com/indigo-web/client/http/method"
	"github.com/indigo-web/client/http/proto"
	"github.com/indigo-web/utils/uf"
	"os"
)

type session interface {
	Send(*Request) (*Response, error)
}

type Request struct {
	method  method.Method
	path    string
	proto   proto.Protocol
	headers headers.Headers
	file    *os.File
	body    []byte
	err     error
}

func NewRequest(hdrs headers.Headers) *Request {
	return &Request{
		proto:   proto.Auto,
		headers: hdrs,
	}
}

func (r *Request) Method(m method.Method) *Request {
	r.method = m
	return r
}

func (r *Request) Path(path string) *Request {
	r.path = path
	return r
}

func (r *Request) Proto(proto proto.Protocol) *Request {
	r.proto = proto
	return r
}

func (r *Request) Header(key string, values ...string) *Request {
	for _, value := range values {
		r.headers.Add(key, value)
	}
	return r
}

// WithFile opens a new file with os.O_RDONLY flag and perm=0
func (r *Request) File(filename string) *Request {
	r.file, r.err = os.OpenFile(filename, os.O_RDONLY, 0)
	return r
}

func (r *Request) String(body string) *Request {
	return r.Bytes(uf.S2B(body))
}

func (r *Request) Bytes(body []byte) *Request {
	r.body = body
	return r
}

// Error returns error, if occurred during request building. This may be caused
// by non-existing filename, passed via File, or BodyFrom, if error occurred during
// reading from it
func (r *Request) Error() error {
	return r.err
}

func (r *Request) Send(session session) (*Response, error) {
	return session.Send(r)
}

func (r *Request) Clear() *Request {
	r.method = method.Unknown
	r.path = ""
	r.proto = proto.Auto
	r.headers.Clear()
	r.file = nil
	r.body = nil
	r.err = nil
	return r
}
