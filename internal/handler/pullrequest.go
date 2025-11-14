package handler

import (
	"avitotest/internal/domain"
	"avitotest/internal/usecase"

	"github.com/labstack/echo/v4"
)

// PullRequestHandler обрабатывает запросы для Pull Request'ов
type PullRequestHandler struct {
	prUseCase *usecase.PullRequestUseCase
	auth      *AuthMiddleware
}

// NewPullRequestHandler создает новый handler для Pull Request'ов
func NewPullRequestHandler(prUseCase *usecase.PullRequestUseCase, auth *AuthMiddleware) *PullRequestHandler {
	return &PullRequestHandler{
		prUseCase: prUseCase,
		auth:      auth,
	}
}

// CreatePullRequest обрабатывает POST /pullRequest/create
func (h *PullRequestHandler) CreatePullRequest(c echo.Context) error {
	var req struct {
		PullRequestID   string `json:"pull_request_id"`
		PullRequestName string `json:"pull_request_name"`
		AuthorID        string `json:"author_id"`
	}

	if err := c.Bind(&req); err != nil {
		return WriteError(c, err, 400)
	}

	pr, err := h.prUseCase.CreatePullRequest(c.Request().Context(), req.PullRequestID, req.PullRequestName, req.AuthorID)
	if err != nil {
		return WriteError(c, err, 0)
	}

	return WriteJSON(c, 201, map[string]interface{}{
		"pr": pr,
	})
}

// MergePullRequest обрабатывает POST /pullRequest/merge
func (h *PullRequestHandler) MergePullRequest(c echo.Context) error {
	var req struct {
		PullRequestID string `json:"pull_request_id"`
	}

	if err := c.Bind(&req); err != nil {
		return WriteError(c, err, 400)
	}

	pr, err := h.prUseCase.MergePullRequest(c.Request().Context(), req.PullRequestID)
	if err != nil {
		return WriteError(c, err, 0)
	}

	return WriteJSON(c, 200, map[string]interface{}{
		"pr": pr,
	})
}

// ReassignReviewer обрабатывает POST /pullRequest/reassign
func (h *PullRequestHandler) ReassignReviewer(c echo.Context) error {
	var req struct {
		PullRequestID string `json:"pull_request_id"`
		OldUserID     string `json:"old_user_id"`
		OldReviewerID string `json:"old_reviewer_id"` // Альтернативное поле для совместимости
	}

	if err := c.Bind(&req); err != nil {
		return WriteError(c, err, 400)
	}

	// Используем old_reviewer_id если old_user_id не указан
	if req.OldUserID == "" {
		req.OldUserID = req.OldReviewerID
	}

	if req.OldUserID == "" {
		return WriteError(c, domain.NewDomainError(domain.ErrorCodeNotFound, "old_user_id is required"), 400)
	}

	pr, newReviewerID, err := h.prUseCase.ReassignReviewer(c.Request().Context(), req.PullRequestID, req.OldUserID)
	if err != nil {
		return WriteError(c, err, 0)
	}

	return WriteJSON(c, 200, map[string]interface{}{
		"pr":          pr,
		"replaced_by": newReviewerID,
	})
}
