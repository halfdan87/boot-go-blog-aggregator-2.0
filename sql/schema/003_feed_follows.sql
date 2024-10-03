-- +goose Up
CREATE TABLE feed_follows (
    id UUID PRIMARY KEY,
    feed_id UUID NOT NULL references feeds(id) ON DELETE CASCADE,
    user_id UUID NOT NULL references users(id) ON DELETE CASCADE,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    unique(feed_id, user_id)
);

-- +goose Down
DROP TABLE feed_follows;
