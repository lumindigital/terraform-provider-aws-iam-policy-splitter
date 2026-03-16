# Ralph Agent Configuration

## Build Instructions

```bash
task local:build
```

## Test Instructions

```bash
task local:test
```

## Run Instructions

```bash
task local:build && ./bin/terraform-aws-iam-policy-splitter
```
## Lint and Formatting Instructions
```bash
task local:fmt && task local:lint
```

## Notes
- This is a a standalone golang application
- Tests use golang's native testing framework
- Update this file when build process changes

