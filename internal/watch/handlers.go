package watch

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"path/filepath"

	"github.com/gjtiquia/txt2ig/internal/renderer"
	"github.com/gjtiquia/txt2ig/internal/watch/templates/pages"
)

func (s *Server) handlePreview(w http.ResponseWriter, r *http.Request) {
	text, err := os.ReadFile(s.watchedFile)
	if err != nil {
		http.Error(w, "Failed to read file", http.StatusInternalServerError)
		return
	}

	rend := renderer.NewRenderer(s.config)
	defer rend.Close()

	img, err := rend.RenderString(string(text))
	if err != nil {
		http.Error(w, "Failed to render", http.StatusInternalServerError)
		return
	}

	base64, err := renderer.EncodeImage(img, "png")
	if err != nil {
		http.Error(w, "Failed to encode image", http.StatusInternalServerError)
		return
	}

	fileName := filepath.Base(s.watchedFile)
	configName := s.getConfigName()

	pages.Preview(fileName, configName, base64).Render(r.Context(), w)
}

func (s *Server) handleSSE(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")

	s.sseHub.AddClient(w)
	defer s.sseHub.RemoveClient(w)

	if f, ok := w.(http.Flusher); ok {
		f.Flush()
	}

	s.sendInitialImage(w)

	<-r.Context().Done()
}

func (s *Server) sendInitialImage(w http.ResponseWriter) {
	text, err := os.ReadFile(s.watchedFile)
	if err != nil {
		return
	}

	rend := renderer.NewRenderer(s.config)
	defer rend.Close()

	img, err := rend.RenderString(string(text))
	if err != nil {
		return
	}

	base64, err := renderer.EncodeImage(img, "png")
	if err != nil {
		return
	}

	data, err := json.Marshal(SSEMessage{
		Type: "image",
		HTML: fmt.Sprintf(`<img src="data:image/png;base64,%s" alt="Preview" class="w-full h-auto border border-gray-700 rounded" id="preview-image">`, base64),
	})
	if err != nil {
		return
	}

	w.Write([]byte("data: "))
	w.Write(data)
	w.Write([]byte("\n\n"))

	if f, ok := w.(http.Flusher); ok {
		f.Flush()
	}
}
