//go:generate ./mockgen.sh

package main

import (
	log "github.com/Sirupsen/logrus"
	"github.com/gorilla/mux"
	"github.com/tarent/iot-eyecatcher-broker/broker"
	"github.com/tarent/iot-eyecatcher-broker/ws"
	"net/http"
	"os"
    "github.com/goji/httpauth"
)

func main() {
	log.Info("Starting websocket message broker")

	listen, found := os.LookupEnv("WS_LISTEN")
	if !found {
		listen = ":8080"
	}

	hub := ws.NewHub()
	go hub.Run()

	b := broker.NewBroker(hub)
	go b.Run()

	r := mux.NewRouter()

	// Root websocket handler
	r.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		ws.Handler(hub, w, r)
	}).Methods(http.MethodGet)

	// "Get last message from topic" handler
	r.HandleFunc("/{topic}", func(w http.ResponseWriter, r *http.Request) {
		broker.GetLastTopicMessageHandler(b, w, r)
	}).Methods(http.MethodGet)

	// CORS-Preflight handler for "get last message from topic" handler
	r.HandleFunc("/{topic}", corsHandler).
		Methods(http.MethodOptions)

    // Post new message handler
    r.Handle("/{topic}",
        httpauth.SimpleBasicAuth("dave", "somepassword")(
		http.HandlerFunc(
			func(w http.ResponseWriter, r *http.Request) {
				broker.MessageHandler(b, w, r)
			}))).
		Methods(http.MethodPost)

    // Start the server
	log.WithField("listenAddress", listen).Info("Start to listen")
	err := http.ListenAndServe(listen, r)
	if err != nil {
		log.WithField("listenAddress", listen).WithField("err", err).Fatal("Error listening")
	}

	log.Info("Exiting websocket message broker")
}

func corsHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Access-Control-Allow-Origin", "*")
	w.Header().Add("Access-Control-Allow-Headers", "content-type,x-authorization")
	w.WriteHeader(http.StatusOK)
}
