package watch

import (
	"encoding/json"
	"net/http"
	"sync"
)

type SSEMessage struct {
	Type  string `json:"type"`
	HTML  string `json:"html,omitempty"`
	Error string `json:"error,omitempty"`
}

type SSEHub struct {
	clients map[http.ResponseWriter]bool
	mu      sync.Mutex
}

func NewSSEHub() *SSEHub {
	return &SSEHub{
		clients: make(map[http.ResponseWriter]bool),
	}
}

func (h *SSEHub) AddClient(w http.ResponseWriter) {
	h.mu.Lock()
	defer h.mu.Unlock()
	h.clients[w] = true
}

func (h *SSEHub) RemoveClient(w http.ResponseWriter) {
	h.mu.Lock()
	defer h.mu.Unlock()
	delete(h.clients, w)
}

func (h *SSEHub) Broadcast(msg SSEMessage) {
	h.mu.Lock()
	defer h.mu.Unlock()

	data, err := json.Marshal(msg)
	if err != nil {
		return
	}

	for client := range h.clients {
		if f, ok := client.(http.Flusher); ok {
			client.Write([]byte("data: "))
			client.Write(data)
			client.Write([]byte("\n\n"))
			f.Flush()
		}
	}
}
