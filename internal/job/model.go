package job

type Job struct {
	ID         string
	Type       string
	Payload    []byte
	Priority   int
	Status     JobStatus
	Attempts   int
	MaxRetries int
}

type JobStatus string

const (
	StatusPending   JobStatus = "pending"
	StatusRunning   JobStatus = "running"
	StatusRetrying  JobStatus = "retrying"
	StatusCompleted JobStatus = "completed"
	StatusFailed    JobStatus = "failed"
)

type Stats struct {
	Total     int `json:"total"`
	Pending   int `json:"pending"`
	Running   int `json:"running"`
	Retrying  int `json:"retrying"`
	Completed int `json:"completed"`
	Failed    int `json:"failed"`
}
