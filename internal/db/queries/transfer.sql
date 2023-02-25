-- name: GetTransfer :one
SELECT * FROM transfers
WHERE id = $1
LIMIT 1;