package service_test

import (
	"testing"

	"quocbui.dev/m/internal/models"
	"quocbui.dev/m/internal/service"
	"quocbui.dev/m/tests/mocks"
)

func setupAnalyticsService() (*service.AnalyticsService, *mocks.MockClickRepository, *mocks.MockLinkRepository) {
	clickRepo := mocks.NewMockClickRepository()
	linkRepo := mocks.NewMockLinkRepository()
	svc := service.NewAnalyticsService(clickRepo, linkRepo)
	return svc, clickRepo, linkRepo
}

func TestAnalyticsService_GetClicksByLinkID_Success(t *testing.T) {
	svc, clickRepo, linkRepo := setupAnalyticsService()

	userID := uint(1)
	linkID := uint(1)

	// Setup link
	linkRepo.Links["test"] = &models.Link{
		ID:        linkID,
		ShortCode: "test",
		UserID:    &userID,
	}

	// Setup clicks
	clickRepo.Clicks = append(clickRepo.Clicks, &models.Click{
		ID:     1,
		LinkID: linkID,
	})
	clickRepo.Clicks = append(clickRepo.Clicks, &models.Click{
		ID:     2,
		LinkID: linkID,
	})

	clicks, total, err := svc.GetClicksByLinkID(linkID, userID, 1, 10)
	if err != nil {
		t.Fatalf("GetClicksByLinkID returned error: %v", err)
	}

	if total != 2 {
		t.Errorf("total = %d, want 2", total)
	}

	if len(clicks) != 2 {
		t.Errorf("len(clicks) = %d, want 2", len(clicks))
	}
}

func TestAnalyticsService_GetClicksByLinkID_LinkNotFound(t *testing.T) {
	svc, _, _ := setupAnalyticsService()

	_, _, err := svc.GetClicksByLinkID(999, 1, 1, 10)
	if err == nil {
		t.Error("Expected error for non-existent link")
	}

	if err != service.ErrLinkNotFound {
		t.Errorf("Expected ErrLinkNotFound, got %v", err)
	}
}

func TestAnalyticsService_GetClicksByLinkID_Unauthorized(t *testing.T) {
	svc, _, linkRepo := setupAnalyticsService()

	ownerID := uint(1)
	otherUserID := uint(2)
	linkID := uint(1)

	linkRepo.Links["test"] = &models.Link{
		ID:        linkID,
		ShortCode: "test",
		UserID:    &ownerID,
	}

	_, _, err := svc.GetClicksByLinkID(linkID, otherUserID, 1, 10)
	if err == nil {
		t.Error("Expected error for unauthorized access")
	}

	if err != service.ErrUnauthorized {
		t.Errorf("Expected ErrUnauthorized, got %v", err)
	}
}

func TestAnalyticsService_GetAnalyticsSummary_Success(t *testing.T) {
	svc, _, linkRepo := setupAnalyticsService()

	userID := uint(1)
	linkID := uint(1)

	linkRepo.Links["test"] = &models.Link{
		ID:        linkID,
		ShortCode: "test",
		UserID:    &userID,
	}

	summary, err := svc.GetAnalyticsSummary(linkID, userID)
	if err != nil {
		t.Fatalf("GetAnalyticsSummary returned error: %v", err)
	}

	if summary == nil {
		t.Fatal("GetAnalyticsSummary returned nil summary")
	}
}

func TestAnalyticsService_GetAnalyticsSummary_LinkNotFound(t *testing.T) {
	svc, _, _ := setupAnalyticsService()

	_, err := svc.GetAnalyticsSummary(999, 1)
	if err == nil {
		t.Error("Expected error for non-existent link")
	}

	if err != service.ErrLinkNotFound {
		t.Errorf("Expected ErrLinkNotFound, got %v", err)
	}
}

func TestAnalyticsService_GetAnalyticsSummary_Unauthorized(t *testing.T) {
	svc, _, linkRepo := setupAnalyticsService()

	ownerID := uint(1)
	otherUserID := uint(2)
	linkID := uint(1)

	linkRepo.Links["test"] = &models.Link{
		ID:        linkID,
		ShortCode: "test",
		UserID:    &ownerID,
	}

	_, err := svc.GetAnalyticsSummary(linkID, otherUserID)
	if err == nil {
		t.Error("Expected error for unauthorized access")
	}

	if err != service.ErrUnauthorized {
		t.Errorf("Expected ErrUnauthorized, got %v", err)
	}
}
