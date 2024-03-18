package middleware

import (
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/raulaguila/go-pass/internal/pkg/dto"
	"github.com/raulaguila/go-pass/internal/pkg/i18n"
	httphelper "github.com/raulaguila/go-pass/pkg/http-helper"
)

func getDTO(c *fiber.Ctx, dto interface{}) error {
	if err := c.BodyParser(dto); err != nil {
		log.Println(err.Error())
		messages := c.Locals(httphelper.LocalLang).(*i18n.Translation)
		return httphelper.NewHTTPResponse(c, fiber.StatusBadRequest, messages.ErrInvalidDatas)
	}

	c.Locals(httphelper.LocalDTO, dto)
	return c.Next()
}

func GetProfileDTO(c *fiber.Ctx) error {
	return getDTO(c, &dto.ProfileInputDTO{})
}

func GetUserDTO(c *fiber.Ctx) error {
	return getDTO(c, &dto.UserInputDTO{})
}

func GetSiteDTO(c *fiber.Ctx) error {
	return getDTO(c, &dto.SiteInputDTO{})
}

func GetPhoneDTO(c *fiber.Ctx) error {
	return getDTO(c, &dto.PhoneInputDTO{})
}

func GetOperatorDTO(c *fiber.Ctx) error {
	return getDTO(c, &dto.OperatorInputDTO{})
}

func GetAccountDTO(c *fiber.Ctx) error {
	return getDTO(c, &dto.AccountInputDTO{})
}

func GetPasswordInputDTO(c *fiber.Ctx) error {
	return getDTO(c, &dto.PasswordInputDTO{})
}