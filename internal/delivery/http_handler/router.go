package http_handler

import (
	"github.com/labstack/echo/v4"
	"test-billing/commons/options"
)

func RegisterEchoRoutes(opt options.Options, handler Handler) (echoRoutes *echo.Echo) {
	echoRoutes = echo.New()

	return echoRoutes
}
