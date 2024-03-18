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

type SiteHandler struct {
	siteService domain.SiteService
}

func (SiteHandler) foreignKeyViolatedMethod(c *fiber.Ctx, translation *i18n.Translation) error {
	switch c.Method() {
	case fiber.MethodPut, fiber.MethodPost, fiber.MethodPatch:
		return httphelper.NewHTTPResponse(c, fiber.StatusBadRequest, translation.ErrSiteNotFound)
	case fiber.MethodDelete:
		return httphelper.NewHTTPResponse(c, fiber.StatusBadRequest, translation.ErrSiteUsed)
	default:
		return httphelper.NewHTTPResponse(c, fiber.StatusInternalServerError, translation.ErrGeneric)
	}
}

func (s SiteHandler) handlerError(c *fiber.Ctx, err error) error {
	translation := c.Locals(httphelper.LocalLang).(*i18n.Translation)

	switch pgerror.HandlerError(err) {
	case pgerror.ErrDuplicatedKey:
		return httphelper.NewHTTPResponse(c, fiber.StatusConflict, translation.ErrSiteRegistered)
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
func NewSiteHandler(route fiber.Router, ps domain.SiteService, mid *middleware.RequesttMiddleware) {
	handler := &SiteHandler{
		siteService: ps,
	}

	route.Use(middleware.MidAccess)

	route.Get("", middleware.GetGenericFilter, handler.getSites)
	route.Post("", middleware.GetSiteDTO, handler.createSite)
	route.Get("/:"+httphelper.ParamID, mid.SiteByID, handler.getSiteBydID)
	route.Put("/:"+httphelper.ParamID, mid.SiteByID, middleware.GetSiteDTO, handler.updateSite)
	route.Delete("/:"+httphelper.ParamID, mid.SiteByID, handler.deleteSite)
}

// getSites godoc
// @Summary      Get sites
// @Description  Get sites
// @Tags         Site
// @Accept       json
// @Produce      json
// @Param        lang query string false "Language responses"
// @Param        filter query filter.Filter false "Optional Filter"
// @Success      200  {array}   dto.ItemsOutputDTO
// @Failure      500  {object}  httphelper.HTTPResponse
// @Router       /site [get]
// @Security	 Bearer
func (h *SiteHandler) getSites(c *fiber.Ctx) error {
	response, err := h.siteService.GetSitesOutputDTO(c.Context(), c.Locals(httphelper.LocalFilter).(*filter.Filter))
	if err != nil {
		return h.handlerError(c, err)
	}

	return c.Status(fiber.StatusOK).JSON(response)
}

// getSiteBydID godoc
// @Summary      Get site by ID
// @Description  Get site by ID
// @Tags         Site
// @Accept       json
// @Produce      json
// @Param        lang query string false "Language responses"
// @Param        id   path			int			true        "Site ID"
// @Success      200  {object}  domain.Site
// @Failure      400  {object}  httphelper.HTTPResponse
// @Failure      404  {object}  httphelper.HTTPResponse
// @Failure      500  {object}  httphelper.HTTPResponse
// @Router       /site/{id} [get]
// @Security	 Bearer
func (h *SiteHandler) getSiteBydID(c *fiber.Ctx) error {
	return c.Status(fiber.StatusOK).JSON(c.Locals(httphelper.LocalObject).(*domain.Site))
}

// createSite godoc
// @Summary      Insert site
// @Description  Insert site
// @Tags         Site
// @Accept       json
// @Produce      json
// @Param        lang query string false "Language responses"
// @Param        site body dto.SiteInputDTO true "Site model"
// @Success      201  {object}  domain.Site
// @Failure      400  {object}  httphelper.HTTPResponse
// @Failure      409  {object}  httphelper.HTTPResponse
// @Failure      500  {object}  httphelper.HTTPResponse
// @Router       /site [post]
// @Security	 Bearer
func (h *SiteHandler) createSite(c *fiber.Ctx) error {
	siteDTO := c.Locals(httphelper.LocalDTO).(*dto.SiteInputDTO)
	site, err := h.siteService.CreateSite(c.Context(), siteDTO)
	if err != nil {
		return h.handlerError(c, err)
	}

	return c.Status(fiber.StatusCreated).JSON(site)
}

// updateSite godoc
// @Summary      Update site by ID
// @Description  Update site by ID
// @Tags         Site
// @Accept       json
// @Produce      json
// @Param        lang query string false "Language responses"
// @Param        id     path    int     true        "Site ID"
// @Param        site body dto.SiteInputDTO true "Site model"
// @Success      200  {object}  domain.Site
// @Failure      400  {object}  httphelper.HTTPResponse
// @Failure      404  {object}  httphelper.HTTPResponse
// @Failure      500  {object}  httphelper.HTTPResponse
// @Router       /site/{id} [put]
// @Security	 Bearer
func (h *SiteHandler) updateSite(c *fiber.Ctx) error {
	siteDTO := c.Locals(httphelper.LocalDTO).(*dto.SiteInputDTO)
	site := c.Locals(httphelper.LocalObject).(*domain.Site)
	if err := h.siteService.UpdateSite(c.Context(), site, siteDTO); err != nil {
		return h.handlerError(c, err)
	}

	return c.Status(fiber.StatusOK).JSON(site)
}

// deleteSite godoc
// @Summary      Delete site by ID
// @Description  Delete site by ID
// @Tags         Site
// @Accept       json
// @Produce      json
// @Param        lang query string false "Language responses"
// @Param        id     path    int     true        "Site ID"
// @Success      204  {object}  nil
// @Failure      404  {object}  httphelper.HTTPResponse
// @Failure      500  {object}  httphelper.HTTPResponse
// @Router       /site/{id} [delete]
// @Security	 Bearer
func (h *SiteHandler) deleteSite(c *fiber.Ctx) error {
	site := c.Locals(httphelper.LocalObject).(*domain.Site)
	if err := h.siteService.DeleteSite(c.Context(), site); err != nil {
		return h.handlerError(c, err)
	}

	return c.Status(fiber.StatusNoContent).Send(nil)
}
