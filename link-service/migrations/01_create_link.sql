CREATE TABLE IF NOT EXISTS links (
    id BIGSERIAL PRIMARY KEY,
    short_code VARCHAR(10) NOT NULL UNIQUE,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    original_url TEXT NOT NULL,
    visits BIGINT NOT NULL DEFAULT 0
);

CREATE INDEX idx_links_short_code ON links(short_code);