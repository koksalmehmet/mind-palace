# golangci-lint Quick Reference

## Installation

```bash
# Install latest version
go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest

# Verify installation
golangci-lint --version
```

## Quick Commands

```bash
# Run all linters
cd apps/cli && golangci-lint run

# Auto-fix issues
golangci-lint run --fix

# Check only new code (since main branch)
golangci-lint run --new-from-rev=origin/main

# Run specific linters
golangci-lint run --disable-all --enable=errcheck,govet,staticcheck

# Fast run (only fast linters)
golangci-lint run --fast

# Generate HTML report
golangci-lint run --out-format=html > lint-report.html
```

## Phased Enablement Priority

### Week 1: Critical (Fix First)

```yaml
- sqlclosecheck # SQL Close() calls
- rowserrcheck # SQL Rows.Err() checks
- errcheck # Unchecked errors
- staticcheck # Comprehensive analysis
- govet # Official Go vet
```

### Week 2: Correctness

```yaml
- bodyclose # HTTP body Close()
- contextcheck # Context propagation
- errorlint # Error wrapping %w
- typecheck # Type correctness
- ineffassign # Dead assignments
```

### Week 3: Quality

```yaml
- gocognit # Cognitive complexity
- revive # Style & best practices
- stylecheck # Go style guide
- gosec # Security issues
```

### Week 4: Polish

```yaml
- prealloc # Slice preallocation
- gocritic # Comprehensive critique
- testifylint # Testify best practices
```

## Common Issues & Fixes

### SQL Resource Leaks

**Issue:**

```go
rows, err := db.Query("SELECT ...")
if err != nil {
    return err
}
for rows.Next() {
    // ...
}
```

**Fix:**

```go
rows, err := db.Query("SELECT ...")
if err != nil {
    return err
}
defer rows.Close()  // Add this

for rows.Next() {
    // ...
}
return rows.Err()   // Add this
```

### Unchecked Errors

**Issue:**

```go
json.Unmarshal(data, &result)
file.Close()
```

**Fix:**

```go
if err := json.Unmarshal(data, &result); err != nil {
    return fmt.Errorf("unmarshal: %w", err)
}
defer func() {
    if err := file.Close(); err != nil {
        log.Printf("failed to close file: %v", err)
    }
}()
```

### HTTP Body Close

**Issue:**

```go
resp, err := http.Get(url)
if err != nil {
    return err
}
body, _ := io.ReadAll(resp.Body)
```

**Fix:**

```go
resp, err := http.Get(url)
if err != nil {
    return err
}
defer resp.Body.Close()  // Add this

body, err := io.ReadAll(resp.Body)
if err != nil {
    return fmt.Errorf("read body: %w", err)
}
```

### Context Propagation

**Issue:**

```go
func fetchData(url string) error {
    req, _ := http.NewRequest("GET", url, nil)
    // ...
}
```

**Fix:**

```go
func fetchData(ctx context.Context, url string) error {
    req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
    if err != nil {
        return err
    }
    // ...
}
```

### Error Wrapping

**Issue:**

```go
if err := doSomething(); err != nil {
    return fmt.Errorf("failed to do something: %v", err)
}
```

**Fix:**

```go
if err := doSomething(); err != nil {
    return fmt.Errorf("failed to do something: %w", err)  // Use %w
}
```

## Suppressing False Positives

```go
// Suppress specific linter for one line
var global = 123 //nolint:gochecknoglobals // Config constant

// Suppress multiple linters
func legacy() { //nolint:funlen,gocognit // Legacy code, refactor later
    // ...
}

// Suppress for entire file (use sparingly)
//nolint:gosec // This file handles user commands, requires command execution
package main
```

## Configuration Sections

### Run Settings

```yaml
run:
  timeout: 5m # Max run time
  tests: true # Include test files
  skip-dirs: # Directories to skip
    - vendor
```

### Linter Settings

```yaml
linters:
  disable-all: true # Start fresh
  enable: # Explicitly enable
    - errcheck
    - govet
```

### Issue Exclusions

```yaml
issues:
  exclude-rules:
    - path: _test\.go # Relax for tests
      linters:
        - errcheck
        - gosec
```

## Integration with Make

Add to your workflow:

```makefile
# In apps/cli/Makefile

.PHONY: lint lint-fix check

# Run linters
lint:
	golangci-lint run --timeout 5m

# Auto-fix what can be fixed
lint-fix:
	golangci-lint run --fix --timeout 5m

# Pre-commit check (fast)
check:
	golangci-lint run --fast --new-from-rev=HEAD~1
```

Then use:

```bash
make lint         # Full check
make lint-fix     # Fix automatically
make check        # Quick pre-commit check
```

## CI/CD Integration

### GitHub Actions

```yaml
- name: golangci-lint
  uses: golangci/golangci-lint-action@v4
  with:
    version: latest
    working-directory: apps/cli
```

### Pre-commit Hook

```bash
#!/bin/bash
# .git/hooks/pre-commit
cd apps/cli
golangci-lint run --new-from-rev=HEAD~1 --fast
```

## Performance Tips

- **Use `--fast`** for quick checks (skips slow linters)
- **Use `--new-from-rev`** to check only changed files
- **Run in parallel** with `--concurrency=4` (default: CPU count)
- **Cache results** with `--out-format=checkstyle` for CI

## Troubleshooting

**"Too many issues found"**

```bash
# See all issues (no limit)
golangci-lint run --max-issues-per-linter=0 --max-same-issues=0
```

**"Linter X is taking too long"**

```bash
# Disable slow linter temporarily
golangci-lint run --disable=gosec
```

**"I want to see what would be fixed"**

```bash
# Dry-run of fixes
golangci-lint run --fix --out-format=colored-line-number | grep "Fixed:"
```

## Resources

- Config: `.golangci.yml`
- Full Guide: `GOLANGCI-LINT-GUIDE.md`
- Linter Docs: https://golangci-lint.run/usage/linters/
- Issue Tracker: https://github.com/golangci/golangci-lint/issues

## Mind Palace Specific Notes

**High-risk areas requiring attention:**

1. `apps/cli/internal/memory/` - SQLite resource management
2. `apps/cli/internal/corridor/` - Global corridor database
3. `apps/cli/internal/index/` - Index database operations
4. `apps/cli/internal/butler/mcp_*.go` - MCP server handlers
5. `apps/cli/internal/llm/` - HTTP client resource cleanup
6. `apps/cli/internal/analysis/lsp_client.go` - Goroutine/context management

**Test coverage:** 77+ test files - ensure test exclusions are appropriate

**Performance-critical paths:**

- Scanner (process thousands of files)
- Parsers (process large files)
- Embedding pipeline (background workers)

Use prealloc and performance linters carefully in these areas.
