package broker

import (
	log "github.com/Sirupsen/logrus"
	"sync"
	"time"
)

type topicList struct {
	topics map[string]Topic
	mutex  *sync.RWMutex
}

func NewTopicList() TopicList {
	return &topicList{make(map[string]Topic), &sync.RWMutex{}}
}

func (t *topicList) UpdateTopic(topic string, message []byte) time.Time {
	t.mutex.Lock()
	defer t.mutex.Unlock()
	_, exists := t.topics[topic]
	if !exists {
		log.WithField("topic", topic).Info("Creating new topic")
		t.topics[topic] = NewTopic(topic)
	}
	return t.topics[topic].UpdateLastMessage(message)
}

func (t *topicList) LastMessage(topic string) (time.Time, []byte) {
	t.mutex.RLock()
	defer t.mutex.RUnlock()
	_, exists := t.topics[topic]
	if !exists {
		return time.Unix(0, 0), nil
	}
	return t.topics[topic].GetLastMessage()
}
