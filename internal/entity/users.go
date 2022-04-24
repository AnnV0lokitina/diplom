package entity

import (
	"crypto/md5"
	"crypto/rand"
	"fmt"
	"time"
)

const (
	sessionIDLength = 16
	SessionLifetime = 10 * time.Minute
)

type User struct {
	Login           string
	PasswordHash    string
	ActiveSessionID string
}

func NewUser(login string, password string) (*User, error) {
	passwordHash := createPasswordHash(password)
	sessionID, err := generateSessionID()
	if err != nil {
		return nil, err
	}
	return &User{
		Login:           login,
		PasswordHash:    passwordHash,
		ActiveSessionID: sessionID,
	}, nil
}

func createPasswordHash(password string) string {
	bytePassword := []byte(password)
	idByte := md5.Sum(bytePassword)
	return fmt.Sprintf("%x", idByte)
}

func generateSessionID() (string, error) {
	b := make([]byte, sessionIDLength)
	_, err := rand.Read(b)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%x", b), nil
}
