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
	ErrRequired            error = errors.New("input is required")
)

type DataType string

const (
	TypeURL   = "url"
	TypeTel   = "tel"
	TypeSMS   = "sms"
	TypeEmail = "email"
)

type CodeRequest struct {
	DataType string `json:"data_type,omitempty"`
	Text     string `json:"text,omitempty"`
}

func GenerateCode(ctx context.Context, req CodeRequest) (string, error) {
	if req.Text == "" {
		return "", ErrRequired
	}
	var textToEncode string
	switch req.DataType {
	case TypeURL:
		textToEncode = req.Text
	case TypeTel:
		textToEncode = "tel:" + req.Text
	case TypeEmail:
		textToEncode = "mailto:" + req.Text
	case TypeSMS:
		textToEncode = "smsto:" + req.Text
	default:
		return "", ErrUnsupportedDataType
	}

	qrcode, err := qr.Encode(textToEncode, qr.H, qr.Auto)
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
