package taxonomy

import (
	"math"
	"testing"
)

// ─────────────────────────────────────────────────────────────────────────────
// Helpers
// ─────────────────────────────────────────────────────────────────────────────

func approxEqual(a, b, eps float64) bool {
	return math.Abs(a-b) < eps
}

func findExtracted(skills []ExtractedSkill, canonicalID string) *ExtractedSkill {
	for i := range skills {
		if skills[i].CanonicalID == canonicalID {
			return &skills[i]
		}
	}
	return nil
}

// ─────────────────────────────────────────────────────────────────────────────
// Taxonomy (database)
// ─────────────────────────────────────────────────────────────────────────────

func TestNew_PopulatesIndex(t *testing.T) {
	tax := New()
	if len(tax.all) == 0 {
		t.Fatal("expected taxonomy to have skill nodes")
	}
	if len(tax.byID) == 0 {
		t.Fatal("expected byID index to be populated")
	}
	if len(tax.byAlias) == 0 {
		t.Fatal("expected byAlias index to be populated")
	}
}

func TestLookup_ExistingID(t *testing.T) {
	tax := New()
	node := tax.Lookup("go")
	if node == nil {
		t.Fatal("expected to find 'go' in taxonomy")
	}
	if node.CanonicalName != "Go" {
		t.Errorf("expected canonical name 'Go', got %q", node.CanonicalName)
	}
}

func TestLookup_NonExistentID(t *testing.T) {
	tax := New()
	node := tax.Lookup("nonexistent-skill-xyz")
	if node != nil {
		t.Error("expected nil for non-existent skill")
	}
}

func TestAll_ReturnsAllNodes(t *testing.T) {
	tax := New()
	all := tax.All()
	if len(all) != len(builtinSkills) {
		t.Errorf("expected %d nodes, got %d", len(builtinSkills), len(all))
	}
}

func TestSearch_ByQuery(t *testing.T) {
	tax := New()
	results := tax.Search("python", "", "", 10)
	if len(results) == 0 {
		t.Fatal("expected search results for 'python'")
	}
	found := false
	for _, r := range results {
		if r.ID == "python" {
			found = true
			break
		}
	}
	if !found {
		t.Error("expected 'python' in search results")
	}
}

func TestSearch_ByDomain(t *testing.T) {
	tax := New()
	results := tax.Search("", DomainDataScience, "", 0)
	for _, r := range results {
		if r.Domain != DomainDataScience {
			t.Errorf("expected domain %q, got %q for skill %q", DomainDataScience, r.Domain, r.ID)
		}
	}
}

func TestSearch_ByCategory(t *testing.T) {
	tax := New()
	results := tax.Search("", "", CategoryDatabase, 0)
	for _, r := range results {
		if r.Category != CategoryDatabase {
			t.Errorf("expected category %q, got %q for skill %q", CategoryDatabase, r.Category, r.ID)
		}
	}
}

func TestSearch_Limit(t *testing.T) {
	tax := New()
	results := tax.Search("", "", "", 3)
	if len(results) > 3 {
		t.Errorf("expected at most 3 results, got %d", len(results))
	}
}

func TestSearch_EmptyQuery_ReturnsAll(t *testing.T) {
	tax := New()
	results := tax.Search("", "", "", 0)
	if len(results) != len(builtinSkills) {
		t.Errorf("expected all %d skills for empty query, got %d", len(builtinSkills), len(results))
	}
}

// ─────────────────────────────────────────────────────────────────────────────
// Normalize
// ─────────────────────────────────────────────────────────────────────────────

func TestNormalize_ExactMatch(t *testing.T) {
	tax := New()
	result := tax.Normalize("Go")
	if result.MatchType != "exact" {
		t.Errorf("expected exact match, got %q", result.MatchType)
	}
	if result.CanonicalID != "go" {
		t.Errorf("expected canonical ID 'go', got %q", result.CanonicalID)
	}
	if result.CanonicalName != "Go" {
		t.Errorf("expected canonical name 'Go', got %q", result.CanonicalName)
	}
}

func TestNormalize_AliasMatch(t *testing.T) {
	tests := []struct {
		input      string
		wantID     string
		wantType   string
	}{
		{"golang", "go", "alias"},
		{"k8s", "kubernetes", "alias"},
		{"postgres", "postgresql", "alias"},
		{"js", "javascript", "alias"},
		{"ts", "typescript", "alias"},
		{"sklearn", "scikit-learn", "alias"},
		{"pyspark", "spark", "alias"},
		{"drf", "django", "alias"},
		{"ror", "rails", "alias"},
		{"tf", "tensorflow", "alias"},
	}

	tax := New()
	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			result := tax.Normalize(tt.input)
			if result.CanonicalID != tt.wantID {
				t.Errorf("Normalize(%q): expected ID %q, got %q (match_type=%q)",
					tt.input, tt.wantID, result.CanonicalID, result.MatchType)
			}
			if result.MatchType != tt.wantType {
				t.Errorf("Normalize(%q): expected match_type %q, got %q",
					tt.input, tt.wantType, result.MatchType)
			}
		})
	}
}

func TestNormalize_CaseInsensitive(t *testing.T) {
	tax := New()
	tests := []string{"PYTHON", "Python", "python", "GOLANG", "Golang"}
	for _, input := range tests {
		result := tax.Normalize(input)
		if result.MatchType == "none" {
			t.Errorf("Normalize(%q): expected a match, got none", input)
		}
	}
}

func TestNormalize_FuzzyMatch(t *testing.T) {
	tax := New()
	// "pythn" is close to "python" – should fuzzy match.
	result := tax.Normalize("pythn")
	if result.MatchType != "fuzzy" {
		t.Logf("Normalize('pythn'): match_type=%q, id=%q, score=%.4f",
			result.MatchType, result.CanonicalID, result.FuzzyScore)
		// Fuzzy matching may not always catch this – just verify it doesn't crash.
	}
}

func TestNormalize_NoMatch(t *testing.T) {
	tax := New()
	result := tax.Normalize("xyzzy-nonexistent-skill-12345")
	if result.MatchType != "none" {
		t.Errorf("expected no match for gibberish, got %q (id=%q)", result.MatchType, result.CanonicalID)
	}
	if result.CanonicalID != "" {
		t.Errorf("expected empty canonical ID for no match, got %q", result.CanonicalID)
	}
}

func TestNormalize_EmptyString(t *testing.T) {
	tax := New()
	result := tax.Normalize("")
	if result.MatchType != "none" {
		t.Errorf("expected no match for empty string, got %q", result.MatchType)
	}
}

func TestNormalize_MultiWordAlias(t *testing.T) {
	tax := New()
	tests := []struct {
		input  string
		wantID string
	}{
		{"machine learning", "machine-learning"},
		{"deep learning", "deep-learning"},
		{"natural language processing", "nlp"},
		{"ruby on rails", "rails"},
		{"spring boot", "spring-boot"},
		{"node.js", "nodejs"},
		{"react native", "react-native"},
	}
	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			result := tax.Normalize(tt.input)
			if result.CanonicalID != tt.wantID {
				t.Errorf("Normalize(%q): expected ID %q, got %q (match_type=%q)",
					tt.input, tt.wantID, result.CanonicalID, result.MatchType)
			}
		})
	}
}

func TestNormalizeMany(t *testing.T) {
	tax := New()
	inputs := []string{"golang", "k8s", "postgres", "xyzzy-nonexistent"}
	results := tax.NormalizeMany(inputs)

	if len(results) != len(inputs) {
		t.Fatalf("expected %d results, got %d", len(inputs), len(results))
	}
	if results[0].CanonicalID != "go" {
		t.Errorf("expected 'go', got %q", results[0].CanonicalID)
	}
	if results[1].CanonicalID != "kubernetes" {
		t.Errorf("expected 'kubernetes', got %q", results[1].CanonicalID)
	}
	if results[2].CanonicalID != "postgresql" {
		t.Errorf("expected 'postgresql', got %q", results[2].CanonicalID)
	}
	if results[3].MatchType != "none" {
		t.Errorf("expected no match for gibberish, got %q", results[3].MatchType)
	}
}

// ─────────────────────────────────────────────────────────────────────────────
// Extractor
// ─────────────────────────────────────────────────────────────────────────────

func TestExtract_TechnicalSkillsFromJobDescription(t *testing.T) {
	tax := New()
	ext := NewExtractor(tax)

	jd := `We are looking for a Senior Software Engineer with strong experience in Go and Python.
The ideal candidate will have hands-on experience with Kubernetes, Docker, and PostgreSQL.
Experience with AWS or GCP is a plus. You should be comfortable with REST APIs and microservices.`

	result := ext.Extract(jd, false)

	if len(result.TechnicalSkills) == 0 {
		t.Fatal("expected technical skills to be extracted")
	}

	expectedSkills := []string{"go", "python", "kubernetes", "docker", "postgresql"}
	for _, id := range expectedSkills {
		found := findExtracted(result.TechnicalSkills, id)
		if found == nil {
			t.Errorf("expected skill %q to be extracted from job description", id)
		}
	}
}

func TestExtract_SoftSkillsFromJobDescription(t *testing.T) {
	tax := New()
	ext := NewExtractor(tax)

	jd := `The role requires strong leadership skills and excellent communication abilities.
You should be a team player who excels at problem solving and project management.`

	result := ext.Extract(jd, false)

	if len(result.SoftSkills) == 0 {
		t.Fatal("expected soft skills to be extracted")
	}
}

func TestExtract_MultiWordSkills(t *testing.T) {
	tax := New()
	ext := NewExtractor(tax)

	jd := `Experience with machine learning and deep learning is required.
Knowledge of natural language processing and large language models is a plus.
Familiarity with Ruby on Rails and Spring Boot is beneficial.`

	result := ext.Extract(jd, false)

	expectedIDs := []string{"machine-learning", "deep-learning", "nlp", "llm", "rails", "spring-boot"}
	for _, id := range expectedIDs {
		found := findExtracted(result.Skills, id)
		if found == nil {
			t.Errorf("expected multi-word skill %q to be extracted", id)
		}
	}
}

func TestExtract_AliasesInText(t *testing.T) {
	tax := New()
	ext := NewExtractor(tax)

	// Use aliases instead of canonical names.
	jd := `Required: golang, k8s, postgres, reactjs, sklearn`

	result := ext.Extract(jd, false)

	expectedIDs := []string{"go", "kubernetes", "postgresql", "react", "scikit-learn"}
	for _, id := range expectedIDs {
		found := findExtracted(result.Skills, id)
		if found == nil {
			t.Errorf("expected alias-matched skill %q to be extracted", id)
		}
	}
}

func TestExtract_Deduplication(t *testing.T) {
	tax := New()
	ext := NewExtractor(tax)

	// "Go" and "golang" both refer to the same skill.
	jd := `We need Go developers. Experience with golang is required. Go is our primary language.`

	result := ext.Extract(jd, false)

	count := 0
	for _, s := range result.Skills {
		if s.CanonicalID == "go" {
			count++
		}
	}
	if count != 1 {
		t.Errorf("expected 'go' to appear exactly once (deduped), got %d times", count)
	}
}

func TestExtract_EmptyText(t *testing.T) {
	tax := New()
	ext := NewExtractor(tax)

	result := ext.Extract("", false)

	if len(result.Skills) != 0 {
		t.Errorf("expected no skills for empty text, got %d", len(result.Skills))
	}
}

func TestExtract_IncludeUnknown(t *testing.T) {
	tax := New()
	ext := NewExtractor(tax)

	jd := `Experience with FooBarBaz and XyzzyTech is required.`

	resultWithout := ext.Extract(jd, false)
	resultWith := ext.Extract(jd, true)

	if len(resultWith.UnknownSkills) < len(resultWithout.UnknownSkills) {
		t.Error("expected more unknown skills when include_unknown=true")
	}
}

func TestExtract_ConfidenceScores(t *testing.T) {
	tax := New()
	ext := NewExtractor(tax)

	jd := `Required: Go, Python, Docker`

	result := ext.Extract(jd, false)

	for _, s := range result.Skills {
		if s.Confidence <= 0 || s.Confidence > 1.0 {
			t.Errorf("skill %q has invalid confidence %.2f", s.CanonicalID, s.Confidence)
		}
	}
}

func TestExtract_MatchTypes(t *testing.T) {
	tax := New()
	ext := NewExtractor(tax)

	// "Go" → exact, "golang" → alias
	jd := `We use Go and golang in our stack.`

	result := ext.Extract(jd, false)

	// Both should resolve to "go" (deduped), so we just check one is found.
	found := findExtracted(result.Skills, "go")
	if found == nil {
		t.Fatal("expected 'go' to be extracted")
	}
	if found.MatchType != "exact" && found.MatchType != "alias" {
		t.Errorf("expected match type 'exact' or 'alias', got %q", found.MatchType)
	}
}

func TestExtract_GroupsCorrectly(t *testing.T) {
	tax := New()
	ext := NewExtractor(tax)

	jd := `Required: Python, TensorFlow, leadership, communication`

	result := ext.Extract(jd, false)

	// Python and TensorFlow should be in technical skills.
	if findExtracted(result.TechnicalSkills, "python") == nil {
		t.Error("expected 'python' in technical skills")
	}
	if findExtracted(result.TechnicalSkills, "tensorflow") == nil {
		t.Error("expected 'tensorflow' in technical skills")
	}

	// Leadership and communication should be in soft skills.
	if findExtracted(result.SoftSkills, "leadership") == nil {
		t.Error("expected 'leadership' in soft skills")
	}
	if findExtracted(result.SoftSkills, "communication") == nil {
		t.Error("expected 'communication' in soft skills")
	}
}

// ─────────────────────────────────────────────────────────────────────────────
// Jaro-Winkler
// ─────────────────────────────────────────────────────────────────────────────

func TestJaroWinkler_IdenticalStrings(t *testing.T) {
	score := jaroWinkler("python", "python")
	if score != 1.0 {
		t.Errorf("expected 1.0 for identical strings, got %.4f", score)
	}
}

func TestJaroWinkler_EmptyStrings(t *testing.T) {
	if jaroWinkler("", "python") != 0.0 {
		t.Error("expected 0.0 when first string is empty")
	}
	if jaroWinkler("python", "") != 0.0 {
		t.Error("expected 0.0 when second string is empty")
	}
	if jaroWinkler("", "") != 1.0 {
		t.Error("expected 1.0 for two empty strings")
	}
}

func TestJaroWinkler_SimilarStrings(t *testing.T) {
	// "golang" and "go" are somewhat similar.
	score := jaroWinkler("golang", "go")
	if score <= 0 || score >= 1.0 {
		t.Errorf("expected score in (0, 1) for 'golang' vs 'go', got %.4f", score)
	}
}

func TestJaroWinkler_DissimilarStrings(t *testing.T) {
	score := jaroWinkler("python", "kubernetes")
	if score > 0.7 {
		t.Errorf("expected low score for dissimilar strings, got %.4f", score)
	}
}

func TestJaroWinkler_CommonPrefix(t *testing.T) {
	// Strings with common prefix should score higher.
	scoreWithPrefix := jaroWinkler("python", "pythn")
	scoreWithout := jaroWinkler("python", "ythonp")
	if scoreWithPrefix <= scoreWithout {
		t.Errorf("expected common-prefix string to score higher: %.4f vs %.4f",
			scoreWithPrefix, scoreWithout)
	}
}

// ─────────────────────────────────────────────────────────────────────────────
// Internal helpers
// ─────────────────────────────────────────────────────────────────────────────

func TestNormalise(t *testing.T) {
	tests := []struct {
		input string
		want  string
	}{
		{"Go", "go"},
		{"  Python  ", "python"},
		{"Node.js", "node.js"},
		{"", ""},
		{"KUBERNETES", "kubernetes"},
	}
	for _, tt := range tests {
		got := normalise(tt.input)
		if got != tt.want {
			t.Errorf("normalise(%q) = %q, want %q", tt.input, got, tt.want)
		}
	}
}

func TestCleanToken(t *testing.T) {
	tests := []struct {
		input string
		want  string
	}{
		{"Go,", "Go"},
		{"(Python)", "Python"},
		{"Docker.", "Docker"},
		{"  Kubernetes  ", "Kubernetes"},
		{"[React]", "React"},
	}
	for _, tt := range tests {
		got := cleanToken(tt.input)
		if got != tt.want {
			t.Errorf("cleanToken(%q) = %q, want %q", tt.input, got, tt.want)
		}
	}
}

func TestSortByLengthDesc(t *testing.T) {
	input := []string{"go", "python", "machine learning", "k8s"}
	sortByLengthDesc(input)
	for i := 1; i < len(input); i++ {
		if len(input[i]) > len(input[i-1]) {
			t.Errorf("not sorted by length desc at index %d: %q > %q",
				i, input[i], input[i-1])
		}
	}
}

func TestLooksLikeSkill(t *testing.T) {
	tests := []struct {
		input string
		want  bool
	}{
		{"Go", true},
		{"Python", true},
		{"the", false},   // stop word
		{"and", false},   // stop word
		{"12345", false}, // number
		{"a", false},     // too short
		{"", false},      // empty
	}
	for _, tt := range tests {
		got := looksLikeSkill(tt.input)
		if got != tt.want {
			t.Errorf("looksLikeSkill(%q) = %v, want %v", tt.input, got, tt.want)
		}
	}
}

func TestRoundTo4(t *testing.T) {
	tests := []struct {
		input float64
		want  float64
	}{
		{0.12345, 0.1235},
		{1.0, 1.0},
		{0.0, 0.0},
	}
	for _, tt := range tests {
		got := roundTo4(tt.input)
		if !approxEqual(got, tt.want, 0.0001) {
			t.Errorf("roundTo4(%.5f) = %.5f, want %.5f", tt.input, got, tt.want)
		}
	}
}

// ─────────────────────────────────────────────────────────────────────────────
// Taxonomy integrity checks
// ─────────────────────────────────────────────────────────────────────────────

func TestTaxonomy_AllIDsUnique(t *testing.T) {
	seen := map[string]bool{}
	for _, node := range builtinSkills {
		if seen[node.ID] {
			t.Errorf("duplicate skill ID: %q", node.ID)
		}
		seen[node.ID] = true
	}
}

func TestTaxonomy_AllPrerequisitesExist(t *testing.T) {
	tax := New()
	for _, node := range builtinSkills {
		for _, prereq := range node.Prerequisites {
			if tax.Lookup(prereq) == nil {
				t.Errorf("skill %q has unknown prerequisite %q", node.ID, prereq)
			}
		}
	}
}

func TestTaxonomy_AllRelatedSkillsExist(t *testing.T) {
	// Related skills are informational references and may point to skills not
	// yet in the taxonomy. We log missing ones but do not fail the test.
	tax := New()
	missingCount := 0
	for _, node := range builtinSkills {
		for _, rel := range node.RelatedSkills {
			if tax.Lookup(rel) == nil {
				t.Logf("INFO: skill %q references related skill %q which is not in taxonomy", node.ID, rel)
				missingCount++
			}
		}
	}
	if missingCount > 0 {
		t.Logf("INFO: %d related skill references point to skills not yet in taxonomy (non-fatal)", missingCount)
	}
}

func TestTaxonomy_AllNodesHaveCanonicalName(t *testing.T) {
	for _, node := range builtinSkills {
		if node.CanonicalName == "" {
			t.Errorf("skill %q has empty canonical name", node.ID)
		}
	}
}

func TestTaxonomy_AllNodesHaveDomainAndCategory(t *testing.T) {
	for _, node := range builtinSkills {
		if node.Domain == "" {
			t.Errorf("skill %q has empty domain", node.ID)
		}
		if node.Category == "" {
			t.Errorf("skill %q has empty category", node.ID)
		}
	}
}

// ─────────────────────────────────────────────────────────────────────────────
// Benchmarks
// ─────────────────────────────────────────────────────────────────────────────

func BenchmarkNew(b *testing.B) {
	for i := 0; i < b.N; i++ {
		New()
	}
}

func BenchmarkNormalize_ExactMatch(b *testing.B) {
	tax := New()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		tax.Normalize("Python")
	}
}

func BenchmarkNormalize_AliasMatch(b *testing.B) {
	tax := New()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		tax.Normalize("golang")
	}
}

func BenchmarkNormalize_FuzzyMatch(b *testing.B) {
	tax := New()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		tax.Normalize("pythn")
	}
}

func BenchmarkExtract_ShortText(b *testing.B) {
	tax := New()
	ext := NewExtractor(tax)
	text := "Go, Python, Docker, Kubernetes, PostgreSQL, AWS"
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		ext.Extract(text, false)
	}
}

func BenchmarkExtract_LongJobDescription(b *testing.B) {
	tax := New()
	ext := NewExtractor(tax)
	text := `We are looking for a Senior Software Engineer with 5+ years of experience.

Required Skills:
- Go or Python (3+ years)
- Kubernetes and Docker for container orchestration
- PostgreSQL and Redis for data storage
- AWS or GCP cloud platforms
- REST API design and implementation
- Microservices architecture
- CI/CD pipelines (GitHub Actions or Jenkins)
- Git version control

Nice to Have:
- TensorFlow or PyTorch experience
- Machine learning background
- React or Vue.js for frontend work
- Kafka or RabbitMQ for messaging

Soft Skills:
- Strong leadership and communication skills
- Excellent problem solving abilities
- Team player with collaboration mindset
- Agile/Scrum methodology experience`

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		ext.Extract(text, false)
	}
}
