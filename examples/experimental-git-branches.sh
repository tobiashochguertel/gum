#!/bin/bash

# Example: Dynamic menu with disabled options based on state
# This demonstrates how to use experimental features for context-aware menus

echo "=== Git Branch Manager with Disabled Items ==="
echo ""
echo "Branches marked as 'current' or 'protected' are disabled"
echo ""

# Simulated branch data (in real scenario, you'd query git)
cat <<EOF | /tmp/gum choose --experimental --option-as-json --limit=1
{"value":"main","label":"main (protected)","config":{"disabled":true}}
{"value":"develop","label":"develop (current)","config":{"disabled":true}}
{"value":"feature/new-ui","label":"feature/new-ui"}
{"value":"feature/api-update","label":"feature/api-update"}
{"value":"hotfix/security-patch","label":"hotfix/security-patch"}
EOF

SELECTED=$?

if [ $SELECTED -eq 0 ]; then
    echo ""
    echo "You selected a branch! (In a real script, you'd check it out here)"
else
    echo ""
    echo "No branch selected or operation cancelled"
fi
