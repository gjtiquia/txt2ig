package web

import (
	"fmt"
	"net/http"

	"github.com/gjtiquia/txt2ig/internal/config"
	"github.com/gjtiquia/txt2ig/internal/renderer"
)

type Server struct {
	config   *config.Config
	renderer *renderer.Renderer
}

func NewServer() *Server {
	return &Server{
		config:   config.DefaultConfig(),
		renderer: nil,
	}
}

func (s *Server) Run(port int) error {
	mux := http.NewServeMux()

	mux.HandleFunc("GET /", s.handleHome)
	mux.HandleFunc("POST /convert", s.handleConvert)

	fs := http.FileServer(http.Dir("static"))
	mux.Handle("GET /static/", http.StripPrefix("/static/", fs))

	handler := loggingMiddleware(mux)

	addr := fmt.Sprintf(":%d", port)
	fmt.Printf("Starting web server on http://localhost%s\n", addr)

	return http.ListenAndServe(addr, handler)
}
