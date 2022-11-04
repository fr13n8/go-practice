package v1

import (
	_ "github.com/fr13n8/go-practice/docs"
	"github.com/fr13n8/go-practice/internal/services"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/swagger"
)

type Handler struct {
	services *services.Services
	router   fiber.Router
}

func NewHandler(svcs *services.Services, router *fiber.Router) *Handler {
	return &Handler{
		services: svcs,
		router:   *router,
	}
}

func (h *Handler) Init() *fiber.Router {
	v1 := h.router.Group("/v1", func(c *fiber.Ctx) error {
		c.Set("Version", "v1")
		return c.Next()
	})
	v1.Get("/swagger/*", swagger.HandlerDefault)
	h.InitTaskRoutes(v1)

	return &v1
}
