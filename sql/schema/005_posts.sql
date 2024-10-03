-- +goose Up
CREATE TABLE posts (
    id UUID PRIMARY KEY,
	created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
	updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    title VARCHAR(255) NOT NULL,
    url VARCHAR(255) NOT NULL UNIQUE,
    description TEXT NULL,
    published_at TIMESTAMP NOT NULL,
    feed_id UUID NOT NULL references feeds(id) ON DELETE CASCADE
);

-- +goose Down
DROP TABLE posts;
