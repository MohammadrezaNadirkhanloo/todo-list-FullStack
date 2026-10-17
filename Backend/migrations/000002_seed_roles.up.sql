BEGIN;

INSERT INTO roles (name, description)
VALUES
    ('admin', 'دسترسی کامل مدیریتی'),
    ('user',  'کاربر عادی')
ON CONFLICT DO NOTHING;

COMMIT;