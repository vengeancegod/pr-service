package team

import (
	"context"
	"fmt"

	"github.com/jackc/pgx"
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
	query := "INSERT INTO teams (team_name) VALUES ($1)"

	_, err := r.pool.Exec(ctx, query, team.TeamName)
	if err != nil {
		return fmt.Errorf("failed to create team: %w", err)
	}
	return nil
}

func (r *repository) GetTeamByName(ctx context.Context, teamName string) (*model.Team, error) {
	query := "SELECT team_name FROM teams WHERE team_name = $1"

	var team model.Team
	err := r.pool.QueryRow(ctx, query, teamName).Scan(&team.TeamName)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, fmt.Errorf("team not found: %s", teamName)
		}
		return nil, fmt.Errorf("failed to get team: %w", err)
	}

	// получаем участников команды
	membersQuery := "SELECT id, username, is_active, team_name FROM users WHERE team_name = $1"
	rows, err := r.pool.Query(ctx, membersQuery, team.TeamName)
	if err != nil {
		return nil, fmt.Errorf("failed to get team members: %w", err)
	}
	defer rows.Close()

	var members []model.User
	for rows.Next() {
		var user model.User
		err := rows.Scan(&user.ID, &user.Username, &user.IsActive, &user.TeamName)
		if err != nil {
			return nil, fmt.Errorf("failed to scan member: %w", err)
		}
		members = append(members, user)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("rows error: %w", err)
	}

	team.Members = members
	return &team, nil
}
