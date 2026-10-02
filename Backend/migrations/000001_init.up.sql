BEGIN;

CREATE TABLE users (
    id            BIGSERIAL PRIMARY KEY,
    username      VARCHAR(64)  NOT NULL UNIQUE,
    password_hash TEXT         NOT NULL,           
    created_at    TIMESTAMPTZ  NOT NULL DEFAULT now()
);

CREATE TABLE sessions (
    token_hash CHAR(64)    PRIMARY KEY,             
    user_id    BIGINT      NOT NULL REFERENCES users (id) ON DELETE CASCADE,
    expires_at TIMESTAMPTZ NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX ix_sessions_expires_at ON sessions (expires_at);

CREATE TABLE todos (
    id         BIGSERIAL PRIMARY KEY,
    user_id    BIGINT       NOT NULL REFERENCES users (id) ON DELETE CASCADE,
    title      VARCHAR(200) NOT NULL,
    done       BOOLEAN      NOT NULL DEFAULT FALSE,
    due_date   TIMESTAMPTZ,                         
    category   VARCHAR(50),                         
    priority   VARCHAR(10)  NOT NULL DEFAULT 'medium',
    tags       TEXT[]       NOT NULL DEFAULT '{}',
    created_at TIMESTAMPTZ  NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ  NOT NULL DEFAULT now(),

    CONSTRAINT ck_todos_priority CHECK (priority IN ('low', 'medium', 'high'))
);


CREATE INDEX ix_todos_user_created ON todos (user_id, created_at DESC);

CREATE INDEX ix_todos_tags ON todos USING GIN (tags);

COMMIT;