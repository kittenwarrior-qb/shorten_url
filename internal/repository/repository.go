package repository

import (
	"quocbui.dev/m/internal/models"
)

type UserRepository interface {
	Create(user *models.User) error
	GetByID(id uint) (*models.User, error)
	GetByEmail(email string) (*models.User, error)
	Update(user *models.User) error
}

type LinkRepository interface {
	Create(link *models.Link) error
	GetByShortCode(shortCode string) (*models.Link, error)
	GetByUserID(userID uint, offset, limit int) ([]models.Link, int64, error)
	IncrementClickCount(id uint) error
	Delete(id uint) error
}

type ClickRepository interface {
	Create(click *models.Click) error
	GetByLinkID(linkID uint, offset, limit int) ([]models.Click, int64, error)
}
