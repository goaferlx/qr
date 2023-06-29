package qr

import (
	"bytes"
	"context"
	"encoding/base64"
	"errors"
	"fmt"
	"image/png"

	"github.com/boombuler/barcode"
	"github.com/boombuler/barcode/qr"
)

// Default size of the QR code - 601px is the breakpoint for small screens on @media queries.
// This will display as full size on mobile devices without compression.
const DefaultQRSize int = 600

var (
	ErrUnsupportedDataType error = errors.New("unsupported data type")
)

func GenerateCode(ctx context.Context, req CodeRequest) (string, error) {
	var text string
	switch req.DataType {
	case "url":
		text = req.Text
	case "tel":
		text = "tel:" + req.Text
	case "email":
		text = "mailto:" + req.Text
	case "sms":
		text = "smsto:" + req.Text
	default:
		return "", ErrUnsupportedDataType
	}

	qrcode, err := qr.Encode(text, qr.H, qr.Auto)
	if err != nil {
		return "", fmt.Errorf("failed to encode qr code: %w", err)
	}

	qrcode, err = barcode.Scale(qrcode, DefaultQRSize, DefaultQRSize)
	if err != nil {
		return "", fmt.Errorf("failed to resize qr code: %w", err)
	}

	var buf bytes.Buffer
	if err := png.Encode(&buf, qrcode); err != nil {
		return "", fmt.Errorf("failed to encode to png: %w", err)
	}

	return base64.RawStdEncoding.EncodeToString(buf.Bytes()), nil
}
