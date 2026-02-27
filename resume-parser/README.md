# Resume Parser

A high-performance Go resume parsing service that extracts structured data from PDF and DOCX files. Part of the [LearnBot AI Career Development Platform](../LearnBot.md).

## Features

- **Multi-format support**: PDF and DOCX file parsing
- **Comprehensive extraction**: Personal info, work experience, education, skills, certifications, and projects
- **Confidence scoring**: Every extracted field includes a 0.0–1.0 confidence score
- **Section detection**: Automatically identifies resume sections using keyword matching and heuristics
- **Error handling**: Structured error responses for malformed or unsupported documents
- **REST API**: Simple HTTP endpoint for file upload and parsing
- **Fast**: Designed for sub-5-second processing per resume

## Architecture

```
resume-parser/
├── cmd/server/          # HTTP server entry point
├── internal/
│   ├── api/             # HTTP handler and routing
│   ├── extractor/       # Field extraction logic
│   │   ├── personal.go      # Name, email, phone, location
│   │   ├── experience.go    # Work experience entries
│   │   ├── education.go     # Education entries
│   │   ├── skills.go        # Technical and soft skills
│   │   ├── certifications.go # Certifications and licenses
│   │   ├── projects.go      # Projects and achievements
│   │   └── sections.go      # Resume section detection
│   ├── parser/          # Document parsing
│   │   ├── pdf.go           # PDF text extraction (dslipak/pdf)
│   │   ├── docx.go          # DOCX text extraction (ZIP/XML)
│   │   └── resume_parser.go # Orchestration pipeline
│   └── schema/          # Data types and structures
│       └── types.go
└── docs/
    └── API.md           # Full API documentation
```

## Quick Start

### Prerequisites

- Go 1.22+

### Build & Run

```bash
cd resume-parser
go mod download
go build -o resume-parser ./cmd/server
./resume-parser -addr :8080
```

### Parse a Resume

```bash
# PDF
curl -X POST http://localhost:8080/api/v1/parse \
  -F "resume=@resume.pdf"

# DOCX with raw text
curl -X POST http://localhost:8080/api/v1/parse \
  -F "resume=@resume.docx" \
  -F "include_raw=true"
```

### Example Response

```json
{
  "success": true,
  "data": {
    "parsed_at": "2024-01-15T10:30:00Z",
    "source_file": "resume.pdf",
    "file_type": "pdf",
    "parser_version": "1.0.0",
    "personal_info": {
      "name": "Jane Doe",
      "email": "jane.doe@email.com",
      "phone": "(555) 987-6543",
      "location": "San Francisco, CA",
      "linkedin": "https://linkedin.com/in/janedoe",
      "confidence": 0.95
    },
    "work_experience": [
      {
        "company": "TechCorp Inc.",
        "title": "Senior Software Engineer",
        "start_date": "March 2021",
        "end_date": "Present",
        "is_current": true,
        "responsibilities": [
          "Architected distributed job scheduling system handling 1M+ tasks/day"
        ],
        "confidence": 0.88
      }
    ],
    "education": [
      {
        "institution": "University of California, Berkeley",
        "degree": "Bachelor of Science",
        "field": "Computer Science",
        "end_date": "2017",
        "gpa": "3.7",
        "confidence": 1.0
      }
    ],
    "skills": [
      { "name": "Go", "category": "technical", "confidence": 0.9 },
      { "name": "Leadership", "category": "soft", "confidence": 0.9 }
    ],
    "certifications": [
      {
        "name": "AWS Certified Solutions Architect",
        "issuer": "aws",
        "date": "2022",
        "confidence": 1.0
      }
    ],
    "projects": [
      {
        "name": "Distributed Task Scheduler",
        "url": "https://github.com/user/scheduler",
        "technologies": ["go", "redis"],
        "confidence": 1.0
      }
    ],
    "overall_confidence": 0.87,
    "sections_found": ["summary", "experience", "education", "skills", "certifications", "projects"],
    "warnings": []
  }
}
```

## Running Tests

```bash
# Run all tests
go test ./...

# Run with coverage
go test -cover ./...

# Run with verbose output
go test -v ./...

# Run specific package
go test ./internal/extractor/...
```

### Coverage Results

| Package | Coverage |
|---------|----------|
| `internal/api` | 86.4% |
| `internal/extractor` | 85.5% |
| `internal/parser` | 85.6% |

## Extraction Logic

### Section Detection

The parser identifies resume sections by scanning for header lines matching known keywords:

- **Summary**: "Summary", "Objective", "Profile", "About"
- **Experience**: "Experience", "Work Experience", "Employment History"
- **Education**: "Education", "Academic Background"
- **Skills**: "Skills", "Technical Skills", "Core Competencies"
- **Certifications**: "Certifications", "Licenses", "Credentials"
- **Projects**: "Projects", "Achievements", "Portfolio"

Both title-case and ALL-CAPS headers are supported.

### Personal Information

Extracted using regex patterns:
- **Email**: RFC 5322-compliant email regex
- **Phone**: US phone number formats (dashes, dots, parentheses)
- **Name**: Heuristic scan of first 10 lines for proper name patterns
- **Location**: "City, State" or "City, Country" patterns
- **LinkedIn/GitHub**: URL pattern matching

### Work Experience

Each job entry is identified by:
1. Lines containing job title keywords (engineer, developer, manager, etc.)
2. Date range patterns (e.g., "Jan 2020 - Present")
3. Bullet points for responsibilities

### Education

Education entries are parsed by:
1. Institution name detection (university, college, institute keywords)
2. Degree keyword matching with word-boundary regex
3. Year range extraction
4. GPA pattern matching

### Skills

Skills are extracted from the skills section and classified as:
- **Technical**: Matched against a curated list of 100+ known technical skills
- **Soft**: Matched against a curated list of soft skills
- **Other**: Unrecognized skills

### Confidence Scoring

Each field's confidence is computed based on:
- Whether required sub-fields were found
- Whether the extracted value matches expected patterns
- Known-good values (e.g., recognized skill names get higher confidence)

## API Documentation

See [`docs/API.md`](docs/API.md) for full API reference including all endpoints, request/response schemas, and error codes.

## Integration with LearnBot

This parser is Phase 1 of the LearnBot platform. The `ParsedResume` output maps directly to the `UserProfile` schema used by the RAG pipeline for:

- Skill gap analysis
- Job matching and acceptance likelihood scoring
- Personalized training recommendations

## Dependencies

| Package | Purpose |
|---------|---------|
| [`github.com/dslipak/pdf`](https://github.com/dslipak/pdf) | PDF text extraction |
| Standard library `archive/zip` | DOCX (ZIP) parsing |
| Standard library `encoding/xml` | DOCX XML parsing |
| Standard library `net/http` | HTTP server |
