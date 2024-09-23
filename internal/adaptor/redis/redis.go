package redis

import (
	"context"
	"encoding/json"
	"sarva/internal/domain"

	"github.com/go-redis/redis/v8"
)

type RedisAdapter struct {
	client *redis.Client
	ctx    context.Context
}
type FileMetadata struct {
	Name string `json:"name"`
	Size int64  `json:"size"`
}

func NewRedisAdapter(client *redis.Client, ctx context.Context) *RedisAdapter {
	return &RedisAdapter{client: client, ctx: ctx}
}

func (r *RedisAdapter) SaveFile(file domain.File) error {
	metadata := FileMetadata{
		Name: file.Name,
		Size: file.Size,
	}

	jsonData, err := json.Marshal(metadata)
	if err != nil {
		return err
	}

	err = r.client.Set(r.ctx, file.Name, jsonData, 0).Err()
	if err != nil {
		return err
	}

	return nil
}
