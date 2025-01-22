package http_handler

import (
	"github.com/labstack/echo/v4"
	"net/http"
	"test-billing/commons/options"
	"test-billing/internal/middleware"
)

func RegisterEchoRoutes(opt options.Options, handler Handler) (echoRoutes *echo.Echo) {
	echoRoutes = echo.New()

	echoRoutes.GET("/health", func(c echo.Context) error {
		return c.String(http.StatusOK, "Service is healthy!")
	})
	echoRoutes.POST("/loans", handler.CreateLoan, middleware.AuthMiddleware())
	echoRoutes.GET("/loans/:id/outstanding", handler.GetOutstanding, middleware.AuthMiddleware())
	echoRoutes.GET("/loans/:id/delinquent", handler.IsDelinquent, middleware.AuthMiddleware())
	echoRoutes.POST("/loans/:id/payments", handler.MakePayment, middleware.AuthMiddleware())
	return
}
