package server

import (
	"context"
	"crypto/tls"
	"net/http"
	"path/filepath"
	"time"

	"ctm_lk/internal/config"
	"ctm_lk/internal/handlers"
	"ctm_lk/internal/middlewares"
	"ctm_lk/internal/models"
	"ctm_lk/pkg/logger"

	"github.com/go-chi/chi/v5"
	"golang.org/x/crypto/acme/autocert"
)

type Server struct {
	http.Server
	models.ServerDB
}

func (s *Server) getSelfSignedOrLetsEncryptCert(certManager *autocert.Manager) func(hello *tls.ClientHelloInfo) (*tls.Certificate, error) {
	return func(hello *tls.ClientHelloInfo) (*tls.Certificate, error) {
		keyFile := filepath.Join(config.Cfg.ProgramPath(), "localhost_cert", "key.pem")
		crtFile := filepath.Join(config.Cfg.ProgramPath(), "localhost_cert", "cert.pem")
		certificate, err := tls.LoadX509KeyPair(crtFile, keyFile)
		if err != nil {
			logger.Infof("%s\nFalling back to Letsencrypt\n", err)
			return certManager.GetCertificate(hello)
		}
		logger.Info("Loaded selfsigned certificate.")
		return &certificate, err
	}
}

func (s *Server) router() http.Handler {
	r := chi.NewRouter()
	r.Use(middlewares.AddAccessAllow, middlewares.ZipHandlerRead, middlewares.ZipHandlerWrite)
	r.Options("/*", handlers.HandlerOptions)
	r.Get("/*", handlers.HandlerStartPage)
	r.Post("/*", handlers.HandlerStartPage)

	r.Post("/api/user/register", handlers.HandlerRegistration(s.NewDBUserRepo()))
	r.Post("/api/user/login", handlers.HandlerLogin(s.NewDBUserRepo()))
	// r.Route("/api", func(r chi.Router) {
	// 	r.Use(middlewares.CheckAuthorization(s.NewDBUserRepo()))
	// 	r.Post("/user/...", ...)
	// 	r.Get("/user/...", ...)
	// })
	return r
}

//Start server with router.
// func (s *Server) Start(ctx context.Context, repo models.Repository, opt models.Options) {
func (s *Server) Start(ctx context.Context) {

	s.router()
	// s.Addr = config.Cfg.ServAddr()
	logger.Info("Старт сервера по адресу", config.Cfg.ServAddr())

	m := &autocert.Manager{
		Cache:      autocert.DirCache("golang-autocert"),
		Prompt:     autocert.AcceptTOS,
		HostPolicy: autocert.HostWhitelist(config.Cfg.ServAddr()),
	}

	tlsConfig := m.TLSConfig()
	tlsConfig.GetCertificate = s.getSelfSignedOrLetsEncryptCert(m)
	s.Addr = config.Cfg.ServAddr()
	s.TLSConfig = tlsConfig
	s.Handler = s.router()
	go s.ListenAndServeTLS("", "")
	// go s.ListenAndServeTLS("localhost/cert.pem", "localhost/key.pem")

	logger.Info("Сервер запущен")

	<-ctx.Done()
	ctx, cancelFunc := context.WithTimeout(context.Background(), time.Second*5)
	defer cancelFunc()
	s.Shutdown(ctx)
}
