package http_handler

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/astaxie/beego/validation"
	"github.com/labstack/echo/v4"
	"test-billing/internal/domain"
)

type Handler struct {
	Opt HandlerOptions
}

// CreateLoan handles the creation of a new loan
func (h Handler) CreateLoan(c echo.Context) (err error) {
	var loan domain.Loan
	err = c.Bind(&loan)
	if err != nil {
		return c.JSON(http.StatusBadRequest, "Invalid request payload")
	}

	//general field validation
	valid := validation.Validation{}
	valid.Required(loan.PrincipalAmount, "principal_amount")
	valid.Min(loan.CustomerID, 1, "customer_id")
	valid.Min(loan.Weeks, 1, "weeks")
	valid.Min(loan.InterestRate, 1, "interest_rate")

	if valid.HasErrors() {
		for _, errValidate := range valid.Errors {
			return c.JSON(http.StatusUnauthorized, fmt.Sprintf("Failed to create loan : %v %v %v", errValidate.Key, errValidate.Value, errValidate.Message))
		}
	}

	err = h.Opt.LoanService.CreateLoan(&loan)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, fmt.Sprintf("Failed to create loan: %v", err))
	}

	return c.JSON(http.StatusCreated, map[string]interface{}{
		"message": "Loan created successfully",
		"loan_id": loan.ID,
	})
}

// GetOutstanding handles fetching the outstanding balance of a loan
func (h Handler) GetOutstanding(c echo.Context) (err error) {

	loanID, _ := strconv.Atoi(c.Param("id"))
	if loanID == 0 {
		return c.JSON(http.StatusBadRequest, "Invalid loan id parameter")
	}

	outstanding, err := h.Opt.LoanService.GetOutstanding(loanID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, fmt.Sprintf("Failed to fetch outstanding: %v", err))
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"outstanding_balance": outstanding,
	})
}

// IsDelinquent checks if a loan is delinquent
func (h Handler) IsDelinquent(c echo.Context) (err error) {
	loanID, _ := strconv.Atoi(c.Param("id"))
	if loanID == 0 {
		return c.JSON(http.StatusBadRequest, "Invalid loan id parameter")
	}

	delinquent, err := h.Opt.LoanService.IsDelinquent(loanID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, fmt.Sprintf("Failed to check delinquency: %v", err))
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"is_delinquent": delinquent,
	})
}

// MakePayment handles making a payment towards a loan
func (h Handler) MakePayment(c echo.Context) (err error) {
	type PaymentRequest struct {
		Amount float64 `json:"amount"`
	}

	loanID, _ := strconv.Atoi(c.Param("id"))
	if loanID == 0 {
		return c.JSON(http.StatusBadRequest, "Invalid loan id parameter")
	}

	var paymentReq PaymentRequest
	err = c.Bind(&paymentReq)
	if err != nil {
		return c.JSON(http.StatusBadRequest, "Invalid request payload")
	}

	//general field validation
	valid := validation.Validation{}
	valid.Min(loanID, 1, "loan_id")

	if valid.HasErrors() {
		for _, errValidate := range valid.Errors {
			return c.JSON(http.StatusUnauthorized, fmt.Sprintf("Failed to make payment : %v %v %v", errValidate.Key, errValidate.Value, errValidate.Message))
		}
	}

	err = h.Opt.LoanService.MakePayment(loanID, paymentReq.Amount)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, fmt.Sprintf("Failed to make payment: %v", err))
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"message": "Payment made successfully",
	})
}
