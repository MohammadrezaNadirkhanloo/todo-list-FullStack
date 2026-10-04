BEGIN;

CREATE TABLE users (
    id            BIGSERIAL PRIMARY KEY,
    username      VARCHAR(64) NOT NULL,
    password_hash TEXT        NOT NULL,       
    created_at    TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at    TIMESTAMPTZ,
    deleted_at    TIMESTAMPTZ,
    created_by    BIGINT,
    updated_by    BIGINT,
    deleted_by    BIGINT
);

-- یکتا فقط بین کاربران حذف‌نشده؛ نام کاربری حذف‌شده برای همیشه اشغال نمی‌ماند.
CREATE UNIQUE INDEX ux_users_username
    ON users (username) WHERE deleted_at IS NULL;

CREATE TABLE todos (
    id         BIGSERIAL PRIMARY KEY,
    user_id    BIGINT       NOT NULL REFERENCES users (id) ON DELETE CASCADE,
    title      VARCHAR(200) NOT NULL,
    done       BOOLEAN      NOT NULL DEFAULT FALSE,
    due_date   TIMESTAMPTZ,                       -- NULL یعنی بدون مهلت (تاریخ + ساعت فرم در یک ستون)
    category   VARCHAR(50),                       -- NULL یعنی دسته انتخاب نشده
    priority   VARCHAR(10)  NOT NULL DEFAULT 'medium',
    tags       TEXT[]       NOT NULL DEFAULT '{}',

    -- فیلدهای BaseModel
    created_at TIMESTAMPTZ  NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ,
    deleted_at TIMESTAMPTZ,
    created_by BIGINT,
    updated_by BIGINT,
    deleted_by BIGINT,

    CONSTRAINT ck_todos_priority CHECK (priority IN ('low', 'medium', 'high'))
);

-- فهرست تسک‌های یک کاربر، جدیدترین اول؛ ردیف‌های حذف‌شده در ایندکس نمی‌آیند.
CREATE INDEX ix_todos_user_created
    ON todos (user_id, created_at DESC)
    WHERE deleted_at IS NULL;

-- جستجو روی تگ‌ها (مثلاً tags @> ARRAY['work'])
CREATE INDEX ix_todos_tags ON todos USING GIN (tags);

COMMIT;