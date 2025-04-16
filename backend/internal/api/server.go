package api

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"radaroficial.app/internal/api/handlers"
)

type Server struct {
	DB     *pgxpool.Pool
	Router *http.ServeMux
	server *http.Server
}

func NewServer(db *pgxpool.Pool) *Server {
	return &Server{
		DB:     db,
		Router: http.NewServeMux(),
	}
}

func (s *Server) RegisterHandlers() {
	// Register handlers
	crawlHandler := handlers.NewCrawlHandler(s.DB)
	s.Router.Handle("/crawl", crawlHandler)
	s.Router.Handle("/reindex", handlers.NewReindexHandler())

}

func (s *Server) Start(port string) error {
	s.server = &http.Server{
		Addr:    ":" + port,
		Handler: s.Router,
	}

	log.Println("Server running on port", port)
	return s.server.ListenAndServe()
}

// Shutdown gracefully shuts down the server with a timeout
func (s *Server) Shutdown(ctx context.Context) error {
	if s.server == nil {
		return nil
	}

	log.Println("Server shutdown initiated...")

	// Create a timeout context for shutdown if one wasn't provided
	shutdownCtx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	// Shutdown the server
	if err := s.server.Shutdown(shutdownCtx); err != nil {
		return err
	}

	log.Println("Server gracefully stopped")
	return nil
}
