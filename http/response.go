package http

import (
	"github.com/indigo-web/client/http/headers"
	"github.com/indigo-web/client/http/proto"
	"github.com/indigo-web/client/http/status"
)

// headersPreAlloc defines a number of headers pairs, space for which will be
// pre-allocated. Don't ask me why exactly 7.
const headersPreAlloc = 7

type Response struct {
	Proto         proto.Protocol
	Code          status.Code
	Status        status.Status
	Headers       headers.Headers
	ContentLength int
	ContentType   string
	Encoding      Encoding
	Body          Body
}

func NewResponse(body Body) *Response {
	return &Response{
		Headers: headers.NewPreAlloc(headersPreAlloc),
		Body:    body,
	}
}

func (r *Response) Clear() {
	r.Code = 0
	r.Headers.Clear()
	r.ContentLength = 0
	r.Encoding = r.Encoding.Clear()
}
