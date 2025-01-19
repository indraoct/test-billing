package http_handler

import (
	"github.com/labstack/echo/v4"
	"test-billing/commons/options"
)

func RegisterEchoRoutes(opt options.Options, handler Handler) (echoRoutes *echo.Echo) {
	echoRoutes = echo.New()

	echoRoutes.POST("/loans", handler.CreateLoan)
	echoRoutes.GET("/loans/:id/outstanding", handler.GetOutstanding)
	echoRoutes.GET("/loans/:id/delinquent", handler.IsDelinquent)
	echoRoutes.POST("/loans/:id/payments", handler.MakePayment)
	return
}
