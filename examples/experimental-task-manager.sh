#!/bin/bash

# Example: Task manager with state-based disabled items
# Demonstrates using config delimiter for simple use cases

echo "=== Task Manager ==="
echo ""
echo "Select a task to work on (disabled tasks are blocked or in progress)"
echo ""

gum choose --experimental \
  --config-delimiter="::" \
  --header="Available Tasks:" \
  "Write documentation" \
  "Update dependencies" \
  'Fix bug #123::{"disabled":true}' \
  'Review PR #456 (in progress)::{"disabled":true}' \
  "Deploy to staging" \
  "Run integration tests" \
  'Update changelog::{"disabled":true}'

SELECTED=$?

if [ $SELECTED -eq 0 ]; then
    echo ""
    echo "Starting work on task!"
else
    echo ""
    echo "No task selected"
fi
