package protocol

import "github.com/indigo-web/client/http"

type Parser interface {
	Parse([]byte) (headersCompleted bool, rest []byte, err error)
}

type Serializer interface {
	Send(r *http.Request) error
}

type Protocol interface {
	Parser
	Serializer
}

var _ Protocol = new(Impl)

type Impl struct {
	Parser
	Serializer
}
