package tests

import (
	"WarpQueue/internal/api"
	"WarpQueue/internal/job"
	"WarpQueue/internal/queue"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

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

func TestGetJobByID(t *testing.T) {
	q := queue.NewMemoryQueue()
	server := api.NewServer(q, 3)

	body := `{"type":"send_email","payload":{"to":"user@example.com"},"priority":1,"max_retries":2}`
	req := httptest.NewRequest(http.MethodPost, "/jobs", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	server.Routes().ServeHTTP(rec, req)

	if rec.Code != http.StatusCreated {
		t.Fatalf("expected status %d, got %d", http.StatusCreated, rec.Code)
	}

	var createResp createJobResponse
	if err := json.Unmarshal(rec.Body.Bytes(), &createResp); err != nil {
		t.Fatalf("failed to decode create job response: %v", err)
	}

	getReq := httptest.NewRequest(http.MethodGet, "/jobs/"+createResp.ID, nil)
	getRec := httptest.NewRecorder()
	server.Routes().ServeHTTP(getRec, getReq)

	if getRec.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, getRec.Code)
	}

	var resp jobResponse
	if err := json.Unmarshal(getRec.Body.Bytes(), &resp); err != nil {
		t.Fatalf("failed to decode get job response: %v", err)
	}

	if resp.ID != createResp.ID {
		t.Fatalf("expected job ID %s, got %s", createResp.ID, resp.ID)
	}

	if resp.Status != job.StatusPending {
		t.Fatalf("expected status %s, got %s", job.StatusPending, resp.Status)
	}

	if resp.Attempts != 0 {
		t.Fatalf("expected attempts 0, got %d", resp.Attempts)
	}
}

func TestGetStats(t *testing.T) {
	q := queue.NewMemoryQueue()

	if err := q.Save(job.Job{ID: "pending", Type: "send_email", Status: job.StatusPending}); err != nil {
		t.Fatalf("save pending job: %v", err)
	}
	if err := q.Save(job.Job{ID: "running", Type: "send_email", Status: job.StatusRunning}); err != nil {
		t.Fatalf("save running job: %v", err)
	}
	if err := q.Save(job.Job{ID: "completed", Type: "send_email", Status: job.StatusCompleted}); err != nil {
		t.Fatalf("save completed job: %v", err)
	}
	if err := q.Save(job.Job{ID: "failed", Type: "send_email", Status: job.StatusFailed}); err != nil {
		t.Fatalf("save failed job: %v", err)
	}

	server := api.NewServer(q, 4)
	req := httptest.NewRequest(http.MethodGet, "/stats", nil)
	rec := httptest.NewRecorder()
	server.Routes().ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, rec.Code)
	}

	var resp statsResponse
	if err := json.Unmarshal(rec.Body.Bytes(), &resp); err != nil {
		t.Fatalf("failed to decode stats response: %v", err)
	}

	if resp.Total != 4 {
		t.Fatalf("expected total 4, got %d", resp.Total)
	}
	if resp.Pending != 1 || resp.Running != 1 || resp.Completed != 1 || resp.Failed != 1 {
		t.Fatalf("unexpected stats counts: %+v", resp)
	}
	if resp.Workers != 4 {
		t.Fatalf("expected workers 4, got %d", resp.Workers)
	}
}
