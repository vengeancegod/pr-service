package handlers

import (
	"encoding/json"
	"net/http"
	"pr-service/internal/dto"
	"pr-service/internal/service"
	httppkg "pr-service/pkg/http"
)

type PullRequestHandler struct {
	prService service.PullRequestService
}

func NewPullRequestHandler(prService service.PullRequestService) *PullRequestHandler {
	return &PullRequestHandler{
		prService: prService,
	}
}

func (h *PullRequestHandler) CreatePR(w http.ResponseWriter, r *http.Request) {
	var req dto.CreatePRRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httppkg.JSONError(w, http.StatusBadRequest, err)
		return
	}

	if req.PullRequestID == "" || req.PullRequestName == "" || req.AuthorID == "" {
		httppkg.JSONError(w, http.StatusBadRequest, &httppkg.ErrorResponse{
			Message: "pull_request_id, pull_request_name and author_id are required",
		})
		return
	}

	createdPR, err := h.prService.CreatePR(r.Context(), req.PullRequestID, req.PullRequestName, req.AuthorID)
	if err != nil {
		if err.Error() == "PR already exists" {
			httppkg.JSONError(w, http.StatusConflict, &httppkg.ErrorResponse{
				Message: "PR id already exists",
			})
			return
		}
		httppkg.JSONError(w, http.StatusInternalServerError, err)
		return
	}

	prResponse := dto.PullRequest{
		PullRequestID:    createdPR.ID,
		PullRequestName:  createdPR.NamePR,
		AuthorID:         createdPR.AuthorID,
		Status:           string(createdPR.Status),
		ReplaceReviewers: createdPR.Reviewers,
	}

	httppkg.JSONResponse(w, http.StatusCreated, map[string]interface{}{
		"pr": prResponse,
	})
}

func (h *PullRequestHandler) MergePR(w http.ResponseWriter, r *http.Request) {
	var req dto.MergePRRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httppkg.JSONError(w, http.StatusBadRequest, err)
		return
	}

	if req.PullRequestID == "" {
		httppkg.JSONError(w, http.StatusBadRequest, &httppkg.ErrorResponse{Message: "pull_request_id is required"})
		return
	}

	mergedPR, err := h.prService.Merge(r.Context(), req.PullRequestID)
	if err != nil {
		httppkg.JSONError(w, http.StatusInternalServerError, err)
		return
	}

	prResponse := dto.PullRequest{
		PullRequestID:    mergedPR.ID,
		PullRequestName:  mergedPR.NamePR,
		AuthorID:         mergedPR.AuthorID,
		Status:           string(mergedPR.Status),
		ReplaceReviewers: mergedPR.Reviewers,
	}

	httppkg.JSONResponse(w, http.StatusOK, map[string]interface{}{
		"pr": prResponse,
	})
}

func (h *PullRequestHandler) ReplaceReviewer(w http.ResponseWriter, r *http.Request) {
	var req dto.ReplaceReviewerRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httppkg.JSONError(w, http.StatusBadRequest, err)
		return
	}

	if req.PullRequestID == "" || req.OldUserID == "" {
		httppkg.JSONError(w, http.StatusBadRequest, &httppkg.ErrorResponse{
			Message: "pull_request_id and old_user_id are required",
		})
		return
	}

	pr, newReviewerID, err := h.prService.ReplaceReviewer(r.Context(), req.PullRequestID, req.OldUserID)
	if err != nil {
		if err.Error() == "PR is merged" {
			httppkg.JSONError(w, http.StatusConflict, &httppkg.ErrorResponse{
				Message: "cannot reassign on merged PR",
			})
			return
		}
		httppkg.JSONError(w, http.StatusInternalServerError, err)
		return
	}

	prResponse := dto.PullRequest{
		PullRequestID:    pr.ID,
		PullRequestName:  pr.NamePR,
		AuthorID:         pr.AuthorID,
		Status:           string(pr.Status),
		ReplaceReviewers: pr.Reviewers,
	}

	response := dto.ReplaceReviewerResponse{
		PR:         prResponse,
		ReplacedBy: newReviewerID,
	}

	httppkg.JSONResponse(w, http.StatusOK, response)
}
