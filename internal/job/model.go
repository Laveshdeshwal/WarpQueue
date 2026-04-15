package job

type Job struct {
	ID       string
	Type     string
	Payload  []byte
	Priority int
	Status   Status
	Attempts int
}

type Status string

const (
	StatusPending   Status = "pending"
	StatusRunning   Status = "running"
	StatusCompleted Status = "completed"
	StatusFailed    Status = "failed"
)
