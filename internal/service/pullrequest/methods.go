package pr

import (
	"context"
	"fmt"
	"math/rand"

	"pr-service/internal/model"
)

func (s *service) CreatePR(ctx context.Context, prID, namePR, authorID string) (*model.PullRequest, error) {
	if namePR == "" {
		return nil, fmt.Errorf("PR name is required")
	}
	if authorID == "" {
		return nil, fmt.Errorf("author ID is required")
	}
	// получаем автора пра
	author, err := s.userRepository.GetUserByID(ctx, authorID)
	if err != nil {
		return nil, fmt.Errorf("failed to get author: %w", err)
	}

	// активные членды команды автора пра
	activeMembers, err := s.userRepository.GetActiveUserFromTeam(ctx, author.TeamName)
	if err != nil {
		return nil, fmt.Errorf("failed to get active team members: %w", err)
	}

	// 2 ревьюера
	reviewers := selectReviewers(activeMembers, authorID, 2)

	pr := &model.PullRequest{
		ID:        prID,
		NamePR:    namePR,
		Status:    model.PRStatusOpen,
		AuthorID:  authorID,
		Reviewers: reviewers,
	}

	if err := s.prRepository.CreatePR(ctx, pr); err != nil {
		return nil, fmt.Errorf("failed to create PR: %w", err)
	}

	return pr, nil
}

func (s *service) GetPRByPRID(ctx context.Context, prID string) (*model.PullRequest, error) {
	if prID == "" {
		return nil, fmt.Errorf("PR ID is required")
	}

	return s.prRepository.GetPRByPRID(ctx, prID)
}

func (s *service) GetPRByReviewerID(ctx context.Context, reviewerID string) ([]model.PullRequest, error) {
	if reviewerID == "" {
		return nil, fmt.Errorf("reviewer ID is required")
	}

	return s.prRepository.GetPRByReviewerID(ctx, reviewerID)
}

func (s *service) ReplaceReviewer(ctx context.Context, prID, oldReviewerID string) (*model.PullRequest, string, error) {
	if prID == "" {
		return nil, "", fmt.Errorf("PR ID is required")
	}
	if oldReviewerID == "" {
		return nil, "", fmt.Errorf("old reviewer ID is required")
	}

	pr, err := s.prRepository.GetPRByPRID(ctx, prID)
	if err != nil {
		return nil, "", fmt.Errorf("failed to get PR: %w", err)
	}

	// проверка на мерж
	if pr.IsMerged() {
		return nil, "", fmt.Errorf("cannot replace reviewer: PR is already merged")
	}

	// проверка что прежний ревьюер действительно ревьювер этого PR
	if !pr.HasReviewer(oldReviewerID) {
		return nil, "", fmt.Errorf("user %s is not a reviewer of PR %s", oldReviewerID, prID)
	}

	oldReviewer, err := s.userRepository.GetUserByID(ctx, oldReviewerID)
	if err != nil {
		return nil, "", fmt.Errorf("failed to get old reviewer: %w", err)
	}

	// активные члены команды предыдущего ревьювера
	activeMembers, err := s.userRepository.GetActiveUserFromTeam(ctx, oldReviewer.TeamName)
	if err != nil {
		return nil, "", fmt.Errorf("failed to get active team members: %w", err)
	}

	// исключение автора и всех текущих ревьюеров
	excludeIDs := append(pr.Reviewers, pr.AuthorID)
	candidates := filterCandidates(activeMembers, excludeIDs)

	if len(candidates) == 0 {
		return nil, "", fmt.Errorf("no available candidates for replacement")
	}

	newReviewer := candidates[rand.Intn(len(candidates))]
	newReviewerID := newReviewer.ID

	if err := s.prRepository.ReplaceReviewer(ctx, prID, oldReviewerID, newReviewerID); err != nil {
		return nil, "", fmt.Errorf("failed to replace reviewer: %w", err)
	}

	updatedPR, err := s.prRepository.GetPRByPRID(ctx, prID)
	if err != nil {
		return nil, "", fmt.Errorf("failed to get updated PR: %w", err)
	}

	return updatedPR, newReviewerID, nil
}

func (s *service) Merge(ctx context.Context, prID string) (*model.PullRequest, error) {
	if prID == "" {
		return nil, fmt.Errorf("PR ID is required")
	}

	return s.prRepository.Merge(ctx, prID)
}

func selectReviewers(users []model.User, excludeID string, maxCount int) []string {
	var candidates []model.User
	for _, u := range users {
		if u.ID != excludeID && u.IsActive {
			candidates = append(candidates, u)
		}
	}

	rand.Shuffle(len(candidates), func(i, j int) {
		candidates[i], candidates[j] = candidates[j], candidates[i]
	})

	count := min(len(candidates), maxCount)
	reviewers := make([]string, count)
	for i := 0; i < count; i++ {
		reviewers[i] = candidates[i].ID
	}

	return reviewers
}

func filterCandidates(users []model.User, excludeIDs []string) []model.User {
	excludeMap := make(map[string]bool)
	for _, id := range excludeIDs {
		excludeMap[id] = true
	}

	var result []model.User
	for _, u := range users {
		if !excludeMap[u.ID] && u.IsActive {
			result = append(result, u)
		}
	}

	return result
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
