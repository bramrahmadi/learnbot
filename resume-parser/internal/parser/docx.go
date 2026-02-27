package parser

import (
	"archive/zip"
	"bytes"
	"encoding/xml"
	"fmt"
	"io"
	"strings"
)

// DOCXParser extracts text from DOCX files.
// DOCX files are ZIP archives containing XML documents.
type DOCXParser struct{}

// NewDOCXParser creates a new DOCXParser instance.
func NewDOCXParser() *DOCXParser {
	return &DOCXParser{}
}

// ExtractText extracts all text content from a DOCX byte slice.
func (d *DOCXParser) ExtractText(data []byte) (string, error) {
	if len(data) == 0 {
		return "", &ParseError{Code: "EMPTY_FILE", Message: "DOCX file is empty"}
	}

	// DOCX files are ZIP archives
	r, err := zip.NewReader(bytes.NewReader(data), int64(len(data)))
	if err != nil {
		return "", &ParseError{
			Code:    "INVALID_FORMAT",
			Message: fmt.Sprintf("file does not appear to be a valid DOCX (ZIP): %v", err),
		}
	}

	// Find word/document.xml
	var docXML []byte
	for _, f := range r.File {
		if f.Name == "word/document.xml" {
			rc, err := f.Open()
			if err != nil {
				return "", &ParseError{
					Code:    "DOCX_READ_ERROR",
					Message: fmt.Sprintf("failed to open word/document.xml: %v", err),
				}
			}
			docXML, err = io.ReadAll(rc)
			rc.Close()
			if err != nil {
				return "", &ParseError{
					Code:    "DOCX_READ_ERROR",
					Message: fmt.Sprintf("failed to read word/document.xml: %v", err),
				}
			}
			break
		}
	}

	if docXML == nil {
		return "", &ParseError{
			Code:    "INVALID_FORMAT",
			Message: "word/document.xml not found in DOCX archive",
		}
	}

	text, err := extractTextFromWordXML(docXML)
	if err != nil {
		return "", &ParseError{
			Code:    "DOCX_PARSE_ERROR",
			Message: fmt.Sprintf("failed to parse DOCX XML: %v", err),
		}
	}

	if strings.TrimSpace(text) == "" {
		return "", &ParseError{
			Code:    "NO_TEXT_CONTENT",
			Message: "DOCX document contains no extractable text",
		}
	}

	return cleanText(text), nil
}

// wordXMLElement represents a simplified Word XML element for text extraction.
type wordXMLElement struct {
	XMLName xml.Name
	Attrs   []xml.Attr `xml:",any,attr"`
	Content string     `xml:",chardata"`
	Inner   []wordXMLElement `xml:",any"`
}

// extractTextFromWordXML parses Word XML and extracts text content,
// preserving paragraph structure.
func extractTextFromWordXML(data []byte) (string, error) {
	decoder := xml.NewDecoder(bytes.NewReader(data))
	var sb strings.Builder

	type stackEntry struct {
		name string
	}
	var stack []stackEntry
	inParagraph := false
	paragraphHasText := false

	for {
		token, err := decoder.Token()
		if err == io.EOF {
			break
		}
		if err != nil {
			return "", fmt.Errorf("XML decode error: %w", err)
		}

		switch t := token.(type) {
		case xml.StartElement:
			localName := t.Name.Local
			stack = append(stack, stackEntry{name: localName})

			switch localName {
			case "p": // w:p - paragraph
				inParagraph = true
				paragraphHasText = false
			case "br": // w:br - line break
				sb.WriteRune('\n')
			case "tab": // w:tab
				sb.WriteRune('\t')
			}

		case xml.EndElement:
			if len(stack) > 0 {
				top := stack[len(stack)-1]
				stack = stack[:len(stack)-1]

				if top.name == "p" {
					inParagraph = false
					if paragraphHasText {
						sb.WriteRune('\n')
					}
				}
			}

		case xml.CharData:
			if inParagraph {
				text := string(t)
				if strings.TrimSpace(text) != "" {
					sb.WriteString(text)
					paragraphHasText = true
				}
			}
		}
	}

	return sb.String(), nil
}

// ParseError represents a document parsing error.
type ParseError struct {
	Code    string
	Message string
}

func (e *ParseError) Error() string {
	return e.Code + ": " + e.Message
}
