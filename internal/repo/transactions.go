package repo

import (
	"context"
	"errors"
	"github.com/AnnV0lokitina/diplom/internal/entity"
	labelError "github.com/AnnV0lokitina/diplom/pkg/error"
	"github.com/jackc/pgx/v4"
	"time"
)

func getUserBalanceFromRows(rows pgx.Rows) (*entity.UserBalance, error) {
	sums := make(map[int]int)
	for rows.Next() {
		var operationType int
		var sum int
		err := rows.Scan(&sum, &operationType)
		if err != nil {
			return nil, err
		}
		sums[operationType] = sum
	}
	if len(sums) == 0 {
		return &entity.UserBalance{
			Current:   0,
			Withdrawn: 0,
		}, nil
	}
	var addSum int
	var subSub int
	subSub, ok := sums[OperationSub]
	if !ok {
		subSub = 0
	}
	addSum, ok = sums[OperationAdd]
	if !ok {
		addSum = 0
	}
	return &entity.UserBalance{
		Current:   entity.PointValue(addSum - subSub),
		Withdrawn: entity.PointValue(subSub),
	}, nil
}

func (r *Repo) GetUserBalance(ctx context.Context, user *entity.User) (*entity.UserBalance, error) {
	ctx, cancel := context.WithTimeout(ctx, 1*time.Second)
	defer cancel()
	sql := "SELECT SUM(t.delta) sum, t.operation_type " +
		"FROM orders o " +
		"JOIN transactions t ON o.id=t.order_id " +
		"WHERE o.user_id=$1 " +
		"GROUP BY t.operation_type"
	rows, _ := r.conn.Query(ctx, sql, user.ID)
	return getUserBalanceFromRows(rows)
}

func (r *Repo) UserOrderWithdraw(
	ctx context.Context,
	user *entity.User,
	orderNumber entity.OrderNumber,
	sum entity.PointValue,
) error {
	ctx, cancel := context.WithTimeout(ctx, 1*time.Second)
	defer cancel()

	tx, err := r.conn.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	sqlCheckBalance := "SELECT SUM(t.delta) sum, t.operation_type " +
		"FROM orders o " +
		"JOIN transactions t ON o.id=t.order_id " +
		"WHERE o.user_id=$1 " +
		"GROUP BY t.operation_type"

	rows, err := tx.Query(ctx, sqlCheckBalance, user.ID)
	if err != nil {
		return err
	}
	userBalance, err := getUserBalanceFromRows(rows)
	if err != nil {
		return err
	}
	if userBalance.Current-sum < 0 {
		return labelError.NewLabelError(labelError.TypeNotEnoughPoints, errors.New("not found"))
	}

	sqlAddOrder := "INSERT INTO orders (num, user_id) " +
		"VALUES ($1, $2) " +
		"ON CONFLICT (num) DO UPDATE SET num=orders.num " +
		"RETURNING id"

	row := tx.QueryRow(ctx, sqlAddOrder, orderNumber, user.ID)
	var orderID int
	err = row.Scan(&orderID)
	if err != nil {
		return err
	}

	sqlAddTransaction := "INSERT INTO transactions (operation_type, delta, order_id) " +
		"VALUES ($1, $2, $3)"

	_, err = tx.Exec(ctx, sqlAddTransaction, OperationSub, sum, orderID)
	if err != nil {
		return err
	}

	tx.Commit(ctx)
	return nil
}

func (r *Repo) GetUserWithdrawals(ctx context.Context, user *entity.User) ([]*entity.Withdrawal, error) {
	ctx, cancel := context.WithTimeout(ctx, 1*time.Second)
	defer cancel()
	sql := "SELECT o.num, t.delta, t.created_at " +
		"FROM orders o " +
		"LEFT JOIN transactions t ON o.id=t.order_id AND t.operation_type=$1 " +
		"WHERE o.user_id=$2 " +
		"ORDER BY t.created_at DESC"
	rows, _ := r.conn.Query(ctx, sql, OperationSub, user.ID)
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
