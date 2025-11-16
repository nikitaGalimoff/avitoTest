package handler

import (
	"avitotest/internal/usecase"
	"log/slog"
	"time"

	"github.com/labstack/echo/v4"
)

type Router struct {
	teamHandler        *TeamHandler
	userHandler        *UserHandler
	pullRequestHandler *PullRequestHandler
	logger             *slog.Logger
}

func NewRouter(
	teamUseCase *usecase.TeamUseCase,
	userUseCase *usecase.UserUseCase,
	prUseCase *usecase.PullRequestUseCase,
	logger *slog.Logger,
) *Router {
	return &Router{
		teamHandler:        NewTeamHandler(teamUseCase),
		userHandler:        NewUserHandler(userUseCase, prUseCase),
		pullRequestHandler: NewPullRequestHandler(prUseCase),
		logger:             logger,
	}
}
func (r *Router) SetupRoutes() *echo.Echo {
	e := echo.New()

	e.Use(r.loggingMiddleware())

	e.POST("/team/add", r.teamHandler.CreateTeam)
	e.GET("/team/get", r.teamHandler.GetTeam)

	e.POST("/users/setIsActive", r.userHandler.SetIsActive)
	e.GET("/users/getReview", r.userHandler.GetReviewPullRequests)

	e.POST("/pullRequest/create", r.pullRequestHandler.CreatePullRequest)
	e.POST("/pullRequest/merge", r.pullRequestHandler.MergePullRequest)
	e.POST("/pullRequest/reassign", r.pullRequestHandler.ReassignReviewer)

	e.GET("/health", func(c echo.Context) error {
		return c.String(200, "OK")
	})

	return e
}
func (r *Router) loggingMiddleware() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			start := time.Now()

			err := next(c)
			if err != nil {
				c.Error(err)
			}

			req := c.Request()
			res := c.Response()

			r.logger.Info("HTTP request",
				"method", req.Method,
				"path", req.URL.Path,
				"status", res.Status,
				"duration", time.Since(start).String(),
				"ip", c.RealIP(),
				"user_agent", req.UserAgent(),
			)

			return err
		}
	}
}
