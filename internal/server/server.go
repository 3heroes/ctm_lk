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
	redirect http.Server
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

func (s *Server) redirectToHTTPS(ctx context.Context) {
	s.redirect = http.Server{
		Addr: config.Cfg.ServAddrHttp(),
		Handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			u := r.URL
			u.Host = config.Cfg.ServAddrHttps()
			u.Scheme = "https"
			logger.Info(u.String())
			http.Redirect(w, r, u.String(), http.StatusMovedPermanently)
		}),
	}
	logger.Info("Старт сервера переадресации по адресу", config.Cfg.ServAddrHttp())

	go s.redirect.ListenAndServe()
}

func (s *Server) router() http.Handler {

	fs := http.FileServer(http.Dir("./html"))
	r := chi.NewRouter()
	r.Use(middlewares.AddAccessAllow, middlewares.ZipHandlerRead, middlewares.ZipHandlerWrite)

	r.Options("/*", handlers.HandlerOptions)
	r.Get("/registration.html", fs.ServeHTTP)
	r.Get("/css/style.css", fs.ServeHTTP)
	r.Get("/js/registration.js", fs.ServeHTTP)
	r.Post("/api/user/register", handlers.HandlerRegistrationCookie(s.NewDBUserRepo()))
	r.Post("/api/user/login", handlers.HandlerLoginCookie(s.NewDBUserRepo()))
	r.Group(func(r chi.Router) {
		r.Use(middlewares.CheckAuthorizationCookie(s.NewDBUserRepo()))
		r.Get("/*", fs.ServeHTTP)
		r.Post("/*", handlers.HandlerStartPage)
	})

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
	logger.Info("Старт сервера по адресу", config.Cfg.ServAddrHttps())

	m := &autocert.Manager{
		Cache:      autocert.DirCache("golang-autocert"),
		Prompt:     autocert.AcceptTOS,
		HostPolicy: autocert.HostWhitelist(config.Cfg.ServAddrHttps()),
	}

	tlsConfig := m.TLSConfig()
	tlsConfig.GetCertificate = s.getSelfSignedOrLetsEncryptCert(m)
	s.Addr = config.Cfg.ServAddrHttps()
	s.TLSConfig = tlsConfig
	s.Handler = s.router()
	go s.ListenAndServeTLS("", "")
	// go s.ListenAndServeTLS("localhost/cert.pem", "localhost/key.pem")

	logger.Info("Сервер запущен")

	<-ctx.Done()
	ctx, cancelFunc := context.WithTimeout(context.Background(), time.Second*5)
	defer cancelFunc()
	s.Shutdown(ctx)
	s.redirect.Shutdown(ctx)
}
