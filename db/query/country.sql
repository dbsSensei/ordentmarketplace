-- name: CreateCountry :one
INSERT INTO countries (code, name)
VALUES ($1, $2) RETURNING *;

-- name: GetCountry :one
SELECT *
FROM countries
WHERE code = $1 LIMIT 1;

-- name: CountCountry :one
SELECT COUNT(*)
FROM countries;