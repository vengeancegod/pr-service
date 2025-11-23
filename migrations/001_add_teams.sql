-- +goose Up
-- +goose StatementBegin
CREATE TABLE teams (
    id VARCHAR(255) PRIMARY KEY,
    team_name VARCHAR(255) UNIQUE NOT NULL
);
CREATE INDEX idx_teams_name ON teams(team_name);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS teams;
-- +goose StatementEnd
