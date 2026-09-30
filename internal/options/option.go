package options

import (
	"encoding/json"
	"fmt"
	"strings"
)

// Config represents the configuration for an option.
type Config struct {
	Disabled bool `json:"disabled,omitempty"`
}

// Option represents a single option with its value, label, and configuration.
type Option struct {
	Value  string `json:"value"`
	Label  string `json:"label,omitempty"`
	Config Config `json:"config,omitempty"`
}

// ParseJSONOptions parses options from JSON format (array or JSONL).
// It supports both JSON array format and JSON Lines format (one JSON object per line).
func ParseJSONOptions(input string) ([]Option, error) {
	input = strings.TrimSpace(input)
	if input == "" {
		return nil, fmt.Errorf("empty input")
	}

	var opts []Option

	// Try parsing as JSON array first
	if strings.HasPrefix(input, "[") {
		if err := json.Unmarshal([]byte(input), &opts); err != nil {
			return nil, fmt.Errorf("failed to parse JSON array: %w", err)
		}
		return validateOptions(opts)
	}

	// Parse as JSONL (one JSON object per line)
	lines := strings.Split(input, "\n")
	for i, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		var opt Option
		if err := json.Unmarshal([]byte(line), &opt); err != nil {
			return nil, fmt.Errorf("failed to parse JSON line %d: %w", i+1, err)
		}
		opts = append(opts, opt)
	}

	return validateOptions(opts)
}

func validateOptions(opts []Option) ([]Option, error) {
	if len(opts) == 0 {
		return nil, fmt.Errorf("at least one option is required")
	}
	for i, opt := range opts {
		if opt.Value == "" {
			return nil, fmt.Errorf("option %d has an empty value", i+1)
		}
	}
	return opts, nil
}

// ParseConfigDelimiter parses a config string using a delimiter.
// Format: {"disabled": true} or other JSON config
func ParseConfigDelimiter(configStr string) (Config, error) {
	var cfg Config
	configStr = strings.TrimSpace(configStr)
	if configStr == "" {
		return cfg, nil
	}

	if err := json.Unmarshal([]byte(configStr), &cfg); err != nil {
		return cfg, fmt.Errorf("failed to parse config JSON: %w", err)
	}

	return cfg, nil
}

// ParseOptionsWithDelimiters parses options using label and config delimiters.
// Supports formats:
// - "label:value::config"
// - "value::config" (if labelDelimiter is empty)
// - "label:value" (if configDelimiter is empty)
// - "value" (if both delimiters are empty)
func ParseOptionsWithDelimiters(optionStr, labelDelimiter, configDelimiter string) (Option, error) {
	var opt Option

	// If no delimiters, the entire string is the value
	if labelDelimiter == "" && configDelimiter == "" {
		opt.Value = optionStr
		return opt, nil
	}

	// Split by config delimiter first (if present)
	var mainPart, configPart string
	if configDelimiter != "" {
		parts := strings.SplitN(optionStr, configDelimiter, 2)
		mainPart = parts[0]
		if len(parts) > 1 {
			configPart = parts[1]
		}
	} else {
		mainPart = optionStr
	}

	// Split main part by label delimiter (if present)
	if labelDelimiter != "" {
		parts := strings.SplitN(mainPart, labelDelimiter, 2)
		if len(parts) == 2 {
			opt.Label = parts[0]
			opt.Value = parts[1]
		} else {
			opt.Value = parts[0]
		}
	} else {
		opt.Value = mainPart
	}

	// Parse config part (if present)
	if configPart != "" {
		cfg, err := ParseConfigDelimiter(configPart)
		if err != nil {
			return opt, err
		}
		opt.Config = cfg
	}

	return opt, nil
}
