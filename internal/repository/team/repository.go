package team

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"

	"pr-service/internal/model"
	rep "pr-service/internal/repository"
)

var _ rep.TeamRepository = (*repository)(nil)

type repository struct {
	pool *pgxpool.Pool
}

func NewRepository(pool *pgxpool.Pool) *repository {
	return &repository{
		pool: pool,
	}
}

func (r *repository) CreateTeam(ctx context.Context, team *model.Team) error {
	query := "INSERT INTO teams (id, team_name) VALUES ($1, $2)"

	_, err := r.pool.Exec(ctx, query, team.ID, team.TeamName)
	if err != nil {
		return fmt.Errorf("failed to create team: %w", err)
	}
	return nil
}

func (r *repository) UpdateTeam(ctx context.Context, team *model.Team) error {
	query := "UPDATE teams SET team_name = $2 WHERE id = $1"

	ct, err := r.pool.Exec(ctx, query, team.ID, team.TeamName)
	if err != nil {
		return fmt.Errorf("failed to update team: %w", err)
	}
	if ct.RowsAffected() == 0 {
		return fmt.Errorf("team not found: %s", team.ID)
	}
	return nil
}

func (r *repository) DeleteTeam(ctx context.Context, id string) error {
	query := `DELETE FROM teams WHERE id = $1`

	ct, err := r.pool.Exec(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete team: %w", err)
	}

	if ct.RowsAffected() == 0 {
		return fmt.Errorf("team not found: %s", id)
	}

	return nil
}

func (r *repository) AddMemberInTeam(ctx context.Context, teamID, userID string) error {
	var exists bool
	checkQuery := "SELECT EXISTS(SELECT 1 FROM teams WHERE id = $1)"

	err := r.pool.QueryRow(ctx, checkQuery, teamID).Scan(&exists)
	if err != nil {
		return fmt.Errorf("failed to check team existence: %w", err)
	}
	if !exists {
		return fmt.Errorf("team not found: %s", teamID)
	}

	query := "UPDATE users SET team_id = $1 WHERE id = $2"

	ct, err := r.pool.Exec(ctx, query, teamID, userID)
	if err != nil {
		return fmt.Errorf("failed to add member to team: %w", err)
	}

	if ct.RowsAffected() == 0 {
		return fmt.Errorf("user not found: %s", userID)
	}

	return nil
}

func (r *repository) RemoveMember(ctx context.Context, teamID, userID string) error {
	query := "UPDATE users SET team_id = '' WHERE id = $1 AND team_id = $2"

	ct, err := r.pool.Exec(ctx, query, userID, teamID)
	if err != nil {
		return fmt.Errorf("failed to remove member from team: %w", err)
	}

	if ct.RowsAffected() == 0 {
		return fmt.Errorf("user not found in team: user_id=%s, team_id=%s", userID, teamID)
	}

	return nil
}
