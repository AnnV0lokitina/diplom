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
	ID              int
	Login           string
	ActiveSessionID string
}

func CreatePasswordHash(password string) string {
	bytePassword := []byte(password)
	idByte := md5.Sum(bytePassword)
	return fmt.Sprintf("%x", idByte)
}

func GenerateSessionID() (string, error) {
	b := make([]byte, sessionIDLength)
	_, err := rand.Read(b)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%x", b), nil
}
