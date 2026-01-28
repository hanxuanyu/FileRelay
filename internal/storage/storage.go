package storage

import (
	"context"
	"io"
)

type Storage interface {
	Save(ctx context.Context, path string, reader io.Reader) error
	Open(ctx context.Context, path string) (io.ReadCloser, error)
	Delete(ctx context.Context, path string) error
	Exists(ctx context.Context, path string) (bool, error)
}

var GlobalStorage Storage
