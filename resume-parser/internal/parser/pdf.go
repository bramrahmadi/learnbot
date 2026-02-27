// Package parser provides document parsing for PDF and DOCX resume files.
package parser

import (
	"bytes"
	"fmt"
	"strings"
	"unicode"

	"github.com/dslipak/pdf"
)

// PDFParser extracts text from PDF files.
type PDFParser struct{}

// NewPDFParser creates a new PDFParser instance.
func NewPDFParser() *PDFParser {
	return &PDFParser{}
}

// ExtractText extracts all text content from a PDF byte slice.
// It returns the raw text and any error encountered.
func (p *PDFParser) ExtractText(data []byte) (string, error) {
	if len(data) == 0 {
		return "", &ParseError{Code: "EMPTY_FILE", Message: "PDF file is empty"}
	}

	// Validate PDF header
	if !bytes.HasPrefix(data, []byte("%PDF")) {
		return "", &ParseError{Code: "INVALID_FORMAT", Message: "file does not appear to be a valid PDF"}
	}

	r, err := pdf.NewReader(bytes.NewReader(data), int64(len(data)))
	if err != nil {
		return "", &ParseError{
			Code:    "PDF_PARSE_ERROR",
			Message: fmt.Sprintf("failed to open PDF: %v", err),
		}
	}

	var sb strings.Builder
	numPages := r.NumPage()

	if numPages == 0 {
		return "", &ParseError{
			Code:    "NO_TEXT_CONTENT",
			Message: "PDF has no pages",
		}
	}

	for i := 1; i <= numPages; i++ {
		page := r.Page(i)
		if page.V.IsNull() {
			continue
		}

		text, err := page.GetPlainText(nil)
		if err != nil {
			// Skip pages that fail, continue with others
			continue
		}
		sb.WriteString(text)
		sb.WriteString("\n")
	}

	result := sb.String()
	if strings.TrimSpace(result) == "" {
		return "", &ParseError{
			Code:    "NO_TEXT_CONTENT",
			Message: "PDF appears to contain no extractable text (may be image-based or encrypted)",
		}
	}

	return cleanText(result), nil
}

// cleanText normalizes extracted text by removing control characters
// and normalizing whitespace while preserving newlines.
func cleanText(text string) string {
	var sb strings.Builder
	prevNewline := false

	for _, r := range text {
		if r == '\n' || r == '\r' {
			if !prevNewline {
				sb.WriteRune('\n')
				prevNewline = true
			}
			continue
		}
		if unicode.IsControl(r) {
			continue
		}
		prevNewline = false
		sb.WriteRune(r)
	}

	// Normalize multiple blank lines to at most two
	result := sb.String()
	for strings.Contains(result, "\n\n\n") {
		result = strings.ReplaceAll(result, "\n\n\n", "\n\n")
	}

	return strings.TrimSpace(result)
}
