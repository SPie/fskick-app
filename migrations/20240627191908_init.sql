-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS "seasons" (
    id INTEGER NOT NULL,
    name VARCHAR(255) NOT NULL,
    active BOOLEAN DEFAULT '0' NOT NULL,
    `updated_at` datetime,
    `deleted_at` datetime,
    `created_at` datetime,
    `uuid` text NOT NULL UNIQUE,
    PRIMARY KEY(id)
);
CREATE INDEX `idx_seasons_deleted_at` ON `seasons`(`deleted_at`);

CREATE TABLE IF NOT EXISTS "players" (
    id INTEGER NOT NULL,
    name VARCHAR(255) NOT NULL,
    created_at DATETIME NOT NULL,
    updated_at DATETIME NOT NULL,
    `deleted_at` datetime,
    `uuid` text NOT NULL UNIQUE,
    PRIMARY KEY(id)
);
CREATE INDEX `idx_players_deleted_at` ON `players`(`deleted_at`);

CREATE TABLE IF NOT EXISTS "games" (
    id INTEGER NOT NULL,
    season_id INTEGER UNSIGNED NOT NULL,
    created_at DATETIME NOT NULL,
    updated_at DATETIME NOT NULL,
    `deleted_at` datetime,
    `uuid` text NOT NULL UNIQUE,
    `played_at` datetime,
    PRIMARY KEY(id)
);

CREATE TABLE IF NOT EXISTS "attendances" (
    id INTEGER NOT NULL,
    game_id INTEGER UNSIGNED NOT NULL,
    player_id INTEGER UNSIGNED NOT NULL,
    win BOOLEAN NOT NULL,
    created_at DATETIME NOT NULL,
    updated_at DATETIME NOT NULL,
    deleted_at datetime null default null,
    uuid text null default null,
    PRIMARY KEY(id)
);
CREATE INDEX IDX_9C6B8FD4E48FD905 ON "attendances" (game_id);
CREATE INDEX IDX_9C6B8FD499E6F5DF ON "attendances" (player_id);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP INDEX IF EXISTS IDX_9C6B8FD4E48FD905;
DROP INDEX IF EXISTS IDX_9C6B8FD499E6F5DF;
DROP TABLE IF EXISTS "attendances";
DROP TABLE IF EXISTS "games";
DROP INDEX IF EXISTS `idx_players_deleted_at`;
DROP TABLE IF EXISTS "players";
DROP INDEX IF EXISTS `idx_seasons_deleted_at`;
DROP TABLE IF EXISTS "seasons";
-- +goose StatementEnd
