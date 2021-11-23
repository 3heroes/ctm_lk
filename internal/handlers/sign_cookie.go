package handlers

import (
	"ctm_lk/internal/models"
	"ctm_lk/pkg/logger"
	"net/http"
)

func HandlerRegistrationCookie(ur models.UsersRepo) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		logger.Info("Обработка запроса регистрации")
		ctx := r.Context()

		user, ok := getUserFromRequest(r)
		if !ok {
			logger.Info(http.StatusBadRequest)
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		finded, err := ur.Get(ctx, user)
		if err != nil {
			logger.Info(http.StatusInternalServerError)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		if finded {
			logger.Info(http.StatusConflict)
			w.WriteHeader(http.StatusConflict)
			return
		}

		if err := ur.Add(ctx, user); err != nil {
			logger.Info(http.StatusInternalServerError)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		c := &http.Cookie{
			Name:  "UserID",
			Value: user.Token,
		}
		http.SetCookie(w, c)
		w.WriteHeader(http.StatusOK)
		logger.Info(http.StatusOK)
	}
}

func HandlerLoginCookie(ur models.UsersRepo) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		logger.Info("Обработка запроса входа")
		ctx := r.Context()

		user, ok := getUserFromRequest(r)
		if !ok {
			logger.Info(http.StatusBadRequest)
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		logPass := user.Password

		finded, err := ur.Get(ctx, user)
		if err != nil {
			logger.Info(http.StatusInternalServerError)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		if !finded || user.Password != logPass {
			logger.Info(http.StatusUnauthorized)
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		c := &http.Cookie{
			Name:  "UserID",
			Value: user.Token,
		}
		http.SetCookie(w, c)
		w.WriteHeader(http.StatusOK)
		logger.Info(http.StatusOK)
	}
}
