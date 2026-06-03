package server

import (
	"Subscribe-service/internal/config"
	"Subscribe-service/internal/handler"
	"context"
	"database/sql"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

type Server struct {
	router *gin.Engine
	cfg    *config.Config
	db     *sql.DB
	http   *http.Server
}

func New(cfg *config.Config, db *sql.DB) *Server {
	r := gin.Default()
	return &Server{
		router: r,
		cfg:    cfg,
		db:     db,
	}
}
func (s *Server) MapRoutes(h *handler.Handler) {
	s.router.POST("/subscriptions", h.CreateSubsc)
	s.router.GET("/subscriptions/{id}", h.GetSubsc)
	s.router.PUT("/subscriptions/{id}", h.UpdateSubsc)
	s.router.DELETE("/subscriptions/{id}", h.DeleteSubsc)
	s.router.GET("/subscriptions", h.GetAllSubsc)
	s.router.GET("/subscriptions/aggregate", h.GetSubscSum)

	s.router.GET("/health", h.Health)

}

func (s *Server) Run(h *handler.Handler) error {
	s.MapRoutes(h)
	addr := s.cfg.Server.Host + s.cfg.Server.Port
	s.http = &http.Server{
		Addr:    addr,
		Handler: s.router,
	}
	go func() {
		if err := s.http.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Printf("server error: %v", err)
		}
	}()

	return nil
}

func (s *Server) Stop(ctx context.Context) error {
	if s.http == nil {
		return nil
	}
	return s.http.Shutdown(ctx)
}
