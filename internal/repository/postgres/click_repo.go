package postgres

import (
	"gorm.io/gorm"

	"quocbui.dev/m/internal/dto"
	"quocbui.dev/m/internal/models"
	"quocbui.dev/m/internal/repository"
)

type clickRepository struct {
	db *gorm.DB
}

func NewClickRepository(db *gorm.DB) repository.ClickRepository {
	return &clickRepository{db: db}
}

func (r *clickRepository) Create(click *models.Click) error {
	return r.db.Create(click).Error
}

func (r *clickRepository) GetByLinkID(linkID uint, page, pageSize int) ([]*models.Click, int64, error) {
	var clicks []*models.Click
	var total int64

	offset := (page - 1) * pageSize

	err := r.db.Model(&models.Click{}).Where("link_id = ?", linkID).Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	err = r.db.Where("link_id = ?", linkID).
		Order("clicked_at DESC").
		Offset(offset).
		Limit(pageSize).
		Find(&clicks).Error

	return clicks, total, err
}

func (r *clickRepository) GetAnalytics(linkID uint) (*dto.AnalyticsSummary, error) {
	var totalClicks int64
	r.db.Model(&models.Click{}).Where("link_id = ?", linkID).Count(&totalClicks)

	summary := &dto.AnalyticsSummary{
		TotalClicks:  totalClicks,
		BrowserStats: make(map[string]int64),
		OSStats:      make(map[string]int64),
		DeviceStats:  make(map[string]int64),
		CountryStats: make(map[string]int64),
		RefererStats: make(map[string]int64),
	}

	// Browser stats
	var browserResults []struct {
		Browser string
		Count   int64
	}
	r.db.Model(&models.Click{}).
		Select("browser, count(*) as count").
		Where("link_id = ?", linkID).
		Group("browser").
		Scan(&browserResults)
	for _, b := range browserResults {
		summary.BrowserStats[b.Browser] = b.Count
	}

	// OS stats
	var osResults []struct {
		OS    string
		Count int64
	}
	r.db.Model(&models.Click{}).
		Select("os, count(*) as count").
		Where("link_id = ?", linkID).
		Group("os").
		Scan(&osResults)
	for _, o := range osResults {
		summary.OSStats[o.OS] = o.Count
	}

	// Device stats
	var deviceResults []struct {
		Device string
		Count  int64
	}
	r.db.Model(&models.Click{}).
		Select("device, count(*) as count").
		Where("link_id = ?", linkID).
		Group("device").
		Scan(&deviceResults)
	for _, d := range deviceResults {
		summary.DeviceStats[d.Device] = d.Count
	}

	// Country stats
	var countryResults []struct {
		Country string
		Count   int64
	}
	r.db.Model(&models.Click{}).
		Select("country, count(*) as count").
		Where("link_id = ?", linkID).
		Group("country").
		Scan(&countryResults)
	for _, c := range countryResults {
		summary.CountryStats[c.Country] = c.Count
	}

	// Referer stats
	var refererResults []struct {
		Referer string
		Count   int64
	}
	r.db.Model(&models.Click{}).
		Select("referer, count(*) as count").
		Where("link_id = ? AND referer != ''", linkID).
		Group("referer").
		Scan(&refererResults)
	for _, ref := range refererResults {
		summary.RefererStats[ref.Referer] = ref.Count
	}

	return summary, nil
}
