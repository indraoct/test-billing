package main

import (
	"context"
	"fmt"
	_ "github.com/lib/pq"
	"github.com/sirupsen/logrus"
	"os"
	"os/signal"
	"runtime"
	"syscall"
	"test-billing/commons/config"
	"test-billing/commons/options"
	"test-billing/internal/delivery/http_handler"
	"test-billing/internal/delivery/job"
	"test-billing/internal/repository"
	"test-billing/internal/service"
	"test-billing/pkg/queue"
	"test-billing/pkg/utils"
)

func main() {

	var err error

	// Read config and check if error occurred
	conf := config.Load()

	//general info
	goVersion := runtime.Version()
	osArch := fmt.Sprintf("%s %s", runtime.GOOS, runtime.GOARCH)
	fmt.Println("App Name	:", conf.Name)
	fmt.Println("Go Version	:", goVersion)
	fmt.Println("OS Arch	:", osArch)

	dbPostgres := &utils.DB{
		Url:             conf.DatabaseUrl,
		MaxOpenConns:    conf.DatabaseMaxOpenConn,
		MaxIdleConns:    conf.DatabaseMaxIdleConn,
		ConnMaxLifetime: conf.DatabaseSetConnMaxIdleTime,
	}
	dbPostgres.Connect()

	// Initialize the notification queue with a buffer size of 100
	q := queue.NewNotificationQueue(100)

	// Start consumer workers (e.g., 3 workers)
	q.StartConsumer(3)

	// Ensure the queue is stopped gracefully on service shutdown
	defer q.Stop()

	//options
	opt := options.Options{
		Config:     conf,
		Logger:     logrus.New(),
		DBPostgres: dbPostgres,
		Queue:      q,
	}

	//repositories
	repo := wiringRepository(repository.RepositoryOptions{Options: opt})

	//services
	services := wiringService(service.ServiceOptions{
		Options:    opt,
		Repository: repo,
	})

	//handlers
	handler := http_handler.Handler{http_handler.HandlerOptions{
		Options: opt,
		Service: services,
	}}

	//set router
	router := http_handler.RegisterEchoRoutes(opt, handler)

	// Start run echo router
	go func() {
		err = router.Start(fmt.Sprintf(":%s", conf.ListenPort))
		if err != nil {
			logrus.WithField("error", err).Error("Echo server has closed")
			os.Exit(1)
		}
	}()

	j := job.RepaymentScheduleJob{job.JobOptions{
		Options:    opt,
		Repository: repo,
	}}

	//running the job repayment notification concurrently
	go func() {
		j.Run()
	}()

	c := make(chan os.Signal, 1)
	// We'll accept graceful shutdowns when quit via SIGINT (Ctrl+C)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)

	// Block until we receive our signal.
	ctx := context.Background()

	// this channel waiting interruption signal
	<-c
	logrus.WithField("service", "stop").Info("server stop gracefully")
	err = router.Server.Shutdown(ctx)
	if err != nil {
		logrus.WithField("error", err).Error("unable to close the server")
		os.Exit(1)
	}

}

func wiringRepository(opt repository.RepositoryOptions) repository.Repository {
	repoLoan := repository.NewRepositoryLoan(opt)
	return repository.Repository{LoanRepo: repoLoan}
}

func wiringService(opt service.ServiceOptions) service.Service {
	svcLoan := service.NewServiceLoan(opt)
	return service.Service{LoanService: svcLoan}
}
