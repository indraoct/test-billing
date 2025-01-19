package service

import (
	"test-billing/commons/options"
	"test-billing/internal/repository"
)

type ServiceOptions struct {
	options.Options
	repository.Repository
}

type Service struct {
	LoanService LoanService
}
