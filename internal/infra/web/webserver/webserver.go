package webserver

import (
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
)

type WebServer struct {
	Router        chi.Router
	Handlers      map[string]http.HandlerFunc
	WebServerPort string
}

func NewWebServer(serverPort string) *WebServer {
	return &WebServer{
		Router:        chi.NewRouter(),
		Handlers:      make(map[string]http.HandlerFunc),
		WebServerPort: serverPort,
	}
}

func (s *WebServer) AddHandler(path string, handler http.HandlerFunc) {
	s.Handlers[path] = handler
}

func (s *WebServer) Start() {
	// Middleware de logging
	s.Router.Use(middleware.Logger)

	// Middleware CORS - ESSENCIAL para permitir acesso do frontend
	s.Router.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"*"}, // Em produção, troque por seu domínio específico
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: false,
		MaxAge:           300, // Cache do preflight
	}))

	// Registra todos os handlers
	for path, handler := range s.Handlers {
		s.Router.Handle(path, handler)
	}

	// Inicia o servidor
	log.Printf("Servidor rodando na porta %s", s.WebServerPort)
	http.ListenAndServe(s.WebServerPort, s.Router)
}