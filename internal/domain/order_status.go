package domain

const (
	OrderStatusPendingMatch       = "pending_match"
	OrderStatusOffered          = "offered"
	OrderStatusAccepted         = "accepted"
	OrderStatusWorkerDeparted   = "worker_departed"
	OrderStatusInProgress       = "in_progress"
	OrderStatusCompleted        = "completed"
	OrderStatusCancelledConsumer = "cancelled_by_consumer"
	OrderStatusCancelledWorker  = "cancelled_by_worker"
	OrderStatusExpired          = "expired"
	OrderStatusDispute          = "dispute"
)

const (
	PaymentPending    = "pending"
	PaymentAuthorized = "authorized"
	PaymentCaptured   = "captured"
	PaymentConfirmed  = "confirmed"
	PaymentRefunded   = "refunded"
	PaymentFailed     = "failed"
)

const (
	WorkerAvailabilityOnline  = "online"
	WorkerAvailabilityOffline = "offline"
	WorkerAvailabilityBusy    = "busy"
)
