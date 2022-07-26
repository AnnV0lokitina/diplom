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

	sqlCheck := "SELECT num, user_id FROM orders WHERE num=$1 LIMIT 1"
	row := tx.QueryRow(ctx, sqlCheck, orderNumber)

	var number string
	var userID int
	numberFind := true
	err = row.Scan(&number, &userID)
	if err != nil {
		if !errors.Is(err, pgx.ErrNoRows) {
			return err
		}
		numberFind = false
	}
	if numberFind {
		if user.ID == userID {
			return labelError.NewLabelError(labelError.TypeCreated, errors.New("number created"))
		}
		return labelError.NewLabelError(labelError.TypeConflict, errors.New("number exists"))
	}

	sqlInsert := "INSERT INTO orders (num, user_id) " +
		"VALUES ($1, $2) " +
		"ON CONFLICT (num) DO NOTHING"

	_, err = tx.Exec(ctx, sqlInsert, orderNumber, user.ID)
	if err != nil {
		return err
	}
	err = tx.Commit(ctx)
	if err != nil {
		return err
	}
	return nil
}

func (r *Repo) GetUserOrders(ctx context.Context, user *entity.User) ([]*entity.Order, error) {
	ctx, cancel := context.WithTimeout(ctx, 1*time.Second)
	defer cancel()
	sqlRequest := "SELECT o.num, o.uploaded_at, o.status, sum(t.delta) accrual " +
		"FROM orders o " +
		"LEFT JOIN transactions t " +
		"ON o.id=t.order_id AND t.operation_type=$1 " +
		"WHERE o.user_id=$2 " +
		"GROUP BY o.num, o.uploaded_at, o.status"
	rows, _ := r.conn.Query(ctx, sqlRequest, OperationAdd, user.ID)
	orders := make([]*entity.Order, 0)
	for rows.Next() {
		order := &entity.Order{}
		var accrual sql.NullInt64
		err := rows.Scan(&order.Number, &order.UploadedAt, &order.Status, &accrual)
		order.Login = user.Login
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
