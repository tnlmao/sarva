package raft

import (
	"encoding/json"
	"fmt"
	"io"
	"strconv"
	"strings"
	"sync"

	"github.com/hashicorp/raft"
)

type FSM struct {
	mu   sync.RWMutex
	data map[string]int64
}

func NewFSM() *FSM {
	return &FSM{
		data: make(map[string]int64),
	}
}

func (f *FSM) Apply(log *raft.Log) interface{} {
	fmt.Println("------------in fsm-----------------")
	command := string(log.Data)
	parts := parseCommand(command)
	i := 0
	switch parts[0] {
	case "SET":
		if len(parts) < 3 {
			return "invalid SET command"
		}
		name := parts[1]
		size := parts[2]
		sint, _ := strconv.Atoi(size)
		f.mu.Lock()
		f.data[name] = int64(sint)
		fmt.Println(i)
		f.mu.Unlock()
		fmt.Printf("Data applied: %s = %s\n", name, size)
		i++
		return nil

	default:
		return "unknown command"
	}
}

func (f *FSM) Snapshot() (raft.FSMSnapshot, error) {
	f.mu.RLock()
	defer f.mu.RUnlock()

	stateCopy := make(map[string]int64)
	for k, v := range f.data {
		stateCopy[k] = v
	}

	return &Snapshot{state: stateCopy}, nil
}

func (f *FSM) Restore(rc io.ReadCloser) error {
	var restoredData map[string]int64
	if err := json.NewDecoder(rc).Decode(&restoredData); err != nil {
		return err
	}

	f.mu.Lock()
	defer f.mu.Unlock()

	f.data = restoredData
	return nil
}

type Snapshot struct {
	state map[string]int64
}

func (s *Snapshot) Persist(sink raft.SnapshotSink) error {
	err := func() error {
		if err := json.NewEncoder(sink).Encode(s.state); err != nil {
			return err
		}
		return sink.Close()
	}()
	if err != nil {
		sink.Cancel()
		return err
	}
	return nil
}

func (s *Snapshot) Release() {
}

func parseCommand(command string) []string {
	return strings.Fields(command)
}
