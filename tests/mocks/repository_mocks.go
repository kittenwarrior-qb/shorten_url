package mocks

import (
	"quocbui.dev/m/internal/dto"
	"quocbui.dev/m/internal/models"

	"gorm.io/gorm"
)

// MockUserRepository is a mock implementation of UserRepository
type MockUserRepository struct {
	Users         map[string]*models.User
	CreateErr     error
	GetByIDErr    error
	GetByEmailErr error
}

func NewMockUserRepository() *MockUserRepository {
	return &MockUserRepository{
		Users: make(map[string]*models.User),
	}
}

func (m *MockUserRepository) Create(user *models.User) error {
	if m.CreateErr != nil {
		return m.CreateErr
	}
	user.ID = uint(len(m.Users) + 1)
	m.Users[user.Email] = user
	return nil
}

func (m *MockUserRepository) GetByID(id uint) (*models.User, error) {
	if m.GetByIDErr != nil {
		return nil, m.GetByIDErr
	}
	for _, user := range m.Users {
		if user.ID == id {
			return user, nil
		}
	}
	return nil, gorm.ErrRecordNotFound
}

func (m *MockUserRepository) GetByEmail(email string) (*models.User, error) {
	if m.GetByEmailErr != nil {
		return nil, m.GetByEmailErr
	}
	if user, ok := m.Users[email]; ok {
		return user, nil
	}
	return nil, gorm.ErrRecordNotFound
}

func (m *MockUserRepository) Update(user *models.User) error {
	m.Users[user.Email] = user
	return nil
}

// MockLinkRepository is a mock implementation of LinkRepository
type MockLinkRepository struct {
	Links     map[string]*models.Link
	CreateErr error
	GetErr    error
	DeleteErr error
	NextID    uint
}

func NewMockLinkRepository() *MockLinkRepository {
	return &MockLinkRepository{
		Links:  make(map[string]*models.Link),
		NextID: 1,
	}
}

func (m *MockLinkRepository) Create(link *models.Link) error {
	if m.CreateErr != nil {
		return m.CreateErr
	}
	link.ID = m.NextID
	m.NextID++
	m.Links[link.ShortCode] = link
	return nil
}

func (m *MockLinkRepository) CreateWithTx(tx *gorm.DB, link *models.Link) error {
	return m.Create(link)
}

func (m *MockLinkRepository) GetByID(id uint) (*models.Link, error) {
	if m.GetErr != nil {
		return nil, m.GetErr
	}
	for _, link := range m.Links {
		if link.ID == id {
			return link, nil
		}
	}
	return nil, gorm.ErrRecordNotFound
}

func (m *MockLinkRepository) GetByShortCode(shortCode string) (*models.Link, error) {
	if m.GetErr != nil {
		return nil, m.GetErr
	}
	if link, ok := m.Links[shortCode]; ok {
		return link, nil
	}
	return nil, gorm.ErrRecordNotFound
}

func (m *MockLinkRepository) GetByShortCodeForUpdate(tx *gorm.DB, shortCode string) (*models.Link, error) {
	return m.GetByShortCode(shortCode)
}

func (m *MockLinkRepository) GetByUserID(userID uint, page, pageSize int) ([]*models.Link, int64, error) {
	var links []*models.Link
	for _, link := range m.Links {
		if link.UserID != nil && *link.UserID == userID {
			links = append(links, link)
		}
	}
	return links, int64(len(links)), nil
}

func (m *MockLinkRepository) IncrementClickCount(id uint) error {
	for _, link := range m.Links {
		if link.ID == id {
			link.ClickCount++
			return nil
		}
	}
	return nil
}

func (m *MockLinkRepository) IncrementClickCountWithTx(tx *gorm.DB, id uint) error {
	return m.IncrementClickCount(id)
}

func (m *MockLinkRepository) Delete(id uint) error {
	if m.DeleteErr != nil {
		return m.DeleteErr
	}
	for code, link := range m.Links {
		if link.ID == id {
			delete(m.Links, code)
			return nil
		}
	}
	return nil
}

// MockClickRepository is a mock implementation of ClickRepository
type MockClickRepository struct {
	Clicks    []*models.Click
	CreateErr error
}

func NewMockClickRepository() *MockClickRepository {
	return &MockClickRepository{
		Clicks: make([]*models.Click, 0),
	}
}

func (m *MockClickRepository) Create(click *models.Click) error {
	if m.CreateErr != nil {
		return m.CreateErr
	}
	click.ID = uint(len(m.Clicks) + 1)
	m.Clicks = append(m.Clicks, click)
	return nil
}

func (m *MockClickRepository) CreateWithTx(tx *gorm.DB, click *models.Click) error {
	return m.Create(click)
}

func (m *MockClickRepository) GetByLinkID(linkID uint, page, pageSize int) ([]*models.Click, int64, error) {
	var clicks []*models.Click
	for _, click := range m.Clicks {
		if click.LinkID == linkID {
			clicks = append(clicks, click)
		}
	}
	return clicks, int64(len(clicks)), nil
}

func (m *MockClickRepository) GetAnalytics(linkID uint) (*dto.AnalyticsSummary, error) {
	return &dto.AnalyticsSummary{}, nil
}

// MockTransactionManager is a mock implementation of TransactionManager
type MockTransactionManager struct {
	ExecuteErr error
}

func NewMockTransactionManager() *MockTransactionManager {
	return &MockTransactionManager{}
}

func (m *MockTransactionManager) ExecuteInTransaction(fn func(tx *gorm.DB) error) error {
	if m.ExecuteErr != nil {
		return m.ExecuteErr
	}
	return fn(nil)
}
