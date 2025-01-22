package domain

import "time"

// RepaymentSchedule represents the repayment details for a specific week of a loan
type RepaymentSchedule struct {
	ID            int       `json:"id"`          // Unique identifier for the repayment schedule
	LoanID        int       `json:"loan_id"`     // The associated loan ID
	WeekNumber    int       `json:"week_number"` // The week number of the repayment
	DueDate       time.Time `json:"due_date"`    // The due date for this repayment
	Status        string    `json:"status"`      // The status of the repayment (e.g., "Paid", "Unpaid")
	CustomerID    int       `json:"customer_id"`
	WeeklyPayment float64   `json:"weekly_payment"`
}
