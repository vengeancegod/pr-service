package handlers

import (
	"encoding/json"
	"net/http"
	"pr-service/internal/dto"
	"pr-service/internal/service"
	httppkg "pr-service/pkg/http"
)

type UserHandler struct {
	userService service.UserService
	prService   service.PullRequestService
}

func NewUserHandler(userService service.UserService, prService service.PullRequestService) *UserHandler {
	return &UserHandler{
		userService: userService,
		prService:   prService,
	}
}

func (h *UserHandler) SetIsActive(w http.ResponseWriter, r *http.Request) {
	var req dto.SetUserActiveRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httppkg.JSONError(w, http.StatusBadRequest, err)
		return
	}

	if req.UserID == "" {
		httppkg.JSONError(w, http.StatusBadRequest, &httppkg.ErrorResponse{Message: "user_id is required"})
		return
	}

	user, err := h.userService.SetIsActive(r.Context(), req.UserID, req.IsActive)
	if err != nil {
		httppkg.JSONError(w, http.StatusInternalServerError, err)
		return
	}

	// Конвертируем в DTO ответ
	userResponse := dto.UserResponse{
		UserID:   user.ID,
		Username: user.Username,
		TeamName: user.TeamName,
		IsActive: user.IsActive,
	}

	httppkg.JSONResponse(w, http.StatusOK, map[string]interface{}{
		"user": userResponse,
	})
}

func (h *UserHandler) GetUserReviewRequests(w http.ResponseWriter, r *http.Request) {
	userID := r.URL.Query().Get("user_id")
	if userID == "" {
		httppkg.JSONError(w, http.StatusBadRequest, &httppkg.ErrorResponse{Message: "user_id is required"})
		return
	}

	prs, err := h.prService.GetPRByReviewerID(r.Context(), userID)
	if err != nil {
		httppkg.JSONError(w, http.StatusInternalServerError, err)
		return
	}

	var prList []dto.PullRequestShort
	for _, pr := range prs {
		prList = append(prList, dto.PullRequestShort{
			PullRequestID:   pr.ID,
			PullRequestName: pr.NamePR,
			AuthorID:        pr.AuthorID,
			Status:          string(pr.Status),
		})
	}

	response := dto.UserReviewResponse{
		UserID:       userID,
		PullRequests: prList,
	}

	httppkg.JSONResponse(w, http.StatusOK, response)
}
