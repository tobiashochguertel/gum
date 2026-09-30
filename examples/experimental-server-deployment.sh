#!/bin/bash

# Example: Server deployment menu with status indicators
# Shows how to use JSON format with disabled items for server management

echo "=== Server Deployment Menu ==="
echo ""
echo "Select servers to deploy to (disabled servers are offline or in maintenance)"
echo ""

cat <<EOF | gum choose --experimental --option-as-json --no-limit
{"value":"prod-web-01","label":"Production Web Server 01 (online)"}
{"value":"prod-web-02","label":"Production Web Server 02 (maintenance)","config":{"disabled":true}}
{"value":"prod-api-01","label":"Production API Server 01 (online)"}
{"value":"prod-api-02","label":"Production API Server 02 (offline)","config":{"disabled":true}}
{"value":"staging-01","label":"Staging Server 01 (online)"}
{"value":"staging-02","label":"Staging Server 02 (online)"}
EOF

SELECTED=$?

if [ $SELECTED -eq 0 ]; then
    echo ""
    echo "Deployment initiated! (servers selected above)"
else
    echo ""
    echo "Deployment cancelled"
fi
