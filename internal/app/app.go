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

	if err := e.Start(":8080"); err != nil {
		panic(err)
	}

}
