package main

import (
	"context"
	"fmt"
	"sarva/internal/adaptor/fileprocessor"
	"sarva/internal/adaptor/logger"
	"sarva/internal/adaptor/raft"
	"sarva/internal/adaptor/redis"
	"sarva/internal/port"
	"sarva/internal/service"

	"github.com/gin-gonic/gin"
	r "github.com/go-redis/redis/v8"
	n "github.com/hashicorp/raft"
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

	nodeMap := make(map[n.ServerID]n.ServerAddress)
	nodeMap[n.ServerID("node1")] = n.ServerAddress("127.0.0.1:5001")
	nodeMap[n.ServerID("node2")] = n.ServerAddress("127.0.0.1:5002")
	nodeMap[n.ServerID("node3")] = n.ServerAddress("127.0.0.1:5003")

	raftNode := raft.SetupRaftCluster(rdb, nodeMap)
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
