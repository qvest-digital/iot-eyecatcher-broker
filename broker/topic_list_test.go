package broker

import (
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"sync"
	"testing"
	"time"
)

func TestNewTopicList(t *testing.T) {

	// given
	a := assert.New(t)

	// when
	testSubject := NewTopicList()

	// then
	a.NotNil(testSubject.(*topicList).topics)
	a.NotNil(testSubject.(*topicList).mutex)
}

func TestUpdateTopicNewTopic(t *testing.T) {

	// given
	a := assert.New(t)
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// and
	testSubject := topicList{
		topics: map[string]Topic{},
		mutex:  &sync.RWMutex{}}

	// when
	before := time.Now()
	ts := testSubject.UpdateTopic("testTopic", []byte("hello"))

	// then
	a.True(before.Before(ts))
	a.True(time.Now().After(ts))
	a.Len(testSubject.topics, 1)
	a.Equal([]byte("hello"), testSubject.topics["testTopic"].(*topic).lastMessage)
	a.Equal("testTopic", testSubject.topics["testTopic"].(*topic).name)
}

func TestUpdateTopicExistingTopic(t *testing.T) {

	// given
	a := assert.New(t)
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// and
	topicMock := NewMockTopic(ctrl)
	topicMock.EXPECT().UpdateLastMessage([]byte("yoo hoo")).Return(time.Unix(50, 0))

	// and
	testSubject := topicList{
		topics: map[string]Topic{"testTopic": topicMock},
		mutex:  &sync.RWMutex{}}

	// when
	ts := testSubject.UpdateTopic("testTopic", []byte("yoo hoo"))

	// then
	a.Equal(time.Unix(50, 0), ts)
}

func TestLastMessage(t *testing.T) {

	// given
	a := assert.New(t)
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// and
	topicMock := NewMockTopic(ctrl)
	topicMock.EXPECT().GetLastMessage().Return(time.Unix(50, 0), []byte("hello"))

	// and
	testSubject := topicList{
		topics: map[string]Topic{"testTopic": topicMock},
		mutex:  &sync.RWMutex{}}

	// when
	ts, msg := testSubject.LastMessage("testTopic")

	// then
	a.Equal(time.Unix(50, 0), ts)
	a.Equal([]byte("hello"), msg)
}

func TestLastMessageNonexistingTopic(t *testing.T) {

	// given
	a := assert.New(t)

	// and
	testSubject := topicList{
		topics: map[string]Topic{},
		mutex:  &sync.RWMutex{}}

	// when
	ts, msg := testSubject.LastMessage("testTopic")

	// then
	a.Equal(time.Unix(0, 0), ts)
	a.Equal([]byte(nil), msg)
}
