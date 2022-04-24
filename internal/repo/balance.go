package repo

import (
	"context"
	"errors"
	"fmt"
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
	fmt.Println(len(sums))
	if len(sums) == 0 {
		return &entity.UserBalance{
			Current:   0,
			Withdrawn: 0,
		}, nil
	}
	var addSum int
	var subSub int
	fmt.Println(addSum)
	fmt.Println(subSub)
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
	sql := "SELECT SUM(balance.delta) sum, balance.operation_type " +
		"FROM orders " +
		"JOIN balance ON orders.num=balance.num " +
		"WHERE orders.login=$1 " +
		"GROUP BY operation_type"
	rows, _ := r.conn.Query(ctx, sql, user.Login)
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

	batch := &pgx.Batch{}

	sql1 := "SELECT SUM(balance.delta) sum, balance.operation_type " +
		"FROM orders " +
		"JOIN balance ON orders.num=balance.num " +
		"WHERE orders.login=$1 " +
		"GROUP BY operation_type"
	_, err = tx.Prepare(ctx, "check_balance", sql1)
	if err != nil {
		return err
	}
	batch.Queue("check_balance", user.Login)

	sql2 := "INSERT INTO balance (operation_type, delta, num) " +
		"VALUES ($1, $2, $3)"

	_, err = tx.Prepare(ctx, "insert_withdraw", sql2)
	if err != nil {
		return err
	}
	batch.Queue("insert_withdraw", OperationSub, sum, orderNumber)

	br := tx.SendBatch(ctx, batch)

	rows, err := br.Query()
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
	_, err = br.Exec()
	if err != nil {
		return err
	}
	br.Close()
	tx.Commit(ctx)
	return nil
}
