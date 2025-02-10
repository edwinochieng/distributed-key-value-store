package hashing

import (
	"hash/crc32"
	"sort"
	"strconv"
	"sync"
)

// ConsistentHash represents the consistent hashing ring
type ConsistentHash struct {
    mu       sync.RWMutex
    nodes    map[uint32]string // Hash ring
    sorted   []uint32          // Sorted hash values
    replicas int               // Number of virtual nodes per real node
}

// NewConsistentHash initializes a consistent hash ring
func NewConsistentHash(replicas int) *ConsistentHash {
    return &ConsistentHash{
        nodes:    make(map[uint32]string),
        sorted:   []uint32{},
        replicas: replicas,
    }
}

// hashKey generates a hash value for a key
func hashKey(key string) uint32 {
    return crc32.ChecksumIEEE([]byte(key))
}

// AddNode adds a node to the hash ring
func (ch *ConsistentHash) AddNode(node string) {
    ch.mu.Lock()
    defer ch.mu.Unlock()

    for i := 0; i < ch.replicas; i++ {
        hash := hashKey(node + strconv.Itoa(i))
        ch.nodes[hash] = node
        ch.sorted = append(ch.sorted, hash)
    }
    sort.Slice(ch.sorted, func(i, j int) bool { return ch.sorted[i] < ch.sorted[j] })
}

// RemoveNode removes a node from the hash ring
func (ch *ConsistentHash) RemoveNode(node string) {
    ch.mu.Lock()
    defer ch.mu.Unlock()

    for i := 0; i < ch.replicas; i++ {
        hash := hashKey(node + strconv.Itoa(i))
        delete(ch.nodes, hash)
    }

    // Rebuild sorted list
    ch.sorted = []uint32{}
    for k := range ch.nodes {
        ch.sorted = append(ch.sorted, k)
    }
    sort.Slice(ch.sorted, func(i, j int) bool { return ch.sorted[i] < ch.sorted[j] })
}

// GetNode returns the closest node for a given key
func (ch *ConsistentHash) GetNode(key string) string {
    ch.mu.RLock()
    defer ch.mu.RUnlock()

    if len(ch.nodes) == 0 {
        return ""
    }

    hash := hashKey(key)

    // Find the closest node
    idx := sort.Search(len(ch.sorted), func(i int) bool { return ch.sorted[i] >= hash })
    if idx == len(ch.sorted) {
        idx = 0
    }

    return ch.nodes[ch.sorted[idx]]
}
