package raft

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"sarva/internal/domain"

	"github.com/go-redis/redis/v8"
	"github.com/hashicorp/raft"
	raftboltdb "github.com/hashicorp/raft-boltdb"
)

type RaftAdapter struct {
	redisClient *redis.Client
	raftNode    *raft.Raft
}

func NewRaftAdapter(redisClient *redis.Client, raftNode *raft.Raft) *RaftAdapter {
	return &RaftAdapter{
		redisClient: redisClient,
		raftNode:    raftNode,
	}
}

func (r *RaftAdapter) UpdateDatabase(file domain.File) error {

	if r.raftNode.State() != raft.Leader {
		return fmt.Errorf("node is not the leader, cannot update database")
	}

	command := fmt.Sprintf("SET %s %d", file.Name, file.Size)
	future := r.raftNode.Apply([]byte(command), 10*time.Second)

	if err := future.Error(); err != nil {
		fmt.Println("Error", "Error applying command to Raft")
		return fmt.Errorf("raft consensus failed: %v", err)
	}

	ctx := context.Background()
	err := r.redisClient.Set(ctx, file.Name, file.Size, 0).Err()
	if err != nil {
		fmt.Println("Error", err.Error())
		return fmt.Errorf("failed to update Redis: %v", err)
	}

	log.Printf("File %s (size: %d) successfully stored in Redis\n", file.Name, file.Size)
	return nil
}

func SetupRaft(redisClient *redis.Client) *raft.Raft {
	raftDir := "./raft_data"
	if err := clearRaftData(raftDir); err != nil {
		log.Fatalf("Failed to clear Raft logs: %v", err)
	}
	config := raft.DefaultConfig()
	config.LocalID = raft.ServerID("node1")
	config.LogLevel = "DEBUG"

	transport, err := raft.NewTCPTransport("127.0.0.1:5000", nil, 3, 10*time.Second, os.Stderr)
	if err != nil {
		log.Fatalf("failed to create TCP transport: %v", err)
	}

	store, err := raftboltdb.NewBoltStore("raft_data/raft-log.bolt")
	if err != nil {
		log.Fatalf("failed to create Bolt store: %v", err)
	}

	snapshotStore, err := raft.NewFileSnapshotStore(".", 1, os.Stderr)
	if err != nil {
		log.Fatalf("failed to create snapshot store: %v", err)
	}
	fsm := NewFSM()
	node, err := raft.NewRaft(config, fsm, store, store, snapshotStore, transport)
	if err != nil {
		log.Fatalf("failed to create Raft node: %v", err)
	}
	configuration := raft.Configuration{
		Servers: []raft.Server{
			{
				ID:      config.LocalID,
				Address: transport.LocalAddr(),
			},
		},
	}

	node.BootstrapCluster(configuration)
	return node
}
func clearRaftData(raftDir string) error {
	err := os.RemoveAll(raftDir)
	if err != nil {
		return err
	}

	return os.MkdirAll(raftDir, 0755)
}

func ClearRedisData(redisClient *redis.Client) error {
	ctx := context.Background()
	err := redisClient.FlushDB(ctx).Err()
	if err != nil {
		return fmt.Errorf("failed to clear Redis: %v", err)
	}
	log.Println("Redis database cleared successfully.")
	return nil
}
