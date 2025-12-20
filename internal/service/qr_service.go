package service

import (
	"encoding/base64"

	"quocbui.dev/m/pkg/utils"
)

// QRService handles QR code generation
type QRService struct {
	logoPath string
}

// NewQRService creates a new QR service
func NewQRService(logoPath string) *QRService {
	return &QRService{
		logoPath: logoPath,
	}
}

// GenerateQRCodeBase64 generates a QR code with logo and returns base64 encoded string
func (s *QRService) GenerateQRCodeBase64(url string) (string, error) {
	qrBytes, err := utils.GenerateQRCode(url, s.logoPath)
	if err != nil {
		return "", err
	}

	return "data:image/png;base64," + base64.StdEncoding.EncodeToString(qrBytes), nil
}
