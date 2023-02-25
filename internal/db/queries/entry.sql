-- name: GetEntry :one
SELECT * FROM entries
WHERE id = $1
LIMIT 1;