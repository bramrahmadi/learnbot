// Package parse provides a public API for resume parsing.
// This package wraps the internal parser and schema packages for use by external modules.
package parse

import (
	"github.com/learnbot/resume-parser/internal/parser"
	"github.com/learnbot/resume-parser/internal/schema"
)

// Re-export types from schema for external use.

// ParsedResume is the top-level structured output of the resume parser.
type ParsedResume = schema.ParsedResume

// ParseRequest is the input to the parser.
type ParseRequest = schema.ParseRequest

// ParseError represents a structured parsing error.
type ParseError = schema.ParseError

// ResumeParser orchestrates document parsing and field extraction.
type ResumeParser struct {
	inner *parser.ResumeParser
}

// New creates a new ResumeParser.
func New() *ResumeParser {
	return &ResumeParser{inner: parser.NewResumeParser()}
}

// Parse accepts a ParseRequest and returns a structured ParsedResume.
func (p *ResumeParser) Parse(req ParseRequest) (*ParsedResume, error) {
	return p.inner.Parse(req)
}
