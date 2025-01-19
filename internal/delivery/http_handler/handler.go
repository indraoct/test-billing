package http_handler

import (
	"test-billing/commons/options"
	"test-billing/internal/service"
)

type HandlerOptions struct {
	options.Options
	service.Service
}
