-- +goose Up
-- +goose StatementBegin
CREATE TABLE users (
    id VARCHAR(255) PRIMARY KEY,
    username VARCHAR(255) NOT NULL,
    is_active BOOLEAN DEFAULT true NOT NULL,
    team_name VARCHAR(255) NOT NULL
    -- CONSTRAINT fk_team FOREIGN KEY (team_name) REFERENCES teams(teams_name) ON DELETE CASCADE
);

CREATE INDEX idx_users_team_name ON users(team_name);
CREATE INDEX idx_users_is_active ON users(is_active);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS users;
-- +goose StatementEnd
