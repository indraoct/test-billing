package repository

import "test-billing/commons/options"

type RepositoryOptions struct {
	options.Options
}

type Repository struct {
	LoanRepo LoanRepository
}
