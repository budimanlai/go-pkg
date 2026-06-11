package base

import (
	"strconv"

	response "github.com/budimanlai/go-pkg/response"
	"github.com/budimanlai/go-pkg/validator"
	"github.com/gofiber/fiber/v2"
	"github.com/jinzhu/copier"
)

// BaseHandler menangani CRUD HTTP standar
// E = Entity, C = Create DTO, U = Update DTO
type BaseHandler[E any, C any, U any] struct {
	Service BaseUsecase[E]
}

func NewBaseHandler[E any, C any, U any](service BaseUsecase[E]) *BaseHandler[E, C, U] {
	return &BaseHandler[E, C, U]{
		Service: service,
	}
}

// Index (GET /)
func (h *BaseHandler[E, C, U]) Index(c *fiber.Ctx) error {
	page, _ := strconv.Atoi(c.Query("page", "1"))
	limit, _ := strconv.Atoi(c.Query("limit", "10"))

	result, err := h.Service.FindAll(c.Context(), page, limit)
	if err != nil {
		return response.ErrorI18n(c, fiber.StatusInternalServerError, err.Error(), nil)
	}

	return response.SuccessWithPagination(c, "app.success", response.PaginationResult{
		Data:      result.Data,
		Total:     result.Total,
		Page:      result.Page,
		Limit:     result.Limit,
		TotalPage: result.TotalPage,
	})
}

// View (GET /:id)
func (h *BaseHandler[E, C, U]) View(c *fiber.Ctx) error {
	id := c.Params("id")

	entity, err := h.Service.FindByID(c.Context(), id)
	if err != nil {
		// Asumsi error karena tidak ketemu
		if entity == nil {
			return response.ErrorI18n(c, fiber.StatusNotFound, "app.error.not_found", nil)
		}
		return response.ErrorI18n(c, fiber.StatusInternalServerError, err.Error(), nil)
	}

	// Safety check double cover
	if entity == nil {
		return response.ErrorI18n(c, fiber.StatusNotFound, "app.error.not_found", nil)
	}

	return response.SuccessI18n(c, "app.success", entity)
}

// Create (POST /)
func (h *BaseHandler[E, C, U]) Create(c *fiber.Ctx) error {
	var req C
	if err := c.BodyParser(&req); err != nil {
		return response.ErrorI18n(c, fiber.StatusBadRequest, "app.error.invalid_request_body", nil)
	}

	if err := validator.ValidateStructWithContext(c, &req); err != nil {
		return response.ValidationErrorI18n(c, err)
	}

	var entity E
	if err := copier.Copy(&entity, &req); err != nil {
		return response.Error(c, fiber.StatusInternalServerError, "Mapping failed")
	}

	if err := h.Service.Create(c.Context(), &entity); err != nil {
		return response.ErrorI18n(c, fiber.StatusInternalServerError, err.Error(), nil)
	}

	return response.SuccessI18n(c, "app.success", entity)
}

// Update (PUT /:id)
func (h *BaseHandler[E, C, U]) Update(c *fiber.Ctx) error {
	id := c.Params("id")

	// 1. Cek eksistensi data
	existing, err := h.Service.FindByID(c.Context(), id)
	if err != nil || existing == nil {
		return response.ErrorI18n(c, fiber.StatusNotFound, "app.error.not_found", nil)
	}

	// 2. Parse Body ke object existing
	var req U
	if err := c.BodyParser(&req); err != nil {
		return response.ErrorI18n(c, fiber.StatusBadRequest, "app.error.invalid_request_body", nil)
	}

	if err := validator.ValidateStructWithContext(c, &req); err != nil {
		return response.ValidationErrorI18n(c, err)
	}

	if err := copier.Copy(existing, &req); err != nil {
		return response.Error(c, fiber.StatusInternalServerError, "Mapping failed")
	}

	// 3. Save Update
	if err := h.Service.Update(c.Context(), existing); err != nil {
		return response.ErrorI18n(c, fiber.StatusInternalServerError, err.Error(), nil)
	}

	return response.SuccessI18n(c, "app.success", existing)
}

// Delete (DELETE /:id)
func (h *BaseHandler[E, C, U]) Delete(c *fiber.Ctx) error {
	id := c.Params("id")

	if err := h.Service.Delete(c.Context(), id); err != nil {
		return response.ErrorI18n(c, fiber.StatusInternalServerError, err.Error(), nil)
	}

	return response.SuccessI18n(c, "app.success", fiber.Map{"deleted": true})
}
