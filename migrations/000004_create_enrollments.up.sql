BEGIN;

CREATE TABLE enrollments (
                             id            SERIAL PRIMARY KEY,
                             user_id       INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
                             course_id     INTEGER NOT NULL REFERENCES courses(id) ON DELETE CASCADE,
                             progress      INTEGER NOT NULL DEFAULT 0,
                             is_completed  BOOLEAN NOT NULL DEFAULT FALSE,
                             enrolled_at   TIMESTAMP NOT NULL DEFAULT NOW(),
                             completed_at  TIMESTAMP NULL,

                             CONSTRAINT unique_user_course UNIQUE (user_id, course_id)
);

CREATE INDEX idx_enrollments_user ON enrollments(user_id);
CREATE INDEX idx_enrollments_course ON enrollments(course_id);

COMMIT;