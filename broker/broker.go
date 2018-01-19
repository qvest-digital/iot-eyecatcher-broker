package broker

import (
	"errors"
	log "github.com/Sirupsen/logrus"
	"github.com/tarent/iot-eyecatcher-broker/json"
	"github.com/tarent/iot-eyecatcher-broker/ws"
	"time"
)

type queuedMessage struct {
	topic     string
	message   []byte
	timestamp time.Time
}

type broker struct {
	messageQueue chan queuedMessage
	hub          ws.Hub
	topics       TopicList
	marshaller   json.Marshaller
}

func NewBroker(hub ws.Hub) Broker {
	return &broker{
		messageQueue: make(chan queuedMessage),
		hub:          hub,
		topics:       NewTopicList(),
		marshaller:   json.NewMarshaller(),
	}
}

func (b *broker) Run() {
	for {
		message := <-b.messageQueue
		updatedTimestamp := b.topics.UpdateTopic(message.topic, message.message)
		j, err := b.marshaller.Marshal(message.topic, updatedTimestamp, message.message)
		if err != nil {
			log.Error("error marshalling json")
			continue
		}
		b.hub.Broadcast(ws.Message{Topic: message.topic, Text: j})
	}
}

func (b *broker) Message(topic string, message []byte) {
	b.messageQueue <- queuedMessage{topic: topic, message: message}
}

func (b *broker) LastMessage(topic string) ([]byte, error) {
	ts, message := b.topics.LastMessage(topic)
	if ts == time.Unix(0, 0) {
		return nil, errors.New("topic not found")
	}
	j, err := b.marshaller.Marshal(topic, ts, message)
	if err != nil {
		return nil, err
	}
	return j, nil
}
