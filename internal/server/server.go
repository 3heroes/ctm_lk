package server

import (
	"context"
	"net/http"
	"time"

	"ctm_lk/internal/config"
	"ctm_lk/internal/handlers"
	"ctm_lk/internal/middlewares"
	"ctm_lk/internal/models"
	"ctm_lk/pkg/logger"

	"github.com/go-chi/chi/v5"
)

type Server struct {
	http.Server
	models.ServerDB
}

//Start server with router.
// func (s *Server) Start(ctx context.Context, repo models.Repository, opt models.Options) {
func (s *Server) Start(ctx context.Context) {
	r := chi.NewRouter()
	r.Use(middlewares.ZipHandlerRead, middlewares.ZipHandlerWrite)
	r.Get("/*", handlers.HandlerStartPage)
	r.Post("/api/user/register", handlers.HandlerRegistration(s.NewDBUserRepo()))
	r.Post("/api/user/login", handlers.HandlerLogin(s.NewDBUserRepo()))
	// r.Route("/api", func(r chi.Router) {
	// 	r.Use(middlewares.CheckAuthorization(s.NewDBUserRepo()))
	// 	r.Post("/user/...", ...)
	// 	r.Get("/user/...", ...)
	// })

	s.Addr = config.Cfg.ServAddr()
	logger.Info("Старт сервера по адресу", config.Cfg.ServAddr())
	s.Handler = r
	go s.ListenAndServe()

	logger.Info("Сервер запущен")

	<-ctx.Done()
	ctx, cancelFunc := context.WithTimeout(context.Background(), time.Second*5)
	defer cancelFunc()
	s.Shutdown(ctx)
}
