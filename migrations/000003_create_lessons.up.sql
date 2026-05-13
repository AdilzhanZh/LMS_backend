BEGIN;

CREATE TABLE lessons (
                         id          SERIAL PRIMARY KEY,
                         course_id   INTEGER NOT NULL REFERENCES courses(id) ON DELETE CASCADE,
                         title       VARCHAR(255) NOT NULL,
                         content     TEXT,
                         video_url   TEXT,
                         duration    INTEGER NOT NULL DEFAULT 0,
                         position    INTEGER NOT NULL DEFAULT 0,
                         is_preview  BOOLEAN NOT NULL DEFAULT FALSE,
                         created_at  TIMESTAMP NOT NULL DEFAULT NOW(),
                         updated_at  TIMESTAMP NOT NULL DEFAULT NOW(),
                         deleted_at  TIMESTAMP NULL
);

CREATE INDEX idx_lessons_course ON lessons(course_id);

COMMIT;