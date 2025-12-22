package config

import (
	"errors"
	"fmt"
	"reflect"
	"regexp"
	"slices"
	"strings"
)

// promptTag represents the parsed prompt tag information
type promptTag struct {
	label        string         // The label is the prompt message shown to the user
	required     bool           // Whether the field is required
	defaultValue string         // The default value for the field
	enums        []string       // Comma-separated list of valid values for the field
	regex        *regexp.Regexp // The regex pattern that the input field value must match
}

// parsePromptTag parses the prompt tag string into a promptTag struct
func parsePromptTag(tag string) (*promptTag, error) {
	pt := &promptTag{}
	for part := range strings.SplitSeq(tag, "|") {
		kv := strings.SplitN(strings.TrimSpace(part), ":", 2)
		if len(kv) != 2 {
			continue
		}
		key := strings.TrimSpace(kv[0])
		value := strings.TrimSpace(kv[1])
		switch key {
		case "label":
			pt.label = value
		case "required":
			pt.required = strings.ToLower(value) == "true" || value == "1"
		case "default":
			pt.defaultValue = value
		case "enums":
			pt.enums = strings.Split(value, ",")
			for i := range pt.enums {
				pt.enums[i] = strings.TrimSpace(pt.enums[i])
			}
		case "regex":
			reg, err := regexp.Compile(value)
			if err != nil {
				return nil, fmt.Errorf("invalid regex %s: %w", value, err)
			}
			pt.regex = reg
		default:
			return nil, errors.New("unknown prompt tag key: " + key)
		}
	}
	return pt, nil
}

// promptValue prompts for a string value
func promptValue(pt *promptTag, alternative string) (string, error) {
	// Build prompt message
	prompt := strings.Builder{}

	// If no label, use the alternative
	if pt.label == "" {
		prompt.WriteString(fmt.Sprintf("Please input %s", alternative))
	} else {
		prompt.WriteString(fmt.Sprintf("Please input %s", pt.label))
		if len(pt.enums) > 0 {
			prompt.WriteString(fmt.Sprintf(" <Options: %s>", strings.Join(pt.enums, ", ")))
		}
		if pt.defaultValue != "" {
			prompt.WriteString(fmt.Sprintf(" (Default: [%s])", pt.defaultValue))
		}
		if pt.required {
			prompt.WriteString(" [Required]")
		} else {
			prompt.WriteString(" [Optional]")
		}
		prompt.WriteString(": ")
	}
	fmt.Print(prompt.String())

	// Read input from user
	var input string
	_, _ = fmt.Scanln(&input)
	input = strings.TrimSpace(input)

	// Use default if empty
	if input == "" && pt.defaultValue != "" {
		input = pt.defaultValue
	}

	// Validate if required
	if (pt.required || len(pt.enums) > 0 || pt.regex != nil) && input == "" {
		return "", fmt.Errorf("%s is required", pt.label)
	}

	// Validate enums
	if len(pt.enums) > 0 {
		if !slices.Contains(pt.enums, input) {
			return "", fmt.Errorf("invalid value for %s, valid options are: %s", pt.label, strings.Join(pt.enums, ", "))
		}
	}

	// Validate the regex
	if pt.regex != nil {
		if !pt.regex.MatchString(input) {
			return "", fmt.Errorf("the input doesn't match the regex %s", pt.regex.String())
		}
	}

	return input, nil
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
				readableName = prefix + " " + readableName
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
