package handler

import (
	"github.com/AnnV0lokitina/diplom/internal/entity"
	"net/http"
	"time"
)

const SessionCookieName = "session"

func addSessionIDToCookie(w http.ResponseWriter, sessionID string) {
	cookie := &http.Cookie{
		Name:    SessionCookieName,
		Value:   sessionID,
		Expires: time.Now().Add(entity.TTL),
	}
	http.SetCookie(w, cookie)
}

func getSessionIDFromCookie(r *http.Request) (string, error) {
	cookie, err := r.Cookie(SessionCookieName)
	if err != nil {
		return "", err
	}
	return cookie.Value, nil
}
