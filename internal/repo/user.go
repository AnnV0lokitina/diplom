package repo

import (
	"context"
	"errors"
	"github.com/AnnV0lokitina/diplom/internal/entity"
	labelError "github.com/AnnV0lokitina/diplom/pkg/error"
	"github.com/jackc/pgx/v4"
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

	sqlInsertUser := `INSERT INTO users (login, password) 
		VALUES ($1, $2) 
		ON CONFLICT (login) DO NOTHING 
		RETURNING id`

	row := tx.QueryRow(ctx, sqlInsertUser, login, passwordHash)
	var userID int
	err = row.Scan(&userID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return labelError.NewLabelError(labelError.TypeConflict, errors.New("login exists"))
		}
		return err
	}
	sqlInsertSession := `INSERT INTO sessions (session_id, created_at, lifetime, user_id) 
		VALUES ($1, $2, $3, $4)`

	timestamp := time.Now().Unix()
	lifetime := entity.TTL.Seconds()
	if _, err = tx.Exec(ctx, sqlInsertSession, sessionID, timestamp, lifetime, userID); err != nil {
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
) (int, error) {
	ctx, cancel := context.WithTimeout(ctx, 1*time.Second)
	defer cancel()
	sqlGetUser := "SELECT id FROM users WHERE login=$1 AND password=$2 LIMIT 1"
	row := r.conn.QueryRow(ctx, sqlGetUser, login, passwordHash)
	var userID int
	err := row.Scan(&userID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return 0, labelError.NewLabelError(labelError.TypeNotFound, errors.New("no registered user"))
		}
		return 0, err
	}
	return userID, nil
}

func (r *Repo) AddUserSession(ctx context.Context, user *entity.User) error {
	ctx, cancel := context.WithTimeout(ctx, 1*time.Second)
	defer cancel()
	sqlInsertSession := `INSERT INTO sessions (session_id, created_at, lifetime, user_id) 
		VALUES ($1, $2, $3, $4)`
	timestamp := time.Now().Unix()
	lifetime := entity.TTL.Seconds()
	_, err := r.conn.Exec(ctx, sqlInsertSession, user.ActiveSessionID, timestamp, lifetime, user.ID)
	if err != nil {
		return err
	}
	return nil
}

func (r *Repo) GetUserBySessionID(ctx context.Context, sessionID string) (*entity.User, error) {
	ctx, cancel := context.WithTimeout(ctx, 1*time.Second)
	defer cancel()
	sqlSelectUser := `SELECT u.id, u.login 
		FROM sessions s 
		JOIN users u ON u.id=s.user_id 
		WHERE session_id=$1 AND created_at > $2 - lifetime 
		LIMIT 1`
	timestamp := time.Now().Unix()
	row := r.conn.QueryRow(ctx, sqlSelectUser, sessionID, timestamp)
	var userID int
	var login string
	err := row.Scan(&userID, &login)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, labelError.NewLabelError(labelError.TypeNotFound, errors.New("no registered user"))
		}
		return nil, err
	}
	user := &entity.User{
		ID:              userID,
		Login:           login,
		ActiveSessionID: sessionID,
	}
	return user, nil
}
