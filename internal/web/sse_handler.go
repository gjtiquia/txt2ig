package web

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sync"
)

type SSEMessage struct {
	Type  string `json:"type"`
	Data  string `json:"data,omitempty"`
	HTML  string `json:"html,omitempty"`
	Error string `json:"error,omitempty"`
}

type SSEHub struct {
	clients   map[http.ResponseWriter]bool
	clientsMu sync.RWMutex
}

func NewSSEHub() *SSEHub {
	return &SSEHub{
		clients: make(map[http.ResponseWriter]bool),
	}
}

func (h *SSEHub) AddClient(w http.ResponseWriter) {
	h.clientsMu.Lock()
	h.clients[w] = true
	h.clientsMu.Unlock()
}

func (h *SSEHub) RemoveClient(w http.ResponseWriter) {
	h.clientsMu.Lock()
	delete(h.clients, w)
	h.clientsMu.Unlock()
}

func (h *SSEHub) Broadcast(message SSEMessage) {
	h.clientsMu.RLock()
	defer h.clientsMu.RUnlock()

	data, err := json.Marshal(message)
	if err != nil {
		log.Printf("error marshaling SSE message: %v", err)
		return
	}

	for client := range h.clients {
		_, err := fmt.Fprintf(client, "data: %s\n\n", data)
		if err != nil {
			continue
		}
		if f, ok := client.(http.Flusher); ok {
			f.Flush()
		}
	}
}
