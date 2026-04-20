package job

type Job struct {
	ID       string
	Type     string
	Payload  []byte
	Priority int
	Status   JobStatus
	Attempts int
}

type JobStatus string

const (
	StatusPending   JobStatus = "pending"
	StatusRunning   JobStatus = "running"
	StatusCompleted JobStatus = "completed"
	StatusFailed    JobStatus = "failed"
)
