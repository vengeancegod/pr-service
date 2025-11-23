package user

import (
	"context"
	"fmt"
	"pr-service/internal/model"
	rep "pr-service/internal/repository"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

var _ rep.UserRepository = (*repository)(nil)

type repository struct {
	pool *pgxpool.Pool
}

func NewRepository(pool *pgxpool.Pool) *repository {
	return &repository{
		pool: pool,
	}
}

func (r *repository) CreateUser(ctx context.Context, user *model.User) error {
	query := "INSERT INTO users (id, user_name, is_active, team_id) VALUES ($1, $2, $3, $4)"

	_, err := r.pool.Exec(ctx, query, user.ID, user.UserName, user.IsActive, user.TeamID)
	if err != nil {
		return fmt.Errorf("failed to create user: %w", err)
	}

	return nil
}

func (r *repository) GetUserByID(ctx context.Context, id string) (*model.User, error) {
	query := "SELECT id, user_name, is_active, team_id FROM users WHERE id = $1"

	var user model.User
	err := r.pool.QueryRow(ctx, query, id).Scan(&user.ID, &user.UserName, &user.IsActive, &user.TeamID)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, fmt.Errorf("user not found: %s", id)
		}
		return nil, fmt.Errorf("failed to get user: %w", err)
	}
	return &user, nil
}

func (r *repository) UpdateUser(ctx context.Context, user *model.User) error {
	query := "UPDATE users SET user_name = $2, is_active = $3, team_id = $4 WHERE id = $1"

	ct, err := r.pool.Exec(ctx, query, user.ID, user.UserName, user.IsActive, user.TeamID)
	if err != nil {
		return fmt.Errorf("failed to update user: %w", err)
	}

	if ct.RowsAffected() == 0 {
		return fmt.Errorf("user not found: %s", user.ID)
	}
	return nil
}

func (r *repository) DeleteUser(ctx context.Context, id string) error {
	query := "DELETE FROM users WHERE id = $1"

	ct, err := r.pool.Exec(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete user: %w", err)
	}

	if ct.RowsAffected() == 0 {
		return fmt.Errorf("user not found: %s", id)
	}
	return nil
}

func (r *repository) GetActiveUserFromTeam(ctx context.Context, teamID string) ([]model.User, error) {
	query := "SELECT id, user_name, is_active, team_id FROM users WHERE team_id = $1 AND is_active = true"
	rows, err := r.pool.Query(ctx, query, teamID)
	if err != nil {
		return nil, fmt.Errorf("failed to get active users: %w", err)
	}
	defer rows.Close()

	var users []model.User
	for rows.Next() {
		var user model.User
		err := rows.Scan(&user.ID, &user.UserName, &user.IsActive, &user.TeamID)
		if err != nil {
			return nil, fmt.Errorf("failed to scan user: %w", err)
		}
		users = append(users, user)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("rows error: %w", err)
	}
	return users, nil
}
