package repo

import (
	"context"
	"errors"
	"github.com/AnnV0lokitina/diplom/internal/entity"
	labelError "github.com/AnnV0lokitina/diplom/pkg/error"
	"github.com/jackc/pgx/v4"
	"time"
)

func (r *Repo) AddOrderInfo(ctx context.Context, orderInfo *entity.OrderUpdateInfo) error {
	ctx, cancel := context.WithTimeout(ctx, 1*time.Second)
	defer cancel()

	tx, err := r.conn.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	sqlSetStatus := `UPDATE orders 
		SET status=$1 
		WHERE num=$2 
		RETURNING id`

	row := r.conn.QueryRow(ctx, sqlSetStatus, orderInfo.Status, orderInfo.Number)
	var orderID int
	err = row.Scan(&orderID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return labelError.NewLabelError(labelError.TypeNotFound, errors.New("no order number"))
		}
		return err
	}

	if orderInfo.Accrual == 0 {
		err = tx.Commit(ctx)
		if err != nil {
			return err
		}
		return nil
	}

	sqlChangeBalance := "INSERT INTO transactions (operation_type, delta, order_id) VALUES ($1, $2, $3)"

	if _, err = tx.Exec(ctx, sqlChangeBalance, OperationAdd, orderInfo.Accrual, orderID); err != nil {
		return err
	}
	err = tx.Commit(ctx)
	if err != nil {
		return err
	}
	return nil
}

func (r *Repo) GetOrdersListForRequest(ctx context.Context) ([]entity.OrderNumber, error) {
	ctx, cancel := context.WithTimeout(ctx, 1*time.Second)
	defer cancel()
	sql := "SELECT num FROM orders WHERE status!=$1 AND status!=$2"
	rows, _ := r.conn.Query(ctx, sql, entity.OrderStatusProcessed, entity.OrderStatusInvalid)
	ordersList := make([]entity.OrderNumber, 0)
	for rows.Next() {
		var order string
		err := rows.Scan(&order)
		if err != nil {
			return nil, err
		}
		ordersList = append(ordersList, entity.OrderNumber(order))
	}
	if len(ordersList) == 0 {
		return nil, labelError.NewLabelError(labelError.TypeNotFound, errors.New("not found"))
	}
	return ordersList, nil
}
