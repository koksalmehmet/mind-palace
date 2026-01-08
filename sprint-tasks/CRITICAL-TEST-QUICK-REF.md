# Critical Tests - Quick Reference

**Goal:** Implement 9 test files for maximum coverage impact  
**Timeline:** 3-5 days  
**Expected Coverage Gain:** +25-30%

---

## Files to Create

### Language Parsers (Days 1-2)

1. **`parser_python_test.go`** - 400 lines, 15 test cases

   - Functions, classes, decorators, imports, async
   - Target: 0% → 75% (+75%)

2. **`parser_javascript_test.go`** - 350 lines, 15 test cases

   - Functions, arrow functions, classes, exports, imports
   - Target: 0% → 75% (+75%)

3. **`parser_typescript_test.go`** - 450 lines, 18 test cases
   - Interfaces, types, enums, generics, decorators
   - Target: 0% → 80% (+80%)

### MCP Tools (Days 3-4)

4. **`mcp_tools_brain_test.go`** - 500 lines, 12+ test cases

   - Store, recall, reflect, forget
   - Kinds: idea, decision, learning
   - Scopes: palace, room, file
   - Target: 15% → 85% (+70%)

5. **`mcp_tools_search_test.go`** - 400 lines, 10+ test cases

   - Explore, impact, symbols, deps, callers
   - Fuzzy search, room filters, limits
   - Target: 25% → 90% (+65%)

6. **`mcp_tools_briefing_test.go`** - 350 lines, 10 test cases
   - Smart briefing (workspace, file, room, task)
   - Styles: summary, detailed, actionable
   - Mock LLM responses
   - Target: 0% → 75% (+75%)

### HTTP Handlers (Day 5)

7. **`handlers_brain_test.go`** - 450 lines, 15 test cases

   - POST /api/remember
   - Validation, scopes, tags, classification
   - Target: 20% → 85% (+65%)

8. **`handlers_search_test.go`** - 350 lines, 12 test cases

   - GET /api/search
   - Query validation, limits, filters
   - Target: 30% → 90% (+60%)

9. **`handlers_context_test.go`** - 400 lines, 12 test cases
   - POST /api/context/preview
   - Auto-injection config, filters
   - Target: 0% → 80% (+80%)

---

## Test Pattern Template

```go
func TestComponent(t *testing.T) {
    tests := []struct {
        name       string
        input      string
        want       expectedType
        wantError  bool
    }{
        {name: "case1", input: "...", want: ...},
        {name: "case2", input: "...", want: ...},
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            got, err := Function(tt.input)
            if (err != nil) != tt.wantError {
                t.Errorf("error = %v, wantError %v", err, tt.wantError)
            }
            if got != tt.want {
                t.Errorf("got %v, want %v", got, tt.want)
            }
        })
    }
}
```

---

## Essential Test Cases by Category

### Parsers (All Languages)

- [ ] Simple declarations (function/class)
- [ ] With parameters/arguments
- [ ] With type annotations
- [ ] Nested structures
- [ ] Exports/imports
- [ ] Comments/docstrings
- [ ] Error cases (empty input, malformed)

### MCP Tools

- [ ] Happy path (valid input)
- [ ] Missing required parameters
- [ ] Empty/invalid parameters
- [ ] Default values
- [ ] Filters (scope, room, tags)
- [ ] Limits and pagination
- [ ] No results found

### HTTP Handlers

- [ ] Valid POST/GET
- [ ] Method not allowed
- [ ] Invalid JSON
- [ ] Missing required fields
- [ ] Empty fields
- [ ] No resources (503)
- [ ] Response structure validation

---

## Quick Commands

```bash
# Create test file
cd apps/cli/internal/analysis
touch parser_python_test.go

# Run single test
go test -v -run TestPythonParser

# Run with coverage
go test -coverprofile=coverage.out ./internal/analysis
go tool cover -html=coverage.out

# Run all tests
make test

# Full coverage report
make test-coverage
```

---

## Common Helpers Needed

```go
// Test server setup
func setupMCPServer(t *testing.T) (*MCPServer, *Butler)
func setupHTTPServer(t *testing.T) *Server

// Assertions
func assertSymbol(t, got Symbol, wantName, wantKind)
func assertHTTPStatus(t, rec *httptest.ResponseRecorder, want int)
func assertContains(t, haystack, needle string)

// Utilities
func mustMarshalJSON(t, v interface{}) string
func decodeJSON(t, body *bytes.Buffer) map[string]any
func seedTestDatabase(t, db *sql.DB, fixtures []Fixture)

// Mocks
type mockLLMClient struct { response string; err error }
```

---

## Priority Order

**P0 - Must Have (Core Functionality):**

1. Python parser (most used)
2. TypeScript parser (this project)
3. Brain MCP tool (core feature)
4. Brain HTTP handler (core API)

**P1 - Should Have (High Value):** 5. JavaScript parser 6. Search MCP tool 7. Search HTTP handler

**P2 - Nice to Have (Additional Coverage):** 8. Briefing MCP tool 9. Context HTTP handler

---

## Coverage Checkpoints

After each file:

```bash
go test -coverprofile=coverage.out ./internal/[package]
go tool cover -func=coverage.out | grep total
```

**Targets:**

- Day 1: Analysis package 25% → 40%
- Day 2: Analysis package 40% → 50%
- Day 3: Butler package 30% → 55%
- Day 4: Butler package 55% → 70%
- Day 5: Dashboard package 40% → 70%

---

## Validation Checklist

After implementation:

- [ ] All tests pass: `go test ./...`
- [ ] Coverage improved: `make test-coverage`
- [ ] No race conditions: `go test -race ./...`
- [ ] Tests run fast: <30s total
- [ ] CI/CD passes
- [ ] No flaky tests (run 3x)
- [ ] Code review ready

---

## Sample Test Counts

| Component         | Functions to Test  | Min Tests | Avg Lines/Test |
| ----------------- | ------------------ | --------- | -------------- |
| Python Parser     | Parse + 8 helpers  | 15        | 25-30          |
| JavaScript Parser | Parse + 7 helpers  | 15        | 23-28          |
| TypeScript Parser | Parse + 10 helpers | 18        | 25-30          |
| Brain MCP         | 4 tools            | 12        | 40-45          |
| Search MCP        | 6 tools            | 10        | 40-45          |
| Briefing MCP      | 1 tool             | 10        | 35-40          |
| Brain Handler     | 3 handlers         | 15        | 30-35          |
| Search Handler    | 1 handler          | 12        | 30-35          |
| Context Handler   | 2 handlers         | 12        | 35-40          |

**Total:** ~119 test cases, ~3,650 lines

---

## Resources

- Full design: `CRITICAL-TEST-DESIGN.md`
- Existing patterns:
  - `parser_cpp_test.go`
  - `mcp_test.go`
  - `handlers_test.go`
- Go testing docs: https://pkg.go.dev/testing
- Table-driven tests: https://go.dev/wiki/TableDrivenTests
