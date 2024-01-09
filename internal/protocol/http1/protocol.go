package http1

import (
	"github.com/indigo-web/client/http"
	"github.com/indigo-web/client/internal/protocol"
	"github.com/indigo-web/client/internal/tcp"
	"github.com/indigo-web/utils/buffer"
)

type Protocol struct {
	parser *Parser
}

func New(
	resp *http.Response, respLineBuff, headersBuff buffer.Buffer,
	client tcp.Client, reqBuff []byte,
) *protocol.Impl {
	return &protocol.Impl{
		Parser:     NewParser(resp, respLineBuff, headersBuff),
		Serializer: NewSerializer(client, reqBuff),
	}
}
