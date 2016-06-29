package docker

import (
	"sync"
)

// MemoryClient is a purely in-memory implementation of Client. It's meant for mocking in unit
// tests. Create one with NewMemoryClient, not by manually creating a struct instance
type MemoryClient struct {
	lck    *sync.Mutex
	Pushes []*Image
	Pulls  []*Image
	Retags []ImageTagPair
}

// NewMemoryClient creates a new MemoryClient. Always call this function instead of manually
// creating a MemoryClient struct
func NewMemoryClient() *MemoryClient {
	return &MemoryClient{lck: new(sync.Mutex)}
}

// Push is the Client interface implementation
func (m *MemoryClient) Push(i *Image) error {
	m.lck.Lock()
	defer m.lck.Unlock()
	m.Pushes = append(m.Pushes, i)
	return nil
}

// Pull is the Client interface implementation
func (m *MemoryClient) Pull(i *Image) error {
	m.lck.Lock()
	defer m.lck.Unlock()
	m.Pulls = append(m.Pulls, i)
	return nil
}

// Retag is the Client interface implementation
func (m *MemoryClient) Retag(src *Image, tar *Image) error {
	m.lck.Lock()
	defer m.lck.Unlock()
	m.Retags = append(m.Retags, ImageTagPair{Source: src, Target: tar})
	return nil
}
