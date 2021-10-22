-- name: CreateAccount :one
INSERT INTO accounts (
  owner, 
  balance, 
  currency
) VALUES (
  $1, $2, $3
)
RETURNING *;

-- name: GetAccount :many
SELECT * FROM accounts
WHERE id = $1 LIMIT 1;

-- name: ListAuthors :many
SELECT * FROM accounts
ORDER BY id
LIMIT $1
OFFSET $2;

-- name: UpdateAccount :one
UPDATE accounts 
SET balance = $2
WHERE id = $1
RETURNING *;

-- name: DeleteAuthor :exec
DELETE FROM accounts WHERE id = $1;