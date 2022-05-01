package repo

import (
	"context"
	"github.com/jackc/pgx/v4"
	"time"
)

type OperationType string

const (
	OperationAdd = "add"
	OperationSub = "sub"
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
