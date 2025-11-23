package user

import (
	"context"
	"encoding/json"
	"net/http"

	"pr-service/internal/model"
	"pr-service/internal/service"
	"pr-service/pkg/api"
)

type UserHandler struct {
	userService service.UserService
}

func NewUserHandler(userService service.UserService) *UserHandler {
	return &UserHandler{
		userService: userService,
	}
}

func (h *UserHandler) UsersSetIsActivePost(ctx context.Context, request *api.UsersSetIsActivePostReq) (api.UsersSetIsActivePostRes, error) {
	err := h.userService.SetUserActive(ctx, request.UserId, request.IsActive)
	if err != nil {
		// Возвращаем ошибку в формате, который ожидает клиент
		return &api.UsersSetIsActivePostBadRequest{
			Error: api.Error{
				Message: err.Error(),
			},
		}, nil
	}

	return &api.UsersSetIsActivePostOK{}, nil
}