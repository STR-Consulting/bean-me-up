package cmd

import (
	"strings"
	"testing"
)

func TestGenerateConfig(t *testing.T) {
	data := configTemplateData{
		ListID:   "123456789",
		ListName: "My Test List",
		Statuses: []string{"to do", "in progress", "complete"},
		CustomFields: []fieldEntry{
			{Name: "Bean ID", Type: "text", ID: "abc-123"},
			{Name: "Due Date", Type: "date", ID: "def-456"},
		},
	}

	result, err := generateConfig(data)
	if err != nil {
		t.Fatalf("generateConfig() error = %v", err)
	}

	// Check for required elements
	checks := []struct {
		name     string
		contains string
	}{
		{"list_id", `list_id: "123456789"`},
		{"list name comment", "# List: My Test List"},
		{"status to do", `- "to do"`},
		{"status in progress", `- "in progress"`},
		{"status complete", `- "complete"`},
		{"field Bean ID", `- "Bean ID" (text): abc-123`},
		{"field Due Date", `- "Due Date" (date): def-456`},
		{"status_mapping comment", "# status_mapping:"},
		{"custom_fields comment", "# custom_fields:"},
	}

	for _, c := range checks {
		t.Run(c.name, func(t *testing.T) {
			if !strings.Contains(result, c.contains) {
				t.Errorf("generateConfig() output missing %q\nGot:\n%s", c.contains, result)
			}
		})
	}
}

func TestGenerateConfig_NoOptionalData(t *testing.T) {
	data := configTemplateData{
		ListID:   "999",
		ListName: "Minimal List",
		Statuses: []string{"open", "closed"},
	}

	result, err := generateConfig(data)
	if err != nil {
		t.Fatalf("generateConfig() error = %v", err)
	}

	// Should have list_id
	if !strings.Contains(result, `list_id: "999"`) {
		t.Error("missing list_id")
	}

	// Should have statuses
	if !strings.Contains(result, `- "open"`) {
		t.Error("missing status 'open'")
	}

	// Should NOT have custom fields section (since no custom fields)
	if strings.Contains(result, "Custom fields: map bean fields") {
		t.Error("should not have custom fields section when no fields provided")
	}
}
