package repo

import (
	"context"
	"errors"
	"github.com/AnnV0lokitina/diplom/internal/entity"
	labelError "github.com/AnnV0lokitina/diplom/pkg/error"
	"github.com/jackc/pgx/v4"
	"time"
)

type IQueryRow interface {
	QueryRow(ctx context.Context, sql string, args ...interface{}) pgx.Row
}

func getUserBalanceFromRows(ctx context.Context, user *entity.User, client IQueryRow) (*entity.UserBalance, error) {
	sqlCheckBalance := `SELECT SUM(CASE 
				WHEN t.operation_type = 'add' then t.delta
				WHEN t.operation_type = 'sub' 0 then t.delta * -1
			END) current,
			SUM(CASE 
				WHEN t.operation_type = 'add' then 0
				WHEN t.operation_type = 'sub' then t.delta
			END) withdrawn 
		FROM orders o 
		JOIN transactions t ON o.id=t.order_id 
		WHERE o.user_id=$1`

	row := client.QueryRow(ctx, sqlCheckBalance, user.ID)

	userBalance := &entity.UserBalance{}
	err := row.Scan(&userBalance.Current, &userBalance.Withdrawn)
	if err != nil && !errors.Is(err, pgx.ErrNoRows) {
		return nil, err
	}
	return userBalance, nil
}

func (r *Repo) GetUserBalance(ctx context.Context, user *entity.User) (*entity.UserBalance, error) {
	ctx, cancel := context.WithTimeout(ctx, 1*time.Second)
	defer cancel()
	return getUserBalanceFromRows(ctx, user, r.conn)
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

	userBalance, err := getUserBalanceFromRows(ctx, user, tx)
	if err != nil {
		return err
	}
	if userBalance.Current-sum < 0 {
		return labelError.NewLabelError(labelError.TypeNotEnoughPoints, errors.New("no points"))
	}

	sqlAddOrder := `INSERT INTO orders (num, user_id) 
		VALUES ($1, $2) 
		ON CONFLICT (num) DO UPDATE SET num=orders.num 
		RETURNING id`

	row := tx.QueryRow(ctx, sqlAddOrder, orderNumber, user.ID)
	var orderID int
	err = row.Scan(&orderID)
	if err != nil {
		return err
	}

	sqlAddTransaction := "INSERT INTO transactions (operation_type, delta, order_id) VALUES ($1, $2, $3)"

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
	sql := `SELECT o.num, t.delta, t.created_at 
		FROM orders o 
		LEFT JOIN transactions t ON o.id=t.order_id AND t.operation_type=$1 
		WHERE o.user_id=$2 
		ORDER BY t.created_at DESC`
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
