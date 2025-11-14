package handler

import (
	"avitotest/internal/usecase"

	"github.com/labstack/echo/v4"
)

// Router настраивает маршруты приложения
type Router struct {
	teamHandler        *TeamHandler
	userHandler        *UserHandler
	pullRequestHandler *PullRequestHandler
	auth               *AuthMiddleware
}

// NewRouter создает новый роутер
func NewRouter(
	teamUseCase *usecase.TeamUseCase,
	userUseCase *usecase.UserUseCase,
	prUseCase *usecase.PullRequestUseCase,
	auth *AuthMiddleware,
) *Router {
	return &Router{
		teamHandler:        NewTeamHandler(teamUseCase, auth),
		userHandler:        NewUserHandler(userUseCase, prUseCase, auth),
		pullRequestHandler: NewPullRequestHandler(prUseCase, auth),
		auth:               auth,
	}
}

// SetupRoutes настраивает все маршруты и возвращает Echo instance
func (r *Router) SetupRoutes() *echo.Echo {
	e := echo.New()

	// Teams
	e.POST("/team/add", r.teamHandler.CreateTeam)
	e.GET("/team/get", r.teamHandler.GetTeam, r.auth.AdminOrUser)

	// Users
	e.POST("/users/setIsActive", r.userHandler.SetIsActive, r.auth.AdminOnly)
	e.GET("/users/getReview", r.userHandler.GetReviewPullRequests, r.auth.AdminOrUser)

	// Pull Requests
	e.POST("/pullRequest/create", r.pullRequestHandler.CreatePullRequest, r.auth.AdminOnly)
	e.POST("/pullRequest/merge", r.pullRequestHandler.MergePullRequest, r.auth.AdminOnly)
	e.POST("/pullRequest/reassign", r.pullRequestHandler.ReassignReviewer, r.auth.AdminOnly)

	// Health check
	e.GET("/health", func(c echo.Context) error {
		return c.String(200, "OK")
	})

	return e
}
