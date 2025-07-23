package test


import (
	"os"
	"path/filepath"
	"testing"

	"github.com/freshpaint/hipaa-tracker/internal/sanitize"
)

func TestLoadSensitiveFields_Success(t *testing.T) {
	// Create a temporary YAML file
	yamlContent := `
		sensitive_fields:
		- Email
		- PHONE
		- zipCode
	`

	tmpFile := filepath.Join(os.TempDir(), "test_sensitive_fields.yaml")
	if err := os.WriteFile(tmpFile, []byte(yamlContent), 0644); err != nil {
		t.Fatalf("Failed to write temp YAML: %v", err)
	}
	defer os.Remove(tmpFile) // Clean up

	if err := sanitize.LoadSensitiveFields(tmpFile); err != nil {
		t.Fatalf("LoadSensitiveFields failed: %v", err)
	}

	testKeys := []string{"email", "Phone", "Zipcode", "ZIP_CODE", "phoNe"}
	for _, key := range testKeys {
		if !isSensitiveHelper(key) {
			t.Errorf("Expected key '%s' to be detected as sensitive", key)
		}
	}
}

func TestLoadSensitiveFields_InvalidFile(t *testing.T) {
	err := sanitize.LoadSensitiveFields("non_existent.yaml")
	if err == nil {
		t.Error("Expected error for missing file, got nil")
	}
}

func TestLoadSensitiveFields_EmptyConfig(t *testing.T) {
	yamlContent := `sensitive_fields: []`
	tmpFile := filepath.Join(os.TempDir(), "test_empty_fields.yaml")
	if err := os.WriteFile(tmpFile, []byte(yamlContent), 0644); err != nil {
		t.Fatalf("Failed to write temp YAML: %v", err)
	}
	defer os.Remove(tmpFile)

	err := sanitize.LoadSensitiveFields(tmpFile)
	if err == nil {
		t.Error("Expected error for empty sensitive_fields, got nil")
	}
}

// local copy of normalize-based check for testing
func isSensitiveHelper(key string) bool {
	// re-use normalize logic inside sanitize
	sanitized, _ := sanitize.Sanitize(map[string]any{key: "value"})
	_, found := sanitized[key]
	return !found
}
