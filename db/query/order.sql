-- name: CreateOrder :one
INSERT INTO orders (
  user_id,
  status
) VALUES (
  $1, $2
) RETURNING *;

-- name: GetOrder :one
SELECT * FROM orders
WHERE id = $1 AND user_id = $2 LIMIT 1;

-- name: GetOrders :many
SELECT * FROM orders
WHERE user_id = $1;

-- name: UpdateOrderStatus :one
UPDATE orders
SET
    status = $1
WHERE id = sqlc.arg(id) RETURNING *;
