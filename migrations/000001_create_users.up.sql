BEGIN;

CREATE TYPE user_role AS ENUM ('student', 'teacher', 'admin');

CREATE TABLE users (
                       id              SERIAL PRIMARY KEY,
                       full_name       VARCHAR(255) NOT NULL,
                       email           VARCHAR(255) NOT NULL UNIQUE,
                       password_hash   TEXT NOT NULL,
                       role            user_role NOT NULL DEFAULT 'student',
                       is_active       BOOLEAN NOT NULL DEFAULT TRUE,
                       created_at      TIMESTAMP NOT NULL DEFAULT NOW(),
                       updated_at      TIMESTAMP NOT NULL DEFAULT NOW()
);

COMMIT;