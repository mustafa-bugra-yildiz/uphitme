-- +goose Up
-- +goose StatementBegin
CREATE OR REPLACE FUNCTION update_modified_column()
RETURNS TRIGGER AS $$
BEGIN
	NEW.updated_at = now();
	RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TYPE task_status AS ENUM (
	'pending',
	'failed',
	'succeeded'
);

CREATE TABLE tasks (
	id         UUID        PRIMARY KEY,
	created_at TIMESTAMPTZ NOT NULL    DEFAULT now(),
	updated_at TIMESTAMPTZ NOT NULL    DEFAULT now(),

	status  task_status NOT NULL DEFAULT 'pending',
	title   TEXT        NOT NULL,
	target  TEXT        NOT NULL,
	payload JSONB       NOT NULL

	CONSTRAINT min_title_length CHECK (length(title) >= 3)
	CONSTRAINT min_target_length CHECK(length(target) >= 7)
);

CREATE TRIGGER update_tasks_modified_time
BEFORE UPDATE ON tasks
FOR EACH ROW
EXECUTE PROCEDURE update_modified_column();
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TRIGGER update_tasks_modified_time ON tasks;

DROP TABLE tasks;

DROP TYPE task_status;

DROP FUNCTION update_modified_column;
-- +goose StatementEnd
