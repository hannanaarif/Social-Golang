-- Down Migration
DROP INDEX IF EXISTS idx_posts_tags;

ALTER TABLE posts
ALTER COLUMN tags TYPE VARCHAR USING tags::VARCHAR;

CREATE INDEX IF NOT EXISTS idx_posts_tags ON posts USING gin (tags gin_trgm_ops);
