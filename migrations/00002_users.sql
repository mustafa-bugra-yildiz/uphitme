-- +goose Up
-- +goose StatementBegin
CREATE EXTENSION citext;

CREATE TABLE users (
	id         UUID        NOT NULL PRIMARY KEY,
	created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
	updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),

	full_name          TEXT               NOT NULL,
	email              CITEXT      UNIQUE NOT NULL,
	email_confirmed_at TIMESTAMPTZ                 DEFAULT NULL,
	hashed_password    TEXT               NOT NULL

	CONSTRAINT min_full_name_length CHECK (length(full_name) >= 3)
	CONSTRAINT min_email_length     CHECK (length(email) >= 3)
);

CREATE TRIGGER update_users_modified_time
BEFORE UPDATE ON users
FOR EACH ROW
EXECUTE PROCEDURE update_modified_column();
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TRIGGER update_users_modified_time ON users;
DROP TABLE users;
DROP EXTENSION citext;
-- +goose StatementEnd
