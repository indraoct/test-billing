package job

import (
	"fmt"
	"test-billing/commons/options"
	"test-billing/internal/repository"
	"test-billing/pkg/queue"
	"time"
)

type JobOptions struct {
	options.Options
	repository.Repository
}

type RepaymentScheduleJob struct {
	Opt JobOptions
}

func (job RepaymentScheduleJob) Run() {
	job.Opt.Logger.Println("Starting repayment schedule notification job...")

	// Fetch upcoming repayments
	repaymentSchedules, err := job.Opt.LoanRepo.GetUpcomingRepayments(7)
	if err != nil {
		job.Opt.Logger.Printf("Error fetching upcoming repayments: %v\n", err)
		return
	}

	// Enqueue notifications
	for _, schedule := range repaymentSchedules {
		message := "Reminder: Your repayment of IDR " +
			formatAmount(schedule.WeeklyPayment) +
			" is due on " +
			schedule.DueDate.Format("2006-01-02") +
			". Please make your payment on time to avoid penalties."
		notification := queue.Notification{
			CustomerID: schedule.CustomerID,
			LoanID:     schedule.LoanID,
			Message:    message,
			Timestamp:  time.Now(),
		}
		job.Opt.Queue.Enqueue(notification)
		job.Opt.Logger.Printf("Notification enqueued for customer %d, loan %d, due on %s\n", schedule.CustomerID, schedule.LoanID, schedule.DueDate.Format("2006-01-02"))
	}

	job.Opt.Logger.Println("Repayment schedule notification job completed.")
}

// Helper function to format the amount.
func formatAmount(amount float64) string {
	return fmt.Sprintf("%.2f", amount)
}
