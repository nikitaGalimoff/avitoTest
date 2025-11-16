package container

import (
	"avitotest/internal/config"
	"avitotest/internal/domain"
	"avitotest/internal/handler"
	"avitotest/internal/repository"
	"avitotest/internal/usecase"
	"avitotest/pkg/logger"
	"database/sql"
	"log/slog"
)

type Container struct {
	Config *config.Config

	DB *sql.DB

	TeamRepo        domain.TeamRepository
	UserRepo        domain.UserRepository
	PullRequestRepo domain.PullRequestRepository

	TeamUseCase        *usecase.TeamUseCase
	UserUseCase        *usecase.UserUseCase
	PullRequestUseCase *usecase.PullRequestUseCase

	Router *handler.Router

	Logger *slog.Logger
}

func NewContainer() (*Container, error) {
	cfg := config.Load()

	postgresDB, err := repository.NewPostgresDB(cfg.GetDBConnectionString())
	if err != nil {
		return nil, err
	}
	db := postgresDB.DB()

	logger := logger.New()

	teamRepo := repository.NewTeamRepository(db)
	userRepo := repository.NewUserRepository(db)
	pullRequestRepo := repository.NewPullRequestRepository(db)

	teamUseCase := usecase.NewTeamUseCase(teamRepo, userRepo)
	userUseCase := usecase.NewUserUseCase(userRepo)
	pullRequestUseCase := usecase.NewPullRequestUseCase(pullRequestRepo, userRepo, teamRepo)

	router := handler.NewRouter(teamUseCase, userUseCase, pullRequestUseCase, logger)

	return &Container{
		Config:             cfg,
		DB:                 db,
		TeamRepo:           teamRepo,
		UserRepo:           userRepo,
		PullRequestRepo:    pullRequestRepo,
		TeamUseCase:        teamUseCase,
		UserUseCase:        userUseCase,
		PullRequestUseCase: pullRequestUseCase,
		Router:             router,
	}, nil
}
