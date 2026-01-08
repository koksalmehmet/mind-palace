# golangci-lint Configuration Analysis for Mind Palace

**Generated:** January 8, 2026  
**Go Version:** 1.25  
**Project:** Mind Palace CLI

## Executive Summary

This document explains the golangci-lint configuration created for the Mind Palace Go CLI project, a sophisticated tool featuring parsers, MCP server integration, SQLite database management, and memory management with embedding pipelines.

## Project Analysis

### Codebase Characteristics

**Project Type:** CLI tool with complex internal architecture

- **Lines of Code:** ~15,000+ across 77+ test files and production code
- **Go Version:** 1.25 (bleeding edge)
- **Key Dependencies:**
  - SQLite (modernc.org/sqlite)
  - Tree-sitter parsers (go-tree-sitter)
  - WebSocket server (gorilla/websocket)
  - JSON schema validation
  - UUID generation

### Code Patterns Identified

1. **Database Access Patterns**

   - Heavy use of `sql.DB`, `sql.Rows`, `sql.Stmt`
   - WAL mode SQLite with pragmas
   - Multiple database instances (index.db, memory.db, corridor.db)
   - Migration system with versioning

2. **Concurrency Patterns**

   - Goroutines in LSP client, embedding pipeline
   - Context usage for cancellation and timeouts
   - Mutex synchronization in LSP client
   - Background workers for embedding queue

3. **Error Handling**

   - Extensive `fmt.Errorf` with `%w` wrapping
   - Custom error types
   - Graceful degradation (non-fatal errors)

4. **Parser Architecture**

   - Multiple language parsers (Bash, C, C++, C#, Go, Python, TypeScript, etc.)
   - Tree-sitter integration
   - Symbol extraction and analysis

5. **HTTP/API Integration**

   - Ollama API client
   - OpenAI API client
   - Anthropic API client
   - Neighbor workspace fetching (HTTP)

6. **Testing Strategy**
   - Uses `testify` assertions
   - E2E tests, integration tests, unit tests
   - Benchmark tests
   - In-memory SQLite for testing

## Linter Configuration Explained

### Phase 1: Critical Linters (Enable Immediately)

These linters catch bugs and critical issues. Enable these first and fix all issues before proceeding.

#### Error Handling & Bug Detection

**`errcheck`** - Unchecked error returns

- **Why:** Mind Palace uses extensive database operations, HTTP calls, and file I/O where unchecked errors could lead to data loss or undefined behavior
- **Configuration:** Checks type assertions and blank identifiers
- **Expected Issues:** 50-100+ (database Close() calls, file operations)
- **Priority:** CRITICAL

**`govet`** - Official Go vet tool

- **Why:** Catches common mistakes (unreachable code, Printf formatting, struct tags)
- **Configuration:** All checks enabled except shadow and fieldalignment
- **Expected Issues:** 10-20
- **Priority:** CRITICAL

**`staticcheck`** - Comprehensive static analysis

- **Why:** Industry-standard linter, catches subtle bugs and inefficiencies
- **Configuration:** All checks enabled, deprecated usage allowed for migration
- **Expected Issues:** 30-50
- **Priority:** CRITICAL

**`gosimple`** - Code simplification

- **Why:** Identifies code that can be simplified
- **Expected Issues:** 20-40
- **Priority:** HIGH

**`ineffassign`** - Ineffectual assignments

- **Why:** Detects assignments to variables that are never read
- **Expected Issues:** 5-15
- **Priority:** HIGH

**`unused`** - Unused code detection

- **Why:** Finds dead code (functions, constants, variables)
- **Expected Issues:** 10-30
- **Priority:** MEDIUM

#### SQL-Specific Linters

**`sqlclosecheck`** - SQL resource cleanup

- **Why:** Critical for Mind Palace's extensive SQLite usage. Ensures `rows.Close()` and `stmt.Close()` are called
- **Configuration:** Default
- **Expected Issues:** 20-40 (multiple databases: index, memory, corridor)
- **Priority:** CRITICAL

**`rowserrcheck`** - SQL rows.Err() checking

- **Why:** Ensures errors during iteration are checked
- **Expected Issues:** 15-30
- **Priority:** CRITICAL

#### Common Mistakes

**`typecheck`** - Type checking

- **Why:** Ensures type correctness
- **Priority:** CRITICAL

**`bodyclose`** - HTTP response body close

- **Why:** Important for Ollama/OpenAI/Anthropic API clients
- **Expected Issues:** 5-10
- **Priority:** HIGH

**`contextcheck`** - Context propagation

- **Why:** Critical for LSP client, LLM API calls with timeouts
- **Expected Issues:** 10-20
- **Priority:** HIGH

**`errname`** - Error naming conventions

- **Why:** Errors should end with "Error" or "Err"
- **Expected Issues:** 5-10
- **Priority:** LOW

**`errorlint`** - Error wrapping with %w

- **Why:** Ensures proper error wrapping for error chain inspection
- **Expected Issues:** 20-30
- **Priority:** MEDIUM

**`nilerr`** - Nil error returns

- **Why:** Catches `if err != nil { return err }` where err is always nil
- **Expected Issues:** 0-5
- **Priority:** MEDIUM

### Phase 2: Important Code Quality

#### Complexity Metrics

**`gocognit`** - Cognitive complexity

- **Why:** Better metric than cyclomatic complexity for readability
- **Configuration:** Threshold of 25 (reasonable for parser/butler logic)
- **Expected Issues:** 15-25 (parsers, butler orchestration, index queries)
- **Priority:** MEDIUM
- **Note:** Some parser functions may legitimately exceed this

**`gocyclo`** - Cyclomatic complexity

- **Why:** Classic complexity metric
- **Configuration:** Threshold of 20
- **Expected Issues:** 10-20
- **Priority:** MEDIUM

**`nestif`** - Nested if depth

- **Why:** Deeply nested ifs harm readability
- **Configuration:** Maximum depth of 5
- **Expected Issues:** 5-10
- **Priority:** MEDIUM

#### Performance

**`prealloc`** - Slice preallocation

- **Why:** Performance optimization for slice operations (important for parsers processing large codebases)
- **Configuration:** Checks simple, range, and for loops
- **Expected Issues:** 20-40
- **Priority:** LOW

**`unconvert`** - Unnecessary type conversions

- **Why:** Cleaner code, micro-optimization
- **Expected Issues:** 5-15
- **Priority:** LOW

#### Style & Best Practices

**`gofmt`** - Go formatting

- **Why:** Standard formatting
- **Expected Issues:** 0 (should be auto-fixed)
- **Priority:** HIGH

**`goimports`** - Import organization

- **Why:** Organize imports consistently
- **Expected Issues:** 0-10
- **Priority:** MEDIUM

**`misspell`** - Spelling

- **Why:** Catches typos in comments, strings
- **Configuration:** US locale
- **Expected Issues:** 10-20
- **Priority:** LOW

**`revive`** - Fast, configurable linter

- **Why:** Modern replacement for golint with many useful rules
- **Configuration:**
  - Exported functions documentation
  - Error naming and strings
  - Context as argument
  - Receiver naming
  - Unused parameters
- **Expected Issues:** 30-60
- **Priority:** MEDIUM

**`stylecheck`** - Style conventions

- **Why:** Enforces Go style guide
- **Configuration:** Most checks enabled, relaxed package comments and underscores (for MCP/JSON fields)
- **Expected Issues:** 20-40
- **Priority:** MEDIUM

### Phase 3: Enhanced Safety & Maintainability

#### Potential Bugs

**`copyloopvar`** - Loop variable semantics

- **Why:** Ensures correct behavior with Go 1.22+ loop variable scoping
- **Expected Issues:** 0-5 (Go 1.25 has fixed loop semantics)
- **Priority:** LOW

**`durationcheck`** - Duration arithmetic

- **Why:** Catches time.Duration multiplication/division errors
- **Expected Issues:** 0-5
- **Priority:** MEDIUM

**`noctx`** - HTTP requests without context

- **Why:** All HTTP requests should have context for cancellation
- **Expected Issues:** 5-10
- **Priority:** MEDIUM

**`nolintlint`** - Nolint directive validation

- **Why:** Ensures nolint directives are valid and necessary
- **Expected Issues:** 0 (no nolint directives found in current code)
- **Priority:** LOW

**`predeclared`** - Predeclared identifier shadowing

- **Why:** Prevents shadowing built-ins like `len`, `error`, `string`
- **Expected Issues:** 5-10
- **Priority:** MEDIUM

**`reassign`** - Function parameter reassignment

- **Why:** Reassigning parameters can be confusing
- **Expected Issues:** 5-15
- **Priority:** LOW

**`unparam`** - Unused function parameters

- **Why:** Finds parameters that are never used
- **Configuration:** Excludes exported functions
- **Expected Issues:** 10-20
- **Priority:** LOW

**`wastedassign`** - Wasted assignments

- **Why:** Assigns to variables that are immediately overwritten
- **Expected Issues:** 5-10
- **Priority:** MEDIUM

#### Security

**`gosec`** - Security analysis

- **Why:** Detects security issues (SQL injection, file permissions, crypto, etc.)
- **Configuration:** Medium severity/confidence, excludes G104 (duplicate of errcheck) and G307
- **Expected Issues:** 10-30 (file operations, temp files, command execution)
- **Priority:** HIGH
- **Note:** Many may be false positives in CLI context

#### Concurrency

**`gocritic`** - Comprehensive critique

- **Why:** Large rule set covering diagnostics, performance, style, and concurrency
- **Configuration:** All tags enabled except whyNoLint, unnamedResult, hugeParam
- **Expected Issues:** 30-60
- **Priority:** MEDIUM
- **Note:** Important for LSP client goroutines, embedding pipeline workers

#### Testing

**`testifylint`** - Testify best practices

- **Why:** Project uses testify extensively (77+ test files)
- **Configuration:** All checks enabled except float-compare
- **Expected Issues:** 10-20
- **Priority:** MEDIUM

**`tparallel`** - t.Parallel() usage

- **Why:** Ensures correct parallel test execution
- **Expected Issues:** 0-5
- **Priority:** LOW

### Phase 4: Optional Strict Linters (Disabled by Default)

These are disabled initially but can be enabled progressively as the codebase matures:

- **`exhaustive`** - Exhaustive switch cases (can be very strict)
- **`exhaustruct`** - Exhaustive struct initialization (too strict for this project)
- **`forcetypeassert`** - Force type assertion checks (covered by errcheck)
- **`funlen`** - Function length (parsers have legitimate long functions)
- **`gochecknoglobals`** - No global variables (some globals are necessary)
- **`gochecknoinits`** - No init functions (some init is used)
- **`godox`** - TODO/FIXME/XXX (no TODOs found currently, can enable later)
- **`goerr113`** - Error wrapping style (too opinionated)
- **`mnd`/`gomnd`** - Magic numbers (too noisy for CLI args, thresholds)
- **`ireturn`** - Interface return types (not always applicable)
- **`nilnil`** - Return nil, nil pattern (sometimes legitimate)
- **`nonamedreturns`** - Named returns (sometimes useful for documentation)
- **`varnamelen`** - Variable name length (short names are fine in limited scope)
- **`wrapcheck`** - Error wrapping (too strict)

## Exclusion Rules Explained

### Test Files

Tests are excluded from complexity and style checks because:

- Test functions can be longer and more complex
- Error handling can be simplified (t.Fatal)
- Security concerns are different

### Generated/Schema Code

- `apps/cli/schemas/` - JSON schema loading (generated/static)
- `apps/cli/starter/` - Template/example code

### Complexity Allowances

- **Parsers** (`apps/cli/internal/analysis/parser_*.go`) - Tree traversal is inherently complex
- **Index/Scan** - Database query building can be complex
- **Butler** - Orchestration layer has legitimate branching logic

### Context in main.go

- `context.Background()` is allowed in main package initialization

## Expected Issues Summary

| Category            | Expected Count | Priority | Effort     |
| ------------------- | -------------- | -------- | ---------- |
| SQL resource leaks  | 30-50          | Critical | Medium     |
| Unchecked errors    | 50-100         | Critical | High       |
| Context propagation | 10-20          | High     | Medium     |
| HTTP body close     | 5-10           | High     | Low        |
| Error wrapping      | 20-30          | Medium   | Medium     |
| Complexity issues   | 20-40          | Medium   | High       |
| Style/naming        | 40-80          | Low      | Medium     |
| Security warnings   | 10-30          | High     | Low-Medium |
| **TOTAL**           | **185-360**    | -        | **High**   |

## Rollout Strategy

### Phase 1: Critical Safety (Week 1)

Enable and fix:

1. `sqlclosecheck`, `rowserrcheck` - Fix all SQL resource leaks
2. `errcheck` - Fix unchecked errors (start with non-test files)
3. `staticcheck`, `govet` - Fix all static analysis issues

**Expected effort:** 2-3 days for 80-150 issues

### Phase 2: Correctness (Week 2)

Enable and fix:

1. `bodyclose`, `contextcheck` - HTTP/context issues
2. `errorlint` - Error wrapping
3. `typecheck`, `ineffassign`, `unused` - Type/assignment issues

**Expected effort:** 1-2 days for 40-70 issues

### Phase 3: Quality (Week 3)

Enable and fix:

1. `gocognit`, `gocyclo`, `nestif` - Refactor complex functions
2. `revive`, `stylecheck` - Style improvements
3. `gosec` - Review and fix legitimate security issues

**Expected effort:** 2-3 days for 60-120 issues

### Phase 4: Polish (Week 4)

Enable and fix:

1. `prealloc`, `unconvert` - Performance optimizations
2. `testifylint` - Improve test quality
3. `gocritic` - Apply critique suggestions

**Expected effort:** 1-2 days for 40-80 issues

## Integration with CI/CD

### Makefile Integration

Add to `Makefile`:

```makefile
.PHONY: lint lint-fix lint-report

# Run golangci-lint
lint:
	@echo "Running golangci-lint..."
	@cd apps/cli && golangci-lint run --timeout 5m

# Run with auto-fix
lint-fix:
	@echo "Running golangci-lint with auto-fix..."
	@cd apps/cli && golangci-lint run --fix --timeout 5m

# Generate detailed report
lint-report:
	@echo "Generating lint report..."
	@cd apps/cli && golangci-lint run --out-format=html > ../../lint-report.html
	@echo "Report saved to lint-report.html"
```

### GitHub Actions

Create `.github/workflows/lint.yml`:

```yaml
name: Lint

on:
  push:
    branches: [main, develop]
  pull_request:

jobs:
  golangci:
    name: golangci-lint
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
      - uses: actions/setup-go@v5
        with:
          go-version: "1.25"
      - name: golangci-lint
        uses: golangci/golangci-lint-action@v4
        with:
          version: latest
          working-directory: apps/cli
          args: --timeout 5m
```

### Pre-commit Hook

Create `.git/hooks/pre-commit`:

```bash
#!/bin/bash
cd apps/cli
golangci-lint run --new-from-rev=HEAD~1 --timeout 2m
```

## Performance Considerations

- **Full run time:** 2-3 minutes (with all linters enabled)
- **Fast linters only:** 30-45 seconds
- **New code only:** 10-20 seconds

Fast linters for quick feedback:

- errcheck
- govet
- staticcheck
- typecheck
- unused

## Customization Tips

### For More Strict Checking

Enable in `.golangci.yml`:

```yaml
linters:
  enable:
    - exhaustive
    - funlen
    - gochecknoglobals
    - godox
```

### For Faster CI

Create `.golangci-ci.yml`:

```yaml
linters:
  disable:
    - gocognit
    - gocyclo
    - gosec
    - gocritic
```

### For Development

```bash
# Check only changed files
golangci-lint run --new-from-rev=main

# Run specific linters
golangci-lint run --disable-all --enable=errcheck,govet
```

## Known False Positives

1. **`gosec` G304** - File path from variable (CLI tool needs this)
2. **`gosec` G204** - Command execution (CLI tool feature)
3. **`gosec` G110** - Decompression bomb (not applicable)
4. **`gocritic` hugeParam** - Tree-sitter nodes are large by design

Use `//nolint:lintername // reason` when necessary.

## Maintenance

- **Monthly:** Update golangci-lint version
- **Quarterly:** Review and enable new linters
- **Per Go release:** Check for deprecated linters
- **Before major release:** Run exhaustive mode

## References

- [golangci-lint Documentation](https://golangci-lint.run/)
- [Enabled Linters List](https://golangci-lint.run/usage/linters/)
- [Go 1.25 Release Notes](https://tip.golang.org/doc/go1.25)
- [Mind Palace Project](https://github.com/koksalmehmet/mind-palace)

---

**Next Steps:**

1. Install golangci-lint: `go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest`
2. Run initial check: `cd apps/cli && golangci-lint run --timeout 5m`
3. Fix Phase 1 issues (critical)
4. Add to CI/CD pipeline
5. Progressively enable Phase 2-4 linters
