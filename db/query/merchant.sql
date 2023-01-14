-- name: CreateMerchant :one
INSERT INTO merchants (admin_id,
                        name,
                        country_code)
VALUES ($1, $2, $3) RETURNING *;

-- name: GetMerchantByAdminID :one
SELECT *
FROM merchants
WHERE admin_id = $1 LIMIT 1;

-- name: UpdateMerchant :one
UPDATE merchants
SET
    name = COALESCE(sqlc.narg(name), name),
    country_code = COALESCE(sqlc.narg(country_code), country_code)
WHERE admin_id = sqlc.arg(admin_id) RETURNING *;