package repo

import (
	"context"
	"errors"
	"github.com/AnnV0lokitina/diplom/internal/entity"
	labelError "github.com/AnnV0lokitina/diplom/pkg/error"
	"time"
)

func (r *Repo) CreateUser(
	ctx context.Context,
	sessionID string,
	login string,
	passwordHash string,
) error {
	ctx, cancel := context.WithTimeout(ctx, 1*time.Second)
	defer cancel()

	tx, err := r.conn.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	sqlInsertUser := "INSERT INTO users (login, password) " +
		"VALUES ($1, $2) " +
		"ON CONFLICT (login) DO NOTHING"

	result, err := tx.Exec(ctx, sqlInsertUser, login, passwordHash)
	if err != nil {
		return err
	}
	if result.RowsAffected() == 0 {
		return labelError.NewLabelError(labelError.TypeConflict, errors.New("login exists"))
	}

	sqlInsertSession := "INSERT INTO sessions (session_id, created_at, lifetime, login) " +
		"VALUES ($1, $2, $3, $4)"

	timestamp := time.Now().Unix()
	lifetime := entity.SessionLifetime.Seconds()
	if _, err = tx.Exec(ctx, sqlInsertSession, sessionID, timestamp, lifetime, login); err != nil {
		return err
	}
	err = tx.Commit(ctx)
	if err != nil {
		return err
	}
	return nil
}

func (r *Repo) AuthUser(
	ctx context.Context,
	login string,
	passwordHash string,
) (bool, error) {
	ctx, cancel := context.WithTimeout(ctx, 1*time.Second)
	defer cancel()
	sql := "SELECT count(*) FROM users WHERE login=$1 AND password=$2"
	rows, _ := r.conn.Query(ctx, sql, login, passwordHash)
	n := 0
	for rows.Next() {
		err := rows.Scan(&n)
		if err != nil {
			return false, err
		}
	}
	return n > 0, nil
}

func (r *Repo) AddUserSession(
	ctx context.Context,
	sessionID string,
	login string,
) error {
	ctx, cancel := context.WithTimeout(ctx, 1*time.Second)
	defer cancel()
	sqlInsertSession := "INSERT INTO sessions (session_id, created_at, lifetime, login) " +
		"VALUES ($1, $2, $3, $4)"
	timestamp := time.Now().Unix()
	lifetime := entity.SessionLifetime.Seconds()
	result, err := r.conn.Exec(ctx, sqlInsertSession, sessionID, timestamp, lifetime, login)
	if err != nil {
		return err
	}
	if result.RowsAffected() == 0 {
		return labelError.NewLabelError(labelError.TypeNotFound, errors.New("no registered user"))
	}
	return nil
}

func (r *Repo) GetUserBySessionID(ctx context.Context, activeSessionID string) (*entity.User, error) {
	ctx, cancel := context.WithTimeout(ctx, 1*time.Second)
	defer cancel()
	sql := "SELECT login FROM sessions WHERE session_id=$1 AND created_at > $2 - lifetime"
	timestamp := time.Now().Unix()
	rows, _ := r.conn.Query(ctx, sql, activeSessionID, timestamp)
	var user *entity.User
	for rows.Next() {
		var login string
		err := rows.Scan(&login)
		if err != nil {
			return nil, err
		}
		user = &entity.User{
			Login:           login,
			ActiveSessionID: activeSessionID,
		}
	}
	if user == nil {
		return nil, labelError.NewLabelError(labelError.TypeNotFound, errors.New("no registered user"))
	}
	return user, nil
}
