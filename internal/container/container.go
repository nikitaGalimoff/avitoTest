package container

import (
	"avitotest/internal/config"
	"avitotest/internal/domain"
	"avitotest/internal/handler"
	"avitotest/internal/repository"
	"avitotest/internal/usecase"
	"database/sql"
)

// Container представляет DI контейнер со всеми зависимостями
type Container struct {
	// Config
	Config *config.Config

	// Database
	DB *sql.DB

	// Repositories
	TeamRepo        domain.TeamRepository
	UserRepo        domain.UserRepository
	PullRequestRepo domain.PullRequestRepository

	// Use Cases
	TeamUseCase        *usecase.TeamUseCase
	UserUseCase        *usecase.UserUseCase
	PullRequestUseCase *usecase.PullRequestUseCase

	// Handlers
	AuthMiddleware *handler.AuthMiddleware
	Router         *handler.Router
}

// NewContainer создает и инициализирует новый контейнер зависимостей
func NewContainer() (*Container, error) {
	cfg := config.Load()

	postgresDB, err := repository.NewPostgresDB(cfg.GetDBConnectionString())
	if err != nil {
		return nil, err
	}
	db := postgresDB.DB()

	// Инициализируем репозитории
	teamRepo := repository.NewTeamRepository(db)
	userRepo := repository.NewUserRepository(db)
	pullRequestRepo := repository.NewPullRequestRepository(db)

	// Инициализируем use cases
	teamUseCase := usecase.NewTeamUseCase(teamRepo, userRepo)
	userUseCase := usecase.NewUserUseCase(userRepo)
	pullRequestUseCase := usecase.NewPullRequestUseCase(pullRequestRepo, userRepo, teamRepo)

	// Инициализируем middleware
	authMiddleware := handler.NewAuthMiddleware(cfg.AdminToken, cfg.UserToken)

	// Инициализируем router
	router := handler.NewRouter(teamUseCase, userUseCase, pullRequestUseCase, authMiddleware)

	return &Container{
		Config:             cfg,
		DB:                 db,
		TeamRepo:           teamRepo,
		UserRepo:           userRepo,
		PullRequestRepo:    pullRequestRepo,
		TeamUseCase:        teamUseCase,
		UserUseCase:        userUseCase,
		PullRequestUseCase: pullRequestUseCase,
		AuthMiddleware:     authMiddleware,
		Router:             router,
	}, nil
}

// Close закрывает все ресурсы контейнера
func (c *Container) Close() error {
	if c.DB != nil {
		return c.DB.Close()
	}
	return nil
}
