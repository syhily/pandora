package config

import (
	"fmt"
	"log"
	"os"
	"reflect"
	"strconv"
	"strings"

	"github.com/syhily/pandora/internal/common"
)

// promptTag represents the parsed prompt tag information
type promptTag struct {
	label        string
	required     bool
	defaultValue string
	validate     string
	optional     bool
}

// parsePromptTag parses the prompt tag string
func parsePromptTag(tag string) promptTag {
	pt := promptTag{}
	parts := strings.Split(tag, ",")
	for _, part := range parts {
		kv := strings.SplitN(strings.TrimSpace(part), "=", 2)
		if len(kv) != 2 {
			continue
		}
		key := strings.TrimSpace(kv[0])
		value := strings.TrimSpace(kv[1])

		switch key {
		case "label":
			pt.label = value
		case "required":
			pt.required = value == "true"
		case "default":
			pt.defaultValue = value
		case "validate":
			pt.validate = value
		case "optional":
			pt.optional = value == "true"
		}
	}
	return pt
}

// promptString prompts for a string value
func promptString(pt promptTag, currentValue string) string {
	label := pt.label
	if label == "" {
		return currentValue
	}

	// Build prompt message
	prompt := fmt.Sprintf("Please input %s", label)
	if pt.defaultValue != "" {
		prompt += fmt.Sprintf(" (Default: [%s])", pt.defaultValue)
	}
	if pt.optional {
		prompt += " [Optional]"
	}
	if pt.required {
		prompt += " [Required]"
	}
	prompt += ": "

	fmt.Print(prompt)

	var input string
	_, _ = fmt.Scanln(&input)

	// Use default if empty
	if input == "" && pt.defaultValue != "" {
		input = pt.defaultValue
	}

	// Validate if required
	if pt.required && input == "" {
		log.Fatalf("%s is required", label)
	}

	// Apply validation
	if input != "" && pt.validate != "" {
		if err := validateValue(input, pt.validate); err != nil {
			log.Fatalf("Validation failed for %s: %v", label, err)
		}
	}

	return input
}

// promptInt prompts for an int value
func promptInt(pt promptTag, currentValue int) int {
	label := pt.label
	if label == "" {
		return currentValue
	}

	// Build prompt message
	prompt := fmt.Sprintf("Please input %s", label)
	if pt.defaultValue != "" {
		prompt += fmt.Sprintf(" (Default: [%s])", pt.defaultValue)
	}
	if pt.optional {
		prompt += " [Optional]"
	}
	if pt.required {
		prompt += " [Required]"
	}
	prompt += ": "

	fmt.Print(prompt)

	var input string
	_, _ = fmt.Scanln(&input)

	// Use default if empty
	if input == "" && pt.defaultValue != "" {
		input = pt.defaultValue
	}

	// If empty and not required, return current value
	if input == "" {
		if pt.required {
			log.Fatalf("%s is required", label)
		}
		return currentValue
	}

	// Parse int
	value, err := strconv.Atoi(input)
	if err != nil {
		log.Fatalf("Invalid integer value for %s: %v", label, err)
	}

	// Apply validation
	if pt.validate != "" {
		if err := validateValue(fmt.Sprintf("%d", value), pt.validate); err != nil {
			log.Fatalf("Validation failed for %s: %v", label, err)
		}
	}

	return value
}

// validateValue applies validation based on the validate tag
func validateValue(value, validateType string) error {
	switch validateType {
	case "http_prefix":
		if !strings.HasPrefix(value, "http") {
			return fmt.Errorf("must start with 'http'")
		}
	case "http_suffix":
		if !strings.HasSuffix(value, "http") {
			return fmt.Errorf("must end with 'http'")
		}
	case "image_format":
		if _, ok := common.SupportExtensions[value]; !ok {
			return fmt.Errorf("unsupported format: %s. Supported: %v", value, common.SupportExtensions)
		}
	}
	return nil
}

// promptStruct recursively prompts for all fields in a struct
func promptStruct(v reflect.Value, t reflect.Type, prefix string) {
	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		fieldValue := v.Field(i)

		// Skip unexported fields
		if !fieldValue.CanSet() {
			continue
		}

		// Handle nested structs
		if field.Type.Kind() == reflect.Struct {
			// Add section header for nested structs
			fieldName := field.Name
			// Convert field name to a more readable format (e.g., ImageS3 -> Image S3)
			readableName := toReadableName(fieldName)
			if prefix != "" {
				fmt.Printf("\n=== %s %s ===\n", prefix, readableName)
			} else {
				fmt.Printf("\n=== %s ===\n", readableName)
			}
			// Pass the field name as prefix for nested prompts
			promptStruct(fieldValue, field.Type, readableName)
			continue
		}

		// Prompt for field with prefix context
		newValue := promptFieldWithPrefix(field, fieldValue, prefix)
		if newValue.IsValid() {
			fieldValue.Set(newValue)
		}
	}
}

// toReadableName converts a field name like "ImageS3" to "Image S3"
func toReadableName(name string) string {
	var result strings.Builder
	for i, r := range name {
		if i > 0 && r >= 'A' && r <= 'Z' {
			result.WriteString(" ")
		}
		result.WriteRune(r)
	}
	return result.String()
}

// promptFieldWithPrefix prompts for a field value with a prefix context
func promptFieldWithPrefix(field reflect.StructField, currentValue reflect.Value, prefix string) reflect.Value {
	tag := field.Tag.Get("prompt")
	if tag == "" {
		return currentValue
	}

	pt := parsePromptTag(tag)

	// Add prefix to label if prefix exists and label doesn't already contain it
	if prefix != "" && !strings.Contains(strings.ToLower(pt.label), strings.ToLower(prefix)) {
		// Extract the base label (e.g., "the s3 bucket" -> "bucket")
		baseLabel := strings.TrimPrefix(pt.label, "the ")
		baseLabel = strings.TrimPrefix(baseLabel, "s3 ")
		pt.label = fmt.Sprintf("the %s %s", strings.ToLower(prefix), baseLabel)
	}

	switch field.Type.Kind() {
	case reflect.String:
		var strValue string
		if currentValue.IsValid() && currentValue.CanInterface() {
			if s, ok := currentValue.Interface().(string); ok {
				strValue = s
			}
		}
		result := promptString(pt, strValue)
		return reflect.ValueOf(result)

	case reflect.Int:
		var intValue int
		if currentValue.IsValid() && currentValue.CanInterface() {
			if i, ok := currentValue.Interface().(int); ok {
				intValue = i
			}
		}
		result := promptInt(pt, intValue)
		return reflect.ValueOf(result)

	default:
		return currentValue
	}
}

// PromptConfig interactively prompts the user for configuration values using struct tags
func PromptConfig() *PandoraConfig {
	cfg := &PandoraConfig{}

	// Handle special case for ProjectRoot (use current directory as default)
	executeRoot, _ := os.Getwd()

	fmt.Println("=== Pandora Configuration ===")
	fmt.Println()

	// Use reflection to prompt for all fields
	v := reflect.ValueOf(cfg).Elem()
	t := reflect.TypeOf(cfg).Elem()

	promptStruct(v, t, "")

	// Post-processing: handle special cases
	// 1. ProjectRoot default to current directory if empty
	if cfg.ProjectRoot == "" || cfg.ProjectRoot == "." {
		cfg.ProjectRoot = executeRoot
	}

	// 2. ImageS3: region or endpoint must have at least one
	if cfg.ImageS3.Region == "" && cfg.ImageS3.Endpoint == "" {
		fmt.Println("\nNote: At least one of region or endpoint must be provided for ImageS3")
		fmt.Println("Please input the s3 region (Optional)")
		_, _ = fmt.Scanln(&cfg.ImageS3.Region)
		fmt.Println("Please input the s3 endpoint (Optional)")
		_, _ = fmt.Scanln(&cfg.ImageS3.Endpoint)
	}
	if cfg.ImageS3.Region == "" {
		cfg.ImageS3.Region = "auto"
	}

	// 3. MusicS3: region defaults to "auto" if empty
	if cfg.MusicS3.Region == "" {
		cfg.MusicS3.Region = "auto"
	}

	// 4. Convert format validation
	if cfg.Convert.DefaultFormat != "" {
		if _, ok := common.SupportExtensions[cfg.Convert.DefaultFormat]; !ok {
			log.Fatalf("Unsupported convert format: %s", cfg.Convert.DefaultFormat)
		}
	}

	return cfg
}
