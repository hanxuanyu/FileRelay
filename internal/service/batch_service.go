package service

import (
	"FileRelay/internal/bootstrap"
	"FileRelay/internal/model"
	"FileRelay/internal/storage"
	"context"
	"errors"
	"log/slog"
	"math/big"
	"strings"
	"time"

	"crypto/rand"

	"gorm.io/gorm"
)

type BatchService struct {
	db *gorm.DB
}

func NewBatchService() *BatchService {
	return &BatchService{db: bootstrap.DB}
}

func (s *BatchService) GetBatchByPickupCode(code string) (*model.FileBatch, error) {
	var batch model.FileBatch
	err := s.db.Preload("FileItems").Where("pickup_code = ? AND status = ?", code, "active").First(&batch).Error
	if err != nil {
		return nil, err
	}

	// 检查是否过期
	if s.IsExpired(&batch) {
		s.MarkAsExpired(&batch)
		return nil, errors.New("batch expired")
	}

	return &batch, nil
}

func (s *BatchService) GetDownloadCountByPickupCode(code string) (int, int, error) {
	var batch model.FileBatch
	// 查询活跃或已过期的批次
	err := s.db.Where("pickup_code = ? AND (status = ? OR status = ?)", code, "active", "expired").First(&batch).Error
	if err != nil {
		return 0, 0, err
	}
	return batch.DownloadCount, batch.MaxDownloads, nil
}

func (s *BatchService) IsExpired(batch *model.FileBatch) bool {
	if batch.Status != "active" {
		return true
	}

	switch batch.ExpireType {
	case "time":
		if batch.ExpireAt != nil && time.Now().After(*batch.ExpireAt) {
			return true
		}
	case "download":
		if batch.MaxDownloads > 0 && batch.DownloadCount >= batch.MaxDownloads {
			return true
		}
	}
	return false
}

func (s *BatchService) MarkAsExpired(batch *model.FileBatch) error {
	slog.Info("Marking batch as expired", "batch_id", batch.ID, "pickup_code", batch.PickupCode)
	return s.db.Model(batch).Update("status", "expired").Error
}

func (s *BatchService) DeleteBatch(ctx context.Context, batchID string) error {
	var batch model.FileBatch
	// 使用 Unscoped 以确保即使是已软删除的批次也能找到并清理其物理文件
	if err := s.db.Unscoped().Preload("FileItems").First(&batch, "id = ?", batchID).Error; err != nil {
		return err
	}

	slog.Info("Deleting batch", "batch_id", batch.ID, "files_count", len(batch.FileItems))

	// 删除物理文件
	for _, item := range batch.FileItems {
		if err := storage.GlobalStorage.Delete(ctx, item.StoragePath); err != nil {
			slog.Error("Failed to delete physical file", "path", item.StoragePath, "error", err)
		}
	}

	// 删除数据库记录 (彻底删除，不再保留元数据以便清理任务不再扫描到它)
	return s.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Where("batch_id = ?", batch.ID).Delete(&model.FileItem{}).Error; err != nil {
			return err
		}
		if err := tx.Unscoped().Delete(&batch).Error; err != nil {
			return err
		}
		return nil
	})
}

func (s *BatchService) IncrementDownloadCount(batchID string) error {
	if batchID == "" {
		return errors.New("batch id is empty")
	}
	result := s.db.Model(&model.FileBatch{}).Where("id = ?", batchID).
		UpdateColumn("download_count", gorm.Expr("download_count + ?", 1))
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return errors.New("batch not found or already deleted")
	}
	return nil
}

func (s *BatchService) GeneratePickupCode(length int) (string, error) {
	const charset = "0123456789"
	b := make([]byte, length)
	for i := range b {
		num, err := rand.Int(rand.Reader, big.NewInt(int64(len(charset))))
		if err != nil {
			return "", err
		}
		b[i] = charset[num.Int64()]
	}
	// 检查是否冲突 (排除已删除的，但包括活跃的和已过期的)
	var count int64
	s.db.Model(&model.FileBatch{}).Where("pickup_code = ?", string(b)).Count(&count)
	if count > 0 {
		return s.GeneratePickupCode(length) // 递归生成
	}
	return string(b), nil
}

func (s *BatchService) UpdateAllPickupCodes(newLength int) error {
	var batches []model.FileBatch
	// 只更新未删除的记录，包括 active 和 expired
	if err := s.db.Find(&batches).Error; err != nil {
		return err
	}

	return s.db.Transaction(func(tx *gorm.DB) error {
		for _, batch := range batches {
			oldCode := batch.PickupCode
			if len(oldCode) == newLength {
				continue
			}

			var newCode string
			if len(oldCode) < newLength {
				// 右侧补零，方便用户输入原码后通过补 0 完成输入
				newCode = oldCode + strings.Repeat("0", newLength-len(oldCode))
			} else {
				// 截取前 newLength 位，保留原码头部
				newCode = oldCode[:newLength]
			}

			// 检查是否冲突 (在事务中检查)
			var count int64
			tx.Model(&model.FileBatch{}).Where("pickup_code = ? AND id != ?", newCode, batch.ID).Count(&count)
			if count > 0 {
				// 如果冲突，生成一个新的随机码
				var err error
				newCode, err = s.generateUniquePickupCodeInTx(tx, newLength)
				if err != nil {
					return err
				}
			}

			if err := tx.Model(&batch).Update("pickup_code", newCode).Error; err != nil {
				return err
			}
		}
		return nil
	})
}

func (s *BatchService) Cleanup(ctx context.Context) error {
	slog.Debug("Starting cleanup scan")
	// 1. 寻找并标记过期的 Active Batches
	var batches []model.FileBatch
	now := time.Now()
	// 同时检查时间过期 and 下载次数过期
	err := s.db.Where("status = ? AND ("+
		"(expire_type = 'time' AND expire_at < ?) OR "+
		"(expire_type = 'download' AND max_downloads > 0 AND download_count >= max_downloads)"+
		")", "active", now).Find(&batches).Error
	if err != nil {
		slog.Error("Failed to query active batches for cleanup", "error", err)
		return err
	}

	if len(batches) > 0 {
		slog.Info("Found expired batches to mark", "count", len(batches))
	}

	for _, batch := range batches {
		_ = s.MarkAsExpired(&batch)
	}

	// 2. 检查异常文件 (物理文件缺失)
	var activeBatches []model.FileBatch
	if err := s.db.Preload("FileItems").Where("status = ?", "active").Find(&activeBatches).Error; err == nil {
		for _, batch := range activeBatches {
			var missingItems []model.FileItem
			for _, item := range batch.FileItems {
				exists, err := storage.GlobalStorage.Exists(ctx, item.StoragePath)
				if err != nil {
					slog.Error("Failed to check file existence", "path", item.StoragePath, "error", err)
					continue
				}
				if !exists {
					missingItems = append(missingItems, item)
				}
			}

			if len(missingItems) > 0 {
				slog.Info("Removing missing files from batch", "batch_id", batch.ID, "count", len(missingItems))
				for _, item := range missingItems {
					if err := s.db.Delete(&item).Error; err != nil {
						slog.Error("Failed to remove missing file record", "item_id", item.ID, "error", err)
					}
				}

				if len(missingItems) == len(batch.FileItems) {
					slog.Warn("All files missing for batch, marking as expired", "batch_id", batch.ID)
					_ = s.MarkAsExpired(&batch)
				}
			}
		}
	}

	// 3. 彻底清理标记为 expired 或 deleted 的批次
	var toDelete []model.FileBatch
	// Unscoped 用于包含已软删除但尚未物理清理的记录
	err = s.db.Unscoped().Where("status IN ? OR deleted_at IS NOT NULL", []string{"expired", "deleted"}).Find(&toDelete).Error
	if err != nil {
		slog.Error("Failed to query batches for physical deletion", "error", err)
		return err
	}

	if len(toDelete) > 0 {
		slog.Info("Found batches for physical deletion", "count", len(toDelete))
	}

	for _, batch := range toDelete {
		if err := s.DeleteBatch(ctx, batch.ID); err != nil {
			slog.Error("Failed to physically delete batch", "batch_id", batch.ID, "error", err)
		}
	}

	return nil
}

func (s *BatchService) generateUniquePickupCodeInTx(tx *gorm.DB, length int) (string, error) {
	const charset = "0123456789"
	for {
		b := make([]byte, length)
		for i := range b {
			num, err := rand.Int(rand.Reader, big.NewInt(int64(len(charset))))
			if err != nil {
				return "", err
			}
			b[i] = charset[num.Int64()]
		}

		var count int64
		tx.Model(&model.FileBatch{}).Where("pickup_code = ?", string(b)).Count(&count)
		if count == 0 {
			return string(b), nil
		}
	}
}
