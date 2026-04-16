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
	queue queue.Queue
}

type createJobRequest struct {
	Type     string          `json:"type" binding:"required"`
	Payload  json.RawMessage `json:"payload" binding:"required"`
	Priority int             `json:"priority" binding:"required"`
	Attempts int             `json:"attempts" binding:"required"`
}

type createJobResponse struct {
	ID string `json:"id"`
}

func NewServer(queue queue.Queue) *Server { return &Server{queue: queue} }

func (s *Server) Routes() http.Handler {
	r := gin.Default()
	r.POST("/jobs", s.createJob)
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
	data := job.Job{ID: id, Type: req.Type, Payload: req.Payload, Priority: req.Priority, Attempts: req.Attempts}
	if err := s.queue.Enqueue(data); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to enqueue job"})
		return
	}
	c.JSON(http.StatusCreated, createJobResponse{ID: id})
}

func newJobID() string {
	b := make([]byte, 8)
	_, _ = rand.Read(b)
	return hex.EncodeToString(b)
}
