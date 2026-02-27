package parser

import (
	"strings"
	"testing"

	"github.com/learnbot/resume-parser/internal/schema"
)

// TestResumeParser_EmptyFile tests error handling for empty files.
func TestResumeParser_EmptyFile(t *testing.T) {
	rp := NewResumeParser()
	_, err := rp.Parse(schema.ParseRequest{
		FileName:    "empty.pdf",
		FileContent: []byte{},
		FileType:    "pdf",
	})
	if err == nil {
		t.Error("expected error for empty file")
	}
}

// TestResumeParser_UnsupportedFormat tests error handling for unsupported formats.
func TestResumeParser_UnsupportedFormat(t *testing.T) {
	rp := NewResumeParser()
	_, err := rp.Parse(schema.ParseRequest{
		FileName:    "resume.txt",
		FileContent: []byte("some content"),
		FileType:    "txt",
	})
	if err == nil {
		t.Error("expected error for unsupported format")
	}
	if pe, ok := err.(*schema.ParseError); ok {
		if pe.Code != "UNSUPPORTED_FORMAT" {
			t.Errorf("expected error code UNSUPPORTED_FORMAT, got %q", pe.Code)
		}
	}
}

// TestResumeParser_InvalidPDF tests error handling for invalid PDF content.
func TestResumeParser_InvalidPDF(t *testing.T) {
	rp := NewResumeParser()
	_, err := rp.Parse(schema.ParseRequest{
		FileName:    "bad.pdf",
		FileContent: []byte("this is not a pdf"),
		FileType:    "pdf",
	})
	if err == nil {
		t.Error("expected error for invalid PDF")
	}
}

// TestResumeParser_InvalidDOCX tests error handling for invalid DOCX content.
func TestResumeParser_InvalidDOCX(t *testing.T) {
	rp := NewResumeParser()
	_, err := rp.Parse(schema.ParseRequest{
		FileName:    "bad.docx",
		FileContent: []byte("this is not a docx"),
		FileType:    "docx",
	})
	if err == nil {
		t.Error("expected error for invalid DOCX")
	}
}

// TestNormalizeFileType tests file type normalization.
func TestNormalizeFileType(t *testing.T) {
	tests := []struct {
		fileType string
		fileName string
		want     string
	}{
		{"pdf", "", "pdf"},
		{"PDF", "", "pdf"},
		{"application/pdf", "", "pdf"},
		{"docx", "", "docx"},
		{"", "resume.pdf", "pdf"},
		{"", "resume.docx", "docx"},
		{"", "resume.DOCX", "docx"},
		{"txt", "resume.txt", ""},
		{"", "resume.txt", ""},
	}

	for _, tt := range tests {
		t.Run(tt.fileType+"/"+tt.fileName, func(t *testing.T) {
			got := normalizeFileType(tt.fileType, tt.fileName)
			if got != tt.want {
				t.Errorf("normalizeFileType(%q, %q) = %q, want %q",
					tt.fileType, tt.fileName, got, tt.want)
			}
		})
	}
}

// TestResumeParser_DOCXParsing tests DOCX text extraction error cases.
func TestResumeParser_DOCXParsing(t *testing.T) {
	dp := NewDOCXParser()

	// Test with empty data
	_, err := dp.ExtractText([]byte{})
	if err == nil {
		t.Error("expected error for empty DOCX")
	}

	// Test with invalid data
	_, err = dp.ExtractText([]byte("not a zip file"))
	if err == nil {
		t.Error("expected error for invalid DOCX")
	}
}

// TestResumeParser_PDFParsing tests PDF text extraction error cases.
func TestResumeParser_PDFParsing(t *testing.T) {
	pp := NewPDFParser()

	// Test with empty data
	_, err := pp.ExtractText([]byte{})
	if err == nil {
		t.Error("expected error for empty PDF")
	}

	// Test with invalid data (no PDF header)
	_, err = pp.ExtractText([]byte("not a pdf file"))
	if err == nil {
		t.Error("expected error for invalid PDF")
	}
}

// TestCleanText tests the text cleaning function.
func TestCleanText(t *testing.T) {
	tests := []struct {
		name  string
		input string
		check func(string) bool
	}{
		{
			name:  "removes control chars",
			input: "Hello\x00World\x01Test",
			check: func(s string) bool {
				return !strings.ContainsAny(s, "\x00\x01")
			},
		},
		{
			name:  "normalizes multiple newlines",
			input: "Line1\n\n\n\nLine2",
			check: func(s string) bool {
				return !strings.Contains(s, "\n\n\n")
			},
		},
		{
			name:  "trims whitespace",
			input: "  \n  Hello World  \n  ",
			check: func(s string) bool {
				return s == "Hello World"
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := cleanText(tt.input)
			if !tt.check(result) {
				t.Errorf("cleanText(%q) = %q, check failed", tt.input, result)
			}
		})
	}
}

// TestParseError_Error tests the ParseError.Error() method.
func TestParseError_Error(t *testing.T) {
	pe := &ParseError{Code: "TEST_CODE", Message: "test message"}
	got := pe.Error()
	if !strings.Contains(got, "TEST_CODE") {
		t.Errorf("expected error string to contain code, got %q", got)
	}
	if !strings.Contains(got, "test message") {
		t.Errorf("expected error string to contain message, got %q", got)
	}
}

// TestSchemaParseError_Error tests the schema.ParseError.Error() method.
func TestSchemaParseError_Error(t *testing.T) {
	pe := &schema.ParseError{Code: "TEST", Message: "msg", Section: "experience"}
	got := pe.Error()
	if !strings.Contains(got, "experience") {
		t.Errorf("expected error string to contain section, got %q", got)
	}

	pe2 := &schema.ParseError{Code: "TEST", Message: "msg"}
	got2 := pe2.Error()
	if got2 == "" {
		t.Error("expected non-empty error string")
	}
}
