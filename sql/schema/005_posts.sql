-- +goose Up
CREATE TABLE posts (
    id UUID PRIMARY KEY NOT NULL,
    created_at TIMESTAMP,
    updated_at TIMESTAMP,
    title TEXT NOT NULL,
    url TEXT UNIQUE NOT NULL,
    description TEXT NOT NULL,
    published_at TIMESTAMP,
    feed_id UUID UNIQUE NOT NULL,
    CONSTRAINT fk_feed_id
        FOREIGN KEY (feed_id)
        REFERENCES feeds(id)
        ON DELETE CASCADE
);

-- +goose Down
DROP TABLE posts;