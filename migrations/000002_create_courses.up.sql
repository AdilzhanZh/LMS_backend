BEGIN;

CREATE TABLE courses (
                         id            SERIAL PRIMARY KEY,
                         title         VARCHAR(255) NOT NULL,
                         description   TEXT,
                         slug          VARCHAR(255) NOT NULL UNIQUE,
                         price         INTEGER NOT NULL DEFAULT 0,
                         duration      INTEGER NOT NULL DEFAULT 0,
                         level         VARCHAR(50),
                         is_active     BOOLEAN NOT NULL DEFAULT FALSE,
                         teacher_id    INTEGER NOT NULL REFERENCES users(id) ON DELETE RESTRICT,
                         created_at    TIMESTAMP NOT NULL DEFAULT NOW(),
                         updated_at    TIMESTAMP NOT NULL DEFAULT NOW(),
                         deleted_at    TIMESTAMP NULL
);

COMMIT;