package handler

import (
	_ "github.com/fuvy/effmob-test/docs"
	"github.com/fuvy/effmob-test/internal/service"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/swagger"
)

type Handler struct {
	service *service.Service
}

func New(srv *service.Service) *Handler {
	return &Handler{service: srv}
}
func (h *Handler) BindRoutes(app *fiber.App) error {
	app.Get("/swagger/*", swagger.HandlerDefault)
	app.Use(commonMiddleware)
	v1 := app.Group("/v1")
	cars := v1.Group("/cars")
	cars.Get("", h.GetCars)
	cars.Delete(":id", h.DeleteCar)
	cars.Patch(":id", h.UpdateCar)
	cars.Post("", h.RegisterCars)
	return nil
}
