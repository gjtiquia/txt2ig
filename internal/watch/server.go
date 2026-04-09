package watch

import (
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"path/filepath"

	"github.com/gjtiquia/txt2ig/internal/config"
	"github.com/gjtiquia/txt2ig/internal/renderer"
)

type Server struct {
	watchedFile string
	outputPath  string
	configFile  string
	config      *config.Config
	watcher     *FileWatcher
	sseHub      *SSEHub
	port        int
}

func NewServer(watchedFile, configFile string) (*Server, error) {
	loader := config.NewConfigLoader()
	if configFile != "" {
		loader.SetCustomPath(configFile)
	}

	cfg, err := loader.Load()
	if err != nil {
		return nil, fmt.Errorf("load config: %w", err)
	}

	return &Server{
		watchedFile: watchedFile,
		configFile:  configFile,
		config:      cfg,
	}, nil
}

func (s *Server) Run(port int) error {
	outputPath := renderer.DetermineOutputPath(s.watchedFile, "")
	s.outputPath = outputPath

	watcher, err := NewFileWatcher()
	if err != nil {
		return fmt.Errorf("create file watcher: %w", err)
	}
	s.watcher = watcher

	s.watcher.Start()

	if err := s.watcher.Watch(s.watchedFile, s.handleFileChange); err != nil {
		return fmt.Errorf("watch file: %w", err)
	}

	fmt.Printf("Watching: %s\n", filepath.Base(s.watchedFile))
	fmt.Printf("Output: %s\n", filepath.Base(s.outputPath))

	if configFile := s.getConfigName(); configFile != "" {
		fmt.Printf("Config: %s\n", configFile)
	}

	if port > 0 {
		fmt.Printf("Preview: http://localhost:%d\n", port)
	}

	fmt.Println("Press Ctrl+C to stop")

	if err := s.regenerateImage(); err != nil {
		log.Printf("initial render: %v", err)
	}

	if port > 0 {
		s.port = port
		return s.startWebServer(port)
	}

	select {}
}

func (s *Server) handleFileChange(path string) {
	log.Printf("File changed: %s, regenerating...\n", path)
	if err := s.regenerateImage(); err != nil {
		log.Printf("regenerate: %v", err)
		if s.sseHub != nil {
			s.sseHub.Broadcast(SSEMessage{
				Type:  "error",
				Error: fmt.Sprintf("Failed to regenerate: %v", err),
			})
		}
	}
}

func (s *Server) regenerateImage() error {
	text, err := os.ReadFile(s.watchedFile)
	if err != nil {
		return fmt.Errorf("read file: %w", err)
	}

	rend := renderer.NewRenderer(s.config)
	defer rend.Close()

	img, err := rend.RenderString(string(text))
	if err != nil {
		return fmt.Errorf("render: %w", err)
	}

	if err := renderer.SaveImage(img, s.outputPath); err != nil {
		return fmt.Errorf("save image: %w", err)
	}

	if s.sseHub != nil {
		base64, err := renderer.EncodeImage(img, "png")
		if err != nil {
			return fmt.Errorf("encode image: %w", err)
		}

		s.sseHub.Broadcast(SSEMessage{
			Type: "image",
			HTML: fmt.Sprintf(`<img src="data:image/png;base64,%s" alt="Preview" class="w-full h-auto border border-gray-700 rounded">`, base64),
		})
	}

	return nil
}

func (s *Server) startWebServer(port int) error {
	addr := fmt.Sprintf(":%d", port)

	listener, err := net.Listen("tcp", addr)
	if err != nil {
		return fmt.Errorf("port %d is already in use", port)
	}
	listener.Close()

	s.sseHub = NewSSEHub()

	mux := http.NewServeMux()
	mux.HandleFunc("GET /", s.handlePreview)
	mux.HandleFunc("GET /sse", s.handleSSE)

	handler := s.loggingMiddleware(mux)

	log.Printf("Starting web server on http://localhost%s\n", addr)

	return http.ListenAndServe(addr, handler)
}

func (s *Server) getConfigName() string {
	if s.configFile != "" {
		return filepath.Base(s.configFile)
	}

	home, err := os.UserHomeDir()
	if err != nil {
		return ""
	}

	xdgConfig := os.Getenv("XDG_CONFIG_HOME")
	if xdgConfig == "" {
		xdgConfig = filepath.Join(home, ".config")
	}

	paths := []string{
		".txt2igconfig.jsonc",
		filepath.Join(xdgConfig, "txt2ig", "config.jsonc"),
		filepath.Join(home, ".txt2ig", "config.jsonc"),
	}

	for _, path := range paths {
		if _, err := os.Stat(path); err == nil {
			return filepath.Base(path)
		}
	}

	return "default"
}

func (s *Server) loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("%s %s", r.Method, r.URL.Path)
		next.ServeHTTP(w, r)
	})
}
