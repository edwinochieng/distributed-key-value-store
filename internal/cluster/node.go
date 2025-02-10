package cluster

import (
	"fmt"
	"sync"

	"github.com/edwinochieng/distributed-key-value-store/internal/hashing"
)

type NodeManager struct {
    mu    sync.RWMutex
    nodes map[string]*Node
    hash  *hashing.ConsistentHash
}

type Node struct {
    ID   string
    Data map[string]string
}

// NewNodeManager initializes a node manager
func NewNodeManager(replicas int) *NodeManager {
    return &NodeManager{
        nodes: make(map[string]*Node),
        hash:  hashing.NewConsistentHash(replicas),
    }
}

// Join adds a new node to the cluster
func (nm *NodeManager) Join(nodeID string) {
    nm.mu.Lock()
    defer nm.mu.Unlock()

    if _, exists := nm.nodes[nodeID]; exists {
        fmt.Println("Node already exists:", nodeID)
        return
    }

    nm.nodes[nodeID] = &Node{ID: nodeID, Data: make(map[string]string)}
    nm.hash.AddNode(nodeID)
    fmt.Println("Node joined:", nodeID)
}

// Leave removes a node from the cluster
func (nm *NodeManager) Leave(nodeID string) {
    nm.mu.Lock()
    defer nm.mu.Unlock()

    if _, exists := nm.nodes[nodeID]; !exists {
        fmt.Println("Node does not exist:", nodeID)
        return
    }

    delete(nm.nodes, nodeID)
    nm.hash.RemoveNode(nodeID)
    fmt.Println("Node left:", nodeID)
}

// GetNodeForKey finds the responsible node for a given key
func (nm *NodeManager) GetNodeForKey(key string) string {
    return nm.hash.GetNode(key)
}
