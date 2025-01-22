package repository

import (
	"database/sql"
	"errors"
	_ "github.com/lib/pq"
	"test-billing/internal/domain"
)

type LoanRepository struct {
	opt RepositoryOptions
}

func NewRepositoryLoan(opt RepositoryOptions) LoanRepository {
	return LoanRepository{opt: opt}
}

// CreateLoan inserts a new loan into the database and returns the loan ID
func (r LoanRepository) CreateLoan(loan *domain.Loan) (int, error) {
	query := `
		INSERT INTO Loans (customer_id, principal_amount, interest_rate, weeks, weekly_payment, outstanding_balance, status)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		RETURNING id
	`

	var loanID int
	err := r.opt.DBPostgres.DB.QueryRow(query, loan.CustomerID, loan.PrincipalAmount, loan.InterestRate, loan.Weeks, loan.WeeklyPayment, loan.OutstandingBalance, loan.Status).Scan(&loanID)
	if err != nil {
		return 0, err
	}

	return loanID, nil
}

// CreateRepaymentSchedule inserts a new repayment schedule into the database
func (r LoanRepository) CreateRepaymentSchedule(schedule *domain.RepaymentSchedule) error {
	query := `
		INSERT INTO RepaymentSchedules (loan_id, week_number, due_date, status)
		VALUES ($1, $2, $3, $4)
	`

	_, err := r.opt.DBPostgres.DB.Exec(query, schedule.LoanID, schedule.WeekNumber, schedule.DueDate, schedule.Status)
	return err
}

// GetLoanByID fetches a loan by its ID
func (r LoanRepository) GetLoanByID(loanID int) (*domain.Loan, error) {
	query := `
		SELECT id, customer_id, principal_amount, interest_rate, weeks, weekly_payment, outstanding_balance, status
		FROM Loans
		WHERE id = $1
	`

	var loan domain.Loan
	err := r.opt.DBPostgres.DB.QueryRow(query, loanID).Scan(
		&loan.ID, &loan.CustomerID, &loan.PrincipalAmount, &loan.InterestRate,
		&loan.Weeks, &loan.WeeklyPayment, &loan.OutstandingBalance, &loan.Status,
	)
	if err == sql.ErrNoRows {
		return nil, errors.New("loan not found")
	} else if err != nil {
		return nil, err
	}

	return &loan, nil
}

// GetRepaymentSchedules fetches all repayment schedules for a given loan
func (r LoanRepository) GetRepaymentSchedules(loanID int) ([]domain.RepaymentSchedule, error) {
	query := `
		SELECT id, loan_id, week_number, due_date, status
		FROM RepaymentSchedules
		WHERE loan_id = $1
		ORDER BY week_number
	`

	rows, err := r.opt.DBPostgres.DB.Query(query, loanID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var schedules []domain.RepaymentSchedule
	for rows.Next() {
		var schedule domain.RepaymentSchedule
		err := rows.Scan(&schedule.ID, &schedule.LoanID, &schedule.WeekNumber, &schedule.DueDate, &schedule.Status)
		if err != nil {
			return nil, err
		}
		schedules = append(schedules, schedule)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return schedules, nil
}

// UpdateRepaymentScheduleStatus updates the status of a repayment schedule
func (r LoanRepository) UpdateRepaymentScheduleStatus(scheduleID int, status string) error {
	query := `
		UPDATE RepaymentSchedules
		SET status = $1, updated_at = NOW()
		WHERE id = $2
	`

	_, err := r.opt.DBPostgres.DB.Exec(query, status, scheduleID)
	return err
}

// UpdateLoanOutstandingBalance updates the outstanding balance of a loan
func (r LoanRepository) UpdateLoanOutstandingBalance(loanID int, outstandingBalance float64) error {
	query := `
		UPDATE Loans
		SET outstanding_balance = $1, updated_at = NOW()
		WHERE id = $2
	`

	_, err := r.opt.DBPostgres.DB.Exec(query, outstandingBalance, loanID)
	return err
}

// GetUpcomingRepayments fetches repayment schedules due within the given number of days.
func (r LoanRepository) GetUpcomingRepayments(days int) ([]domain.RepaymentSchedule, error) {
	query := `
		SELECT 
			rs.loan_id, 
			l.customer_id, 
			rs.due_date, 
			l.weekly_payment 
		FROM 
			repaymentschedules rs
		INNER JOIN 
			loans l 
		ON 
			rs.loan_id = l.id
		WHERE 
			rs.due_date BETWEEN CURRENT_DATE AND CURRENT_DATE + $1::INTERVAL
	`
	rows, err := r.opt.DBPostgres.DB.Query(query, days)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var schedules []domain.RepaymentSchedule
	for rows.Next() {
		var schedule domain.RepaymentSchedule
		if err = rows.Scan(&schedule.LoanID, &schedule.CustomerID, &schedule.DueDate, &schedule.WeeklyPayment); err != nil {
			return nil, err
		}
		schedules = append(schedules, schedule)
	}

	return schedules, nil
}
