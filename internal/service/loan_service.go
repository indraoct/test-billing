package service

import (
	"errors"
	"fmt"
	"test-billing/commons/constants"
	"test-billing/pkg/queue"
	"time"

	"test-billing/internal/domain"
)

type LoanService struct {
	opt ServiceOptions
}

func NewServiceLoan(opt ServiceOptions) LoanService {
	return LoanService{opt: opt}
}

// CreateLoan creates a new loan and its repayment schedule
func (s LoanService) CreateLoan(loan *domain.Loan) (err error) {
	loan.OutstandingBalance = loan.PrincipalAmount + (loan.PrincipalAmount * loan.InterestRate / 100.0)
	loan.WeeklyPayment = loan.OutstandingBalance / float64(loan.Weeks)
	loan.Status = constants.LoanStatusActive

	loanID, err := s.opt.Repository.LoanRepo.CreateLoan(loan)
	if err != nil {
		return err
	}

	// Generate repayment schedule
	startDate := time.Now()
	for week := 1; week <= loan.Weeks; week++ {
		schedule := domain.RepaymentSchedule{
			LoanID:     loanID,
			WeekNumber: week,
			DueDate:    startDate.AddDate(0, 0, week*7),
			Status:     constants.ScheduleUnpaid,
		}
		err = s.opt.Repository.LoanRepo.CreateRepaymentSchedule(&schedule)
		if err != nil {
			return err
		}
	}

	return nil
}

// GetOutstanding returns the outstanding balance for a loan
func (s LoanService) GetOutstanding(loanID int) (float64, error) {
	loan, err := s.opt.Repository.LoanRepo.GetLoanByID(loanID)
	if err != nil {
		return 0, err
	}
	return loan.OutstandingBalance, nil
}

// IsDelinquent checks if a loan is delinquent
func (s LoanService) IsDelinquent(loanID int) (bool, error) {
	schedules, err := s.opt.Repository.LoanRepo.GetRepaymentSchedules(loanID)
	if err != nil {
		return false, err
	}

	// Count consecutive unpaid weeks
	consecutiveUnpaid := 0
	for _, schedule := range schedules {
		if schedule.Status == constants.ScheduleUnpaid {
			consecutiveUnpaid++
			if consecutiveUnpaid >= 2 {
				return true, nil
			}
		} else {
			consecutiveUnpaid = 0
		}
	}
	return false, nil
}

// MakePayment processes a payment for a loan
func (s LoanService) MakePayment(loanID int, amount float64) (err error) {
	loan, err := s.opt.Repository.LoanRepo.GetLoanByID(loanID)
	if err != nil {
		return err
	}

	if amount != loan.WeeklyPayment {
		return errors.New("payment amount must match the weekly payment")
	}

	// Find the next unpaid schedule
	schedules, err := s.opt.Repository.LoanRepo.GetRepaymentSchedules(loanID)
	if err != nil {
		return err
	}

	for _, schedule := range schedules {
		if schedule.Status == constants.ScheduleUnpaid {
			// Mark the schedule as paid
			err = s.opt.Repository.LoanRepo.UpdateRepaymentScheduleStatus(schedule.ID, constants.SchedulePaid)
			if err != nil {
				return err
			}

			// Update the loan's outstanding balance
			loan.OutstandingBalance -= amount
			err = s.opt.Repository.LoanRepo.UpdateLoanOutstandingBalance(loanID, loan.OutstandingBalance)
			if err != nil {
				return err
			}

			// Enqueue a notification for successful payment
			notification := queue.Notification{
				CustomerID: loan.CustomerID,
				LoanID:     loan.ID,
				Message:    fmt.Sprintf("Your payment of IDR %v has been received. Thank you!", amount),
				Timestamp:  time.Now(),
			}
			s.opt.Queue.Enqueue(notification)
			s.opt.Logger.WithField("notification", notification).Info("Success sending notification")

			return nil
		}
	}

	return errors.New("no unpaid schedules available for payment")
}
