package postgres

import (
	"errors"

	"gorm.io/gorm"

	"quocbui.dev/m/internal/models"
	"quocbui.dev/m/internal/repository"
)

type linkRepository struct {
	db *gorm.DB
}

func NewLinkRepository(db *gorm.DB) repository.LinkRepository {
	return &linkRepository{db: db}
}

func (r *linkRepository) Create(link *models.Link) error {
	return r.db.Create(link).Error
}

func (r *linkRepository) GetByID(id uint) (*models.Link, error) {
	var link models.Link
	err := r.db.First(&link, id).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}
	return &link, err
}

func (r *linkRepository) GetByShortCode(shortCode string) (*models.Link, error) {
	var link models.Link
	err := r.db.Where("short_code = ?", shortCode).First(&link).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}
	return &link, err
}

func (r *linkRepository) GetByUserID(userID uint, page, pageSize int) ([]*models.Link, int64, error) {
	var links []*models.Link
	var total int64

	offset := (page - 1) * pageSize

	err := r.db.Model(&models.Link{}).Where("user_id = ?", userID).Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	err = r.db.Where("user_id = ?", userID).
		Order("created_at DESC").
		Offset(offset).
		Limit(pageSize).
		Find(&links).Error

	return links, total, err
}

func (r *linkRepository) IncrementClickCount(id uint) error {
	return r.db.Model(&models.Link{}).Where("id = ?", id).
		UpdateColumn("click_count", gorm.Expr("click_count + ?", 1)).Error
}

func (r *linkRepository) Delete(id uint) error {
	return r.db.Delete(&models.Link{}, id).Error
}
