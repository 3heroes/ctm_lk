package middlewares

import (
	"context"
	"ctm_lk/internal/config"
	"ctm_lk/internal/models"
	"ctm_lk/pkg/logger"
	"net/http"
)

func CheckAuthorizationCookie(ur models.UsersRepo) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			c, err := r.Cookie("UserID")
			key := ""
			if err == nil {
				key = c.Value
			}

			if len(key) == 0 {
				w.Header().Add("Location", "https://"+config.Cfg.ServAddrHttps()+"/registration.html")
				w.WriteHeader(http.StatusTemporaryRedirect)
				return
			}

			u := new(models.User)
			u.Token = key

			finded, err := ur.Get(r.Context(), u)
			if err != nil {
				logger.Info(http.StatusInternalServerError)
				w.WriteHeader(http.StatusInternalServerError)
				return
			}

			if !finded {
				logger.Info(http.StatusUnauthorized)
				w.Header().Add("Location", "https://"+config.Cfg.ServAddrHttps()+"/registration.html")
				w.WriteHeader(http.StatusTemporaryRedirect)
				return
			}

			c = &http.Cookie{
				Name:  "UserID",
				Value: key,
			}
			http.SetCookie(w, c)
			ctx := context.WithValue(r.Context(), models.UKeyName, u.ID)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
