package database

import (
	"FilesToTarBot/predictor"
	"archive/tar"
	"sync"
	"time"
)

// MemoryCache is an in memory cache to store files which want to be downloaded temporary
// The cache is deleted every hour
type MemoryCache struct {
	m  map[int64]memoryCacheValue
	mu sync.RWMutex
}

// memoryCacheValue contains a file and the time which it has been inserted
type memoryCacheValue struct {
	// The files to hold
	files []File
	// Predict the result tar size
	tarPredictor predictor.Tar
	// When was this file modification in unix epoch
	lastModification int64
}

// NewMemoryCache creates a new memory cache and setups a cleanup goroutine
func NewMemoryCache() *MemoryCache {
	m := &MemoryCache{m: make(map[int64]memoryCacheValue)}
	go m.cleanupGoroutine()
	return m
}

// cleanupGoroutine deletes old files from memory
func (m *MemoryCache) cleanupGoroutine() {
	for {
		time.Sleep(time.Hour)
		m.mu.Lock()
		start := time.Now().Unix()
		for k, v := range m.m {
			// Delete entries older than a day
			if start-v.lastModification > 3600*24 {
				delete(m.m, k)
			}
		}
		m.mu.Unlock()
	}
}

func (m *MemoryCache) AddFile(userID int64, file File) error {
	m.mu.Lock()
	files := m.m[userID]
	// Test the size (note that files is just a copy)
	files.tarPredictor.AddFile(tar.Header{
		Name: file.Name,
		Size: file.Size,
	})
	if files.tarPredictor.Total() >= MaxFileSize {
		m.mu.Unlock()
		return TooBigFileError
	}
	// Check other files
	for i := range files.files {
		if files.files[i].Name == file.Name {
			m.mu.Unlock()
			return FileAlreadyExistsError
		}
	}
	// Add the file
	files.files = append(files.files, file)
	files.lastModification = time.Now().Unix()
	m.m[userID] = files
	m.mu.Unlock()
	return nil
}

func (m *MemoryCache) Reset(userID int64) {
	m.mu.Lock()
	delete(m.m, userID)
	m.mu.Unlock()
}

func (m *MemoryCache) GetFiles(userID int64) ([]File, int64) {
	m.mu.RLock()
	files := m.m[userID]
	m.mu.RUnlock()
	return files.files, files.tarPredictor.Total()
}
