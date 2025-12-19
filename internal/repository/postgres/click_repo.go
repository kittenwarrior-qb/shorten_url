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

// CreateWithTx creates a click record within a transaction
func (r *clickRepository) CreateWithTx(tx *gorm.DB, click *models.Click) error {
	return tx.Create(click).Error
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
		TotalClicks: totalClicks,
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
	if len(browserResults) > 0 {
		summary.Browsers = make(map[string]int64)
		for _, b := range browserResults {
			summary.Browsers[b.Browser] = b.Count
		}
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
	if len(osResults) > 0 {
		summary.OS = make(map[string]int64)
		for _, o := range osResults {
			summary.OS[o.OS] = o.Count
		}
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
	if len(deviceResults) > 0 {
		summary.Devices = make(map[string]int64)
		for _, d := range deviceResults {
			summary.Devices[d.Device] = d.Count
		}
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
	if len(countryResults) > 0 {
		summary.Countries = make(map[string]int64)
		for _, c := range countryResults {
			summary.Countries[c.Country] = c.Count
		}
	}

	// Referer source stats (Facebook, Google, Direct...)
	var sourceResults []struct {
		RefererSource string
		Count         int64
	}
	r.db.Model(&models.Click{}).
		Select("referer_source, count(*) as count").
		Where("link_id = ?", linkID).
		Group("referer_source").
		Scan(&sourceResults)
	if len(sourceResults) > 0 {
		summary.RefererSources = make(map[string]int64)
		for _, s := range sourceResults {
			summary.RefererSources[s.RefererSource] = s.Count
		}
	}

	// Referer domain stats (chi tiết)
	var domainResults []struct {
		RefererDomain string
		Count         int64
	}
	r.db.Model(&models.Click{}).
		Select("referer_domain, count(*) as count").
		Where("link_id = ? AND referer_domain != ''", linkID).
		Group("referer_domain").
		Order("count DESC").
		Limit(10).
		Scan(&domainResults)
	if len(domainResults) > 0 {
		summary.RefererDomains = make(map[string]int64)
		for _, d := range domainResults {
			summary.RefererDomains[d.RefererDomain] = d.Count
		}
	}

	return summary, nil
}
