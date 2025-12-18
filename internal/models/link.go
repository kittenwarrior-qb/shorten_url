package models

import (
	"time"

	"gorm.io/gorm"
)

type Link struct {
	ID          uint           `gorm:"primaryKey"`
	UserID      *uint          `gorm:"index"`
	ShortCode   string         `gorm:"uniqueIndex;size:20;not null"`
	OriginalURL string         `gorm:"size:2048;not null"`
	CustomAlias *string        `gorm:"size:20"`
	ClickCount  int64          `gorm:"default:0"`
	ExpiresAt   *time.Time     `gorm:"index"`
	CreatedAt   time.Time      `gorm:"autoCreateTime"`
	UpdatedAt   time.Time      `gorm:"autoUpdateTime"`
	DeletedAt   gorm.DeletedAt `gorm:"index"`
	User        *User          `gorm:"foreignKey:UserID"`
	Clicks      []Click        `gorm:"foreignKey:LinkID"`
}
