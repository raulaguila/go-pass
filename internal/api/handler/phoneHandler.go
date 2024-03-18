package handler

import (
	"errors"
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/raulaguila/go-pass/internal/api/middleware"
	"github.com/raulaguila/go-pass/internal/pkg/domain"
	"github.com/raulaguila/go-pass/internal/pkg/dto"
	"github.com/raulaguila/go-pass/internal/pkg/i18n"
	"github.com/raulaguila/go-pass/pkg/filter"
	httphelper "github.com/raulaguila/go-pass/pkg/http-helper"
	"github.com/raulaguila/go-pass/pkg/pgerror"
	"github.com/raulaguila/go-pass/pkg/validator"
)

type PhoneHandler struct {
	phoneService domain.PhoneService
}

func (PhoneHandler) foreignKeyViolatedMethod(c *fiber.Ctx, translation *i18n.Translation) error {
	switch c.Method() {
	case fiber.MethodPut, fiber.MethodPost, fiber.MethodPatch:
		return httphelper.NewHTTPResponse(c, fiber.StatusBadRequest, translation.ErrPhoneNotFound)
	case fiber.MethodDelete:
		return httphelper.NewHTTPResponse(c, fiber.StatusBadRequest, translation.ErrPhoneUsed)
	default:
		return httphelper.NewHTTPResponse(c, fiber.StatusInternalServerError, translation.ErrGeneric)
	}
}

func (s PhoneHandler) handlerError(c *fiber.Ctx, err error) error {
	translation := c.Locals(httphelper.LocalLang).(*i18n.Translation)

	switch pgerror.HandlerError(err) {
	case pgerror.ErrDuplicatedKey:
		return httphelper.NewHTTPResponse(c, fiber.StatusConflict, translation.ErrPhoneRegistered)
	case pgerror.ErrForeignKeyViolated:
		return s.foreignKeyViolatedMethod(c, translation)
	case pgerror.ErrUndefinedColumn:
		return httphelper.NewHTTPResponse(c, fiber.StatusBadRequest, translation.ErrUndefinedColumn)
	}

	if errors.As(err, &validator.ErrValidator) {
		return httphelper.NewHTTPResponse(c, fiber.StatusBadRequest, err)
	}

	log.Println(err.Error())
	return httphelper.NewHTTPResponse(c, fiber.StatusInternalServerError, translation.ErrGeneric)
}

// Creates a new handler.
func NewPhoneHandler(route fiber.Router, ps domain.PhoneService, mid *middleware.RequesttMiddleware) {
	handler := &PhoneHandler{
		phoneService: ps,
	}

	route.Use(middleware.MidAccess)

	route.Get("", middleware.GetGenericFilter, handler.getPhones)
	route.Post("", middleware.GetPhoneDTO, handler.createPhone)
	route.Get("/:"+httphelper.ParamID, mid.PhoneByID, handler.getPhoneBydID)
	route.Put("/:"+httphelper.ParamID, mid.PhoneByID, middleware.GetPhoneDTO, handler.updatePhone)
	route.Delete("/:"+httphelper.ParamID, mid.PhoneByID, handler.deletePhone)
}

// getPhones godoc
// @Summary      Get phones
// @Description  Get phones
// @Tags         Phone
// @Accept       json
// @Produce      json
// @Param        lang query string false "Language responses"
// @Param        filter query filter.Filter false "Optional Filter"
// @Success      200  {array}   dto.ItemsOutputDTO
// @Failure      500  {object}  httphelper.HTTPResponse
// @Router       /phone [get]
// @Security	 Bearer
func (h *PhoneHandler) getPhones(c *fiber.Ctx) error {
	response, err := h.phoneService.GetPhonesOutputDTO(c.Context(), c.Locals(httphelper.LocalFilter).(*filter.Filter))
	if err != nil {
		return h.handlerError(c, err)
	}

	return c.Status(fiber.StatusOK).JSON(response)
}

// getPhoneBydID godoc
// @Summary      Get phone by ID
// @Description  Get phone by ID
// @Tags         Phone
// @Accept       json
// @Produce      json
// @Param        lang query string false "Language responses"
// @Param        id   path			int			true        "Phone ID"
// @Success      200  {object}  domain.Phone
// @Failure      400  {object}  httphelper.HTTPResponse
// @Failure      404  {object}  httphelper.HTTPResponse
// @Failure      500  {object}  httphelper.HTTPResponse
// @Router       /phone/{id} [get]
// @Security	 Bearer
func (h *PhoneHandler) getPhoneBydID(c *fiber.Ctx) error {
	return c.Status(fiber.StatusOK).JSON(c.Locals(httphelper.LocalObject).(*domain.Phone))
}

// createPhone godoc
// @Summary      Insert phone
// @Description  Insert phone
// @Tags         Phone
// @Accept       json
// @Produce      json
// @Param        lang query string false "Language responses"
// @Param        phone body dto.PhoneInputDTO true "Phone model"
// @Success      201  {object}  domain.Phone
// @Failure      400  {object}  httphelper.HTTPResponse
// @Failure      409  {object}  httphelper.HTTPResponse
// @Failure      500  {object}  httphelper.HTTPResponse
// @Router       /phone [post]
// @Security	 Bearer
func (h *PhoneHandler) createPhone(c *fiber.Ctx) error {
	phoneDTO := c.Locals(httphelper.LocalDTO).(*dto.PhoneInputDTO)
	phone, err := h.phoneService.CreatePhone(c.Context(), phoneDTO)
	if err != nil {
		return h.handlerError(c, err)
	}

	return c.Status(fiber.StatusCreated).JSON(phone)
}

// updatePhone godoc
// @Summary      Update phone by ID
// @Description  Update phone by ID
// @Tags         Phone
// @Accept       json
// @Produce      json
// @Param        lang query string false "Language responses"
// @Param        id     path    int     true        "Phone ID"
// @Param        phone body dto.PhoneInputDTO true "Phone model"
// @Success      200  {object}  domain.Phone
// @Failure      400  {object}  httphelper.HTTPResponse
// @Failure      404  {object}  httphelper.HTTPResponse
// @Failure      500  {object}  httphelper.HTTPResponse
// @Router       /phone/{id} [put]
// @Security	 Bearer
func (h *PhoneHandler) updatePhone(c *fiber.Ctx) error {
	phoneDTO := c.Locals(httphelper.LocalDTO).(*dto.PhoneInputDTO)
	phone := c.Locals(httphelper.LocalObject).(*domain.Phone)
	if err := h.phoneService.UpdatePhone(c.Context(), phone, phoneDTO); err != nil {
		return h.handlerError(c, err)
	}

	return c.Status(fiber.StatusOK).JSON(phone)
}

// deletePhone godoc
// @Summary      Delete phone by ID
// @Description  Delete phone by ID
// @Tags         Phone
// @Accept       json
// @Produce      json
// @Param        lang query string false "Language responses"
// @Param        id     path    int     true        "Phone ID"
// @Success      204  {object}  nil
// @Failure      404  {object}  httphelper.HTTPResponse
// @Failure      500  {object}  httphelper.HTTPResponse
// @Router       /phone/{id} [delete]
// @Security	 Bearer
func (h *PhoneHandler) deletePhone(c *fiber.Ctx) error {
	phone := c.Locals(httphelper.LocalObject).(*domain.Phone)
	if err := h.phoneService.DeletePhone(c.Context(), phone); err != nil {
		return h.handlerError(c, err)
	}

	return c.Status(fiber.StatusNoContent).Send(nil)
}
