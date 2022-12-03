package v1

import (
	"github.com/fr13n8/go-practice/internal/domain"
	"github.com/gofiber/fiber/v2"
	"github.com/opentracing/opentracing-go"
	"gorm.io/gorm"
)

func (h *Handler) InitTaskRoutes(router fiber.Router) *fiber.Router {
	routes := router.Group("/task")
	routes.Get("/all", h.GetAll)
	routes.Get("/:id", h.Get)
	routes.Post("/", h.Create)
	routes.Put("/:id", h.Update)
	routes.Delete("/:id", h.Delete)

	return &routes
}

// Get ShowAccount godoc
// @Summary      Get a task
// @Description  get task by ID
// @Tags         tasks
// @Accept       json
// @Produce      json
// @Param        id   path      int  true  "Task ID"
// @Success      200  {object}  domain.Task
// @Router       /api/v1/task/:id [get]
func (h *Handler) Get(ctx *fiber.Ctx) error {
	span, jCtx := opentracing.StartSpanFromContext(ctx.Context(), "task.Get")
	defer span.Finish()

	id := ctx.Params("id")
	task, err := h.services.Task.Get(jCtx, id)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return ctx.Status(404).JSON(err)
		}
		return ctx.Status(400).JSON(
			fiber.Map{
				"error": "Something went wrong",
			},
		)
	}

	return ctx.Status(200).JSON(task)
}

// GetAll ShowAccount godoc
// @Summary      Get all tasks
// @Description  get all tasks
// @Tags         tasks
// @Accept       json
// @Produce      json
// @Success      200  {object}  []domain.Task
// @Router       /api/v1/task [get]
func (h *Handler) GetAll(ctx *fiber.Ctx) error {
	span, jCtx := opentracing.StartSpanFromContext(ctx.Context(), "task.GetAll")
	defer span.Finish()

	envs, err := h.services.Task.GetAll(jCtx)
	if err != nil {
		return ctx.Status(400).JSON(
			fiber.Map{
				"error": "Something went wrong",
			},
		)
	}

	return ctx.Status(200).JSON(envs)
}

// @Summary Create task
// @Tags tasks
// @Description create task
// @ID create-task
// @Accept  json
// @Produce  json
// @Param input body domain.TaskCreate true "task info"
// @Success 200 {object}  domain.Task
// @Router /api/v1/task [post]
func (h *Handler) Create(ctx *fiber.Ctx) error {
	span, jCtx := opentracing.StartSpanFromContext(ctx.Context(), "task.Create")
	defer span.Finish()

	reqBody := domain.TaskCreate{}
	if err := ctx.BodyParser(&reqBody); err != nil {
		return ctx.Status(400).JSON(
			fiber.Map{
				"error": "Something went wrong",
			},
		)
	}

	task, err := h.services.Task.Create(jCtx, reqBody)
	if err != nil {
		return ctx.Status(400).JSON(
			fiber.Map{
				"error": "Something went wrong",
			},
		)
	}
	return ctx.Status(200).JSON(task)
}

// @Summary Update task
// @Tags tasks
// @Description update task
// @ID update-task
// @Accept  json
// @Produce  json
// @Param input body domain.TaskUpdate true "task info"
// @Param id path string true "Task id"
// @Success 200 {object} domain.Task
// @Router /api/v1/task/:id [put]
func (h *Handler) Update(ctx *fiber.Ctx) error {
	span, jCtx := opentracing.StartSpanFromContext(ctx.Context(), "task.Update")
	defer span.Finish()

	id := ctx.Params("id")
	_, err := h.services.Task.Get(jCtx, id)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return ctx.Status(404).JSON(err)
		}
		return ctx.Status(400).JSON(
			fiber.Map{
				"error": "Something went wrong",
			},
		)
	}

	reqBody := domain.TaskUpdate{}
	if err := ctx.BodyParser(&reqBody); err != nil {
		return ctx.Status(400).JSON(
			fiber.Map{
				"error": "Something went wrong",
			},
		)
	}

	task, err := h.services.Task.Update(jCtx, reqBody, id)
	if err != nil {
		return ctx.Status(400).JSON(
			fiber.Map{
				"error": "Something went wrong",
			},
		)
	}
	return ctx.Status(200).JSON(task)
}

// @Summary Delete task
// @Security ApiKeyAuth
// @Tags tasks
// @Description delete task by id
// @ID delete-task-by-id
// @Accept  json
// @Produce  json
// @Param id path string true "Task id"
// @Success 200 {string} Ok
// @Router /api/v1/task/:id [delete]
func (h *Handler) Delete(ctx *fiber.Ctx) error {
	span, jCtx := opentracing.StartSpanFromContext(ctx.Context(), "task.Delete")
	defer span.Finish()

	id := ctx.Params("id")
	_, err := h.services.Task.Get(jCtx, id)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return ctx.Status(404).JSON(err)
		}
		return ctx.Status(400).JSON(
			fiber.Map{
				"error": "Something went wrong",
			},
		)
	}

	err = h.services.Task.Delete(jCtx, id)
	if err != nil {
		return ctx.Status(400).JSON(
			fiber.Map{
				"error": "Something went wrong",
			},
		)
	}
	return ctx.Status(200).JSON(nil)
}
