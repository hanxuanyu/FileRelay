package storage

import (
	"context"
	"io"
	"os"
	"path/filepath"
)

type LocalStorage struct {
	RootPath string
}

func NewLocalStorage(rootPath string) *LocalStorage {
	// 确保根目录存在
	if _, err := os.Stat(rootPath); os.IsNotExist(err) {
		os.MkdirAll(rootPath, 0755)
	}
	return &LocalStorage{RootPath: rootPath}
}

func (s *LocalStorage) Save(ctx context.Context, path string, reader io.Reader) error {
	fullPath := filepath.Join(s.RootPath, path)
	dir := filepath.Dir(fullPath)
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		os.MkdirAll(dir, 0755)
	}

	file, err := os.Create(fullPath)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = io.Copy(file, reader)
	return err
}

func (s *LocalStorage) Open(ctx context.Context, path string) (io.ReadCloser, error) {
	fullPath := filepath.Join(s.RootPath, path)
	return os.Open(fullPath)
}

func (s *LocalStorage) Delete(ctx context.Context, path string) error {
	fullPath := filepath.Join(s.RootPath, path)
	return os.Remove(fullPath)
}

func (s *LocalStorage) Exists(ctx context.Context, path string) (bool, error) {
	fullPath := filepath.Join(s.RootPath, path)
	_, err := os.Stat(fullPath)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}
