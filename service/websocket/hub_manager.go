package websocket

import (
	"ChiragKr04/go-backend/types"
	"log"
	"sync"
)

type HubManager struct {
	hubs map[string]*HubFile
	mu   sync.RWMutex
}

func NewHubManager() *HubManager {
	return &HubManager{
		hubs: make(map[string]*HubFile),
	}
}

func newHub() *HubFile {
	return &HubFile{
		HubType: &types.Hub{
			Broadcast:  make(chan []byte),
			Register:   make(chan *types.Client),
			Unregister: make(chan *types.Client),
			Clients:    make(map[*types.Client]bool),
		},
	}
}

func (m *HubManager) GetHub(roomId string) *HubFile {
	m.mu.RLock()
	hub, exists := m.hubs[roomId]
	m.mu.RUnlock()

	if exists {
		return hub
	}

	m.mu.Lock()
	defer m.mu.Unlock()

	hub, exists = m.hubs[roomId]
	if exists {
		return hub
	}

	hub = newHub()
	m.hubs[roomId] = hub
	go hub.Run()
	log.Printf("Created new hub for room: %s", roomId)
	return hub
}

func (m *HubManager) RemoveHub(roomId string) {
	m.mu.Lock()
	defer m.mu.Unlock()
	delete(m.hubs, roomId)
	log.Printf("Removed hub for room: %s", roomId)
}
