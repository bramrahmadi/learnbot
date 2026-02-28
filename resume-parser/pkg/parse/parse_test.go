package parse_test

import (
	"testing"

	"github.com/learnbot/resume-parser/pkg/parse"
)

func TestNew(t *testing.T) {
	p := parse.New()
	if p == nil {
		t.Fatal("expected non-nil ResumeParser")
	}
}

func TestParse_EmptyContent(t *testing.T) {
	p := parse.New()
	_, err := p.Parse(parse.ParseRequest{
		FileName:    "empty.pdf",
		FileContent: []byte{},
		FileType:    "pdf",
	})
	if err == nil {
		t.Error("expected error for empty file content")
	}
}

func TestParse_UnsupportedFormat(t *testing.T) {
	p := parse.New()
	_, err := p.Parse(parse.ParseRequest{
		FileName:    "resume.txt",
		FileContent: []byte("some content"),
		FileType:    "txt",
	})
	if err == nil {
		t.Error("expected error for unsupported format")
	}
}
