package pullrequest

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"pr-service/internal/model"
	rep "pr-service/internal/repository"
)

var _ rep.PullRequestRepository = (*repository)(nil)

type repository struct {
	pool *pgxpool.Pool
}

func NewRepository(pool *pgxpool.Pool) *repository {
	return &repository{
		pool: pool,
	}
}

func (r *repository) CreatePR(ctx context.Context, pr *model.PullRequest) error {
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback(ctx)

	query := "INSERT INTO pull_requests (id, name_pr, status, author_id) VALUES ($1, $2, $3, $4)"
	_, err = tx.Exec(ctx, query, pr.ID, pr.NamePR, pr.Status, pr.AuthorID)
	if err != nil {
		return fmt.Errorf("failed to create PR: %w", err)
	}

	if len(pr.Reviewers) > 0 {
		reviewerQuery := "INSERT INTO pr_reviewers (pr_id, reviewer_id) VALUES ($1, $2)"

		for _, reviewerID := range pr.Reviewers {
			_, err = tx.Exec(ctx, reviewerQuery, pr.ID, reviewerID)
			if err != nil {
				return fmt.Errorf("failed to add reviewer: %w", err)
			}
		}
	}

	if err = tx.Commit(ctx); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}

func (r *repository) GetPRByPRID(ctx context.Context, id string) (*model.PullRequest, error) {
	query := "SELECT id, name_pr, status, author_id FROM pull_requests WHERE id = $1"

	var pr model.PullRequest
	err := r.pool.QueryRow(ctx, query, id).Scan(
		&pr.ID,
		&pr.NamePR,
		&pr.Status,
		&pr.AuthorID,
	)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, fmt.Errorf("PR not found: %s", id)
		}
		return nil, fmt.Errorf("failed to get PR: %w", err)
	}

	reviewerGetQuery := "SELECT reviewer_id FROM pr_reviewers WHERE pr_id = $1"

	rows, err := r.pool.Query(ctx, reviewerGetQuery, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get reviewers: %w", err)
	}
	defer rows.Close()

	var reviewers []string
	for rows.Next() {
		var reviewerID string
		if err := rows.Scan(&reviewerID); err != nil {
			return nil, fmt.Errorf("failed to scan reviewer: %w", err)
		}
		reviewers = append(reviewers, reviewerID)
	}

	pr.Reviewers = reviewers

	return &pr, nil
}

func (r *repository) GetPRByReviewerID(ctx context.Context, reviewerID string) ([]model.PullRequest, error) {
	query := "SELECT DISTINCT pr.id, pr.name_pr, pr.status, pr.author_id FROM pull_requests pr INNER JOIN pr_reviewers prr ON pr.id = prr.pr_id WHERE prr.reviewer_id = $1"

	rows, err := r.pool.Query(ctx, query, reviewerID)
	if err != nil {
		return nil, fmt.Errorf("failed to get PRs by reviewer: %w", err)
	}
	defer rows.Close()

	var prs []model.PullRequest
	for rows.Next() {
		var pr model.PullRequest
		err := rows.Scan(&pr.ID, &pr.NamePR, &pr.Status, &pr.AuthorID)
		if err != nil {
			return nil, fmt.Errorf("failed to scan PR: %w", err)
		}

		reviewerGetQuery := "SELECT reviewer_id FROM pr_reviewers WHERE pr_id = $1"

		reviewerRows, err := r.pool.Query(ctx, reviewerGetQuery, pr.ID)
		if err != nil {
			return nil, fmt.Errorf("failed to get reviewers: %w", err)
		}

		var reviewers []string
		for reviewerRows.Next() {
			var revID string
			if err := reviewerRows.Scan(&revID); err != nil {
				reviewerRows.Close()
				return nil, fmt.Errorf("failed to scan reviewer: %w", err)
			}
			reviewers = append(reviewers, revID)
		}
		reviewerRows.Close()

		pr.Reviewers = reviewers
		prs = append(prs, pr)
	}

	return prs, nil
}

func (r *repository) ReplaceReviewer(ctx context.Context, prID, oldReviewerID, newReviewerID string) error {
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback(ctx)

	var status string
	statusQuery := "SELECT status FROM pull_requests WHERE id = $1"
	err = tx.QueryRow(ctx, statusQuery, prID).Scan(&status)
	if err != nil {
		if err == pgx.ErrNoRows {
			return fmt.Errorf("PR not found: %s", prID)
		}
		return fmt.Errorf("failed to check PR status: %w", err)
	}

	if status == string(model.PRStatusMerged) {
		return fmt.Errorf("cannot modify reviewers: PR is already merged")
	}

	reviewerDeleteQuery := "DELETE FROM pr_reviewers WHERE pr_id = $1 AND reviewer_id = $2"

	cmdTag, err := tx.Exec(ctx, reviewerDeleteQuery, prID, oldReviewerID)
	if err != nil {
		return fmt.Errorf("failed to remove old reviewer: %w", err)
	}

	if cmdTag.RowsAffected() == 0 {
		return fmt.Errorf("reviewer not found in PR: reviewer_id=%s, pr_id=%s", oldReviewerID, prID)
	}

	reviewerInsertQuery := "INSERT INTO pr_reviewers (pr_id, reviewer_id) VALUES ($1, $2)"

	_, err = tx.Exec(ctx, reviewerInsertQuery, prID, newReviewerID)
	if err != nil {
		return fmt.Errorf("failed to add new reviewer: %w", err)
	}

	if err = tx.Commit(ctx); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}

func (r *repository) Merge(ctx context.Context, prID string) (*model.PullRequest, error) {
	// Идемпотентност обновляем только если статус open
	query := "UPDATE pull_requests SET status = $2 WHERE id = $1 AND status = $3"

	_, err := r.pool.Exec(ctx, query, prID, model.PRStatusMerged, model.PRStatusOpen)
	if err != nil {
		return nil, fmt.Errorf("failed to merge PR: %w", err)
	}

	return r.GetPRByPRID(ctx, prID)
}
