-- name: CreateChirp :one
INSERT INTO chirps (body, user_id, created_at, updated_at)
VALUES ($1, $2, NOW(), NOW())
RETURNING id,
    created_at,
    updated_at,
    body,
    user_id;
-- name: GetAllChirps :many
SELECT *
FROM chirps
ORDER BY created_at ASC;
-- name: GetChirpById :one
SELECT *
FROM chirps
WHERE id = $1;
-- name: DeleteChirp :exec
DELETE FROM chirps
Where id = $1;
-- name: GetChirpsByUserId :many
SELECT *
FROM chirps
WHERE user_id = $1;