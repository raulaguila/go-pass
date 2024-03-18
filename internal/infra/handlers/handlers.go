package handlers

import (
	"log"
	"os"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/swagger"

	"github.com/raulaguila/go-pass/docs"
	"github.com/raulaguila/go-pass/internal/api/handler"
	"github.com/raulaguila/go-pass/internal/api/middleware"
	"github.com/raulaguila/go-pass/internal/api/service"
	"github.com/raulaguila/go-pass/internal/pkg/domain"
	"github.com/raulaguila/go-pass/internal/pkg/i18n"
	"github.com/raulaguila/go-pass/internal/pkg/repository"
	httphelper "github.com/raulaguila/go-pass/pkg/http-helper"

	"gorm.io/gorm"
)

var (
	profileRepository  domain.ProfileRepository
	userRepository     domain.UserRepository
	authRepository     domain.AuthRepository
	siteRepository     domain.SiteRepository
	operatorRepository domain.OperatorRepository
	phoneRepository    domain.PhoneRepository
	accountRepository  domain.AccountRepository

	profileService  domain.ProfileService
	userService     domain.UserService
	authService     domain.AuthService
	siteService     domain.SiteService
	operatorService domain.OperatorService
	phoneService    domain.PhoneService
	accountService  domain.AccountService
)

func initRepositories(postgresdb *gorm.DB) {
	// Create repositories.
	profileRepository = repository.NewProfileRepository(postgresdb)
	userRepository = repository.NewUserRepository(postgresdb)
	authRepository = repository.NewAuthRepository(userRepository)
	siteRepository = repository.NewSiteRepository(postgresdb)
	operatorRepository = repository.NewOperatorRepository(postgresdb)
	phoneRepository = repository.NewPhoneRepository(postgresdb)
	accountRepository = repository.NewAccountRepository(postgresdb)
}

func initServices() {
	// Create services.
	profileService = service.NewProfileService(profileRepository)
	userService = service.NewUserService(userRepository)
	authService = service.NewAuthService(authRepository)
	siteService = service.NewSiteService(siteRepository)
	operatorService = service.NewOperatorService(operatorRepository)
	phoneService = service.NewPhoneService(phoneRepository)
	accountService = service.NewAccountService(accountRepository)
}

func initHandelrs(app *fiber.App, postgresdb *gorm.DB) {
	reqMid := middleware.NewRequesttMiddleware(postgresdb)

	// Initialize access middleares
	middleware.MidAccess = middleware.Auth(os.Getenv("ACCESS_TOKEN_PUBLIC"), authService)
	middleware.MidRefresh = middleware.Auth(os.Getenv("RFRESH_TOKEN_PUBLIC"), authService)

	// Prepare endpoints for the API.
	handler.NewMiscHandler(app.Group(""))
	handler.NewAuthHandler(app.Group("/auth"), authService)
	handler.NewProfileHandler(app.Group("/profile"), profileService, reqMid)
	handler.NewUserHandler(app.Group("/user"), userService, reqMid)
	handler.NewSiteHandler(app.Group("/site"), siteService, reqMid)
	handler.NewOperatorHandler(app.Group("/operator"), operatorService, reqMid)
	handler.NewPhoneHandler(app.Group("/phone"), phoneService, reqMid)
	handler.NewAccountHandler(app.Group("/account"), accountService, reqMid)

	// Prepare an endpoint for 'Not Found'.
	app.All("*", func(c *fiber.Ctx) error {
		messages := c.Locals(httphelper.LocalLang).(*i18n.Translation)
		return httphelper.NewHTTPResponse(c, fiber.StatusNotFound, messages.ErrorNonexistentRoute)
	})
}

func HandleRequests(app *fiber.App, postgresdb *gorm.DB) {
	if strings.ToLower(os.Getenv("API_SWAGGO")) == "true" {
		docs.SwaggerInfo.Version = os.Getenv("SYS_VERSION")

		// 	// Config swagger
		app.Get("/swagger/*", swagger.New(swagger.Config{
			DisplayRequestDuration: true,
			DocExpansion:           "none",
			ValidatorUrl:           "none",
		}))
	}

	initRepositories(postgresdb)
	initServices()
	initHandelrs(app, postgresdb)

	log.Fatal(app.Listen(":" + os.Getenv("API_PORT")))
}
