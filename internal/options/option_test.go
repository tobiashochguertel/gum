package options

import (
	"reflect"
	"testing"
)

func TestParseJSONOptions(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		want    []Option
		wantErr bool
	}{
		{
			name:  "JSON array format",
			input: `[{"value":"opt1","label":"Option 1"},{"value":"opt2","label":"Option 2","config":{"disabled":true}}]`,
			want: []Option{
				{Value: "opt1", Label: "Option 1"},
				{Value: "opt2", Label: "Option 2", Config: Config{Disabled: true}},
			},
			wantErr: false,
		},
		{
			name: "JSONL format",
			input: `{"value":"opt1","label":"Option 1"}
{"value":"opt2","label":"Option 2","config":{"disabled":true}}`,
			want: []Option{
				{Value: "opt1", Label: "Option 1"},
				{Value: "opt2", Label: "Option 2", Config: Config{Disabled: true}},
			},
			wantErr: false,
		},
		{
			name: "JSONL with blank lines",
			input: `{"value":"opt1"}

{"value":"opt2"}`,
			want: []Option{
				{Value: "opt1"},
				{Value: "opt2"},
			},
			wantErr: false,
		},
		{
			name:    "empty input",
			input:   "",
			want:    nil,
			wantErr: true,
		},
		{
			name:    "invalid JSON",
			input:   `{invalid}`,
			want:    nil,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ParseJSONOptions(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseJSONOptions() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ParseJSONOptions() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestParseConfigDelimiter(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		want    Config
		wantErr bool
	}{
		{
			name:    "disabled true",
			input:   `{"disabled":true}`,
			want:    Config{Disabled: true},
			wantErr: false,
		},
		{
			name:    "disabled false",
			input:   `{"disabled":false}`,
			want:    Config{Disabled: false},
			wantErr: false,
		},
		{
			name:    "empty config",
			input:   "",
			want:    Config{},
			wantErr: false,
		},
		{
			name:    "invalid JSON",
			input:   `{invalid}`,
			want:    Config{},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ParseConfigDelimiter(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseConfigDelimiter() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ParseConfigDelimiter() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestParseOptionsWithDelimiters(t *testing.T) {
	tests := []struct {
		name            string
		optionStr       string
		labelDelimiter  string
		configDelimiter string
		want            Option
		wantErr         bool
	}{
		{
			name:            "label and value with config",
			optionStr:       `Option 1:opt1::{"disabled":true}`,
			labelDelimiter:  ":",
			configDelimiter: "::",
			want:            Option{Value: "opt1", Label: "Option 1", Config: Config{Disabled: true}},
			wantErr:         false,
		},
		{
			name:            "value with config only",
			optionStr:       `opt1::{"disabled":true}`,
			labelDelimiter:  "",
			configDelimiter: "::",
			want:            Option{Value: "opt1", Config: Config{Disabled: true}},
			wantErr:         false,
		},
		{
			name:            "label and value only",
			optionStr:       "Option 1:opt1",
			labelDelimiter:  ":",
			configDelimiter: "",
			want:            Option{Value: "opt1", Label: "Option 1"},
			wantErr:         false,
		},
		{
			name:            "value only",
			optionStr:       "opt1",
			labelDelimiter:  "",
			configDelimiter: "",
			want:            Option{Value: "opt1"},
			wantErr:         false,
		},
		{
			name:            "value without label delimiter match",
			optionStr:       "opt1",
			labelDelimiter:  ":",
			configDelimiter: "",
			want:            Option{Value: "opt1"},
			wantErr:         false,
		},
		{
			name:            "invalid config JSON",
			optionStr:       `opt1::{invalid}`,
			labelDelimiter:  "",
			configDelimiter: "::",
			want:            Option{},
			wantErr:         true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ParseOptionsWithDelimiters(tt.optionStr, tt.labelDelimiter, tt.configDelimiter)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseOptionsWithDelimiters() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ParseOptionsWithDelimiters() = %v, want %v", got, tt.want)
			}
		})
	}
}
