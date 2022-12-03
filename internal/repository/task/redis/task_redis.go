package redis_repo

import (
	"bytes"
	"context"
	"encoding/gob"
	"fmt"
	"time"

	"github.com/fr13n8/go-practice/internal/domain"
	"github.com/go-redis/redis"
	"github.com/opentracing/opentracing-go"
)

const (
	taskKey = "task"
)

type Repo struct {
	redis  *redis.Client
	prefix string
}

func NewTask(redis *redis.Client) *Repo {
	return &Repo{redis, taskKey}
}

func (e *Repo) Set(ctx context.Context, task domain.Task, expire int) error {
	span, _ := opentracing.StartSpanFromContext(ctx, "task.redis.Set")
	defer span.Finish()

	var b bytes.Buffer
	if err := gob.NewEncoder(&b).Encode(task); err != nil {
		return err
	}

	key := e.createKey(task.ID)
	return e.redis.Set(key, b.Bytes(), time.Second*time.Duration(expire)).Err()
}

func (e *Repo) Delete(ctx context.Context, id string) error {
	span, _ := opentracing.StartSpanFromContext(ctx, "task.redis.Delete")
	defer span.Finish()

	if err := e.redis.HDel(e.prefix, id).Err(); err != nil {
		return err
	}

	return nil
}

func (e *Repo) Get(ctx context.Context, id string) (domain.Task, error) {
	span, _ := opentracing.StartSpanFromContext(ctx, "task.redis.Get")
	defer span.Finish()

	var task domain.Task
	key := e.createKey(id)
	cmd := e.redis.Get(key)
	cmdb, err := cmd.Bytes()
	if err != nil {
		return domain.Task{}, err
	}
	b := bytes.NewReader(cmdb)
	if err := gob.NewDecoder(b).Decode(&task); err != nil {
		return domain.Task{}, err
	}
	return task, nil
}

func (e *Repo) Created(ctx context.Context, task domain.Task) error {
	span, _ := opentracing.StartSpanFromContext(ctx, "task.redis.PublishCreated")
	defer span.Finish()

	return e.publish(task, "created")
}

func (e *Repo) publish(task domain.Task, event string) error {
	var b bytes.Buffer
	if err := gob.NewEncoder(&b).Encode(task); err != nil {
		return err
	}

	return e.redis.Publish(fmt.Sprintf("%s.%s", e.prefix, event), b.Bytes()).Err()
}

func (e *Repo) createKey(id string) string {
	return fmt.Sprintf("%s:%s", e.prefix, id)
}
