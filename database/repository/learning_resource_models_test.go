package repository

import (
	"testing"

	"github.com/google/uuid"
)

// ─────────────────────────────────────────────────────────────────────────────
// Enum constant tests
// ─────────────────────────────────────────────────────────────────────────────

func TestResourceTypeConstants(t *testing.T) {
	types := []ResourceType{
		ResourceTypeCourse,
		ResourceTypeCertification,
		ResourceTypeDocumentation,
		ResourceTypeVideo,
		ResourceTypeBook,
		ResourceTypePractice,
		ResourceTypeArticle,
		ResourceTypeProject,
		ResourceTypeOther,
	}

	for _, rt := range types {
		if string(rt) == "" {
			t.Errorf("ResourceType constant should not be empty")
		}
	}
}

func TestResourceDifficultyConstants(t *testing.T) {
	difficulties := []ResourceDifficulty{
		ResourceDifficultyBeginner,
		ResourceDifficultyIntermediate,
		ResourceDifficultyAdvanced,
		ResourceDifficultyExpert,
		ResourceDifficultyAllLevels,
	}

	for _, d := range difficulties {
		if string(d) == "" {
			t.Errorf("ResourceDifficulty constant should not be empty")
		}
	}
}

func TestResourceCostTypeConstants(t *testing.T) {
	costTypes := []ResourceCostType{
		ResourceCostFree,
		ResourceCostFreemium,
		ResourceCostPaid,
		ResourceCostSubscription,
		ResourceCostFreeAudit,
		ResourceCostEmployerSponsored,
	}

	for _, ct := range costTypes {
		if string(ct) == "" {
			t.Errorf("ResourceCostType constant should not be empty")
		}
	}
}

func TestUserResourceStatusConstants(t *testing.T) {
	statuses := []UserResourceStatus{
		UserResourceStatusSaved,
		UserResourceStatusInProgress,
		UserResourceStatusCompleted,
		UserResourceStatusAbandoned,
	}

	for _, s := range statuses {
		if string(s) == "" {
			t.Errorf("UserResourceStatus constant should not be empty")
		}
	}
}

// ─────────────────────────────────────────────────────────────────────────────
// ResourceQueryFilter tests
// ─────────────────────────────────────────────────────────────────────────────

func TestResourceQueryFilter_Defaults(t *testing.T) {
	// A zero-value filter should have sensible defaults.
	var filter ResourceQueryFilter

	if filter.Limit != 0 {
		t.Errorf("expected default Limit=0, got %d", filter.Limit)
	}
	if filter.Offset != 0 {
		t.Errorf("expected default Offset=0, got %d", filter.Offset)
	}
	if filter.IsFree {
		t.Error("expected IsFree=false by default")
	}
	if filter.HasCertificate {
		t.Error("expected HasCertificate=false by default")
	}
	if filter.HasHandsOn {
		t.Error("expected HasHandsOn=false by default")
	}
}

func TestResourceQueryFilter_SkillName(t *testing.T) {
	filter := ResourceQueryFilter{
		SkillName: "Python",
	}
	if filter.SkillName != "Python" {
		t.Errorf("expected SkillName=Python, got %s", filter.SkillName)
	}
}

func TestResourceQueryFilter_MultipleFilters(t *testing.T) {
	filter := ResourceQueryFilter{
		SkillName:      "Go",
		ResourceType:   ResourceTypeCourse,
		Difficulty:     ResourceDifficultyIntermediate,
		CostType:       ResourceCostFree,
		IsFree:         true,
		HasCertificate: true,
		HasHandsOn:     true,
		MinRating:      4.5,
		Limit:          20,
		Offset:         40,
	}

	if filter.SkillName != "Go" {
		t.Errorf("expected SkillName=Go, got %s", filter.SkillName)
	}
	if filter.ResourceType != ResourceTypeCourse {
		t.Errorf("expected ResourceType=course, got %s", filter.ResourceType)
	}
	if filter.Difficulty != ResourceDifficultyIntermediate {
		t.Errorf("expected Difficulty=intermediate, got %s", filter.Difficulty)
	}
	if !filter.IsFree {
		t.Error("expected IsFree=true")
	}
	if !filter.HasCertificate {
		t.Error("expected HasCertificate=true")
	}
	if !filter.HasHandsOn {
		t.Error("expected HasHandsOn=true")
	}
	if filter.MinRating != 4.5 {
		t.Errorf("expected MinRating=4.5, got %f", filter.MinRating)
	}
	if filter.Limit != 20 {
		t.Errorf("expected Limit=20, got %d", filter.Limit)
	}
	if filter.Offset != 40 {
		t.Errorf("expected Offset=40, got %d", filter.Offset)
	}
}

// ─────────────────────────────────────────────────────────────────────────────
// CreateResourceInput tests
// ─────────────────────────────────────────────────────────────────────────────

func TestCreateResourceInput_SkillsSlice(t *testing.T) {
	level := ResourceDifficultyIntermediate
	input := CreateResourceInput{
		Title:        "Test Course",
		Slug:         "test-course",
		URL:          "https://example.com",
		ResourceType: ResourceTypeCourse,
		Difficulty:   ResourceDifficultyIntermediate,
		CostType:     ResourceCostFree,
		Skills: []ResourceSkillInput{
			{SkillName: "Python", IsPrimary: true, CoverageLevel: &level},
			{SkillName: "Django", IsPrimary: false},
		},
	}

	if len(input.Skills) != 2 {
		t.Errorf("expected 2 skills, got %d", len(input.Skills))
	}
	if input.Skills[0].SkillName != "Python" {
		t.Errorf("expected first skill=Python, got %s", input.Skills[0].SkillName)
	}
	if !input.Skills[0].IsPrimary {
		t.Error("expected first skill to be primary")
	}
	if input.Skills[0].CoverageLevel == nil {
		t.Error("expected coverage level to be set")
	}
	if *input.Skills[0].CoverageLevel != ResourceDifficultyIntermediate {
		t.Errorf("expected coverage level=intermediate, got %s", *input.Skills[0].CoverageLevel)
	}
}

// ─────────────────────────────────────────────────────────────────────────────
// UpsertProgressInput tests
// ─────────────────────────────────────────────────────────────────────────────

func TestUpsertProgressInput_StatusValues(t *testing.T) {
	tests := []struct {
		status UserResourceStatus
		valid  bool
	}{
		{UserResourceStatusSaved, true},
		{UserResourceStatusInProgress, true},
		{UserResourceStatusCompleted, true},
		{UserResourceStatusAbandoned, true},
		{UserResourceStatus("invalid"), false},
	}

	validStatuses := map[UserResourceStatus]bool{
		UserResourceStatusSaved:      true,
		UserResourceStatusInProgress: true,
		UserResourceStatusCompleted:  true,
		UserResourceStatusAbandoned:  true,
	}

	for _, tt := range tests {
		isValid := validStatuses[tt.status]
		if isValid != tt.valid {
			t.Errorf("status %q: expected valid=%v, got valid=%v", tt.status, tt.valid, isValid)
		}
	}
}

// ─────────────────────────────────────────────────────────────────────────────
// LearningPathWithResources tests
// ─────────────────────────────────────────────────────────────────────────────

func TestLearningPathWithResources_ResourceCount(t *testing.T) {
	path := LearningPathWithResources{
		ResourceCount:         5,
		RequiredResourceCount: 3,
		Resources: []LearningPathStep{
			{StepOrder: 1, IsRequired: true},
			{StepOrder: 2, IsRequired: true},
			{StepOrder: 3, IsRequired: false},
		},
	}

	if path.ResourceCount != 5 {
		t.Errorf("expected ResourceCount=5, got %d", path.ResourceCount)
	}
	if path.RequiredResourceCount != 3 {
		t.Errorf("expected RequiredResourceCount=3, got %d", path.RequiredResourceCount)
	}
	if len(path.Resources) != 3 {
		t.Errorf("expected 3 resources, got %d", len(path.Resources))
	}
}

func TestLearningPathStep_Ordering(t *testing.T) {
	steps := []LearningPathStep{
		{StepOrder: 3, IsRequired: false},
		{StepOrder: 1, IsRequired: true},
		{StepOrder: 2, IsRequired: true},
	}

	// Verify step orders are distinct.
	seen := map[int16]bool{}
	for _, s := range steps {
		if seen[s.StepOrder] {
			t.Errorf("duplicate step order: %d", s.StepOrder)
		}
		seen[s.StepOrder] = true
	}
}

// ─────────────────────────────────────────────────────────────────────────────
// LearningResourceWithSkills tests
// ─────────────────────────────────────────────────────────────────────────────

func TestLearningResourceWithSkills_ProviderName(t *testing.T) {
	res := LearningResourceWithSkills{
		ProviderName: "Coursera",
	}
	if res.ProviderName != "Coursera" {
		t.Errorf("expected ProviderName=Coursera, got %s", res.ProviderName)
	}
}

// ─────────────────────────────────────────────────────────────────────────────
// CreateLearningPathInput tests
// ─────────────────────────────────────────────────────────────────────────────

func TestCreateLearningPathInput_Resources(t *testing.T) {
	resourceID := mustParseUUID("20000000-0000-0000-0000-000000000001")
	notes := "Start here"
	input := CreateLearningPathInput{
		Title:      "Python Path",
		Slug:       "python-path",
		Difficulty: ResourceDifficultyBeginner,
		Resources: []LearningPathResourceInput{
			{
				ResourceID: resourceID,
				StepOrder:  1,
				IsRequired: true,
				Notes:      &notes,
			},
		},
	}

	if len(input.Resources) != 1 {
		t.Errorf("expected 1 resource, got %d", len(input.Resources))
	}
	if input.Resources[0].StepOrder != 1 {
		t.Errorf("expected StepOrder=1, got %d", input.Resources[0].StepOrder)
	}
	if !input.Resources[0].IsRequired {
		t.Error("expected IsRequired=true")
	}
	if input.Resources[0].Notes == nil || *input.Resources[0].Notes != "Start here" {
		t.Error("expected Notes='Start here'")
	}
}

// ─────────────────────────────────────────────────────────────────────────────
// Helpers
// ─────────────────────────────────────────────────────────────────────────────

func mustParseUUID(s string) uuid.UUID {
	id, err := uuid.Parse(s)
	if err != nil {
		panic(err)
	}
	return id
}
