package ws

import (
	"io"
	"net"
	"time"
)

type Client interface {
	Send() chan Message
	Close()
	RemoteAddr() string
}

type Hub interface {
	Run()
	Unregister(c Client)
	Register(c Client)
	Broadcast(message Message)
}

type ConnI interface {
	RemoteAddr() net.Addr
	Close() error
	SetReadLimit(int64)
	SetReadDeadline(time.Time) error
	SetWriteDeadline(time.Time) error
	SetPongHandler(func(string) error)
	ReadMessage() (int, []byte, error)
	WriteMessage(int, []byte) error
	NextWriter(int) (io.WriteCloser, error)
}
