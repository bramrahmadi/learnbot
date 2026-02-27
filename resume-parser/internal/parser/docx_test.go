package parser

import (
	"archive/zip"
	"bytes"
	"fmt"
	"testing"
)

// buildDOCXWithContent creates a DOCX with specific XML content.
func buildDOCXWithContent(bodyXML string) []byte {
	var buf bytes.Buffer
	w := zip.NewWriter(&buf)

	docXML := fmt.Sprintf(`<?xml version="1.0" encoding="UTF-8" standalone="yes"?>
<w:document xmlns:w="http://schemas.openxmlformats.org/wordprocessingml/2006/main">
  <w:body>%s</w:body>
</w:document>`, bodyXML)

	f, _ := w.Create("word/document.xml")
	f.Write([]byte(docXML))

	ct, _ := w.Create("[Content_Types].xml")
	ct.Write([]byte(`<?xml version="1.0" encoding="UTF-8" standalone="yes"?>
<Types xmlns="http://schemas.openxmlformats.org/package/2006/content-types">
  <Default Extension="rels" ContentType="application/vnd.openxmlformats-package.relationships+xml"/>
  <Default Extension="xml" ContentType="application/xml"/>
  <Override PartName="/word/document.xml" ContentType="application/vnd.openxmlformats-officedocument.wordprocessingml.document.main+xml"/>
</Types>`))

	w.Close()
	return buf.Bytes()
}

// TestDOCXParser_ValidDocument tests successful DOCX text extraction.
func TestDOCXParser_ValidDocument(t *testing.T) {
	dp := NewDOCXParser()

	bodyXML := `<w:p><w:r><w:t>John Doe</w:t></w:r></w:p>
<w:p><w:r><w:t>john@example.com</w:t></w:r></w:p>
<w:p><w:r><w:t>Software Engineer</w:t></w:r></w:p>`

	docxData := buildDOCXWithContent(bodyXML)
	text, err := dp.ExtractText(docxData)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if text == "" {
		t.Error("expected non-empty text")
	}
	if !bytes.Contains([]byte(text), []byte("John Doe")) {
		t.Errorf("expected text to contain 'John Doe', got: %q", text)
	}
}

// TestDOCXParser_MissingDocumentXML tests DOCX without word/document.xml.
func TestDOCXParser_MissingDocumentXML(t *testing.T) {
	dp := NewDOCXParser()

	var buf bytes.Buffer
	w := zip.NewWriter(&buf)
	f, _ := w.Create("other/file.xml")
	f.Write([]byte("<root/>"))
	w.Close()

	_, err := dp.ExtractText(buf.Bytes())
	if err == nil {
		t.Error("expected error for DOCX without word/document.xml")
	}
}

// TestDOCXParser_EmptyDocument tests DOCX with empty body.
func TestDOCXParser_EmptyDocument(t *testing.T) {
	dp := NewDOCXParser()
	docxData := buildDOCXWithContent("")

	_, err := dp.ExtractText(docxData)
	if err == nil {
		t.Error("expected error for DOCX with no text content")
	}
}

// TestDOCXParser_WithLineBreaks tests DOCX with line break elements.
func TestDOCXParser_WithLineBreaks(t *testing.T) {
	dp := NewDOCXParser()

	bodyXML := `<w:p><w:r><w:t>Line 1</w:t></w:r><w:br/><w:r><w:t>Line 2</w:t></w:r></w:p>`
	docxData := buildDOCXWithContent(bodyXML)

	text, err := dp.ExtractText(docxData)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if text == "" {
		t.Error("expected non-empty text")
	}
}

// TestDOCXParser_WithTabs tests DOCX with tab elements.
func TestDOCXParser_WithTabs(t *testing.T) {
	dp := NewDOCXParser()

	bodyXML := `<w:p><w:r><w:t>Col1</w:t></w:r><w:tab/><w:r><w:t>Col2</w:t></w:r></w:p>`
	docxData := buildDOCXWithContent(bodyXML)

	text, err := dp.ExtractText(docxData)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if text == "" {
		t.Error("expected non-empty text")
	}
}
