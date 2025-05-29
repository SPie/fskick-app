-- +goose Up
-- +goose StatementBegin
SELECT 'up SQL query';
CREATE TABLE IF NOT EXISTS "users" (
    `id` INTEGER NOT NULL,
    `uuid` text NOT NULL UNIQUE,
    `email` VARCHAR(255) NOT NULL UNIQUE,
    `password` text NOT NULL,
    `player_id` INTEGER NOT NULL,
    `created_at` datetime,
    `updated_at` datetime,
    `deleted_at` datetime,
    PRIMARY KEY(`id`),

    FOREIGN KEY (`player_id`) REFERENCES `players` (`id`)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
DROP TABLE IF EXISTS "users";
-- +goose StatementEnd
