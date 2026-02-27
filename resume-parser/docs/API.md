# Resume Parser API Documentation

## Overview

The Resume Parser API extracts structured information from PDF and DOCX resume files. It uses regex-based and heuristic NLP techniques to identify and extract personal information, work experience, education, skills, certifications, and projects. Each extracted field includes a confidence score (0.0–1.0).

**Base URL:** `http://localhost:8080`  
**API Version:** `v1`  
**Parser Version:** `1.0.0`

---

## Endpoints

### POST `/api/v1/parse`

Parse a resume file and return structured JSON data.

#### Request

**Content-Type:** `multipart/form-data`

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| `resume` | file | ✅ | The resume file (PDF or DOCX, max 10 MB) |
| `include_raw` | string | ❌ | Set to `"true"` to include raw extracted text in the response |

#### Supported File Types

| Extension | MIME Type |
|-----------|-----------|
| `.pdf` | `application/pdf` |
| `.docx` | `application/vnd.openxmlformats-officedocument.wordprocessingml.document` |

#### Example Request

```bash
# Parse a PDF resume
curl -X POST http://localhost:8080/api/v1/parse \
  -F "resume=@/path/to/resume.pdf"

# Parse a DOCX resume with raw text included
curl -X POST http://localhost:8080/api/v1/parse \
  -F "resume=@/path/to/resume.docx" \
  -F "include_raw=true"
```

#### Response

**Content-Type:** `application/json`

**Success Response (200 OK):**

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
      "github": "https://github.com/janedoe",
      "website": "https://janedoe.dev",
      "confidence": 0.95
    },
    "work_experience": [
      {
        "company": "TechCorp Inc.",
        "title": "Senior Software Engineer",
        "start_date": "March 2021",
        "end_date": "Present",
        "is_current": true,
        "location": "",
        "responsibilities": [
          "Architected distributed job scheduling system handling 1M+ tasks/day",
          "Reduced infrastructure costs by 35% through optimization"
        ],
        "confidence": 0.88
      }
    ],
    "education": [
      {
        "institution": "University of California, Berkeley",
        "degree": "Bachelor of Science",
        "field": "Computer Science",
        "start_date": "2013",
        "end_date": "2017",
        "gpa": "3.7",
        "honors": "",
        "confidence": 1.0
      }
    ],
    "skills": [
      {
        "name": "Go",
        "category": "technical",
        "confidence": 0.9
      },
      {
        "name": "Leadership",
        "category": "soft",
        "confidence": 0.9
      }
    ],
    "certifications": [
      {
        "name": "AWS Certified Solutions Architect - Associate",
        "issuer": "aws",
        "date": "2022",
        "expiry_date": "",
        "id": "AWS-SAA-12345",
        "confidence": 1.0
      }
    ],
    "projects": [
      {
        "name": "Distributed Task Scheduler",
        "description": "Open-source job scheduling library",
        "technologies": ["go", "redis", "docker"],
        "url": "https://github.com/janedoe/scheduler",
        "date": "",
        "confidence": 1.0
      }
    ],
    "summary": "Experienced software engineer with 6+ years building scalable backend systems.",
    "overall_confidence": 0.87,
    "sections_found": ["summary", "experience", "education", "skills", "certifications", "projects"],
    "warnings": [],
    "raw_text": null
  }
}
```

**Error Response:**

```json
{
  "success": false,
  "error": {
    "code": "INVALID_FORMAT",
    "message": "file does not appear to be a valid PDF",
    "section": ""
  }
}
```

---

### GET `/api/v1/health`

Health check endpoint.

#### Example Request

```bash
curl http://localhost:8080/api/v1/health
```

#### Response (200 OK)

```json
{
  "status": "ok",
  "version": "1.0.0",
  "time": "2024-01-15T10:30:00Z"
}
```

---

## Response Schema

### `ParsedResume`

| Field | Type | Description |
|-------|------|-------------|
| `parsed_at` | string (ISO 8601) | Timestamp when the resume was parsed |
| `source_file` | string | Original filename |
| `file_type` | string | `"pdf"` or `"docx"` |
| `parser_version` | string | Parser version string |
| `personal_info` | `PersonalInfo` | Extracted contact information |
| `work_experience` | `[]WorkExperience` | List of work experience entries |
| `education` | `[]Education` | List of education entries |
| `skills` | `[]Skill` | List of extracted skills |
| `certifications` | `[]Certification` | List of certifications |
| `projects` | `[]Project` | List of projects/achievements |
| `summary` | string | Professional summary/objective |
| `overall_confidence` | float (0.0–1.0) | Weighted average confidence score |
| `sections_found` | `[]string` | List of section types detected |
| `warnings` | `[]string` | Warnings about missing or low-confidence data |
| `raw_text` | string | Raw extracted text (only if `include_raw=true`) |

### `PersonalInfo`

| Field | Type | Description |
|-------|------|-------------|
| `name` | string | Full name |
| `email` | string | Email address |
| `phone` | string | Phone number |
| `location` | string | City, State/Country |
| `linkedin` | string | LinkedIn profile URL |
| `github` | string | GitHub profile URL |
| `website` | string | Personal website URL |
| `confidence` | float (0.0–1.0) | Extraction confidence |

### `WorkExperience`

| Field | Type | Description |
|-------|------|-------------|
| `company` | string | Company name |
| `title` | string | Job title |
| `start_date` | string | Start date (e.g., "January 2020") |
| `end_date` | string | End date or "Present" |
| `is_current` | boolean | Whether this is the current position |
| `location` | string | Job location (if available) |
| `responsibilities` | `[]string` | List of responsibilities/achievements |
| `confidence` | float (0.0–1.0) | Extraction confidence |

### `Education`

| Field | Type | Description |
|-------|------|-------------|
| `institution` | string | School/university name |
| `degree` | string | Degree type (e.g., "Bachelor of Science") |
| `field` | string | Field of study |
| `start_date` | string | Start year |
| `end_date` | string | Graduation year |
| `gpa` | string | GPA (if available) |
| `honors` | string | Honors (e.g., "Magna Cum Laude") |
| `confidence` | float (0.0–1.0) | Extraction confidence |

### `Skill`

| Field | Type | Description |
|-------|------|-------------|
| `name` | string | Skill name |
| `category` | string | `"technical"`, `"soft"`, or `"other"` |
| `confidence` | float (0.0–1.0) | Extraction confidence |

### `Certification`

| Field | Type | Description |
|-------|------|-------------|
| `name` | string | Certification name |
| `issuer` | string | Issuing organization |
| `date` | string | Issue date |
| `expiry_date` | string | Expiry date (if available) |
| `id` | string | Credential ID (if available) |
| `confidence` | float (0.0–1.0) | Extraction confidence |

### `Project`

| Field | Type | Description |
|-------|------|-------------|
| `name` | string | Project name |
| `description` | string | Project description |
| `technologies` | `[]string` | Technologies used |
| `url` | string | Project URL (if available) |
| `date` | string | Project date (if available) |
| `confidence` | float (0.0–1.0) | Extraction confidence |

---

## Error Codes

| Code | HTTP Status | Description |
|------|-------------|-------------|
| `EMPTY_REQUEST` | 400 | File content is empty |
| `EMPTY_FILE` | 400 | Uploaded file is empty |
| `INVALID_FORMAT` | 400 | File is not a valid PDF or DOCX |
| `UNSUPPORTED_FORMAT` | 400 | File type is not supported |
| `MISSING_FILE` | 400 | The `resume` form field is missing |
| `INVALID_REQUEST` | 400 | Malformed multipart request |
| `NO_TEXT_CONTENT` | 422 | Document contains no extractable text (e.g., image-based PDF) |
| `PDF_PARSE_ERROR` | 500 | Internal error parsing PDF |
| `DOCX_PARSE_ERROR` | 500 | Internal error parsing DOCX |
| `PARSE_ERROR` | 500 | General parsing error |
| `INTERNAL_ERROR` | 500 | Unexpected server error |

---

## Confidence Scores

Each extracted field includes a `confidence` score between 0.0 and 1.0:

| Range | Interpretation |
|-------|----------------|
| 0.9–1.0 | High confidence — field clearly identified |
| 0.7–0.9 | Medium confidence — field likely correct |
| 0.5–0.7 | Low confidence — field may be inaccurate |
| < 0.5 | Very low confidence — treat with caution |

The `overall_confidence` is a weighted average across all extracted sections.

---

## Running the Server

```bash
# Build and run
cd resume-parser
go build -o resume-parser ./cmd/server
./resume-parser -addr :8080

# Or run directly
go run ./cmd/server -addr :8080
```

**Flags:**

| Flag | Default | Description |
|------|---------|-------------|
| `-addr` | `:8080` | HTTP server listen address |

---

## Limitations

1. **Image-based PDFs**: PDFs that contain scanned images without embedded text cannot be parsed. The API returns a `NO_TEXT_CONTENT` error.
2. **Complex layouts**: Multi-column or heavily formatted resumes may have reduced extraction accuracy.
3. **Non-English resumes**: The parser is optimized for English-language resumes.
4. **File size**: Maximum upload size is 10 MB.
5. **Processing time**: Target is under 5 seconds per resume; complex documents may take longer.
