package web

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gjtiquia/txt2ig/internal/config"
	"github.com/gjtiquia/txt2ig/internal/renderer"
)

type Server struct {
	config      *config.Config
	renderer    *renderer.Renderer
	watcher     *FileWatcher
	sseHub      *SSEHub
	watchedFile string
}

func NewServer() *Server {
	return &Server{
		config: config.DefaultConfig(),
		sseHub: NewSSEHub(),
	}
}

func (s *Server) WithWatch(file string) *Server {
	s.watchedFile = file
	return s
}

func (s *Server) Run(port int) error {
	if s.watchedFile != "" {
		if err := s.startWatcher(); err != nil {
			return fmt.Errorf("start watcher: %w", err)
		}
		defer s.watcher.Close()
	}

	mux := http.NewServeMux()

	mux.HandleFunc("GET /", s.handleHome)
	mux.HandleFunc("POST /convert", s.handleConvert)
	mux.HandleFunc("GET /watch", s.handleWatchPage)
	mux.HandleFunc("GET /sse/watch", s.handleWatchSSE)

	fs := http.FileServer(http.Dir("static"))
	mux.Handle("GET /static/", http.StripPrefix("/static/", fs))

	handler := loggingMiddleware(mux)

	addr := fmt.Sprintf(":%d", port)
	log.Printf("Starting web server on http://localhost%s\n", addr)
	if s.watchedFile != "" {
		log.Printf("Watching file: %s\n", s.watchedFile)
	}

	return http.ListenAndServe(addr, handler)
}

func (s *Server) startWatcher() error {
	var err error
	s.watcher, err = NewFileWatcher()
	if err != nil {
		return err
	}

	s.watcher.Start()

	return s.watcher.Watch(s.watchedFile, func(path string) {
		log.Printf("File changed: %s, regenerating...\n", path)
		s.regenerateAndBroadcast()
	})
}

func (s *Server) regenerateAndBroadcast() {
	cfg := s.config

	data, err := os.ReadFile(s.watchedFile)
	if err != nil {
		s.sseHub.Broadcast(SSEMessage{
			Type:  "error",
			Error: fmt.Sprintf("Failed to read file: %v", err),
		})
		return
	}

	text := string(data)

	rend := renderer.NewRenderer(cfg)
	defer rend.Close()

	img, err := rend.RenderString(text)
	if err != nil {
		s.sseHub.Broadcast(SSEMessage{
			Type:  "error",
			Error: fmt.Sprintf("Failed to render: %v", err),
		})
		return
	}

	base64, err := renderer.EncodeImage(img, "png")
	if err != nil {
		s.sseHub.Broadcast(SSEMessage{
			Type:  "error",
			Error: fmt.Sprintf("Failed to encode: %v", err),
		})
		return
	}

	s.sseHub.Broadcast(SSEMessage{
		Type: "image",
		HTML: fmt.Sprintf(`<img src="data:image/png;base64,%s" alt="Preview" class="w-full h-auto">`, base64),
	})
}
