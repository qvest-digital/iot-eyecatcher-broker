package json

import (
	"bytes"
	log "github.com/Sirupsen/logrus"
	"text/template"
	"time"
)

type JsonData struct {
	Topic     string
	Timestamp int64
	Message   string
}

type marshaller struct {
	jsonTemplate TemplateI
}

func NewMarshaller() Marshaller {
	return &marshaller{
		template.Must(template.New("jsonMessage").Parse(jsonTemplate))}
}

func (m *marshaller) Marshal(topic string, timestamp time.Time, message []byte) ([]byte, error) {
	data := JsonData{
		Topic:     topic,
		Timestamp: timestamp.Unix(),
		Message:   string(message),
	}

	var buf bytes.Buffer
	err := m.jsonTemplate.Execute(&buf, data)
	if err != nil {
		log.WithField("err", err).Error("error creating JSON")
		return nil, err
	}
	return buf.Bytes(), nil
}

const jsonTemplate = `{
    "topic": "{{.Topic}}",
    "timestamp": {{.Timestamp}},
    "message": {{.Message}}
}`
