CREATE TABLE IF NOT EXISTS highlights (
    id uuid PRIMARY KEY DEFAULT uuid_generate_v4(),
    url VARCHAR(3096),
    title VARCHAR(512),
    created_at TIMESTAMP,

    reddit_submission_id VARCHAR(10),
    reddit_permalink VARCHAR(3096),
    reddit_author VARCHAR(64),
    reddit_created_at TIMESTAMP,

    CONSTRAINT highlights_reddit_submission_id_unique UNIQUE (reddit_submission_id)
);

CREATE TABLE IF NOT EXISTS highlight_mirrors (
    id uuid PRIMARY KEY DEFAULT uuid_generate_v4(),
    highlight_id uuid,
    url VARCHAR(3096),
    created_at TIMESTAMP,

    reddit_permalink VARCHAR(3096),
    reddit_author VARCHAR(64),
    reddit_created_at TIMESTAMP,

    CONSTRAINT highlight_mirrors_highlight_id_url_unique UNIQUE (highlight_id, url),

    CONSTRAINT highlight_mirrors_highlight_id_fk FOREIGN KEY (highlight_id) REFERENCES highlights (id)
);
