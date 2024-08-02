package server

import (
	"github.com/gofiber/fiber/v2"
)

func (s *FiberServer) RegisterFiberRoutes() {
	s.App.Get("/", s.rootHandler)
	s.App.Get("/health", s.healthHandler)
	s.App.Get("/dragons", s.dragonsHandler)

}

func (s *FiberServer) rootHandler(c *fiber.Ctx) error {
	resp := fiber.Map{
		"message": "success!",
	}

	return c.JSON(resp)
}

func (s *FiberServer) healthHandler(c *fiber.Ctx) error {
	return c.JSON(s.db.Health())
}

func (s *FiberServer) dragonsHandler(c *fiber.Ctx) error {
	dragons, err := s.db.GetAllDragons()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to retrieve dragons",
		})
	}

	return c.JSON(dragons)
}
