-- name: CreateChirp :one
INSERT INTO chirps (id, created_at, updated_at, body, user_id)
VALUES ($1, $2, $3, $4, $5)
RETURNING *;

-- name: GetChirps :many
SELECT * FROM chirps
WHERE (
    sqlc.narg('user_id')::uuid IS NULL 
    OR user_id = sqlc.narg('user_id')::uuid
)
ORDER BY
  CASE
    WHEN sqlc.narg('sort_dir') = 'asc'  THEN created_at
  END ASC,
  CASE
    WHEN sqlc.narg('sort_dir') = 'desc' THEN created_at
  END DESC;

-- name: GetChirpById :one
SELECT * FROM chirps WHERE id = $1;

-- name: DeleteChirpById :exec
DELETE FROM chirps WHERE id = $1;
