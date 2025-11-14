package app

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"avitotest/internal/container"
)

// Run запускает приложение: подключается к БД, инициализирует зависимости и запускает сервер
func Run() error {
	// Создаем и инициализируем DI контейнер
	ctn, err := container.NewContainer()
	if err != nil {
		return fmt.Errorf("failed to create container: %w", err)
	}
	defer func() {
		if err := ctn.Close(); err != nil {
			log.Printf("error closing container: %v", err)
		}
	}()

	// Настраиваем маршруты
	e := ctn.Router.SetupRoutes()

	// Канал для получения сигналов завершения
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)

	// Запускаем сервер в отдельной горутине
	serverErr := make(chan error, 1)
	go func() {
		addr := ":" + ctn.Config.ServerPort
		log.Printf("Server starting on port %s", ctn.Config.ServerPort)
		if err := e.Start(addr); err != nil && err != http.ErrServerClosed {
			serverErr <- err
		}
	}()

	// Ждем сигнала завершения или ошибки сервера
	select {
	case err := <-serverErr:
		return fmt.Errorf("server error: %w", err)
	case sig := <-stop:
		log.Printf("Received signal: %v. Shutting down gracefully...", sig)
	}

	// Graceful shutdown
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := e.Shutdown(ctx); err != nil {
		return fmt.Errorf("server shutdown error: %w", err)
	}

	log.Println("Server stopped gracefully")
	return nil
}
