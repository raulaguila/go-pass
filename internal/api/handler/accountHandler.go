package handler

import (
	"errors"
	"fmt"
	"log"
	"strings"

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

type AccountHandler struct {
	accountService domain.AccountService
}

func (AccountHandler) foreignKeyViolatedMethod(c *fiber.Ctx, err error, translation *i18n.Translation) error {
	switch c.Method() {
	case fiber.MethodPut, fiber.MethodPost, fiber.MethodPatch:
		if strings.Contains(err.Error(), "site") {
			return httphelper.NewHTTPResponse(c, fiber.StatusBadRequest, translation.ErrSiteNotFound)
		}
		if strings.Contains(err.Error(), "phone") {
			return httphelper.NewHTTPResponse(c, fiber.StatusBadRequest, translation.ErrPhoneNotFound)
		}
		return httphelper.NewHTTPResponse(c, fiber.StatusBadRequest, translation.ErrAccountNotFound)
	case fiber.MethodDelete:
		return httphelper.NewHTTPResponse(c, fiber.StatusBadRequest, translation.ErrAccountUsed)
	default:
		return httphelper.NewHTTPResponse(c, fiber.StatusInternalServerError, translation.ErrGeneric)
	}
}

func (s AccountHandler) handlerError(c *fiber.Ctx, err error) error {
	translation := c.Locals(httphelper.LocalLang).(*i18n.Translation)
	fmt.Printf("Error: %\n", err.Error())

	switch pgerror.HandlerError(err) {
	case pgerror.ErrDuplicatedKey:
		return httphelper.NewHTTPResponse(c, fiber.StatusConflict, translation.ErrAccountRegistered)
	case pgerror.ErrForeignKeyViolated:
		return s.foreignKeyViolatedMethod(c, err, translation)
	case pgerror.ErrUndefinedColumn:
		return httphelper.NewHTTPResponse(c, fiber.StatusBadRequest, translation.ErrUndefinedColumn)
	case translation.ErrInvalidId:
		return httphelper.NewHTTPResponse(c, fiber.StatusBadRequest, translation.ErrInvalidId)
	}

	if errors.As(err, &validator.ErrValidator) {
		return httphelper.NewHTTPResponse(c, fiber.StatusBadRequest, err)
	}

	log.Println(err.Error())
	return httphelper.NewHTTPResponse(c, fiber.StatusInternalServerError, translation.ErrGeneric)
}

// Creates a new handler.
func NewAccountHandler(route fiber.Router, ps domain.AccountService, mid *middleware.RequesttMiddleware) {
	handler := &AccountHandler{
		accountService: ps,
	}

	route.Use(middleware.MidAccess)

	route.Get("", middleware.GetGenericFilter, handler.getAccounts)
	route.Post("", middleware.GetAccountDTO, handler.createAccount)
	route.Get("/:"+httphelper.ParamID, mid.AccountByID, handler.getAccountBydID)
	route.Get("/:"+httphelper.ParamID+"/pass", mid.AccountByID, handler.getAccountPasswordBydID)
	route.Put("/:"+httphelper.ParamID, mid.AccountByID, middleware.GetAccountDTO, handler.updateAccount)
	route.Delete("/:"+httphelper.ParamID, mid.AccountByID, handler.deleteAccount)
}

// getAccounts godoc
// @Summary      Get accounts
// @Description  Get accounts
// @Tags         Account
// @Accept       json
// @Produce      json
// @Param        lang query string false "Language responses"
// @Param        filter query filter.Filter false "Optional Filter"
// @Success      200  {array}   dto.ItemsOutputDTO
// @Failure      500  {object}  httphelper.HTTPResponse
// @Router       /account [get]
// @Security	 Bearer
func (h *AccountHandler) getAccounts(c *fiber.Ctx) error {
	user := c.Locals(httphelper.LocalUser).(*domain.User)
	response, err := h.accountService.GetAccountsOutputDTO(c.Context(), c.Locals(httphelper.LocalFilter).(*filter.Filter), user.Id)
	if err != nil {
		return h.handlerError(c, err)
	}

	return c.Status(fiber.StatusOK).JSON(response)
}

// getAccountBydID godoc
// @Summary      Get account by ID
// @Description  Get account by ID
// @Tags         Account
// @Accept       json
// @Produce      json
// @Param        lang query string false "Language responses"
// @Param        id   path			int			true        "Account ID"
// @Success      200  {object}  domain.Account
// @Failure      400  {object}  httphelper.HTTPResponse
// @Failure      404  {object}  httphelper.HTTPResponse
// @Failure      500  {object}  httphelper.HTTPResponse
// @Router       /account/{id} [get]
// @Security	 Bearer
func (h *AccountHandler) getAccountBydID(c *fiber.Ctx) error {
	return c.Status(fiber.StatusOK).JSON(c.Locals(httphelper.LocalObject).(*domain.Account))
}

// getAccountPasswordBydID godoc
// @Summary      Get password account by ID
// @Description  Get password account by ID
// @Tags         Account
// @Accept       json
// @Produce      json
// @Param        lang query string false "Language responses"
// @Param        id   path			int			true        "Account ID"
// @Success      200  {object}  map[string]string
// @Failure      400  {object}  httphelper.HTTPResponse
// @Failure      404  {object}  httphelper.HTTPResponse
// @Failure      500  {object}  httphelper.HTTPResponse
// @Router       /account/{id}/pass [get]
// @Security	 Bearer
func (h *AccountHandler) getAccountPasswordBydID(c *fiber.Ctx) error {
	account := c.Locals(httphelper.LocalObject).(*domain.Account)
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"password": account.DecodePass(),
	})
}

// createAccount godoc
// @Summary      Insert account
// @Description  Insert account
// @Tags         Account
// @Accept       json
// @Produce      json
// @Param        lang query string false "Language responses"
// @Param        account body dto.AccountInputDTO true "Account model"
// @Success      201  {object}  domain.Account
// @Failure      400  {object}  httphelper.HTTPResponse
// @Failure      409  {object}  httphelper.HTTPResponse
// @Failure      500  {object}  httphelper.HTTPResponse
// @Router       /account [post]
// @Security	 Bearer
func (h *AccountHandler) createAccount(c *fiber.Ctx) error {
	user := c.Locals(httphelper.LocalUser).(*domain.User)
	accountDTO := c.Locals(httphelper.LocalDTO).(*dto.AccountInputDTO)
	account, err := h.accountService.CreateAccount(c.Context(), accountDTO, user.Id)
	if err != nil {
		return h.handlerError(c, err)
	}

	return c.Status(fiber.StatusCreated).JSON(account)
}

// updateAccount godoc
// @Summary      Update account by ID
// @Description  Update account by ID
// @Tags         Account
// @Accept       json
// @Produce      json
// @Param        lang query string false "Language responses"
// @Param        id     path    int     true        "Account ID"
// @Param        account body dto.AccountInputDTO true "Account model"
// @Success      200  {object}  domain.Account
// @Failure      400  {object}  httphelper.HTTPResponse
// @Failure      404  {object}  httphelper.HTTPResponse
// @Failure      500  {object}  httphelper.HTTPResponse
// @Router       /account/{id} [put]
// @Security	 Bearer
func (h *AccountHandler) updateAccount(c *fiber.Ctx) error {
	user := c.Locals(httphelper.LocalUser).(*domain.User)
	accountDTO := c.Locals(httphelper.LocalDTO).(*dto.AccountInputDTO)
	account := c.Locals(httphelper.LocalObject).(*domain.Account)
	account, err := h.accountService.UpdateAccount(c.Context(), account, accountDTO, user.Id)
	if err != nil {
		return h.handlerError(c, err)
	}

	return c.Status(fiber.StatusOK).JSON(account)
}

// deleteAccount godoc
// @Summary      Delete account by ID
// @Description  Delete account by ID
// @Tags         Account
// @Accept       json
// @Produce      json
// @Param        lang query string false "Language responses"
// @Param        id     path    int     true        "Account ID"
// @Success      204  {object}  nil
// @Failure      404  {object}  httphelper.HTTPResponse
// @Failure      500  {object}  httphelper.HTTPResponse
// @Router       /account/{id} [delete]
// @Security	 Bearer
func (h *AccountHandler) deleteAccount(c *fiber.Ctx) error {
	user := c.Locals(httphelper.LocalUser).(*domain.User)
	account := c.Locals(httphelper.LocalObject).(*domain.Account)
	if err := h.accountService.DeleteAccount(c.Context(), account, user.Id); err != nil {
		return h.handlerError(c, err)
	}

	return c.Status(fiber.StatusNoContent).Send(nil)
}
