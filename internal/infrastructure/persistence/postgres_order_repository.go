package persistence

import (
	"context"
	"errors"
	"fmt"
	"time"

	"api-ddd/internal/domain/order"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type PostgresOrderRepository struct {
	db *pgxpool.Pool
}

func NewPostgresOrderRepository(
	db *pgxpool.Pool,
) *PostgresOrderRepository {
	return &PostgresOrderRepository{
		db: db,
	}
}

func (r *PostgresOrderRepository) Save(
	ctx context.Context,
	o *order.Order,
) error {

	if o.ID() != 0 {
		return fmt.Errorf("el pedido ya está guardado")
	}

	tx, err := r.db.Begin(ctx)
	if err != nil {
		return fmt.Errorf("iniciando transacción: %w", err)
	}

	defer tx.Rollback(ctx)

	var orderID int64

	err = tx.QueryRow(ctx, `
		INSERT INTO orders (status, created_at)
		VALUES ($1, $2)
		RETURNING id
	`, o.Status(), o.CreatedAt()).Scan(&orderID)

	if err != nil {
		return fmt.Errorf("guardando pedido: %w", err)
	}

	for _, item := range o.Items() {
		_, err = tx.Exec(ctx, `
			INSERT INTO order_items (order_id, product_id, quantity)
			VALUES ($1, $2, $3)
		`, orderID, item.ProductID(), item.Quantity())

		if err != nil {
			return fmt.Errorf("guardando producto: %w", err)
		}
	}

	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("confirmando transacción: %w", err)
	}

	return o.AssignID(orderID)
}

func (r *PostgresOrderRepository) FindByID(
	ctx context.Context,
	id int64,
) (*order.Order, error) {

	var (
		orderID   int64
		status    string
		createdAt time.Time
	)

	err := r.db.QueryRow(ctx, `
		SELECT id, status, created_at
		FROM orders
		WHERE id = $1
	`, id).Scan(&orderID, &status, &createdAt)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, order.ErrOrderNotFound
		}

		return nil, fmt.Errorf("consultando pedido: %w", err)
	}

	rows, err := r.db.Query(ctx, `
		SELECT product_id, quantity
		FROM order_items
		WHERE order_id = $1
	`, orderID)

	if err != nil {
		return nil, fmt.Errorf("consultando productos: %w", err)
	}

	defer rows.Close()

	var items []order.OrderItem

	for rows.Next() {
		var (
			productID int64
			quantity  int
		)

		if err := rows.Scan(&productID, &quantity); err != nil {
			return nil, fmt.Errorf("leyendo producto: %w", err)
		}

		item, err := order.NewOrderItem(productID, quantity)
		if err != nil {
			return nil, fmt.Errorf("reconstruyendo producto: %w", err)
		}

		items = append(items, item)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("recorriendo productos: %w", err)
	}

	return order.RehydrateOrder(
		orderID,
		items,
		order.OrderStatus(status),
		createdAt,
	)
}

func (r *PostgresOrderRepository) Update(
	ctx context.Context,
	o *order.Order,
) error {

	result, err := r.db.Exec(ctx, `
		UPDATE orders
		SET status = $1
		WHERE id = $2
	`, o.Status(), o.ID())

	if err != nil {
		return fmt.Errorf("actualizando pedido: %w", err)
	}

	if result.RowsAffected() == 0 {
		return order.ErrOrderNotFound
	}

	return nil
}

// Verifica en compilación que implementamos la interfaz.
var _ order.OrderRepository = (*PostgresOrderRepository)(nil)