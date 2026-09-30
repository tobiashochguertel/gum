#!/bin/bash

# Example: File filter with permissions
# Shows how to use filter with disabled items for files with restricted access

echo "=== File Selector (some files are read-only) ==="
echo ""

cat <<EOF | /tmp/gum filter --experimental --option-as-json
{"value":"README.md","label":"README.md"}
{"value":"config.yaml","label":"config.yaml (read-only)","config":{"disabled":true}}
{"value":"main.go","label":"main.go"}
{"value":"system.conf","label":"system.conf (protected)","config":{"disabled":true}}
{"value":"script.sh","label":"script.sh"}
{"value":"data.json","label":"data.json"}
EOF

SELECTED=$?

if [ $SELECTED -eq 0 ]; then
    echo ""
    echo "File selected for editing!"
else
    echo ""
    echo "No file selected"
fi
