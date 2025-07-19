-- name: InsertUser :one
INSERT INTO users (
  first_name,
  last_name,
  email,
  phone,
  password,
  created_at,
  updated_at
) VALUES (
  $1, $2, $3, $4, $5, $6, $7
) RETURNING *;

-- name: UpdateUser :execresult
UPDATE users
SET
  first_name = $2,
  last_name = $3,
  email = $4,
  phone = $5,
  updated_at = $6
WHERE id = $1;

-- name: UpdateUserPassword :execresult
UPDATE users
SET
  password = $2,
  updated_at = $3
WHERE id = $1;

-- name: GetUserByID :one
SELECT * FROM users
WHERE id = $1;

-- name: GetUserByEmail :one
SELECT * FROM users
WHERE email = $1;

-- name: GetPublicProfileByIds :many
SELECT id, first_name, last_name FROM users
WHERE id = ANY(sqlc.arg(user_ids)::string[]);
