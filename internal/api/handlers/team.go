package handlers

import (
	"encoding/json"
	"net/http"
	"pr-service/internal/dto"
	"pr-service/internal/model"
	"pr-service/internal/service"
	httppkg "pr-service/pkg/http"
)

type TeamHandler struct {
	teamService service.TeamService
}

func NewTeamHandler(teamService service.TeamService) *TeamHandler {
	return &TeamHandler{
		teamService: teamService,
	}
}

func (h *TeamHandler) AddTeam(w http.ResponseWriter, r *http.Request) {
	var req dto.CreateTeamRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httppkg.JSONError(w, http.StatusBadRequest, err)
		return
	}

	if req.TeamName == "" {
		httppkg.JSONError(w, http.StatusBadRequest, &httppkg.ErrorResponse{Message: "team name is required"})
		return
	}

	var users []model.User
	for _, member := range req.Members {
		users = append(users, model.User{
			ID:       member.UserID,
			Username: member.Username,
			IsActive: member.IsActive,
			TeamName: req.TeamName,
		})
	}

	team := &model.Team{
		TeamName: req.TeamName,
		Members:  users,
	}

	err := h.teamService.CreateTeam(r.Context(), team)
	if err != nil {
		httppkg.JSONError(w, http.StatusInternalServerError, err)
		return
	}

	httppkg.EmptyResponse(w, http.StatusCreated)
}

func (h *TeamHandler) GetTeam(w http.ResponseWriter, r *http.Request) {
	teamName := r.URL.Query().Get("team_name")
	if teamName == "" {
		httppkg.JSONError(w, http.StatusBadRequest, &httppkg.ErrorResponse{Message: "team_name is required"})
		return
	}

	team, err := h.teamService.GetTeamByName(r.Context(), teamName)
	if err != nil {
		httppkg.JSONError(w, http.StatusInternalServerError, err)
		return
	}

	httppkg.JSONResponse(w, http.StatusOK, team)
}
