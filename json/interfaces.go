package json

import (
	"io"
	"time"
)

type Marshaller interface {
	Marshal(topic string, timestamp time.Time, message []byte) ([]byte, error)
}

type TemplateI interface {
	Execute(io.Writer, interface{}) error
}
