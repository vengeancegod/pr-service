-- +goose Up
-- +goose StatementBegin
CREATE TABLE pull_requests (
    id VARCHAR(255) PRIMARY KEY,
    name_pr VARCHAR(255) NOT NULL,
    status VARCHAR(10) NOT NULL CHECK (status IN('OPEN', 'MERGED')),
    author_id VARCHAR(255) NOT NULL,
    CONSTRAINT fk_author FOREIGN KEY (author_id) REFERENCES users(id) ON DELETE CASCADE
);

CREATE INDEX idx_pr_author ON pull_requests(author_id);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS pull_requests;
-- +goose StatementEnd
