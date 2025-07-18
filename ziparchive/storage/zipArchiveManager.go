package storage

import (
	"archive/zip"
	"bytes"
	"errors"
	"fmt"
	"sync"
	"time"
)

type zipArchiveManager struct {
	mu          sync.Mutex
	ZipArchives map[string]*zipArchive
	InProgress  int
}

func NewZipArchiveManager() *zipArchiveManager {
	ZAManager := make(map[string]*zipArchive)
	return &zipArchiveManager{
		ZipArchives: ZAManager,
		InProgress:  0,
	}
}

func (zam *zipArchiveManager) CreateZipArchive() *zipArchive {
	zam.mu.Lock()
	defer zam.mu.Unlock()

	id := fmt.Sprintf("%d", time.Now().UnixNano())
	za := zipArchive{
		ID:        id,
		Status:    ZAStatusCreated,
		Tasks:     make(map[string]*task),
		ZipBuffer: new(bytes.Buffer),
	}

	za.ZipWriter = zip.NewWriter(za.ZipBuffer)

	zam.ZipArchives[id] = &za

	return &za
}

func (zam *zipArchiveManager) FindTask(id string) (*task, error) {
	zam.mu.Lock()
	defer zam.mu.Unlock()

	for _, v := range zam.ZipArchives {
		for _, vv := range v.Tasks {
			if vv.ID == id {
				return vv, nil
			}
		}
	}
	return nil, errors.New("invalid task ID")
}
