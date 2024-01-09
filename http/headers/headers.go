package headers

import (
	"github.com/indigo-web/client/internal/keyvalue"
)

type (
	Header  = keyvalue.Pair
	Headers = *keyvalue.Storage
)

func New() Headers {
	return keyvalue.New()
}

func NewPreAlloc(n int) Headers {
	return keyvalue.NewPreAlloc(n)
}

func NewFromMap(m map[string][]string) Headers {
	return keyvalue.NewFromMap(m)
}
