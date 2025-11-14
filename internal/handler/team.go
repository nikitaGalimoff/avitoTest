package handler

import (
	"avitotest/internal/domain"
	"avitotest/internal/usecase"
	"github.com/labstack/echo/v4"
)

// TeamHandler обрабатывает запросы для команд
type TeamHandler struct {
	teamUseCase *usecase.TeamUseCase
	auth        *AuthMiddleware
}

// NewTeamHandler создает новый handler для команд
func NewTeamHandler(teamUseCase *usecase.TeamUseCase, auth *AuthMiddleware) *TeamHandler {
	return &TeamHandler{
		teamUseCase: teamUseCase,
		auth:        auth,
	}
}

// CreateTeam обрабатывает POST /team/add
func (h *TeamHandler) CreateTeam(c echo.Context) error {
	var req struct {
		TeamName string                `json:"team_name"`
		Members  []domain.TeamMember `json:"members"`
	}

	if err := c.Bind(&req); err != nil {
		return WriteError(c, err, 400)
	}

	team := &domain.Team{
		TeamName: req.TeamName,
		Members:  req.Members,
	}

	result, err := h.teamUseCase.CreateTeam(c.Request().Context(), team)
	if err != nil {
		return WriteError(c, err, 0)
	}

	return WriteJSON(c, 201, map[string]interface{}{
		"team": result,
	})
}

// GetTeam обрабатывает GET /team/get
func (h *TeamHandler) GetTeam(c echo.Context) error {
	teamName := c.QueryParam("team_name")
	if teamName == "" {
		return WriteError(c, domain.NewDomainError(domain.ErrorCodeNotFound, "team_name is required"), 400)
	}

	team, err := h.teamUseCase.GetTeam(c.Request().Context(), teamName)
	if err != nil {
		return WriteError(c, err, 0)
	}

	return WriteJSON(c, 200, team)
}
