# Experimental Features

This document describes the experimental features available in `gum` for the `choose` and `filter` commands. These features are designed to provide more control over options, including the ability to disable or inactivate specific items.

## Overview

The experimental features introduce:

1. **Disabled/Inactive Options**: Options can be marked as disabled, preventing them from being selected while still displaying them to users
2. **Config Delimiter**: Allows options to be defined with configuration metadata using a delimiter
3. **JSON/JSONL Input**: Enables defining options as JSON objects with separate value, label, and config fields

## Enabling Experimental Features

To enable experimental features, use the `--experimental` flag:

```bash
gum choose --experimental [other flags]
gum filter --experimental [other flags]
```

Or set the environment variable:

```bash
export GUM_EXPERIMENTAL=true
```

## Feature 1: Config Delimiter

The config delimiter allows you to attach configuration metadata to options using a delimiter (default: `::`)

### Syntax

```
option_text::{"disabled": true}
```

### Examples

#### Choose with Config Delimiter

```bash
gum choose --experimental \
  --config-delimiter="::" \
  "Available Option 1" \
  "Available Option 2" \
  "Disabled Option::{"disabled":true}" \
  "Available Option 3"
```

#### Combined with Label Delimiter

```bash
gum choose --experimental \
  --label-delimiter=":" \
  --config-delimiter="::" \
  "Display Text:actual_value::{"disabled":true}" \
  "Another Option:value2"
```

### Configuration Options

- `--config-delimiter`: Set the delimiter for config (default: `::`)
- Environment: `GUM_CHOOSE_CONFIG_DELIMITER` or `GUM_FILTER_CONFIG_DELIMITER`

### Config JSON Format

The configuration is a JSON object that supports:

```json
{
  "disabled": true
}
```

Currently supported config fields:
- `disabled` (boolean): Marks the option as disabled/inactive

## Feature 2: JSON/JSONL Input

Define options as JSON objects to have full control over value, label, and configuration.

### JSON Array Format

```json
[
  {
    "value": "actual_value",
    "label": "Display Text",
    "config": {
      "disabled": true
    }
  }
]
```

### JSONL (JSON Lines) Format

One JSON object per line:

```jsonl
{"value":"value1","label":"Option 1"}
{"value":"value2","label":"Option 2","config":{"disabled":true}}
{"value":"value3","label":"Option 3"}
```

### Examples

#### Choose with JSON Array

```bash
echo '[
  {"value":"opt1","label":"Available Option"},
  {"value":"opt2","label":"Disabled Option","config":{"disabled":true}},
  {"value":"opt3","label":"Another Option"}
]' | gum choose --experimental --option-as-json
```

#### Choose with JSONL

```bash
cat options.jsonl | gum choose --experimental --option-as-json
```

Where `options.jsonl` contains:

```jsonl
{"value":"item1","label":"First Item"}
{"value":"item2","label":"Second (Disabled)","config":{"disabled":true}}
{"value":"item3","label":"Third Item"}
```

#### Filter with JSON

```bash
cat <<EOF | gum filter --experimental --option-as-json
{"value":"file1.txt","label":"File 1"}
{"value":"file2.txt","label":"File 2 (Read-only)","config":{"disabled":true}}
{"value":"file3.txt","label":"File 3"}
EOF
```

### JSON Field Descriptions

- `value` (required): The actual value that will be output when the option is selected
- `label` (optional): The text displayed to the user. If omitted, `value` is used for display
- `config` (optional): Configuration object
  - `disabled` (boolean): If true, the option cannot be selected

## Feature 3: Disabled Items

Disabled items have the following characteristics:

- **Visual Indication**: Rendered with distinct styling (faded/dimmed by default)
- **Navigation**: Cursor automatically skips over disabled items
- **Selection**: Cannot be selected or toggled
- **Display**: Still visible to show users what options exist but are currently unavailable

### Styling Disabled Items

For `gum choose`:

```bash
gum choose --experimental \
  --disabled.foreground="240" \
  --disabled.faint=true \
  --disabled-prefix="✗ " \
  [options...]
```

For `gum filter`:

```bash
gum filter --experimental \
  --disabled-text.foreground="240" \
  --disabled-text.faint=true \
  --disabled-prefix=" - " \
  [options...]
```

Available style flags for disabled items:
- `--disabled.foreground` / `--disabled-text.foreground`: Text color
- `--disabled.background` / `--disabled-text.background`: Background color
- `--disabled.faint` / `--disabled-text.faint`: Use faint text
- `--disabled-prefix`: Prefix character(s) for disabled items

Environment variables:
- `GUM_CHOOSE_DISABLED_FOREGROUND`
- `GUM_CHOOSE_DISABLED_PREFIX`
- `GUM_FILTER_DISABLED_TEXT_FOREGROUND`
- `GUM_FILTER_DISABLED_PREFIX`

## Use Cases

### 1. Dynamic Shell Menus

Show all options but disable those that aren't currently applicable:

```bash
#!/bin/bash

# Get list of services with status
services=$(systemctl list-units --type=service --all --no-pager --no-legend | awk '{print $1,$3}')

# Build JSON options
echo "$services" | while read name state; do
  if [ "$state" = "running" ]; then
    echo "{\"value\":\"$name\",\"label\":\"$name (running)\"}"
  else
    echo "{\"value\":\"$name\",\"label\":\"$name (stopped)\",\"config\":{\"disabled\":true}}"
  fi
done | gum choose --experimental --option-as-json
```

### 2. Feature Flags in Development

```bash
#!/bin/bash

cat <<EOF | gum choose --experimental --option-as-json --no-limit
{"value":"feature-a","label":"Feature A (Stable)"}
{"value":"feature-b","label":"Feature B (Beta)","config":{"disabled":true}}
{"value":"feature-c","label":"Feature C (Stable)"}
{"value":"feature-d","label":"Feature D (Alpha)","config":{"disabled":true}}
EOF
```

### 3. Workflow State Management

```bash
#!/bin/bash

# Show tasks with their current state
cat <<EOF | gum filter --experimental --option-as-json
{"value":"task-1","label":"Write documentation"}
{"value":"task-2","label":"Review code (blocked)","config":{"disabled":true}}
{"value":"task-3","label":"Deploy to staging"}
{"value":"task-4","label":"Run tests (in progress)","config":{"disabled":true}}
EOF
```

## Backwards Compatibility

All experimental features are opt-in and do not affect existing functionality:

- Without `--experimental` flag, commands work exactly as before
- Config delimiter and JSON parsing are only active when experimental mode is enabled
- Disabled items are only recognized in experimental mode

## Future Enhancements

Planned additions to the config object:

- `selected`: Pre-select specific items
- `style`: Custom styling per option
- `group`: Grouping options together
- `description`: Additional descriptive text
- `shortcut`: Keyboard shortcuts for specific options

## Environment Variables

All experimental features can be configured via environment variables:

```bash
# Enable experimental features globally
export GUM_EXPERIMENTAL=true

# Configure delimiters
export GUM_CHOOSE_CONFIG_DELIMITER="::"
export GUM_FILTER_CONFIG_DELIMITER="::"

# Enable JSON parsing
export GUM_CHOOSE_OPTION_AS_JSON=true
export GUM_FILTER_OPTION_AS_JSON=true

# Style disabled items
export GUM_CHOOSE_DISABLED_FOREGROUND="240"
export GUM_CHOOSE_DISABLED_PREFIX="✗ "
export GUM_FILTER_DISABLED_PREFIX=" - "
```

## Troubleshooting

### Disabled items can still be selected

Make sure you've enabled experimental mode with `--experimental`.

### JSON parsing fails

- Check that your JSON is valid
- For JSON arrays, ensure the input starts with `[`
- For JSONL, ensure each line is a valid JSON object
- Check for proper escaping of quotes and special characters

### Config delimiter not working

- Ensure `--experimental` flag is set
- Verify the delimiter matches what you're using (default is `::`)
- Check that the config portion is valid JSON

### All options appear disabled

- Verify your JSON config syntax
- Check that `disabled` is a boolean, not a string: `{"disabled":true}` not `{"disabled":"true"}`

## Feedback

These features are experimental and subject to change. Please provide feedback on:
- API design and usability
- Additional config options needed
- Performance with large option lists
- Integration with other tools
