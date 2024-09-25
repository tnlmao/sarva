package raft

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os"
	"time"

	"sarva/internal/domain"

	"github.com/go-redis/redis/v8"
	"github.com/hashicorp/raft"
	raftboltdb "github.com/hashicorp/raft-boltdb"
)

type RaftNode struct {
	ID   raft.ServerID
	Addr raft.ServerAddress
	Raft *raft.Raft
}
type RaftAdapter struct {
	redisClient *redis.Client
	raftNode    map[raft.ServerID]*RaftNode
}

func NewRaftAdapter(redisClient *redis.Client, raftNodes map[raft.ServerID]*RaftNode) *RaftAdapter {
	return &RaftAdapter{
		redisClient: redisClient,
		raftNode:    raftNodes,
	}
}

func (r *RaftAdapter) UpdateDatabase(file domain.File) error {

	leader := GetLeader(r.raftNode)
	fmt.Println("leader", leader)
	if leader.State() != raft.Leader {
		log.Println("This  is not the leader.")
		return errors.New("this is not the leader node")
	}
	command := fmt.Sprintf("SET %s %d", file.Name, file.Size)
	future := leader.Apply([]byte(command), 10*time.Second)
	if err := future.Error(); err != nil {
		return fmt.Errorf("raft consensus failed: %v", err)
	}
	log.Printf("File %s (size: %d) successfully stored in Redis\n", file.Name, file.Size)
	return nil
}

func SetupRaftCluster(redisClient *redis.Client, nodeMap map[raft.ServerID]raft.ServerAddress) map[raft.ServerID]*RaftNode {
	raftDir := "./raft_data"
	nodes := make(map[raft.ServerID]*RaftNode)
	if err := clearRaftData(raftDir); err != nil {
		log.Fatalf("Failed to clear Raft logs: %v", err)
	}

	for id, addr := range nodeMap {
		nodeDir := fmt.Sprintf("%s/%s", raftDir, id)
		if err := os.MkdirAll(nodeDir, 0755); err != nil {
			log.Fatalf("Failed to create directory for node %s: %v", id, err)
		}
		config := raft.DefaultConfig()
		config.LocalID = id
		config.LogLevel = "DEBUG"

		transport, err := raft.NewTCPTransport(string(addr), nil, 3, 10*time.Second, os.Stderr)
		if err != nil {
			log.Fatalf("Failed to create TCP transport for node %s: %v", id, err)
		}

		store, err := raftboltdb.NewBoltStore(fmt.Sprintf("%s/raft-log.bolt", nodeDir))
		if err != nil {
			log.Fatalf("Failed to create Bolt store for node %s: %v", id, err)
		}

		snapshotStore, err := raft.NewFileSnapshotStore(nodeDir, 1, os.Stderr)
		if err != nil {
			log.Fatalf("Failed to create snapshot store for node %s: %v", id, err)
		}

		fsm := NewFSM()

		raftNode, err := raft.NewRaft(config, fsm, store, store, snapshotStore, transport)
		if err != nil {
			log.Fatalf("Failed to create Raft node %s: %v", id, err)
		}

		nodes[id] = &RaftNode{
			ID:   id,
			Addr: addr,
			Raft: raftNode,
		}
	}

	configuration := raft.Configuration{
		Servers: make([]raft.Server, 0, len(nodeMap)),
	}
	for id, addr := range nodeMap {
		configuration.Servers = append(configuration.Servers, raft.Server{
			ID:      id,
			Address: addr,
		})
	}
	for _, node := range nodes {
		if err := node.Raft.BootstrapCluster(configuration).Error(); err != nil {
			return nil
		}
	}
	return nodes
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
func GetLeader(nodes map[raft.ServerID]*RaftNode) *raft.Raft {
	var leaderNode *RaftNode
	for _, node := range nodes {
		leaderID := node.Raft.Leader()
		fmt.Println("Current Node ID:", node.ID)
		fmt.Println("Leader ID:", leaderID)

		if leaderID == node.Addr {
			leaderNode = node
			break
		}
	}

	if leaderNode != nil {
		return leaderNode.Raft
	}

	fmt.Println("No leader found.")
	return nil
}
