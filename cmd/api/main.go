package main

import (
	"context"
	"fmt"
	"sarva/internal/adaptor/logger"
	"sarva/internal/adaptor/raft"
	"sarva/internal/adaptor/redis"
	"sarva/internal/port"
	"sarva/internal/service"

	"sarva/internal/adaptor/fileprocessor"

	"github.com/gin-gonic/gin"
	r "github.com/go-redis/redis/v8"
)

func main() {
	rdb := r.NewClient(&r.Options{
		Addr: "localhost:6379",
	})
	ctx := context.Background()
	err := raft.ClearRedisData(rdb)
	if err != nil {
		fmt.Printf("Failed to clear Redis on startup: %v\n", err)
		return
	}

	raftNode := raft.SetupRaft(rdb)
	consensus := raft.NewRaftAdapter(rdb, raftNode)
	logger := logger.NewHclogAdapter()

	processor := &fileprocessor.DummyFileProcessor{}
	repository := redis.NewRedisAdapter(rdb, ctx)
	fileService := service.NewFileService(processor, repository, consensus, logger)

	r := gin.Default()
	handler := port.NewHandler(fileService)
	r.POST("/upload", handler.UploadFile)
	r.Run(":8080")
}
