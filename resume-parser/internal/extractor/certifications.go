package extractor

import (
	"regexp"
	"strings"

	"github.com/learnbot/resume-parser/internal/schema"
)

var (
	// Certification ID patterns: "Credential ID: ABC123", "License #: 12345"
	certIDRe = regexp.MustCompile(`(?i)(?:credential\s+id|license\s+(?:no|#|number)|cert(?:ificate)?\s+(?:id|no|#))[\s:]+([A-Z0-9\-]+)`)

	// Expiry date patterns
	expiryRe = regexp.MustCompile(`(?i)(?:expires?|expiry|valid\s+(?:until|through))[\s:]+([A-Za-z]+\s+\d{4}|\d{1,2}/\d{4}|\d{4})`)

	// Known certification issuers
	knownIssuers = []string{
		"aws", "amazon", "google", "microsoft", "azure", "oracle",
		"cisco", "comptia", "pmi", "isaca", "isc2", "ec-council",
		"red hat", "linux foundation", "cncf", "hashicorp",
		"salesforce", "servicenow", "databricks", "snowflake",
		"coursera", "udemy", "edx", "pluralsight",
	}
)

// ExtractCertifications parses certification entries from the certifications section text.
func ExtractCertifications(text string) []schema.Certification {
	if strings.TrimSpace(text) == "" {
		return nil
	}

	lines := strings.Split(text, "\n")
	var certs []schema.Certification
	var current []string

	flush := func() {
		if len(current) == 0 {
			return
		}
		block := strings.TrimSpace(strings.Join(current, "\n"))
		if block != "" {
			if cert := parseCertBlock(block); cert != nil {
				certs = append(certs, *cert)
			}
		}
		current = nil
	}

	for _, line := range lines {
		trimmed := strings.TrimSpace(line)
		if trimmed == "" {
			flush()
			continue
		}
		// Bullet points start a new cert
		if isBulletLine(trimmed) && len(current) > 0 {
			flush()
		}
		current = append(current, trimmed)
	}
	flush()

	return certs
}

// parseCertBlock extracts a single certification from a text block.
func parseCertBlock(block string) *schema.Certification {
	cert := &schema.Certification{}
	lines := strings.Split(block, "\n")

	for i, line := range lines {
		line = strings.TrimSpace(line)
		line = strings.TrimLeft(line, "•-*>·")
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}

		// Credential ID
		if m := certIDRe.FindStringSubmatch(line); len(m) > 1 {
			cert.ID = m[1]
			continue
		}

		// Expiry date
		if m := expiryRe.FindStringSubmatch(line); len(m) > 1 {
			cert.ExpiryDate = m[1]
			continue
		}

		// Date range or single date
		if m := dateRangeRe.FindStringSubmatch(line); len(m) >= 3 {
			cert.Date = strings.TrimSpace(m[1])
			continue
		}
		if years := yearRe.FindAllString(line, -1); len(years) > 0 && cert.Date == "" {
			cert.Date = years[len(years)-1]
		}

		// First line is the cert name
		if i == 0 || cert.Name == "" {
			cert.Name = line
			// Try to extract issuer from name line
			cert.Issuer = extractIssuer(line)
			continue
		}

		// Subsequent lines may be issuer
		if cert.Issuer == "" {
			cert.Issuer = extractIssuer(line)
			if cert.Issuer == "" {
				cert.Issuer = line
			}
		}
	}

	if cert.Name == "" {
		return nil
	}

	// Confidence
	score := 0.0
	total := 2.0
	if cert.Name != "" {
		score += 1.0
	}
	if cert.Issuer != "" || cert.Date != "" {
		score += 1.0
	}
	cert.Confidence = schema.ConfidenceScore(score / total)

	return cert
}

// extractIssuer attempts to find a known issuer in the text.
func extractIssuer(s string) string {
	lower := strings.ToLower(s)
	for _, issuer := range knownIssuers {
		if strings.Contains(lower, issuer) {
			// Return the properly-cased version
			idx := strings.Index(lower, issuer)
			return strings.TrimSpace(s[idx : idx+len(issuer)])
		}
	}
	return ""
}

// isBulletLine returns true if the line starts with a bullet character.
func isBulletLine(s string) bool {
	if len(s) == 0 {
		return false
	}
	// Check ASCII bullet chars
	switch s[0] {
	case '-', '*', '>':
		return true
	}
	// Check Unicode bullet chars using rune
	r := []rune(s)[0]
	switch r {
	case '•', '·', '‣', '◦', '▪', '▸':
		return true
	}
	return false
}
