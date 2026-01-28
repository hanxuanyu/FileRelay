package task

import (
	"FileRelay/internal/service"
	"context"
	"log/slog"
	"time"
)

type Cleaner struct {
	batchService *service.BatchService
}

func NewCleaner() *Cleaner {
	return &Cleaner{
		batchService: service.NewBatchService(),
	}
}

func (c *Cleaner) Start(ctx context.Context) {
	ticker := time.NewTicker(1 * time.Hour)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			c.Clean()
		}
	}
}

func (c *Cleaner) Clean() {
	slog.Info("Running cleanup task")
	if err := c.batchService.Cleanup(context.Background()); err != nil {
		slog.Error("Error during cleanup", "error", err)
	}
}
