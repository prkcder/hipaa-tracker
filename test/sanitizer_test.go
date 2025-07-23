package test


import (
	"reflect"
	"testing"

	"github.com/freshpaint/hipaa-tracker/internal/sanitize"
)

func TestSanitize_RemovesSensitiveFields(t *testing.T) {
	sanitizeConfig := []string{"password", "token", "email"}
	sanitize.InitTestConfig(sanitizeConfig)

	input := map[string]interface{}{
		"username": "john",
		"email":    "john@example.com",
		"token":    "abcd1234",
		"password": "secret",
		"age":      30,
	}

	expected := map[string]interface{}{
		"username": "john",
		"age":      30,
	}

	result, sanitized := sanitize.Sanitize(input)

	if !sanitized {
		t.Errorf("Expected sanitized = true, got false")
	}

	if !reflect.DeepEqual(result, expected) {
		t.Errorf("Expected %v, got %v", expected, result)
	}
}

func TestSanitize_NoSensitiveFieldsPresent(t *testing.T) {
	sanitizeConfig := []string{"ssn", "credit_card"}
	sanitize.InitTestConfig(sanitizeConfig)

	input := map[string]interface{}{
		"username": "alice",
		"age":      25,
	}

	expected := map[string]interface{}{
		"username": "alice",
		"age":      25,
	}

	result, sanitized := sanitize.Sanitize(input)

	if sanitized {
		t.Errorf("Expected sanitized = false, got true")
	}

	if !reflect.DeepEqual(result, expected) {
		t.Errorf("Expected %v, got %v", expected, result)
	}
}

func TestSanitize_CaseInsensitiveMatch(t *testing.T) {
	sanitizeConfig := []string{"Email", "TOKEN"}
	sanitize.InitTestConfig(sanitizeConfig)

	input := map[string]interface{}{
		"Email": "user@example.com",
		"token": "xyz789",
		"name":  "Test User",
	}

	expected := map[string]interface{}{
		"name": "Test User",
	}

	result, sanitized := sanitize.Sanitize(input)

	if !sanitized {
		t.Errorf("Expected sanitized = true, got false")
	}

	if !reflect.DeepEqual(result, expected) {
		t.Errorf("Expected %v, got %v", expected, result)
	}
}