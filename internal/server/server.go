package server

import (
	"Subscribe-service/internal/config"
	"Subscribe-service/internal/handler"
	"context"
	"database/sql"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	_ "Subscribe-service/docs"
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
	api := s.router.Group("/subscriptions")

	api.POST("", h.CreateSubs)
	api.GET("", h.GetAllSubs)
	api.GET("/item", h.GetSubs)
	api.PUT("/item", h.UpdateSubs)
	api.DELETE("/item", h.DeleteSubs)
	api.GET("/aggregate", h.GetSubsSum)

	s.router.GET("/health", h.Health)

	s.router.GET("/swagger/*any",
		ginSwagger.WrapHandler(swaggerFiles.Handler),
	)
}

func (s *Server) Run(h *handler.Handler) error {
	s.MapRoutes(h)

	addr := s.cfg.Server.Host + s.cfg.Server.Port

	s.http = &http.Server{
		Addr:    addr,
		Handler: s.router,
	}

	log.Printf("server starting on %s", addr)

	go func() {
		if err := s.http.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Printf("server runtime error %v", err)
		}
	}()

	return nil
}

func (s *Server) Stop(ctx context.Context) error {
	log.Println("server shutdown initiated")

	if s.http == nil {
		return nil
	}

	err := s.http.Shutdown(ctx)
	if err != nil {
		log.Printf("server shutdown error %v", err)
		return err
	}

	log.Println("server stopped gracefully")

	return nil
}
