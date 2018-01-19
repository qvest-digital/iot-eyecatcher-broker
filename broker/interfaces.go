package broker

import "time"

type Broker interface {
	Run()
	Message(topic string, message []byte)
	LastMessage(topic string) ([]byte, error)
}

type Topic interface {
	UpdateLastMessage(message []byte) time.Time
	GetLastMessage() (time.Time, []byte)
}

type TopicList interface {
	UpdateTopic(topic string, message []byte) time.Time
	LastMessage(topic string) (time.Time, []byte)
}
