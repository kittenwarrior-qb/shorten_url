package utils

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/yeqown/go-qrcode/v2"
	"github.com/yeqown/go-qrcode/writer/standard"
)

// GenerateQRCode generates a QR code PNG with optional logo
func GenerateQRCode(content string, logoPath string) ([]byte, error) {
	qrc, err := qrcode.NewWith(content,
		qrcode.WithErrorCorrectionLevel(qrcode.ErrorCorrectionHighest),
	)
	if err != nil {
		return nil, err
	}

	// Create temp file
	tmpFile, err := os.CreateTemp("", "qr-*.png")
	if err != nil {
		return nil, err
	}
	tmpPath := tmpFile.Name()
	tmpFile.Close()
	defer os.Remove(tmpPath)

	// Find logo file
	logoFound := findLogoFile(logoPath)
	fmt.Printf("Logo path: %s, Found: %s\n", logoPath, logoFound)

	options := []standard.ImageOption{
		standard.WithBgColorRGBHex("#ffffff"),
		standard.WithFgColorRGBHex("#000000"),
		standard.WithQRWidth(30),
		standard.WithBuiltinImageEncoder(standard.PNG_FORMAT),
	}

	if logoFound != "" {
		options = append(options, standard.WithLogoImageFilePNG(logoFound))
		options = append(options, standard.WithLogoSizeMultiplier(3))
	}

	w, err := standard.New(tmpPath, options...)
	if err != nil {
		fmt.Printf("Writer error: %v\n", err)
		return nil, err
	}

	if err := qrc.Save(w); err != nil {
		fmt.Printf("Save error: %v\n", err)
		return nil, err
	}

	return os.ReadFile(tmpPath)
}

func findLogoFile(logoPath string) string {
	// Try original path
	if _, err := os.Stat(logoPath); err == nil {
		return logoPath
	}

	// Try from working directory
	wd, _ := os.Getwd()
	absPath := filepath.Join(wd, logoPath)
	if _, err := os.Stat(absPath); err == nil {
		return absPath
	}

	return ""
}
