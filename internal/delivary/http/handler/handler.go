package handler

import (
	v1 "github.com/fr13n8/go-practice/internal/delivary/http/handler/v1"
	"github.com/fr13n8/go-practice/internal/delivary/http/middlewares"
	"github.com/fr13n8/go-practice/internal/services"
	"github.com/gofiber/adaptor/v2"
	"github.com/gofiber/fiber/v2"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

type Handler struct {
	services *services.Services
}

func NewHandler(svcs *services.Services) *Handler {
	return &Handler{
		services: svcs,
	}
}

// @title           Go Practice API
// @version         1.0
// @description     This is a sample server celler server.
// @termsOfService  http://swagger.io/terms/

// @contact.name   API Support
// @contact.url    http://www.swagger.io/support
// @contact.email  support@swagger.io

// @license.name  Apache 2.0
// @license.url   http://www.apache.org/licenses/LICENSE-2.0.html

// @host      localhost:8080
// @BasePath  /api/

// @securityDefinitions.basic  BasicAuth
func (h *Handler) Init(srv *fiber.App) {
	api := srv.Group("/api", middlewares.RecrdRequestLatency)
	handler := v1.NewHandler(h.services, &api)

	// srv.Static("/", "static")
	// srv.Get("*", func(c *fiber.Ctx) error {
	// 	return c.SendFile("static/index.html")
	// })
	h.RegisterPromethesRoutes(srv)
	handler.Init()
}

func (h *Handler) RegisterPromethesRoutes(srv *fiber.App) {
	srv.Get("/metrics", adaptor.HTTPHandler(promhttp.Handler()))
}
