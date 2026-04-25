package api

import (
	"WarpQueue/internal/job"
	"WarpQueue/internal/queue"
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"net/http"

	"github.com/gin-gonic/gin"
)

type Server struct {
	queue       queue.Queue
	workerCount int
}

type createJobRequest struct {
	Type       string          `json:"type" binding:"required"`
	Payload    json.RawMessage `json:"payload" binding:"required"`
	Priority   int             `json:"priority" binding:"required"`
	MaxRetries int             `json:"max_retries" binding:"required"`
}

type createJobResponse struct {
	ID string `json:"id"`
}

type jobResponse struct {
	ID       string        `json:"id"`
	Type     string        `json:"type"`
	Status   job.JobStatus `json:"status"`
	Attempts int           `json:"attempts"`
}

type statsResponse struct {
	Total     int `json:"total"`
	Pending   int `json:"pending"`
	Running   int `json:"running"`
	Retrying  int `json:"retrying"`
	Completed int `json:"completed"`
	Failed    int `json:"failed"`
	Workers   int `json:"workers"`
}

func NewServer(queue queue.Queue, workerCount int) *Server {
	return &Server{queue: queue, workerCount: workerCount}
}

func (s *Server) Routes() http.Handler {
	r := gin.Default()
	r.POST("/jobs", s.createJob)
	r.GET("/jobs/:id", s.getJob)
	r.GET("/stats", s.getStats)
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"Server": "go-warp-queue-server",
			"Status": "running",
		})
	})
	r.GET("/size", func(c *gin.Context) {
		c.JSON(200, gin.H{"size": s.queue.Size()})
	})
	return r
}

func (s *Server) createJob(c *gin.Context) {
	var req createJobRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}
	id := newJobID()
	data := job.Job{
		ID:         id,
		Type:       req.Type,
		Payload:    req.Payload,
		Priority:   req.Priority,
		MaxRetries: req.MaxRetries,
	}
	if err := s.queue.Enqueue(data); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to enqueue job"})
		return
	}
	c.JSON(http.StatusCreated, createJobResponse{ID: id})
}

func (s *Server) getJob(c *gin.Context) {
	storedJob, err := s.queue.Get(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "job not found"})
		return
	}

	c.JSON(http.StatusOK, jobResponse{
		ID:       storedJob.ID,
		Type:     storedJob.Type,
		Status:   storedJob.Status,
		Attempts: storedJob.Attempts,
	})
}

func (s *Server) getStats(c *gin.Context) {
	stats := s.queue.Stats()

	c.JSON(http.StatusOK, statsResponse{
		Total:     stats.Total,
		Pending:   stats.Pending,
		Running:   stats.Running,
		Retrying:  stats.Retrying,
		Completed: stats.Completed,
		Failed:    stats.Failed,
		Workers:   s.workerCount,
	})
}

func newJobID() string {
	b := make([]byte, 8)
	_, _ = rand.Read(b)
	return hex.EncodeToString(b)
}
