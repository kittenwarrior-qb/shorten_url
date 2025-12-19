package repository

import (
	"quocbui.dev/m/internal/dto"
	"quocbui.dev/m/internal/models"

	"gorm.io/gorm"
)

// TransactionManager handles database transactions
type TransactionManager interface {
	// ExecuteInTransaction runs the given function within a transaction
	ExecuteInTransaction(fn func(tx *gorm.DB) error) error
}

type UserRepository interface {
	Create(user *models.User) error
	GetByID(id uint) (*models.User, error)
	GetByEmail(email string) (*models.User, error)
	Update(user *models.User) error
}

type LinkRepository interface {
	Create(link *models.Link) error
	CreateWithTx(tx *gorm.DB, link *models.Link) error
	GetByID(id uint) (*models.Link, error)
	GetByShortCode(shortCode string) (*models.Link, error)
	GetByShortCodeForUpdate(tx *gorm.DB, shortCode string) (*models.Link, error)
	GetByUserID(userID uint, page, pageSize int) ([]*models.Link, int64, error)
	IncrementClickCount(id uint) error
	IncrementClickCountWithTx(tx *gorm.DB, id uint) error
	Delete(id uint) error
}

type ClickRepository interface {
	Create(click *models.Click) error
	CreateWithTx(tx *gorm.DB, click *models.Click) error
	GetByLinkID(linkID uint, page, pageSize int) ([]*models.Click, int64, error)
	GetAnalytics(linkID uint) (*dto.AnalyticsSummary, error)
}
