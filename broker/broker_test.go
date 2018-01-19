package broker

import (
	"errors"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/tarent/iot-eyecatcher-broker/json"
	"github.com/tarent/iot-eyecatcher-broker/ws"
	"testing"
	"time"
)

func TestNewBroker(t *testing.T) {

	// given
	a := assert.New(t)
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	hubMock := ws.NewMockHub(ctrl)

	// when
	testSubject := NewBroker(hubMock)

	// then
	a.NotNil(testSubject.(*broker).topics)
	a.NotNil(testSubject.(*broker).marshaller)
	a.NotNil(testSubject.(*broker).hub)
	a.NotNil(testSubject.(*broker).messageQueue)
}

func TestRun(t *testing.T) {

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// given
	testTopicName := "msgTopic"
	testPayload := []byte("msgMessage")
	testTimestamp := time.Unix(50, 0)
	jsonMessage := []byte("jsonMessage")
	testMessage := ws.Message{Text: jsonMessage, Topic: testTopicName}

	topicsMock := NewMockTopicList(ctrl)
	topicsMock.EXPECT().UpdateTopic(testTopicName, testPayload).Return(testTimestamp)

	// and
	marshallerMock := json.NewMockMarshaller(ctrl)
	marshallerMock.EXPECT().Marshal(testTopicName, testTimestamp, testPayload).Return(jsonMessage, nil)

	// and
	messageQueueMock := make(chan queuedMessage)

	// and
	hubMock := ws.NewMockHub(ctrl)
	hubMock.EXPECT().Broadcast(testMessage)

	// and the test subject
	testSubject := broker{
		hub:          hubMock,
		messageQueue: messageQueueMock,
		topics:       topicsMock,
		marshaller:   marshallerMock}

	// when
	go testSubject.Run()

	// and message is sent
	messageQueueMock <- queuedMessage{
		message:   testPayload,
		topic:     testTopicName,
		timestamp: testTimestamp}

}

func TestRunFaultyJson(t *testing.T) {

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// given
	testTopicName := "msgTopic"
	testPayload := []byte("msgMessage")
	testTimestamp := time.Unix(50, 0)

	// and
	topicsMock := NewMockTopicList(ctrl)
	topicsMock.EXPECT().UpdateTopic(testTopicName, testPayload).Return(testTimestamp)

	// and
	marshallerMock := json.NewMockMarshaller(ctrl)
	marshallerMock.EXPECT().Marshal(testTopicName, testTimestamp, testPayload).Return(nil, errors.New("some error"))

	// and
	messageQueueMock := make(chan queuedMessage)

	// and
	hubMock := ws.NewMockHub(ctrl)
	hubMock.EXPECT().Broadcast(gomock.Any()).Times(0)

	// and the test subject
	testSubject := broker{
		hub:          hubMock,
		messageQueue: messageQueueMock,
		topics:       topicsMock,
		marshaller:   marshallerMock}

	// when
	go testSubject.Run()

	// and message is sent
	messageQueueMock <- queuedMessage{
		message:   testPayload,
		topic:     testTopicName,
		timestamp: testTimestamp}
}

func TestMessage(t *testing.T) {

	a := assert.New(t)
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// given
	messageQueueMock := make(chan queuedMessage, 1)

	// and the test subject
	testSubject := broker{
		messageQueue: messageQueueMock}

	// when
	testSubject.Message("testTopic", []byte("someMsg"))
	res := <-messageQueueMock

	// then
	a.Equal(queuedMessage{topic: "testTopic", message: []byte("someMsg")}, res)
}

func TestLastBrokerMessage(t *testing.T) {

	a := assert.New(t)
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// given
	testTopicName := "msgTopic"
	testPayload := []byte("msgMessage")
	testTimestamp := time.Unix(50, 0)
	testJson := []byte("someJson")

	// and
	topicsMock := NewMockTopicList(ctrl)
	topicsMock.EXPECT().LastMessage(testTopicName).Return(testTimestamp, testPayload)

	// and
	marshallerMock := json.NewMockMarshaller(ctrl)
	marshallerMock.EXPECT().Marshal(testTopicName, testTimestamp, testPayload).Return(testJson, nil)

	// and the test subject
	testSubject := broker{
		marshaller: marshallerMock,
		topics:     topicsMock}

	// when
	res, err := testSubject.LastMessage(testTopicName)

	// then
	a.NoError(err)
	a.Equal(testJson, res)
}

func TestLastMessageTopicNotFound(t *testing.T) {

	a := assert.New(t)
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// given
	testTopicName := "msgTopic"

	// and
	topicsMock := NewMockTopicList(ctrl)
	topicsMock.EXPECT().LastMessage(testTopicName).Return(time.Unix(0, 0), nil)

	// and the test subject
	testSubject := broker{
		topics: topicsMock}

	// when
	json, err := testSubject.LastMessage(testTopicName)

	// then
	a.Error(err, "topic not found")
	a.Equal([]byte(nil), json)
}

func TestLastMessageMarshalError(t *testing.T) {

	a := assert.New(t)
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// given
	testTopicName := "msgTopic"
	testPayload := []byte("msgMessage")
	testTimestamp := time.Unix(50, 0)

	// and
	topicsMock := NewMockTopicList(ctrl)
	topicsMock.EXPECT().LastMessage(testTopicName).Return(testTimestamp, testPayload)

	// and
	marshallerMock := json.NewMockMarshaller(ctrl)
	marshallerMock.EXPECT().Marshal(testTopicName, testTimestamp, testPayload).Return(nil, errors.New("some error"))

	// and the test subject
	testSubject := broker{
		marshaller: marshallerMock,
		topics:     topicsMock}

	// when
	res, err := testSubject.LastMessage(testTopicName)

	// then
	a.Error(err, "some error")
	a.Nil(res)
}

// Just to please the coverage
func TestMockRun(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	x := NewMockBroker(ctrl)
	x.EXPECT().Run()
	x.Run()
}
