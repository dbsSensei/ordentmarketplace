package db

import (
	"context"
	"database/sql"
	"fmt"
	"time"
)

// Store defines all functions to execute db queries and transactions
type Store interface {
	Querier
	OrderTx(ctx context.Context, arg OrderTxParams) (OrderTxResult, error)
}

// SQLStore provides all functions to execute SQL queries and transactions
type SQLStore struct {
	db *sql.DB
	*Queries
}

// NewStore creates a new store
func NewStore(db *sql.DB) Store {
	return &SQLStore{
		db:      db,
		Queries: New(db),
	}
}

// ExecTx executes a function within a database transaction
func (store *SQLStore) execTx(ctx context.Context, fn func(*Queries) error) error {
	tx, err := store.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	q := New(tx)
	err = fn(q)
	if err != nil {
		if rbErr := tx.Rollback(); rbErr != nil {
			return fmt.Errorf("tx err: %v, rb err: %v", err, rbErr)
		}
		return err
	}

	return tx.Commit()
}

// OrderTxParams contains the input parameters of the order
type OrderTxParams struct {
	UserID int32               `json:"user_id"`
	Status string              `json:"status"`
	Items  []OrderItemTxParams `json:"items" binding:"required"`
}

type OrderItemTxParams struct {
	OrderID   int32 `json:"order_id"`
	ProductID int32 `json:"product_id" binding:"required"`
	Quantity  int32 `json:"quantity" binding:"required"`
}

// OrderTxResult is the result of the transfer transaction
type OrderTxResult struct {
	ID        int32       `json:"id"`
	Status    string      `json:"status"`
	Items     []OrderItem `json:"items"`
	CreatedAt time.Time   `json:"created_at"`
}

// OrderTx performs a money transfer from one account to the other.
// It creates the transfer, add account entries, and update accounts' balance within a database transaction
func (store *SQLStore) OrderTx(ctx context.Context, arg OrderTxParams) (OrderTxResult, error) {
	var result OrderTxResult

	err := store.execTx(ctx, func(q *Queries) error {
		var err error

		orderResult, err := q.CreateOrder(ctx, CreateOrderParams{
			UserID: arg.UserID,
			Status: arg.Status,
		})
		if err != nil {
			return err
		}

		var items []OrderItem
		for _, item := range arg.Items {
			orderItemResult, err := q.CreateOrderItem(ctx, CreateOrderItemParams{
				OrderID:   orderResult.ID,
				ProductID: item.ProductID,
				Quantity:  item.Quantity,
			})
			if err != nil {
				return err
			}
			items = append(items, orderItemResult)
		}

		result = OrderTxResult{
			ID:        orderResult.ID,
			Status:    orderResult.Status,
			Items:     items,
			CreatedAt: orderResult.CreatedAt,
		}
		return err
	})

	return result, err
}
