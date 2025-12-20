# Experimental Features

> **Note**: These features are experimental and may change in future releases.

Gum now supports experimental features for `choose` and `filter` commands that allow you to:

- **Disable specific options** while keeping them visible
- **Use JSON/JSONL format** to define options with metadata
- **Configure options** with delimiters for advanced control

## Quick Example

```bash
# Enable experimental mode and use JSON options with disabled items
echo '[
  {"value":"option1","label":"Available Option"},
  {"value":"option2","label":"Disabled Option","config":{"disabled":true}},
  {"value":"option3","label":"Another Available"}
]' | gum choose --experimental --option-as-json
```

## Use Cases

- **State-aware menus**: Show all options but disable unavailable ones based on current state
- **Permission-based selection**: Display restricted items that users can see but not select
- **Workflow management**: Indicate blocked or in-progress tasks

## Documentation

For complete documentation on experimental features, see [EXPERIMENTAL_FEATURES.md](./EXPERIMENTAL_FEATURES.md).

## Examples

Check out the example scripts in the `examples/` directory:

- [`experimental-git-branches.sh`](./examples/experimental-git-branches.sh) - Branch selector with protected branches
- [`experimental-server-deployment.sh`](./examples/experimental-server-deployment.sh) - Server deployment with status
- [`experimental-task-manager.sh`](./examples/experimental-task-manager.sh) - Task list with blocked items
- [`experimental-file-filter.sh`](./examples/experimental-file-filter.sh) - File filter with permissions

## Enabling Experimental Features

```bash
# Via flag
gum choose --experimental [options]

# Via environment variable
export GUM_EXPERIMENTAL=true
gum choose [options]
```
