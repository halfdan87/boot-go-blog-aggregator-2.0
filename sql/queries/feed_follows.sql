-- name: CreateFeedFollow :one
-- we need to return also the names of the user and the feed
WITH inserted_feed_follow AS (
    INSERT INTO feed_follows (id, feed_id, user_id, created_at, updated_at)
    VALUES (
        $1,
        $2,
        $3,
        now(),
        now()
    )
    RETURNING feed_id, user_id, created_at, updated_at
)
SELECT
    inserted_feed_follow.*,
    f.name as feed_name,
    u.name as user_name
FROM inserted_feed_follow
INNER JOIN feeds f ON f.id = inserted_feed_follow.feed_id
INNER JOIN users u ON u.id = inserted_feed_follow.user_id;

-- name: GetFeedFollowsForUser :many
SELECT
    ff.*,
    f.name as feed_name,
    u.name as user_name
FROM feed_follows ff
INNER JOIN users u ON u.id = ff.user_id
INNER JOIN feeds f ON f.id = ff.feed_id
WHERE ff.user_id = $1;


-- name: DeleteFeedFollowByUserIdAndFeedUrl :exec
DELETE FROM feed_follows ff
WHERE ff.user_id = $1 AND ff.feed_id = (SELECT f.id FROM feeds f WHERE f.url = $2);
