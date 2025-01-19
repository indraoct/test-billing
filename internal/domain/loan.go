package domain

type Loan struct {
	ID                 int
	CustomerID         int
	PrincipalAmount    float64
	InterestRate       float64
	Weeks              int
	WeeklyPayment      float64
	OutstandingBalance float64
	Status             string // Active, Paid, Delinquent
}
