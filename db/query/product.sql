-- name: CreateProduct :one
INSERT INTO products (name,
                      merchant_id,
                      price,
                      status,
                      category_id)
VALUES ($1, $2, $3, $4, $5) RETURNING *;

-- name: GetProducts :many
SELECT *
FROM products;

-- name: GetProduct :one
SELECT *
FROM products
WHERE id = $1 LIMIT 1;

-- name: UpdateProduct :one
UPDATE products
SET
    name = COALESCE(sqlc.narg(name), name),
    price = COALESCE(sqlc.narg(price), price),
    status = COALESCE(sqlc.narg(status), status),
    category_id = COALESCE(sqlc.narg(category_id), category_id)
WHERE id = sqlc.arg(id) AND merchant_id = sqlc.arg(merchant_id) RETURNING *;

-- name: DeleteProduct :exec
DELETE FROM products
WHERE id = $1;