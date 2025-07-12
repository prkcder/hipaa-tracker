package sanitize

import (
	"fmt"
	"os"
	"strings"

	"gopkg.in/yaml.v3"
)

// Config represents the structure of the YAML config
type Config struct {
	SensitiveFields []string `yaml:"sensitive_fields"`
}

var config Config

// LoadSensitiveFields loads the YAML config from disk
func LoadSensitiveFields(path string) error {

	data, err := os.ReadFile(path)
	
	if err != nil {
		return fmt.Errorf("failed to read YAML file: %w", err)
	}

	if err := yaml.Unmarshal(data, &config); err != nil {
		return fmt.Errorf("failed to parse YAML: %w", err)
	}

	if len(config.SensitiveFields) == 0 {
		return fmt.Errorf("no sensitive_fields defined in YAML")
	}

	return nil
}

// Sanitize removes keys from the payload that match sensitive fields
func Sanitize(payload map[string]any) (map[string]any, bool) {

	sanitized := false

	clean := make(map[string]any, len(payload))

	for key, value := range payload {
		if isSensitive(key) {
			sanitized = true
			continue // skip sensitive key
		}
		clean[key] = value
	}

	return clean, sanitized
}

// isSensitive performs case-insensitive matching against loaded keys
func isSensitive(key string) bool {

	for _, s := range config.SensitiveFields {
		if strings.EqualFold(s, key) {
			return true
		}
	}

	return false
}