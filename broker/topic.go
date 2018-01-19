package broker

import "time"

type topic struct {
	name        string
	lastMessage []byte
	created     time.Time
	lastUpdated time.Time
}

func NewTopic(name string) Topic {
	return &topic{
		name:        name,
		lastMessage: make([]byte, 0),
		created:     time.Now(),
		lastUpdated: time.Unix(0, 0),
	}
}

func (t *topic) GetLastMessage() (time.Time, []byte) {
	if t.lastUpdated != time.Unix(0, 0) {
		return t.lastUpdated, t.lastMessage
	}
	return t.lastUpdated, nil
}

func (t *topic) UpdateLastMessage(message []byte) time.Time {
	t.lastMessage = message
	t.lastUpdated = time.Now()
	return t.lastUpdated
}
