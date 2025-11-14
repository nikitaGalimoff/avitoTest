package app

import (
	"avitotest/internal/container"
)

// Run запускает приложение: подключается к БД, инициализирует зависимости и запускает сервер
func Run() {

	ctn, err := container.NewContainer()
	if err != nil {
		panic(err)
	}

	e := ctn.Router.SetupRoutes()

	e.Start(ctn.Config.ServerPort)

}
