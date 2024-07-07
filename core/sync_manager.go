package core

// SyncManager handles the synchronization logic
type SyncManager struct {
	// services maps service names to their respective CloudService implementations
	services map[string]CloudService
}

func NewSyncManager() *SyncManager {
	return &SyncManager{
		services: make(map[string]CloudService),
	}
}

// RegisterService adds a new cloud service to the manager
func (sm *SyncManager) RegisterService(name string, service CloudService) {
	sm.services[name] = service
}

func (sm *SyncManager) GetService(name string) (CloudService, bool) {
	service, exists := sm.services[name]
	return service, exists
}
