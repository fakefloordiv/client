package client

import (
	"github.com/indigo-web/chunkedbody"
	"github.com/indigo-web/client/http"
	"github.com/indigo-web/client/http/headers"
	"github.com/indigo-web/client/http/method"
	"github.com/indigo-web/client/internal/protocol"
	http1 "github.com/indigo-web/client/internal/protocol/http1"
	"github.com/indigo-web/client/internal/tcp"
	"github.com/indigo-web/utils/buffer"
	"net"
	"time"
)

const (
	readTimeout         = 90 * time.Second
	writeTimeout        = 90 * time.Second
	tcpBuffSize         = 4 * 1024
	respLineBuffInitial = 256
	respLineBuffMax     = 1024
	headersBuffInitial  = 2 * 1024
	headersBuffMax      = 32 * 1024
	renderBuffDefault   = 2 * 1024
	preAllocHeaders     = 10
)

type Session struct {
	client   tcp.Client
	protocol protocol.Protocol
	request  *http.Request
	response *http.Response
}

func NewSession(host string) (*Session, error) {
	conn, err := net.Dial("tcp", host)
	if err != nil {
		return nil, err
	}

	respLineBuff := buffer.New(respLineBuffInitial, respLineBuffMax)
	headersBuff := buffer.New(headersBuffInitial, headersBuffMax)
	buff := make([]byte, tcpBuffSize)
	client := tcp.NewClient(conn, readTimeout, writeTimeout, buff)
	body := http1.NewBody(client, chunkedbody.NewParser(chunkedbody.DefaultSettings()))
	resp := http.NewResponse(body)
	serializerBuff := make([]byte, 0, renderBuffDefault)

	return &Session{
		client:   client,
		protocol: http1.New(resp, *respLineBuff, *headersBuff, client, serializerBuff),
		request:  http.NewRequest(headers.NewPreAlloc(preAllocHeaders)),
		response: resp,
	}, nil
}

func (s *Session) Send(request *http.Request) (*http.Response, error) {
	if err := s.response.Body.Reset(); err != nil {
		return nil, err
	}

	if err := s.protocol.Send(request); err != nil {
		return nil, err
	}

	for {
		data, err := s.client.Read()
		if err != nil {
			return nil, err
		}

		headersCompleted, rest, err := s.protocol.Parse(data)
		if err != nil {
			// TODO: we should be more error-tolerant. Keep reading till the end (if the error isn't too hard)
			return nil, err
		}

		s.client.Unread(rest)

		if headersCompleted {
			s.response.Body.Init(s.response)

			return s.response, nil
		}
	}
}

func (s *Session) GET(path string) *http.Request {
	return s.request.WithMethod(method.GET).WithPath(path)
}

func (s *Session) HEAD(path string) *http.Request {
	return s.request.WithMethod(method.HEAD).WithPath(path)
}

func (s *Session) POST(path string) *http.Request {
	return s.request.WithMethod(method.POST).WithPath(path)
}

func (s *Session) PUT(path string) *http.Request {
	return s.request.WithMethod(method.PUT).WithPath(path)
}

func (s *Session) DELETE(path string) *http.Request {
	return s.request.WithMethod(method.DELETE).WithPath(path)
}

func (s *Session) CONNECT(path string) *http.Request {
	return s.request.WithMethod(method.CONNECT).WithPath(path)
}

func (s *Session) OPTIONS(path string) *http.Request {
	return s.request.WithMethod(method.OPTIONS).WithPath(path)
}

func (s *Session) TRACE(path string) *http.Request {
	return s.request.WithMethod(method.TRACE).WithPath(path)
}

func (s *Session) PATCH(path string) *http.Request {
	return s.request.WithMethod(method.PATCH).WithPath(path)
}
