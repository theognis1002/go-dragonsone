package server

import (
	"github.com/gofiber/fiber/v2"

	"go-dragonstone/internal/database"
)

type FiberServer struct {
	*fiber.App

	db database.Service
}

func New() *FiberServer {
	server := &FiberServer{
		App: fiber.New(fiber.Config{
			ServerHeader: "go-dragonstone",
			AppName:      "go-dragonstone",
		}),

		db: database.New(),
	}

	return server
}
