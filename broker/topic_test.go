package broker

import (
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestNewTopic(t *testing.T) {

	// given
	a := assert.New(t)

	// when
	testSubject := NewTopic("testName")

	// then
	a.Equal("testName", testSubject.(*topic).name)
	a.Equal(time.Unix(0, 0), testSubject.(*topic).lastUpdated)
	a.Equal([]byte{}, testSubject.(*topic).lastMessage)
	a.NotNil(testSubject.(*topic).created)
}

func TestGetLastMessage(t *testing.T) {

	// given
	a := assert.New(t)

	// and
	testSubject := topic{
		lastUpdated: time.Unix(50, 0),
		lastMessage: []byte("hello")}

	// when
	ts, msg := testSubject.GetLastMessage()

	// then
	a.Equal([]byte("hello"), msg)
	a.Equal(time.Unix(50, 0), ts)
}

func TestGetLastMessageNoMessageAvailable(t *testing.T) {

	// given
	a := assert.New(t)

	// and
	testSubject := topic{
		lastUpdated: time.Unix(0, 0)}

	// when
	ts, msg := testSubject.GetLastMessage()

	// then
	a.Equal([]byte(nil), msg)
	a.Equal(time.Unix(0, 0), ts)
}

func TestUpdateLastMessage(t *testing.T) {

	// given
	a := assert.New(t)

	// and
	testSubject := topic{
		lastUpdated: time.Unix(50, 0),
		lastMessage: []byte("hello")}

	before := time.Now()

	// when
	ts := testSubject.UpdateLastMessage([]byte("yoo hoo"))

	// then
	a.True(ts.Before(time.Now()))
	a.True(ts.After(before))
	a.Equal([]byte("yoo hoo"), testSubject.lastMessage)
	a.Equal(ts, testSubject.lastUpdated)
}
