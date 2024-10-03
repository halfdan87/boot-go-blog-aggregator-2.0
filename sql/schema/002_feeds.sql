-- +goose Up
CREATE TABLE feeds (
    id UUID PRIMARY KEY,
	name VARCHAR(255) NOT NULL,
	url VARCHAR(255) NOT NULL UNIQUE,
	created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
	updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
	user_id UUID NOT NULL references users(id) ON DELETE CASCADE
);

-- +goose Down
DROP TABLE feeds;
