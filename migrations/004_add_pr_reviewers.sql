-- +goose Up
-- +goose StatementBegin
CREATE TABLE pr_reviewers (
    pr_id VARCHAR(255) NOT NULL,
    reviewer_id VARCHAR(255) NOT NULL,
    PRIMARY KEY (pr_id, reviewer_id),
    CONSTRAINT fk_pr FOREIGN KEY (pr_id) REFERENCES pull_requests(id) ON DELETE CASCADE,
    CONSTRAINT fk_reviewer FOREIGN KEY (reviewer_id) REFERENCES users(id) ON DELETE CASCADE
)
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS pr_reviewers;
-- +goose StatementEnd
