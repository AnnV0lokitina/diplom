package repo

import (
	"context"
	"errors"
	"github.com/AnnV0lokitina/diplom/internal/entity"
	labelError "github.com/AnnV0lokitina/diplom/pkg/error"
	"github.com/jackc/pgx/v4"
	"time"
)

const (
	OperationAdd = 0
	OperationSub = 1
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

func (r *Repo) GetUserOrderByNumber(
	ctx context.Context,
	user *entity.User,
	orderNumber entity.OrderNumber,
) (*entity.Order, error) {
	ctx, cancel := context.WithTimeout(ctx, 1*time.Second)
	defer cancel()
	sql := "SELECT id, num, login, uploaded_at, status " +
		"FROM orders " +
		"WHERE login=$1 AND num=$2 " +
		"LIMIT 1"
	rows, _ := r.conn.Query(ctx, sql, user.Login, orderNumber)
	orders := make([]*entity.Order, 0)
	for rows.Next() {
		order := &entity.Order{}
		err := rows.Scan(&order.ID, &order.Number, &order.Login, &order.UploadedAt, &order.Status)
		if err != nil {
			return nil, err
		}
		orders = append(orders, order)
	}
	if len(orders) == 0 {
		return nil, labelError.NewLabelError(labelError.TypeNotFound, errors.New("not found"))
	}
	return orders[0], nil
}

func (r *Repo) GetUserWithdrawals(ctx context.Context, user *entity.User) ([]*entity.Withdrawal, error) {
	ctx, cancel := context.WithTimeout(ctx, 1*time.Second)
	defer cancel()
	sql := "SELECT o.num, b.delta, b.created_at " +
		"FROM orders o " +
		"LEFT JOIN balance b ON o.num=b.num AND b.operation_type=$1 " +
		"WHERE o.login=$2 " +
		"ORDER BY b.created_at DESC"
	rows, _ := r.conn.Query(ctx, sql, OperationSub, user.Login)
	withdrawals := make([]*entity.Withdrawal, 0)
	for rows.Next() {
		w := &entity.Withdrawal{}
		err := rows.Scan(&w.OrderNumber, &w.Sum, &w.ProcessedAt)
		if err != nil {
			return nil, err
		}
		withdrawals = append(withdrawals, w)
	}
	if len(withdrawals) == 0 {
		return nil, labelError.NewLabelError(labelError.TypeNotFound, errors.New("not found"))
	}
	return withdrawals, nil
}

func (r *Repo) AddOrderInfo(
	ctx context.Context,
	orderNumber entity.OrderNumber,
	status entity.OrderStatus,
	accrual entity.PointValue,
) error {
	ctx, cancel := context.WithTimeout(ctx, 1*time.Second)
	defer cancel()

	tx, err := r.conn.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	sqlSetStatus := "UPDATE orders " +
		"SET status=$1 " +
		"WHERE num=$2 "
	//"RETURNING id"

	_, err = tx.Exec(ctx, sqlSetStatus, status, orderNumber)
	if err != nil {
		return err
	}
	//id := 0
	//for rows.Next() {
	//	err := rows.Scan(&id)
	//	if err != nil {
	//		return err
	//	}
	//}

	//if id == 0 {
	//	return labelError.NewLabelError(labelError.TypeNotFound, errors.New("no order number"))
	//}

	if accrual == 0 {
		err = tx.Commit(ctx)
		if err != nil {
			return err
		}
		return nil
	}

	sqlChangeBalance := "INSERT INTO balance (operation_type, delta, num) " +
		"VALUES ($1, $2, $3)"

	if _, err = tx.Exec(ctx, sqlChangeBalance, OperationAdd, accrual, orderNumber); err != nil {
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
