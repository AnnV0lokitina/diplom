package repo

import (
	"context"
	"database/sql"
	"errors"
	"github.com/AnnV0lokitina/diplom/internal/entity"
	labelError "github.com/AnnV0lokitina/diplom/pkg/error"
	"github.com/jackc/pgx/v4"
	"time"
)

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

func (r *Repo) GetUserOrders(ctx context.Context, user *entity.User) ([]*entity.Order, error) {
	ctx, cancel := context.WithTimeout(ctx, 1*time.Second)
	defer cancel()
	sqlRequest := "SELECT o.num, o.login, o.uploaded_at, o.status, sum(b.delta) accrual " +
		"FROM orders o " +
		"LEFT JOIN balance b " +
		"ON o.id=b.order_id AND operation_type=$1 " +
		"WHERE o.login=$2 " +
		"GROUP BY o.num, o.login, o.uploaded_at, o.status"
	rows, _ := r.conn.Query(ctx, sqlRequest, OperationAdd, user.Login)
	orders := make([]*entity.Order, 0)
	for rows.Next() {
		order := &entity.Order{}
		var accrual sql.NullInt64
		err := rows.Scan(&order.Number, &order.Login, &order.UploadedAt, &order.Status, &accrual)
		if accrual.Valid {
			order.Accrual = entity.PointValue(accrual.Int64)
		}
		if err != nil {
			return nil, err
		}
		orders = append(orders, order)
	}
	if len(orders) == 0 {
		return nil, labelError.NewLabelError(labelError.TypeNotFound, errors.New("not found"))
	}
	l := len(orders)
	return orders[:l:l], nil
}
