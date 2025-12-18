package models

import "time"

type Click struct {
	ID          uint      `gorm:"primaryKey"`
	LinkID      uint      `gorm:"index;not null"`
	IPAddress   string    `gorm:"size:45"`
	UserAgent   string    `gorm:"size:512"`
	Browser     string    `gorm:"size:50"`
	BrowserVer  string    `gorm:"size:20"`
	OS          string    `gorm:"size:50"`
	Device      string    `gorm:"size:50"` // Desktop, Mobile, Tablet
	Country     string    `gorm:"size:100;index"`
	CountryCode string    `gorm:"size:2"`
	City        string    `gorm:"size:100"`
	Referer     string    `gorm:"size:2048"`
	ClickedAt   time.Time `gorm:"autoCreateTime;index"`
	Link        *Link     `gorm:"foreignKey:LinkID"`
}
