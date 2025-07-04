-- +goose Up
-- +goose StatementBegin
SELECT "up SQL query";
ALTER TABLE players ADD COLUMN email TEXT NULL;
ALTER TABLE players ADD COLUMN password TEXT NULL;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SELECT "down SQL query";
ALTER TABLE players DROP COLUMN email;
ALTER TABLE players DROP COLUMN password;
-- +goose StatementEnd
