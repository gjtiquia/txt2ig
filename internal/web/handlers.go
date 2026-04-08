package web

import (
	"encoding/json"
	"fmt"
	"net/http"

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
