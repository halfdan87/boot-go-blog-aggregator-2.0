-- name: CreateFeed :one
INSERT INTO feeds (id, created_at, updated_at, name, url, user_id)
VALUES (
    $1,
	now(),
	now(),
    $2,
    $3,
	$4
)
RETURNING *;

-- name: GetAllFeeds :many
SELECT id, name, url, created_at, updated_at, user_id FROM feeds;

-- name: GetFeedByUrl :one
SELECT id, name, url, created_at, updated_at, user_id FROM feeds WHERE url = $1;

-- name: MarkFeedAsFetched :exec
UPDATE feeds SET last_fetched_at = now(), updated_at = now() WHERE id = $1;

-- name: GetNextFeedToFetch :one
SELECT id, name, url, created_at, updated_at, last_fetched_at, user_id FROM feeds 
ORDER BY last_fetched_at IS NULL, last_fetched_at ASC;


