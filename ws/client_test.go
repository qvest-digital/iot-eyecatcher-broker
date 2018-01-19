package ws

import (
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/tarent/iot-eyecatcher-broker/mocks/net"
	"reflect"
	"testing"
	"unsafe"
)

func TestSend(t *testing.T) {

	// given
	a := assert.New(t)

	// and
	ch := make(chan Message)

	testSubject := client{
		send: ch}

	// when
	res := testSubject.Send()

	// then
	a.Equal(ch, res)
}

func TestClose(t *testing.T) {

	// given
	a := assert.New(t)

	// and
	ch := make(chan Message)

	testSubject := client{
		send: ch}

	// when
	testSubject.Close()

	// then
	a.True(isChanClosed(ch))

}

func TestRemoteAddr(t *testing.T) {

	// given
	a := assert.New(t)
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	addrMock := net.NewMockAddr(ctrl)
	addrMock.EXPECT().String().Return("address")

	connMock := NewMockConnI(ctrl)
	connMock.EXPECT().RemoteAddr().Return(addrMock)

	testSubject := client{
		conn: connMock}

	// when
	addr := testSubject.RemoteAddr()

	// then
	a.Equal("address", addr)

}

// see http://stackoverflow.com/questions/16105325/how-to-check-a-channel-is-closed-or-not-without-reading-it
func isChanClosed(ch interface{}) bool {
	if reflect.TypeOf(ch).Kind() != reflect.Chan {
		panic("only channels!")
	}

	// get interface value pointer, from cgo_export
	// typedef struct { void *t; void *v; } GoInterface;
	// then get channel real pointer
	cptr := *(*uintptr)(unsafe.Pointer(
		unsafe.Pointer(uintptr(unsafe.Pointer(&ch)) + unsafe.Sizeof(uint(0))),
	))

	// this function will return true if chan.closed > 0
	// see hchan on https://github.com/golang/go/blob/master/src/runtime/chan.go
	// type hchan struct {
	// qcount   uint           // total data in the queue
	// dataqsiz uint           // size of the circular queue
	// buf      unsafe.Pointer // points to an array of dataqsiz elements
	// elemsize uint16
	// closed   uint32
	// **

	cptr += unsafe.Sizeof(uint(0)) * 2
	cptr += unsafe.Sizeof(unsafe.Pointer(uintptr(0)))
	cptr += unsafe.Sizeof(uint16(0))
	return *(*uint32)(unsafe.Pointer(cptr)) > 0
}
