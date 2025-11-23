package model

type PRStatus string

const (
	PRStatusOpen   PRStatus = "OPEN"
	PRStatusMerged PRStatus = "MERGED"
)

type PullRequest struct {
	ID       string   
	NamePR   string   
	Status   PRStatus 
	AuthorID string   

	Reviewers []string 
}

func (pr *PullRequest) IsMerged() bool {
	return pr.Status == PRStatusMerged
}

func (pr *PullRequest) ReviewerCount() int {
	return len(pr.Reviewers)
}

func (pr *PullRequest) HasReviewer(userID string) bool {
	for _, reviewerID := range pr.Reviewers {
		if reviewerID == userID {
			return true
		}
	}
	return false
}

func (pr *PullRequest) CanModifyReviewers() bool {
	return !pr.IsMerged()
}
