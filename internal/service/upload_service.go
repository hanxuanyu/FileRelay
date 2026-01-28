package service

import (
	"FileRelay/internal/bootstrap"
	"FileRelay/internal/config"
	"FileRelay/internal/model"
	"FileRelay/internal/storage"
	"context"
	"fmt"
	"log/slog"
	"mime/multipart"
	"path/filepath"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type UploadService struct {
	db *gorm.DB
}

func NewUploadService() *UploadService {
	return &UploadService{db: bootstrap.DB}
}

func (s *UploadService) CreateBatch(ctx context.Context, files []*multipart.FileHeader, remark string, expireType string, expireValue interface{}) (*model.FileBatch, error) {
	// 1. 生成取件码
	batchService := NewBatchService()
	pickupCode, err := batchService.GeneratePickupCode(config.GlobalConfig.Security.PickupCodeLength)
	if err != nil {
		return nil, err
	}

	// 2. 准备 Batch
	batch := &model.FileBatch{
		ID:         uuid.New().String(),
		PickupCode: pickupCode,
		Remark:     remark,
		ExpireType: expireType,
		Status:     "active",
		Type:       "file",
	}

	s.applyExpire(batch, expireType, expireValue)

	// 3. 处理文件上传
	err = s.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(batch).Error; err != nil {
			return err
		}

		for _, fileHeader := range files {
			fileItem, err := s.processFile(ctx, tx, batch.ID, fileHeader)
			if err != nil {
				return err
			}
			batch.FileItems = append(batch.FileItems, *fileItem)
		}
		return nil
	})

	if err == nil {
		slog.Info("File batch created", "batch_id", batch.ID, "pickup_code", batch.PickupCode, "files_count", len(files))
	} else {
		slog.Error("Failed to create file batch", "error", err)
	}

	return batch, err
}

func (s *UploadService) CreateTextBatch(ctx context.Context, content string, remark string, expireType string, expireValue interface{}) (*model.FileBatch, error) {
	// 1. 生成取件码
	batchService := NewBatchService()
	pickupCode, err := batchService.GeneratePickupCode(config.GlobalConfig.Security.PickupCodeLength)
	if err != nil {
		return nil, err
	}

	// 2. 准备 Batch
	batch := &model.FileBatch{
		ID:         uuid.New().String(),
		PickupCode: pickupCode,
		Remark:     remark,
		ExpireType: expireType,
		Status:     "active",
		Type:       "text",
		Content:    content,
	}

	s.applyExpire(batch, expireType, expireValue)

	// 3. 保存
	if err := s.db.Create(batch).Error; err != nil {
		slog.Error("Failed to create text batch", "error", err)
		return nil, err
	}

	slog.Info("Text batch created", "batch_id", batch.ID, "pickup_code", batch.PickupCode)
	return batch, nil
}

func (s *UploadService) applyExpire(batch *model.FileBatch, expireType string, expireValue interface{}) {
	switch expireType {
	case "time":
		if days, ok := expireValue.(int); ok {
			expireAt := time.Now().Add(time.Duration(days) * 24 * time.Hour)
			batch.ExpireAt = &expireAt
		}
	case "download":
		if max, ok := expireValue.(int); ok {
			batch.MaxDownloads = max
		}
	}
}

func (s *UploadService) processFile(ctx context.Context, tx *gorm.DB, batchID string, fileHeader *multipart.FileHeader) (*model.FileItem, error) {
	file, err := fileHeader.Open()
	if err != nil {
		return nil, err
	}
	defer file.Close()

	// 生成唯一存储路径
	ext := filepath.Ext(fileHeader.Filename)
	fileID := uuid.New().String()
	storagePath := fmt.Sprintf("%s/%s%s", batchID, fileID, ext)

	// 保存到存储层
	if err := storage.GlobalStorage.Save(ctx, storagePath, file); err != nil {
		return nil, err
	}

	// 创建数据库记录
	item := &model.FileItem{
		ID:           fileID,
		BatchID:      batchID,
		OriginalName: fileHeader.Filename,
		StoragePath:  storagePath,
		Size:         fileHeader.Size,
		MimeType:     fileHeader.Header.Get("Content-Type"),
	}

	if err := tx.Create(item).Error; err != nil {
		return nil, err
	}

	return item, nil
}
