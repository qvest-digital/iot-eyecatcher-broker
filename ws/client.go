package ws

import (
	"encoding/json"
	log "github.com/Sirupsen/logrus"
	"github.com/gorilla/websocket"
	"net/http"
	"time"
)

const (
	// Time allowed to write a message to the peer.
	writeWait = 10 * time.Second

	// Time allowed to read the message form the peer
	pongWait = 10 * time.Hour

	// Send pings to peer with this period. Must be less than pongWait.
	pingPeriod = (pongWait * 9) / 10
)

type SubscribeMsg struct {
	Operation string `json:"operation"`
	Topic     string `json:"topic"`
}

type client struct {
	hub           Hub
	conn          ConnI
	send          chan Message
	subscriptions map[string]interface{}
	stop          bool
	quit          chan bool
}

func (c *client) Send() chan Message {
	return c.send
}

func (c *client) Close() {
	close(c.send)
}

func (c *client) RemoteAddr() string {
	return c.conn.RemoteAddr().String()
}

func (c *client) readPump() {
	for c.stop != true {
		_, msg, err := c.conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway) {
				log.WithField("err", err).Info("Client unexpectedly closed connection")
			}
			c.quit <- true
			return
		}

		var subscribeMsg SubscribeMsg
		err = json.Unmarshal(msg, &subscribeMsg)
		if err == nil {
			switch subscribeMsg.Operation {
			case "subscribe":
				log.WithField("topic", subscribeMsg.Topic).Info("Subscribe")
				c.subscribe(subscribeMsg.Topic)
				continue
			case "unsubscribe":
				log.WithField("topic", subscribeMsg.Topic).Info("Unsubscribe")
				c.unsubscribe(subscribeMsg.Topic)
				continue
			}
		}
		log.Warn("Unidentified message received")
	}
	c.quit <- true
}

func (c *client) writePump() {
	for c.stop != true {
		select {
		case message, ok := <-c.send:

			if c.isSubscribed(message.Topic) {
				c.conn.SetWriteDeadline(time.Now().Add(writeWait))
				if !ok {
					// The hub closed the channel.
					c.conn.WriteMessage(websocket.CloseMessage, []byte{})
					c.quit <- true
					return
				}

				w, err := c.conn.NextWriter(websocket.TextMessage)
				if err != nil {
					log.WithField("err", err).Info("NextWriter error")
					c.quit <- true
					return
				}
				w.Write(message.Text)
				if err := w.Close(); err != nil {
					log.WithField("err", err).Info("Close error")
					c.quit <- true
					return
				}
			}
		}
	}
	c.quit <- true
}

func (c *client) subscribe(topic string) {
	c.subscriptions[topic] = true
}

func (c *client) unsubscribe(topic string) {
	delete(c.subscriptions, topic)
}

func (c *client) isSubscribed(topic string) bool {
	for k := range c.subscriptions {
		if k == topic {
			return true
		}
	}
	return false
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  0,
	WriteBufferSize: 4096,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func Handler(hub Hub, w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.WithField("err", err).WithField("remoteAddr", r.RemoteAddr).Error("error upgrading connection")
		return
	}
	client := &client{
		hub:           hub,
		conn:          conn,
		send:          make(chan Message, 32),
		subscriptions: make(map[string]interface{}, 0),
		stop:          false,
		quit:          make(chan bool, 1)}

	client.hub.Register(client)
	go client.writePump()
	go client.readPump()

	<-client.quit

	client.stop = true
	client.hub.Unregister(client)
	client.conn.Close()
}
