package queue

import (
	"fmt"
	"log"
	"sync"
	"time"
)

// Notification represents a notification message.
type Notification struct {
	CustomerID int
	LoanID     int
	Message    string
	Timestamp  time.Time
}

// NotificationQueue handles enqueuing and processing notifications.
type NotificationQueue struct {
	queue    chan Notification
	wg       sync.WaitGroup
	shutdown chan struct{}
}

// NewNotificationQueue creates a new notification queue.
func NewNotificationQueue(bufferSize int) *NotificationQueue {
	return &NotificationQueue{
		queue:    make(chan Notification, bufferSize),
		shutdown: make(chan struct{}),
	}
}

// Enqueue adds a notification to the queue.
func (nq *NotificationQueue) Enqueue(notification Notification) {
	select {
	case nq.queue <- notification:
		log.Printf("Notification enqueued: %+v", notification)
	default:
		log.Println("Notification queue is full, dropping notification")
	}
}

// StartConsumer starts a worker to process notifications.
func (nq *NotificationQueue) StartConsumer(workerCount int) {
	for i := 0; i < workerCount; i++ {
		nq.wg.Add(1)
		go nq.consumerWorker(i)
	}
}

// consumerWorker processes notifications from the queue.
func (nq *NotificationQueue) consumerWorker(workerID int) {
	defer nq.wg.Done()

	for {
		select {
		case notification := <-nq.queue:
			// Simulate processing a notification
			log.Printf("[Worker %d] Processing notification: %+v", workerID, notification)
			nq.processNotification(notification)
		case <-nq.shutdown:
			log.Printf("[Worker %d] Shutting down", workerID)
			return
		}
	}
}

// processNotification handles the notification (e.g., send an email or SMS).
func (nq *NotificationQueue) processNotification(notification Notification) {
	// Simulate sending a notification (e.g., email or SMS)
	fmt.Printf("Sending notification to customer %d: %s\n", notification.CustomerID, notification.Message)
	time.Sleep(1 * time.Second) // Simulate delay in processing
}

// Stop gracefully shuts down the notification queue.
func (nq *NotificationQueue) Stop() {
	close(nq.shutdown)
	nq.wg.Wait()
	close(nq.queue)
	log.Println("Notification queue stopped")
}
