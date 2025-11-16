package app

import (
	"avitotest/internal/container"
)

func Run() {
	ctn, err := container.NewContainer()
	if err != nil {
		panic(err)
	}
	e := ctn.Router.SetupRoutes()
	if err := e.Start(":" + ctn.Config.ServerPort); err != nil {
		panic(err)
	}
}
