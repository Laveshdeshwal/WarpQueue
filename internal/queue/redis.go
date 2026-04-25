package queue

import (
	"WarpQueue/internal/job"
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/redis/go-redis/v9"
)

type RedisQueue struct {
	client     *redis.Client
	listKey    string
	jobKeyBase string
}

func NewRedisQueue(addr, password string, db int, keyPrefix string) (*RedisQueue, error) {
	client := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: password,
		DB:       db,
	})

	if err := client.Ping(context.Background()).Err(); err != nil {
		return nil, err
	}

	prefix := strings.TrimSpace(keyPrefix)
	if prefix == "" {
		prefix = "warpqueue"
	}

	return &RedisQueue{
		client:     client,
		listKey:    prefix + ":queue",
		jobKeyBase: prefix + ":job:",
	}, nil
}

func (q *RedisQueue) Save(newJob job.Job) error {
	return q.saveJob(context.Background(), newJob)
}

func (q *RedisQueue) Get(id string) (job.Job, error) {
	return q.getJob(context.Background(), id)
}

func (q *RedisQueue) Update(updatedJob job.Job) error {
	return q.saveJob(context.Background(), updatedJob)
}

func (q *RedisQueue) ListByStatus(status job.JobStatus) []job.Job {
	ctx := context.Background()
	keys, err := q.scanJobKeys(ctx)
	if err != nil {
		return []job.Job{}
	}

	jobs := make([]job.Job, 0)
	for _, key := range keys {
		payload, err := q.client.Get(ctx, key).Result()
		if err != nil {
			continue
		}

		var storedJob job.Job
		if err := json.Unmarshal([]byte(payload), &storedJob); err != nil {
			continue
		}

		if storedJob.Status == status {
			jobs = append(jobs, storedJob)
		}
	}

	return jobs
}

func (q *RedisQueue) Stats() job.Stats {
	ctx := context.Background()
	keys, err := q.scanJobKeys(ctx)
	if err != nil {
		return job.Stats{}
	}

	stats := job.Stats{
		Total: len(keys),
	}

	for _, key := range keys {
		payload, err := q.client.Get(ctx, key).Result()
		if err != nil {
			continue
		}

		var storedJob job.Job
		if err := json.Unmarshal([]byte(payload), &storedJob); err != nil {
			continue
		}

		switch storedJob.Status {
		case job.StatusPending:
			stats.Pending++
		case job.StatusRunning:
			stats.Running++
		case job.StatusRetrying:
			stats.Retrying++
		case job.StatusCompleted:
			stats.Completed++
		case job.StatusFailed:
			stats.Failed++
		}
	}

	return stats
}

func (q *RedisQueue) Enqueue(newJob job.Job) error {
	ctx := context.Background()
	newJob.Status = job.StatusPending
	if err := q.saveJob(ctx, newJob); err != nil {
		return err
	}

	return q.client.RPush(ctx, q.listKey, newJob.ID).Err()
}

func (q *RedisQueue) Dequeue() (job.Job, error) {
	ctx := context.Background()
	jobID, err := q.client.LPop(ctx, q.listKey).Result()
	if err != nil {
		if err == redis.Nil {
			return job.Job{}, ErrQueueEmpty
		}
		return job.Job{}, err
	}

	storedJob, err := q.getJob(ctx, jobID)
	if err != nil {
		return job.Job{}, err
	}

	storedJob.Attempts++
	storedJob.Status = job.StatusRunning
	if err := q.saveJob(ctx, storedJob); err != nil {
		return job.Job{}, err
	}

	return storedJob, nil
}

func (q *RedisQueue) Size() int {
	size, err := q.client.LLen(context.Background(), q.listKey).Result()
	if err != nil {
		return 0
	}

	return int(size)
}

func (q *RedisQueue) saveJob(ctx context.Context, j job.Job) error {
	payload, err := json.Marshal(j)
	if err != nil {
		return err
	}

	return q.client.Set(ctx, q.jobKey(j.ID), payload, 0).Err()
}

func (q *RedisQueue) getJob(ctx context.Context, id string) (job.Job, error) {
	payload, err := q.client.Get(ctx, q.jobKey(id)).Result()
	if err != nil {
		if err == redis.Nil {
			return job.Job{}, fmt.Errorf("job %s not found", id)
		}
		return job.Job{}, err
	}

	var storedJob job.Job
	if err := json.Unmarshal([]byte(payload), &storedJob); err != nil {
		return job.Job{}, err
	}

	return storedJob, nil
}

func (q *RedisQueue) scanJobKeys(ctx context.Context) ([]string, error) {
	keys := make([]string, 0)
	iter := q.client.Scan(ctx, 0, q.jobKeyBase+"*", 0).Iterator()
	for iter.Next(ctx) {
		keys = append(keys, iter.Val())
	}

	if err := iter.Err(); err != nil {
		return nil, err
	}

	return keys, nil
}

func (q *RedisQueue) jobKey(id string) string {
	return q.jobKeyBase + id
}
