package ws

import log "github.com/Sirupsen/logrus"

type hub struct {
	clients    map[Client]bool
	broadcast  chan Message
	register   chan Client
	unregister chan Client
}

type Message struct {
	Text  []byte
	Topic string
}

func NewHub() Hub {
	return &hub{
		broadcast:  make(chan Message),
		register:   make(chan Client),
		unregister: make(chan Client),
		clients:    make(map[Client]bool),
	}
}

func (h *hub) Broadcast(message Message) {
	h.broadcast <- message
}

func (h *hub) Unregister(c Client) {
	h.unregister <- c
}

func (h *hub) Register(c Client) {
	h.register <- c
}

func (h *hub) Run() {
	for {
		select {
		case client := <-h.register:
			h.clients[client] = true
			log.WithField("remoteAddr", client.RemoteAddr()).WithField("numClients", len(h.clients)).Info("Registered new client")

		case client := <-h.unregister:
			if _, ok := h.clients[client]; ok {
				delete(h.clients, client)
				client.Close()
				log.WithField("remoteAddr", client.RemoteAddr()).WithField("numClients", len(h.clients)).Info("Deregistered client")
			}
		case message := <-h.broadcast:
			if len(h.clients) > 0 {
				log.WithField("topic", message.Topic).
					WithField("messageLength", len(message.Text)).
					WithField("numClients", len(h.clients)).
					Info("Broadcasting message")
				for client := range h.clients {
					select {
					case client.Send() <- message:
					default:
						client.Close()
						delete(h.clients, client)
					}
				}
			}
		}
	}
}
