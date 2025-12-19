package utils

import (
	"bytes"
	"os"

	"github.com/yeqown/go-qrcode/v2"
	"github.com/yeqown/go-qrcode/writer/standard"
)

// GenerateQRCode generates a QR code PNG with optional logo
func GenerateQRCode(content string, logoPath string) ([]byte, error) {
	qrc, err := qrcode.New(content)
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

	options := []standard.ImageOption{
		standard.WithBgColorRGBHex("#ffffff"),
		standard.WithFgColorRGBHex("#000000"),
		standard.WithQRWidth(10),
	}

	// Add logo if file exists
	if logoPath != "" {
		if _, err := os.Stat(logoPath); err == nil {
			options = append(options, standard.WithLogoImageFilePNG(logoPath))
			options = append(options, standard.WithLogoSizeMultiplier(3))
		}
	}

	w, err := standard.New(tmpPath, options...)
	if err != nil {
		return nil, err
	}

	if err := qrc.Save(w); err != nil {
		return nil, err
	}

	// Read file to bytes
	var buf bytes.Buffer
	data, err := os.ReadFile(tmpPath)
	if err != nil {
		return nil, err
	}
	buf.Write(data)

	return buf.Bytes(), nil
}
