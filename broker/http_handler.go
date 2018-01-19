package broker

import (
	log "github.com/Sirupsen/logrus"
	"github.com/gorilla/mux"
	"io/ioutil"
	"net/http"
)

// Send new message
func MessageHandler(broker Broker, w http.ResponseWriter, r *http.Request) {

	l := log.WithField("remoteAddr", r.RemoteAddr)

	w.Header().Add("Access-Control-Allow-Origin", "*")
	w.Header().Add("Access-Control-Allow-Headers", "content-type,x-authorization")

	message, err := ioutil.ReadAll(r.Body)
	if err != nil {
		l.WithField("err", err).Error("error reading http request body")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	vars := mux.Vars(r)
	topic := vars["topic"]

	if len(topic) < 1 {
		l.Error("missing topic")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	l = l.WithField("topic", topic)

	if len(message) < 1 {
		l.Error("missing message")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	l = l.WithField("messageLength", len(message))

	l.Info("Received new message")

	broker.Message(topic, message)

	w.WriteHeader(http.StatusCreated)
}

// Get newest message in topic
func GetLastTopicMessageHandler(broker Broker, w http.ResponseWriter, r *http.Request) {

	l := log.WithField("remoteAddr", r.RemoteAddr)

	vars := mux.Vars(r)
	topic := vars["topic"]

	if len(topic) < 1 {
		l.Error("missing topic")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	l = l.WithField("topic", topic)

	l.Info("Requested last topic message")

	msg, err := broker.LastMessage(topic)
	if err != nil {
		if err.Error() == "topic not found" {
			w.WriteHeader(http.StatusNotFound)
			return
		}
		w.WriteHeader(http.StatusInternalServerError)
		log.WithField("err", err).Error("error getting last message")
		return
	}

	w.Header().Add("Access-Control-Allow-Origin", "*")
	w.Header().Add("Access-Control-Allow-Headers", "content-type,x-authorization")

	w.WriteHeader(http.StatusOK)
	w.Write(msg)
}
