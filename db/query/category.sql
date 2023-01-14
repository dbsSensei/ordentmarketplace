-- name: CreateCategory :one
INSERT INTO categories (name,
                        parent_id)
VALUES ($1, $2) RETURNING *;

-- name: GetCategories :many
SELECT *
FROM categories;

-- name: GetCategory :one
SELECT *
FROM categories
WHERE id = $1 LIMIT 1;

-- name: DeleteCategory :exec
DELETE FROM categories
WHERE id = $1;