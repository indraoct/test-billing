package domain

type Loan struct {
	ID                 int     `json:"id"`
	CustomerID         int     `json:"customer_id"`
	PrincipalAmount    float64 `json:"principal_amount"`
	InterestRate       float64 `json:"interest_rate"`
	Weeks              int     `json:"weeks"`
	WeeklyPayment      float64 `json:"weekly_payment"`
	OutstandingBalance float64 `json:"outstanding_balance"`
	Status             string  `json:"status"` // Active, Paid, Delinquent
}
