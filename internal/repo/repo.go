package repo

import (
	"context"
	"errors"
	"github.com/AnnV0lokitina/diplom/internal/entity"
	labelError "github.com/AnnV0lokitina/diplom/pkg/error"
	"github.com/jackc/pgx/v4"
	"time"
)

type Repo struct {
	conn *pgx.Conn
}

func NewRepo(ctx context.Context, dsn string) (*Repo, error) {
	ctx, cancel := context.WithTimeout(ctx, 2*time.Second)
	defer cancel()
	conn, err := pgx.Connect(ctx, dsn)
	if err != nil {
		return nil, err
	}
	return &Repo{
		conn: conn,
	}, nil
}

func (r *Repo) Close(ctx context.Context) error {
	return r.conn.Close(ctx)
}

func (r *Repo) CreateUser(
	ctx context.Context,
	sessionID string,
	login string,
	passwordHash string,
) error {
	ctx, cancel := context.WithTimeout(ctx, 1*time.Second)
	defer cancel()
	sql := "INSERT INTO users (login, password, active_session_id) " +
		"VALUES ($1, $2, $3) " +
		"ON CONFLICT (login) DO NOTHING"
	result, err := r.conn.Exec(ctx, sql, login, passwordHash, sessionID)
	if err != nil {
		return err
	}
	if result.RowsAffected() == 0 {
		return labelError.NewLabelError(labelError.TypeConflict, errors.New("login exists"))
	}
	return nil
}

func (r *Repo) UpdateUserSession(
	ctx context.Context,
	sessionID string,
	login string,
	passwordHash string,
) error {
	ctx, cancel := context.WithTimeout(ctx, 1*time.Second)
	defer cancel()
	sql := "UPDATE users " +
		"SET active_session_id=$1 " +
		"WHERE login=$2 AND password=$3"
	result, err := r.conn.Exec(ctx, sql, sessionID, login, passwordHash)
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
	sql := "select login from users where active_session_id=$1"
	rows, _ := r.conn.Query(ctx, sql, activeSessionID)
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
