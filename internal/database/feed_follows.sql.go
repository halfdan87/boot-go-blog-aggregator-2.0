// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.27.0
// source: feed_follows.sql

package database

import (
	"context"
	"time"

	"github.com/google/uuid"
)

const createFeedFollow = `-- name: CreateFeedFollow :one
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
    inserted_feed_follow.feed_id, inserted_feed_follow.user_id, inserted_feed_follow.created_at, inserted_feed_follow.updated_at,
    f.name as feed_name,
    u.name as user_name
FROM inserted_feed_follow
INNER JOIN feeds f ON f.id = inserted_feed_follow.feed_id
INNER JOIN users u ON u.id = inserted_feed_follow.user_id
`

type CreateFeedFollowParams struct {
	ID     uuid.UUID
	FeedID uuid.UUID
	UserID uuid.UUID
}

type CreateFeedFollowRow struct {
	FeedID    uuid.UUID
	UserID    uuid.UUID
	CreatedAt time.Time
	UpdatedAt time.Time
	FeedName  string
	UserName  string
}

// we need to return also the names of the user and the feed
func (q *Queries) CreateFeedFollow(ctx context.Context, arg CreateFeedFollowParams) (CreateFeedFollowRow, error) {
	row := q.db.QueryRowContext(ctx, createFeedFollow, arg.ID, arg.FeedID, arg.UserID)
	var i CreateFeedFollowRow
	err := row.Scan(
		&i.FeedID,
		&i.UserID,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.FeedName,
		&i.UserName,
	)
	return i, err
}

const getFeedFollowsForUser = `-- name: GetFeedFollowsForUser :many
SELECT
    ff.id, ff.feed_id, ff.user_id, ff.created_at, ff.updated_at,
    f.name as feed_name,
    u.name as user_name
FROM feed_follows ff
INNER JOIN users u ON u.id = ff.user_id
INNER JOIN feeds f ON f.id = ff.feed_id
WHERE ff.user_id = $1
`

type GetFeedFollowsForUserRow struct {
	ID        uuid.UUID
	FeedID    uuid.UUID
	UserID    uuid.UUID
	CreatedAt time.Time
	UpdatedAt time.Time
	FeedName  string
	UserName  string
}

func (q *Queries) GetFeedFollowsForUser(ctx context.Context, userID uuid.UUID) ([]GetFeedFollowsForUserRow, error) {
	rows, err := q.db.QueryContext(ctx, getFeedFollowsForUser, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []GetFeedFollowsForUserRow
	for rows.Next() {
		var i GetFeedFollowsForUserRow
		if err := rows.Scan(
			&i.ID,
			&i.FeedID,
			&i.UserID,
			&i.CreatedAt,
			&i.UpdatedAt,
			&i.FeedName,
			&i.UserName,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}
