BEGIN;

CREATE TABLE roles (
    id          BIGSERIAL PRIMARY KEY,
    name        VARCHAR(32)  NOT NULL,
    description VARCHAR(255) NOT NULL DEFAULT '',
    created_at  TIMESTAMPTZ  NOT NULL DEFAULT now(),
    updated_at  TIMESTAMPTZ,
    deleted_at  TIMESTAMPTZ,
    created_by  BIGINT,
    updated_by  BIGINT,
    deleted_by  BIGINT
);

CREATE UNIQUE INDEX ux_roles_name ON roles (name) WHERE deleted_at IS NULL;
CREATE INDEX ix_roles_deleted_at ON roles (deleted_at);

CREATE TABLE users (
    id            BIGSERIAL PRIMARY KEY,
    username      VARCHAR(64) NOT NULL,
    password_hash TEXT        NOT NULL,
    enabled       BOOLEAN     NOT NULL DEFAULT TRUE,
    created_at    TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at    TIMESTAMPTZ,
    deleted_at    TIMESTAMPTZ,
    created_by    BIGINT,
    updated_by    BIGINT,
    deleted_by    BIGINT
);

CREATE UNIQUE INDEX ux_users_username
    ON users (username) WHERE deleted_at IS NULL;
CREATE INDEX ix_users_deleted_at ON users (deleted_at);

CREATE TABLE user_roles (
    user_id BIGINT NOT NULL REFERENCES users (id) ON DELETE CASCADE,
    role_id BIGINT NOT NULL REFERENCES roles (id) ON DELETE CASCADE,
    PRIMARY KEY (user_id, role_id)
);

CREATE INDEX ix_user_roles_role_id ON user_roles (role_id);

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
    updated_at TIMESTAMPTZ,
    deleted_at TIMESTAMPTZ,
    created_by BIGINT,
    updated_by BIGINT,
    deleted_by BIGINT,
    CONSTRAINT ck_todos_priority CHECK (priority IN ('low', 'medium', 'high'))
);

CREATE INDEX ix_todos_user_created
    ON todos (user_id, created_at DESC)
    WHERE deleted_at IS NULL;

CREATE INDEX ix_todos_tags ON todos USING GIN (tags);

COMMIT;