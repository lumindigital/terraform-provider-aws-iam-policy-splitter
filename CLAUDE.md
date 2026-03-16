# CLAUDE.md

This file provides guidance to Claude Code when working with this repository.

## Project Overview

`terraform-aws-iam-policy-splitter` is a Terraform provider that distributes AWS IAM policy JSON strings across multiple bins using a best-fit decreasing bin packing algorithm. It respects AWS's 6144-character policy size limit and excludes whitespace from size calculations (matching AWS IAM behavior).

## Build, Test, and Lint Commands

All commands use [Taskfile](https://taskfile.dev):

```bash
task local:build      # Build the provider
task local:test       # Run unit tests (timeout: 120s, parallel: 10)
task local:testacc    # Run acceptance tests (sets TF_ACC=1, timeout: 120m)
task local:lint       # Run golangci-lint with auto-fix
task local:fmt        # Format Go code with gofmt
task local:generate   # Run code generation (from tools/)
task local:install    # Install the provider binary
task local:all        # Run all tasks in sequence
```

### Running a single test

```bash
go test -v -run TestBestFit_BasicDistribution ./internal/binpack/
```

## Project Structure

```
main.go                              # Terraform provider entry point
internal/
  binpack/
    binpack.go                       # Core bin packing algorithm (BestFit, BestFitWithSizes)
    binpack_test.go                  # Unit tests (13 tests, 100% coverage)
  provider/
    provider.go                      # Terraform provider definition
    provider_test.go                 # Provider test factories
    aws_iam_policy_splitter_data_source.go      # Data source implementation
    aws_iam_policy_splitter_data_source_test.go # Acceptance tests (8 tests, require TF_ACC=1)
docs/                                # Auto-generated provider documentation
examples/                            # Terraform usage examples
```

## Key Architecture Details

- **Bin packing algorithm**: `internal/binpack` implements best-fit decreasing. Items are sorted largest-first, then placed in the tightest-fitting bin.
- **Whitespace handling**: The data source strips whitespace before calculating effective sizes but preserves original policy strings in output. This matches AWS IAM behavior.
- **Deterministic output**: The data source ID is a truncated SHA-256 hash of all input policies.
- **Error safety**: Error messages truncate policy content to 50 characters to avoid leaking sensitive data.

## Code Conventions

- Uses Terraform Plugin Framework v1 (not the older SDKv2)
- Linting enforced via `golangci-lint` (config: `.golangci.yml`)
- Struct fields use Go naming conventions (e.g., `JSONPolicies` not `JsonPolicies`) with `tfsdk` tags for Terraform mapping
- Magic numbers should be extracted as named constants
- Test globals require `//nolint:gochecknoglobals` with justification
