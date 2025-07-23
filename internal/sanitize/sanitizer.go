package sanitize

import (
	"fmt"
	"os"
	"strings"
	"log"
	"regexp"


	"gopkg.in/yaml.v3"
)

// Config represents the structure of the YAML config
type Config struct {
	SensitiveFields []string `yaml:"sensitive_fields"`
}

var (
	config Config
	normalizedSet map[string]struct{}
	normalizerExpr = regexp.MustCompile(`[^a-zA-Z0-9]`)
	sanitizeMode = "remove"

)

func normalizeKey(s string) string {
	return normalizerExpr.ReplaceAllString(strings.ToLower(s), "")
}

//  loads the YAML config from disk
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

	normalizedSet = make(map[string]struct{}, len(config.SensitiveFields))

	for _, f := range config.SensitiveFields {
		normalizedSet[normalizeKey(f)] = struct{}{}
	}

	if mode := os.Getenv("SANITIZE_MODE"); mode == "redact" || mode == "remove" {
		sanitizeMode = mode
	}

	log.Printf("Loaded %d sensitive fields", len(config.SensitiveFields))

	return nil
}

// removes keys from the payload that match sensitive fields
func Sanitize(payload map[string]any) (map[string]any, bool) {

	sanitized := false

	clean := make(map[string]any, len(payload))

	for key, value := range payload {
		if isSensitive(key) {
			sanitized = true
			if sanitizeMode == "redact" {
				clean[key] = "[REDACTED]"
			}
			continue // skip sensitive key
		}
		clean[key] = value
	}

	return clean, sanitized
}

// performs case-insensitive matching against loaded keys
func isSensitive(key string) bool {

	// for _, s := range config.SensitiveFields {
	// 	if strings.EqualFold(s, key) {
	// 		return true
	// 	}
	// }

	// return false
	_, exists := normalizedSet[normalizeKey(key)]
	return exists
}

func ReloadSensitiveFields() error {
	return LoadSensitiveFields("sensitive_fields.yaml")
}


func InitTestConfig(fields []string) {
	config.SensitiveFields = fields

	normalizedSet = make(map[string]struct{}, len(fields))

	for _, f := range fields {
		normalizedSet[normalizeKey(f)] = struct{}{}
	}
}