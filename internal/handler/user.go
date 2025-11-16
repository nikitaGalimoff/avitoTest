package handler

import (
	"avitotest/internal/domain"
	"avitotest/internal/usecase"
	"github.com/labstack/echo/v4"
)

type UserHandler struct {
	userUseCase *usecase.UserUseCase
	prUseCase   *usecase.PullRequestUseCase
}

func NewUserHandler(
	userUseCase *usecase.UserUseCase,
	prUseCase *usecase.PullRequestUseCase,
) *UserHandler {
	return &UserHandler{
		userUseCase: userUseCase,
		prUseCase:   prUseCase,
	}
}

func (h *UserHandler) SetIsActive(c echo.Context) error {
	var req struct {
		UserID   string `json:"user_id"`
		IsActive bool   `json:"is_active"`
	}

	if err := c.Bind(&req); err != nil {
		return WriteError(c, err, 400)
	}

	user, err := h.userUseCase.SetIsActive(c.Request().Context(), req.UserID, req.IsActive)
	if err != nil {
		return WriteError(c, err, 0)
	}

	return WriteJSON(c, 200, map[string]interface{}{
		"user": user,
	})
}

func (h *UserHandler) GetReviewPullRequests(c echo.Context) error {
	userID := c.QueryParam("user_id")
	if userID == "" {
		return WriteError(c, domain.NewDomainError(domain.ErrorCodeNotFound, "user_id is required"), 400)
	}

	prs, err := h.prUseCase.GetPullRequestsByReviewer(c.Request().Context(), userID)
	if err != nil {
		return WriteError(c, err, 0)
	}

	return WriteJSON(c, 200, map[string]interface{}{
		"user_id":       userID,
		"pull_requests": prs,
	})
}
