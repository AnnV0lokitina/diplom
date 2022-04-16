package handler

import (
	"fmt"
	"net/http"
)

func (h *Handler) GetBalance() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		sessionID, err := getSessionIDFromCookie(r)
		if err != nil {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}
		fmt.Println(sessionID)
		w.WriteHeader(http.StatusOK)
	}
}

func (h *Handler) Withdraw() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		sessionID, err := getSessionIDFromCookie(r)
		if err != nil {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}
		fmt.Println(sessionID)
		w.WriteHeader(http.StatusOK)
	}
}

func (h *Handler) GetWithdrawals() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		sessionID, err := getSessionIDFromCookie(r)
		if err != nil {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}
		fmt.Println(sessionID)
		w.WriteHeader(http.StatusOK)
	}
}
