CREATE TABLE posts (
    id BIGSERIAL PRIMARY KEY,
    content TEXT NOT NULL,
    title TEXT NOT NULL,
    user_id BIGINT NOT NULL,
    tags TEXT[],
    created_at TIMESTAMP NOT NULL DEFAULT NOW()
);
