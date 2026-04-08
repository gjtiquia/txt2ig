package web

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"github.com/gjtiquia/txt2ig/internal/config"
	"github.com/gjtiquia/txt2ig/internal/renderer"
	"github.com/gjtiquia/txt2ig/internal/web/templates/components"
	"github.com/gjtiquia/txt2ig/internal/web/templates/pages"
)

func (s *Server) handleHome(w http.ResponseWriter, r *http.Request) {
	defaultConfig := config.DefaultConfig()
	configBytes, err := json.MarshalIndent(defaultConfig, "", "    ")
	if err != nil {
		http.Error(w, "Failed to generate default config", http.StatusInternalServerError)
		return
	}

	pages.Home(string(configBytes)).Render(r.Context(), w)
}

func (s *Server) handleConvert(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		components.ErrorAlert("Failed to parse form data").Render(r.Context(), w)
		return
	}

	text := r.FormValue("text")
	configStr := r.FormValue("config")

	if text == "" {
		components.ErrorAlert("Text is required").Render(r.Context(), w)
		return
	}

	cfg, err := config.ParseConfig([]byte(configStr))
	if err != nil {
		components.ErrorAlert(fmt.Sprintf("Invalid config: %v", err)).Render(r.Context(), w)
		return
	}

	rend := renderer.NewRenderer(cfg)
	defer rend.Close()

	img, err := rend.RenderString(text)
	if err != nil {
		components.ErrorAlert(fmt.Sprintf("Failed to render: %v", err)).Render(r.Context(), w)
		return
	}

	format := "png"
	base64, err := renderer.EncodeImage(img, format)
	if err != nil {
		components.ErrorAlert(fmt.Sprintf("Failed to encode image: %v", err)).Render(r.Context(), w)
		return
	}

	components.ImageResult(base64, format).Render(r.Context(), w)
}

func (s *Server) handleWatchPage(w http.ResponseWriter, r *http.Request) {
	initialImage := ""
	if s.watchedFile != "" {
		if data, err := os.ReadFile(s.watchedFile); err == nil {
			cfg := s.config
			rend := renderer.NewRenderer(cfg)
			if img, err := rend.RenderString(string(data)); err == nil {
				if base64, err := renderer.EncodeImage(img, "png"); err == nil {
					initialImage = base64
				}
			}
			rend.Close()
		}
	}

	pages.Watch(s.watchedFile, initialImage).Render(r.Context(), w)
}

func (s *Server) handleWatchSSE(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")

	s.sseHub.AddClient(w)
	defer s.sseHub.RemoveClient(w)

	if f, ok := w.(http.Flusher); ok {
		f.Flush()
	}

	<-r.Context().Done()
}
