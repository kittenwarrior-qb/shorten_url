package service_test

import (
	"testing"
	"time"

	"quocbui.dev/m/internal/models"
	"quocbui.dev/m/internal/service"
	"quocbui.dev/m/tests/mocks"
)

func setupLinkService() (*service.LinkService, *mocks.MockLinkRepository, *mocks.MockClickRepository) {
	linkRepo := mocks.NewMockLinkRepository()
	clickRepo := mocks.NewMockClickRepository()
	txManager := mocks.NewMockTransactionManager()
	geoIP := service.NewGeoIPService()

	svc := service.NewLinkService(linkRepo, clickRepo, txManager, geoIP)
	return svc, linkRepo, clickRepo
}

func TestLinkService_CreateLink_Success(t *testing.T) {
	svc, _, _ := setupLinkService()

	userID := uint(1)
	link, err := svc.CreateLink("https://example.com/long/url", nil, &userID, nil, 6)
	if err != nil {
		t.Fatalf("CreateLink returned error: %v", err)
	}

	if link == nil {
		t.Fatal("CreateLink returned nil link")
	}

	if link.OriginalURL != "https://example.com/long/url" {
		t.Errorf("link.OriginalURL = %s, want https://example.com/long/url", link.OriginalURL)
	}

	if len(link.ShortCode) != 6 {
		t.Errorf("link.ShortCode length = %d, want 6", len(link.ShortCode))
	}

	if link.UserID == nil || *link.UserID != userID {
		t.Error("link.UserID should be set to provided userID")
	}
}

func TestLinkService_CreateLink_WithCustomAlias(t *testing.T) {
	svc, _, _ := setupLinkService()

	userID := uint(1)
	alias := "my-custom-link"
	link, err := svc.CreateLink("https://example.com", &alias, &userID, nil, 6)
	if err != nil {
		t.Fatalf("CreateLink returned error: %v", err)
	}

	if link.ShortCode != alias {
		t.Errorf("link.ShortCode = %s, want %s", link.ShortCode, alias)
	}

	if link.CustomAlias == nil || *link.CustomAlias != alias {
		t.Error("link.CustomAlias should be set")
	}
}

func TestLinkService_CreateLink_InvalidURL(t *testing.T) {
	svc, _, _ := setupLinkService()

	invalidURLs := []string{
		"",
		"not-a-url",
		"ftp://example.com",
		"javascript:alert(1)",
	}

	for _, url := range invalidURLs {
		_, err := svc.CreateLink(url, nil, nil, nil, 6)
		if err == nil {
			t.Errorf("CreateLink(%q) expected error for invalid URL", url)
		}
	}
}

func TestLinkService_CreateLink_InvalidAlias(t *testing.T) {
	svc, _, _ := setupLinkService()

	invalidAliases := []string{
		"ab",                    // too short
		"abcdefghij12345678901", // too long (21 chars)
		"my link",               // contains space
		"my@link",               // contains special char
	}

	for _, alias := range invalidAliases {
		_, err := svc.CreateLink("https://example.com", &alias, nil, nil, 6)
		if err == nil {
			t.Errorf("CreateLink with alias %q expected error for invalid alias", alias)
		}
	}
}

func TestLinkService_CreateLink_DuplicateAlias(t *testing.T) {
	svc, linkRepo, _ := setupLinkService()

	alias := "existing"
	linkRepo.Links[alias] = &models.Link{
		ID:        1,
		ShortCode: alias,
	}

	_, err := svc.CreateLink("https://example.com", &alias, nil, nil, 6)
	if err == nil {
		t.Error("Expected error for duplicate alias")
	}
}

func TestLinkService_CreateLink_WithExpiration(t *testing.T) {
	svc, _, _ := setupLinkService()

	expiresAt := time.Now().Add(24 * time.Hour)
	link, err := svc.CreateLink("https://example.com", nil, nil, &expiresAt, 6)
	if err != nil {
		t.Fatalf("CreateLink returned error: %v", err)
	}

	if link.ExpiresAt == nil {
		t.Fatal("link.ExpiresAt should be set")
	}

	if !link.ExpiresAt.Equal(expiresAt) {
		t.Errorf("link.ExpiresAt = %v, want %v", link.ExpiresAt, expiresAt)
	}
}

func TestLinkService_Redirect_Success(t *testing.T) {
	svc, linkRepo, _ := setupLinkService()

	linkRepo.Links["abc123"] = &models.Link{
		ID:          1,
		ShortCode:   "abc123",
		OriginalURL: "https://example.com/original",
	}

	clickInfo := &service.ClickInfo{
		IPAddress: "127.0.0.1",
		UserAgent: "Mozilla/5.0",
		Referer:   "",
	}

	originalURL, err := svc.Redirect("abc123", clickInfo)
	if err != nil {
		t.Fatalf("Redirect returned error: %v", err)
	}

	if originalURL != "https://example.com/original" {
		t.Errorf("originalURL = %s, want https://example.com/original", originalURL)
	}
}

func TestLinkService_Redirect_NotFound(t *testing.T) {
	svc, _, _ := setupLinkService()

	clickInfo := &service.ClickInfo{
		IPAddress: "127.0.0.1",
		UserAgent: "Mozilla/5.0",
	}

	_, err := svc.Redirect("nonexistent", clickInfo)
	if err == nil {
		t.Error("Expected error for non-existent link")
	}
}

func TestLinkService_Redirect_Expired(t *testing.T) {
	svc, linkRepo, _ := setupLinkService()

	expiredTime := time.Now().Add(-1 * time.Hour)
	linkRepo.Links["expired"] = &models.Link{
		ID:          1,
		ShortCode:   "expired",
		OriginalURL: "https://example.com",
		ExpiresAt:   &expiredTime,
	}

	clickInfo := &service.ClickInfo{
		IPAddress: "127.0.0.1",
		UserAgent: "Mozilla/5.0",
	}

	_, err := svc.Redirect("expired", clickInfo)
	if err == nil {
		t.Error("Expected error for expired link")
	}
}

func TestLinkService_GetUserLinks_Success(t *testing.T) {
	svc, linkRepo, _ := setupLinkService()

	userID := uint(1)
	linkRepo.Links["link1"] = &models.Link{ID: 1, ShortCode: "link1", UserID: &userID}
	linkRepo.Links["link2"] = &models.Link{ID: 2, ShortCode: "link2", UserID: &userID}

	links, total, err := svc.GetUserLinks(userID, 1, 10)
	if err != nil {
		t.Fatalf("GetUserLinks returned error: %v", err)
	}

	if total != 2 {
		t.Errorf("total = %d, want 2", total)
	}

	if len(links) != 2 {
		t.Errorf("len(links) = %d, want 2", len(links))
	}
}

func TestLinkService_GetUserLinks_Empty(t *testing.T) {
	svc, _, _ := setupLinkService()

	links, total, err := svc.GetUserLinks(999, 1, 10)
	if err != nil {
		t.Fatalf("GetUserLinks returned error: %v", err)
	}

	if total != 0 {
		t.Errorf("total = %d, want 0", total)
	}

	if len(links) != 0 {
		t.Errorf("len(links) = %d, want 0", len(links))
	}
}

func TestLinkService_GetLinkWithAnalytics_Success(t *testing.T) {
	svc, linkRepo, _ := setupLinkService()

	userID := uint(1)
	linkRepo.Links["mylink"] = &models.Link{
		ID:          1,
		ShortCode:   "mylink",
		OriginalURL: "https://example.com",
		UserID:      &userID,
	}

	link, err := svc.GetLinkWithAnalytics("mylink", userID)
	if err != nil {
		t.Fatalf("GetLinkWithAnalytics returned error: %v", err)
	}

	if link == nil {
		t.Fatal("GetLinkWithAnalytics returned nil link")
	}
}

func TestLinkService_GetLinkWithAnalytics_NotFound(t *testing.T) {
	svc, _, _ := setupLinkService()

	_, err := svc.GetLinkWithAnalytics("nonexistent", 1)
	if err == nil {
		t.Error("Expected error for non-existent link")
	}
}

func TestLinkService_GetLinkWithAnalytics_Unauthorized(t *testing.T) {
	svc, linkRepo, _ := setupLinkService()

	ownerID := uint(1)
	otherUserID := uint(2)
	linkRepo.Links["mylink"] = &models.Link{
		ID:        1,
		ShortCode: "mylink",
		UserID:    &ownerID,
	}

	_, err := svc.GetLinkWithAnalytics("mylink", otherUserID)
	if err == nil {
		t.Error("Expected error for unauthorized access")
	}
}

func TestLinkService_DeleteLink_Success(t *testing.T) {
	svc, linkRepo, _ := setupLinkService()

	userID := uint(1)
	linkRepo.Links["todelete"] = &models.Link{
		ID:        1,
		ShortCode: "todelete",
		UserID:    &userID,
	}

	err := svc.DeleteLink("todelete", userID)
	if err != nil {
		t.Fatalf("DeleteLink returned error: %v", err)
	}

	if _, exists := linkRepo.Links["todelete"]; exists {
		t.Error("Link should be deleted from repository")
	}
}

func TestLinkService_DeleteLink_NotFound(t *testing.T) {
	svc, _, _ := setupLinkService()

	err := svc.DeleteLink("nonexistent", 1)
	if err == nil {
		t.Error("Expected error for non-existent link")
	}
}

func TestLinkService_DeleteLink_Unauthorized(t *testing.T) {
	svc, linkRepo, _ := setupLinkService()

	ownerID := uint(1)
	otherUserID := uint(2)
	linkRepo.Links["mylink"] = &models.Link{
		ID:        1,
		ShortCode: "mylink",
		UserID:    &ownerID,
	}

	err := svc.DeleteLink("mylink", otherUserID)
	if err == nil {
		t.Error("Expected error for unauthorized delete")
	}
}
