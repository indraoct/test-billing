package http_handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"test-billing/internal/domain"
)

type Handler struct {
	Opt HandlerOptions
}

// CreateLoan handles the creation of a new loan
func (h Handler) CreateLoan(w http.ResponseWriter, r *http.Request) {
	var loan domain.Loan
	err := json.NewDecoder(r.Body).Decode(&loan)
	if err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	err = h.Opt.LoanService.CreateLoan(&loan)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to create loan: %v", err), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"message": "Loan created successfully",
		"loan_id": loan.ID,
	})
}

// GetOutstanding handles fetching the outstanding balance of a loan
func (h Handler) GetOutstanding(w http.ResponseWriter, r *http.Request) {
	loanID, err := strconv.Atoi(r.URL.Query().Get("loan_id"))
	if err != nil {
		http.Error(w, "Invalid loan_id parameter", http.StatusBadRequest)
		return
	}

	outstanding, err := h.Opt.LoanService.GetOutstanding(loanID)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to fetch outstanding: %v", err), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(map[string]interface{}{
		"outstanding_balance": outstanding,
	})
}

// IsDelinquent checks if a loan is delinquent
func (h Handler) IsDelinquent(w http.ResponseWriter, r *http.Request) {
	loanID, err := strconv.Atoi(r.URL.Query().Get("loan_id"))
	if err != nil {
		http.Error(w, "Invalid loan_id parameter", http.StatusBadRequest)
		return
	}

	delinquent, err := h.Opt.LoanService.IsDelinquent(loanID)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to check delinquency: %v", err), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(map[string]interface{}{
		"is_delinquent": delinquent,
	})
}

// MakePayment handles making a payment towards a loan
func (h Handler) MakePayment(w http.ResponseWriter, r *http.Request) {
	type PaymentRequest struct {
		LoanID int     `json:"loan_id"`
		Amount float64 `json:"amount"`
	}

	var paymentReq PaymentRequest
	err := json.NewDecoder(r.Body).Decode(&paymentReq)
	if err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	err = h.Opt.LoanService.MakePayment(paymentReq.LoanID, paymentReq.Amount)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to make payment: %v", err), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"message": "Payment made successfully",
	})
}
