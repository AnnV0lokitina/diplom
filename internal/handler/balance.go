package handler

import (
	log "github.com/sirupsen/logrus"
	"net/http"
)

func (h *Handler) GetBalance() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		sessionID, err := getSessionIDFromCookie(r)
		if err != nil {
			log.WithFields(log.Fields{
				"session ID": sessionID,
			}).Info("authorization failed")
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}
		log.WithFields(log.Fields{
			"session ID": sessionID,
		}).Info("authorization success")
		w.WriteHeader(http.StatusOK)
	}
}

func (h *Handler) Withdraw() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		sessionID, err := getSessionIDFromCookie(r)
		if err != nil {
			log.WithFields(log.Fields{
				"session ID": sessionID,
			}).Info("authorization failed")
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}
		log.WithFields(log.Fields{
			"session ID": sessionID,
		}).Info("authorization success")
		w.WriteHeader(http.StatusOK)
	}
}

func (h *Handler) GetWithdrawals() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		sessionID, err := getSessionIDFromCookie(r)
		if err != nil {
			log.WithFields(log.Fields{
				"session ID": sessionID,
			}).Info("authorization failed")
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}
		log.WithFields(log.Fields{
			"session ID": sessionID,
		}).Info("authorization success")
		w.WriteHeader(http.StatusOK)
	}
}
