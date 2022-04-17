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

func (r *Repo) AddOrder(ctx context.Context, user *entity.User, orderNumber entity.OrderNumber) error {
	ctx, cancel := context.WithTimeout(ctx, 1*time.Second)
	defer cancel()

	tx, err := r.conn.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	batch := &pgx.Batch{}

	sql1 := "SELECT num, login FROM orders WHERE num=$1 LIMIT 1"
	_, err = tx.Prepare(ctx, "check", sql1)
	if err != nil {
		return err
	}
	batch.Queue("check", orderNumber)

	sql2 := "INSERT INTO orders (num, login) " +
		"VALUES ($1, $2) " +
		"ON CONFLICT (num) DO NOTHING"

	_, err = tx.Prepare(ctx, "insert", sql2)
	if err != nil {
		return err
	}
	batch.Queue("insert", orderNumber, user.Login)

	br := tx.SendBatch(ctx, batch)

	var number string
	var login string
	numberFind := true
	err = br.QueryRow().Scan(&number, &login)
	if err != nil {
		if !errors.Is(err, pgx.ErrNoRows) {
			return err
		}
		numberFind = false
	}
	if numberFind {
		if user.Login == login {
			return labelError.NewLabelError(labelError.TypeCreated, errors.New("number created"))
		}
		return labelError.NewLabelError(labelError.TypeConflict, errors.New("number exists"))
	}
	_, err = br.Exec()
	if err != nil {
		return err
	}
	br.Close()
	tx.Commit(ctx)
	return nil
}
