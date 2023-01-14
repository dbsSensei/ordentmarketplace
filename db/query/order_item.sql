-- name: CreateOrderItem :one
INSERT INTO order_items (order_id,
                    product_id,
                    quantity)
VALUES ($1, $2, $3) RETURNING *;

-- name: GetOrderItems :many
SELECT *
FROM order_items
WHERE order_id = $1;
