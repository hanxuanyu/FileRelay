package storage

import (
	"context"
	"io"
	"path/filepath"
	"strings"

	"github.com/studio-b12/gowebdav"
)

type WebDAVStorage struct {
	client *gowebdav.Client
	root   string
}

func NewWebDAVStorage(url, username, password, root string) *WebDAVStorage {
	client := gowebdav.NewClient(url, username, password)
	return &WebDAVStorage{
		client: client,
		root:   root,
	}
}

func (s *WebDAVStorage) getFullPath(path string) string {
	return filepath.ToSlash(filepath.Join(s.root, path))
}

func (s *WebDAVStorage) Save(ctx context.Context, path string, reader io.Reader) error {
	fullPath := s.getFullPath(path)

	// Ensure directory exists
	dir := filepath.Dir(fullPath)
	if dir != "." && dir != "/" {
		parts := strings.Split(strings.Trim(dir, "/"), "/")
		current := ""
		for _, part := range parts {
			current += "/" + part
			_ = s.client.Mkdir(current, 0755)
		}
	}

	return s.client.WriteStream(fullPath, reader, 0644)
}

func (s *WebDAVStorage) Open(ctx context.Context, path string) (io.ReadCloser, error) {
	fullPath := s.getFullPath(path)
	return s.client.ReadStream(fullPath)
}

func (s *WebDAVStorage) Delete(ctx context.Context, path string) error {
	fullPath := s.getFullPath(path)
	return s.client.Remove(fullPath)
}

func (s *WebDAVStorage) Exists(ctx context.Context, path string) (bool, error) {
	fullPath := s.getFullPath(path)
	_, err := s.client.Stat(fullPath)
	if err != nil {
		// gowebdav's Stat returns error if not found
		// We could check for 404 but gowebdav doesn't export error types easily
		// Usually we check if it's a 404
		if strings.Contains(err.Error(), "404") || strings.Contains(err.Error(), "not found") {
			return false, nil
		}
		return false, err
	}
	return true, nil
}
