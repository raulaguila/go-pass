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

type OperatorHandler struct {
	operatorService domain.OperatorService
}

func (OperatorHandler) foreignKeyViolatedMethod(c *fiber.Ctx, translation *i18n.Translation) error {
	switch c.Method() {
	case fiber.MethodPut, fiber.MethodPost, fiber.MethodPatch:
		return httphelper.NewHTTPResponse(c, fiber.StatusBadRequest, translation.ErrOperatorNotFound)
	case fiber.MethodDelete:
		return httphelper.NewHTTPResponse(c, fiber.StatusBadRequest, translation.ErrOperatorUsed)
	default:
		return httphelper.NewHTTPResponse(c, fiber.StatusInternalServerError, translation.ErrGeneric)
	}
}

func (s OperatorHandler) handlerError(c *fiber.Ctx, err error) error {
	translation := c.Locals(httphelper.LocalLang).(*i18n.Translation)

	switch pgerror.HandlerError(err) {
	case pgerror.ErrDuplicatedKey:
		return httphelper.NewHTTPResponse(c, fiber.StatusConflict, translation.ErrOperatorRegistered)
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
func NewOperatorHandler(route fiber.Router, ps domain.OperatorService, mid *middleware.RequesttMiddleware) {
	handler := &OperatorHandler{
		operatorService: ps,
	}

	route.Use(middleware.MidAccess)

	route.Get("", middleware.GetGenericFilter, handler.getOperators)
	route.Post("", middleware.GetOperatorDTO, handler.createOperator)
	route.Get("/:"+httphelper.ParamID, mid.OperatorByID, handler.getOperatorBydID)
	route.Put("/:"+httphelper.ParamID, mid.OperatorByID, middleware.GetOperatorDTO, handler.updateOperator)
	route.Delete("/:"+httphelper.ParamID, mid.OperatorByID, handler.deleteOperator)
}

// getOperators godoc
// @Summary      Get operators
// @Description  Get operators
// @Tags         Operator
// @Accept       json
// @Produce      json
// @Param        lang query string false "Language responses"
// @Param        filter query filter.Filter false "Optional Filter"
// @Success      200  {array}   dto.ItemsOutputDTO
// @Failure      500  {object}  httphelper.HTTPResponse
// @Router       /operator [get]
// @Security	 Bearer
func (h *OperatorHandler) getOperators(c *fiber.Ctx) error {
	response, err := h.operatorService.GetOperatorsOutputDTO(c.Context(), c.Locals(httphelper.LocalFilter).(*filter.Filter))
	if err != nil {
		return h.handlerError(c, err)
	}

	return c.Status(fiber.StatusOK).JSON(response)
}

// getOperatorBydID godoc
// @Summary      Get operator by ID
// @Description  Get operator by ID
// @Tags         Operator
// @Accept       json
// @Produce      json
// @Param        lang query string false "Language responses"
// @Param        id   path			int			true        "Operator ID"
// @Success      200  {object}  domain.Operator
// @Failure      400  {object}  httphelper.HTTPResponse
// @Failure      404  {object}  httphelper.HTTPResponse
// @Failure      500  {object}  httphelper.HTTPResponse
// @Router       /operator/{id} [get]
// @Security	 Bearer
func (h *OperatorHandler) getOperatorBydID(c *fiber.Ctx) error {
	return c.Status(fiber.StatusOK).JSON(c.Locals(httphelper.LocalObject).(*domain.Operator))
}

// createOperator godoc
// @Summary      Insert operator
// @Description  Insert operator
// @Tags         Operator
// @Accept       json
// @Produce      json
// @Param        lang query string false "Language responses"
// @Param        operator body dto.OperatorInputDTO true "Operator model"
// @Success      201  {object}  domain.Operator
// @Failure      400  {object}  httphelper.HTTPResponse
// @Failure      409  {object}  httphelper.HTTPResponse
// @Failure      500  {object}  httphelper.HTTPResponse
// @Router       /operator [post]
// @Security	 Bearer
func (h *OperatorHandler) createOperator(c *fiber.Ctx) error {
	operatorDTO := c.Locals(httphelper.LocalDTO).(*dto.OperatorInputDTO)
	operator, err := h.operatorService.CreateOperator(c.Context(), operatorDTO)
	if err != nil {
		return h.handlerError(c, err)
	}

	return c.Status(fiber.StatusCreated).JSON(operator)
}

// updateOperator godoc
// @Summary      Update operator by ID
// @Description  Update operator by ID
// @Tags         Operator
// @Accept       json
// @Produce      json
// @Param        lang query string false "Language responses"
// @Param        id     path    int     true        "Operator ID"
// @Param        operator body dto.OperatorInputDTO true "Operator model"
// @Success      200  {object}  domain.Operator
// @Failure      400  {object}  httphelper.HTTPResponse
// @Failure      404  {object}  httphelper.HTTPResponse
// @Failure      500  {object}  httphelper.HTTPResponse
// @Router       /operator/{id} [put]
// @Security	 Bearer
func (h *OperatorHandler) updateOperator(c *fiber.Ctx) error {
	operatorDTO := c.Locals(httphelper.LocalDTO).(*dto.OperatorInputDTO)
	operator := c.Locals(httphelper.LocalObject).(*domain.Operator)
	if err := h.operatorService.UpdateOperator(c.Context(), operator, operatorDTO); err != nil {
		return h.handlerError(c, err)
	}

	return c.Status(fiber.StatusOK).JSON(operator)
}

// deleteOperator godoc
// @Summary      Delete operator by ID
// @Description  Delete operator by ID
// @Tags         Operator
// @Accept       json
// @Produce      json
// @Param        lang query string false "Language responses"
// @Param        id     path    int     true        "Operator ID"
// @Success      204  {object}  nil
// @Failure      404  {object}  httphelper.HTTPResponse
// @Failure      500  {object}  httphelper.HTTPResponse
// @Router       /operator/{id} [delete]
// @Security	 Bearer
func (h *OperatorHandler) deleteOperator(c *fiber.Ctx) error {
	operator := c.Locals(httphelper.LocalObject).(*domain.Operator)
	if err := h.operatorService.DeleteOperator(c.Context(), operator); err != nil {
		return h.handlerError(c, err)
	}

	return c.Status(fiber.StatusNoContent).Send(nil)
}
