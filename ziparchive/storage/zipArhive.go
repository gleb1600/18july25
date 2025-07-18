package storage

import (
	"archive/zip"
	"bytes"
	"errors"
	"fmt"
	"io"
	"net/http"
	"sync"
	"time"
)

type zipArchiveStatus string

const (
	ZAStatusCreated               zipArchiveStatus = "ZA Created"
	ZAStatusCompletedSuccessfully zipArchiveStatus = "ZA Completed Successfully"
	//ZAStatusCompletedUnsuccessfully zipArchiveStatus = "ZA Completed Unsuccessfully"
)

type zipArchive struct {
	mu        sync.Mutex
	ID        string
	Status    zipArchiveStatus
	Tasks     map[string]*task
	URL       string
	ZipBuffer *bytes.Buffer
	ZipWriter *zip.Writer
}

func (za *zipArchive) CreateTask(r *http.Request) *task {
	za.mu.Lock()
	defer za.mu.Unlock()

	id := fmt.Sprintf("%d", time.Now().UnixNano())
	tsk := task{
		ID:     id,
		Status: TaskStatusCreated,
	}
	za.Tasks[id] = &tsk

	return &tsk
}

func (za *zipArchive) AddFileToZA(url, filename string) error {

	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return errors.New("the material on the link is unavailable")
	}

	zipEntry, err := za.ZipWriter.Create(filename)
	if err != nil {
		return err
	}
	if _, err := io.Copy(zipEntry, resp.Body); err != nil {
		return err
	}
	return nil
}
