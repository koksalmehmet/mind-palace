# Test Implementation Checklist

**Goal:** Implement 9 critical test files for +25-30% coverage  
**Timeline:** January 8-12, 2026 (5 days)  
**Status:** 🔴 Not Started

---

## Phase 1: Language Parsers (Days 1-2)

### Day 1 - January 8, 2026

- [ ] **Task 1.1: Python Parser Tests** (4 hours)

  - [ ] Create `apps/cli/internal/analysis/parser_python_test.go`
  - [ ] Copy template from `SAMPLE-parser_python_test.go`
  - [ ] Implement 15 test cases:
    - [ ] Simple function
    - [ ] Function with params & return type
    - [ ] Function with docstring
    - [ ] Empty class
    - [ ] Class with methods
    - [ ] Class with `__init__`
    - [ ] Class inheritance
    - [ ] Decorated function
    - [ ] Decorated class
    - [ ] Async function
    - [ ] Global variables
    - [ ] Import statements
    - [ ] Method calls (relationships)
    - [ ] Nested functions
    - [ ] Property methods
  - [ ] Run tests: `go test -v -run TestPythonParser ./internal/analysis`
  - [ ] Verify coverage: 0% → 75%+
  - [ ] Commit: "test: add comprehensive Python parser tests"

- [ ] **Task 1.2: JavaScript Parser Tests** (4 hours)
  - [ ] Create `apps/cli/internal/analysis/parser_javascript_test.go`
  - [ ] Implement 15 test cases:
    - [ ] Function declaration
    - [ ] Arrow function (const)
    - [ ] Arrow function (let)
    - [ ] Class declaration
    - [ ] Class with methods
    - [ ] Class with static methods
    - [ ] Class inheritance
    - [ ] Named export (function)
    - [ ] Named export (class)
    - [ ] Default export
    - [ ] Variable declarations
    - [ ] Object method shorthand
    - [ ] Import statements
    - [ ] Async function
    - [ ] Error handling
  - [ ] Run tests: `go test -v -run TestJavaScriptParser ./internal/analysis`
  - [ ] Verify coverage: 0% → 75%+
  - [ ] Commit: "test: add comprehensive JavaScript parser tests"

**End of Day 1 Checkpoint:**

- [ ] 2 test files created
- [ ] ~750 lines of test code
- [ ] Analysis package coverage: 25% → 45%
- [ ] All tests passing

---

### Day 2 - January 9, 2026

- [ ] **Task 1.3: TypeScript Parser Tests** (6 hours)

  - [ ] Create `apps/cli/internal/analysis/parser_typescript_test.go`
  - [ ] Implement 18 test cases:
    - [ ] Function with type annotations
    - [ ] Arrow function with types
    - [ ] Interface declaration
    - [ ] Interface with methods
    - [ ] Interface inheritance
    - [ ] Type alias (simple)
    - [ ] Type alias (object)
    - [ ] Type alias (union)
    - [ ] Enum declaration
    - [ ] Const enum
    - [ ] Class with TypeScript features
    - [ ] Class implementing interface
    - [ ] Generic function
    - [ ] Generic class
    - [ ] Exported types
    - [ ] Namespace declaration
    - [ ] Decorators
    - [ ] Async/await with types
  - [ ] Run tests: `go test -v -run TestTypeScriptParser ./internal/analysis`
  - [ ] Verify coverage: 0% → 80%+
  - [ ] Commit: "test: add comprehensive TypeScript parser tests"

- [ ] **Task 1.4: Integration & Review** (2 hours)
  - [ ] Run all analysis tests: `go test ./internal/analysis/...`
  - [ ] Generate coverage report: `go test -coverprofile=coverage.out ./internal/analysis`
  - [ ] View HTML coverage: `go tool cover -html=coverage.out`
  - [ ] Fix any failing tests
  - [ ] Update documentation if needed

**End of Day 2 Checkpoint:**

- [ ] 3 parser test files complete
- [ ] ~1,200 total test lines
- [ ] Analysis package coverage: 25% → 55%+
- [ ] Zero failing tests
- [ ] Commit: "test: complete language parser test suite"

---

## Phase 2: MCP Tools (Days 3-4)

### Day 3 - January 10, 2026

- [ ] **Task 2.1: Brain MCP Tool Tests** (6 hours)

  - [ ] Create `apps/cli/internal/butler/mcp_tools_brain_test.go`
  - [ ] Copy template from `SAMPLE-mcp_tools_brain_test.go`
  - [ ] Implement TestToolStore (12 cases):
    - [ ] Store idea (explicit)
    - [ ] Store decision with tags
    - [ ] Store learning
    - [ ] Auto-classification
    - [ ] Room scope
    - [ ] File scope
    - [ ] Multiple tags
    - [ ] Default palace scope
    - [ ] Backward compatibility
    - [ ] Missing content error
    - [ ] Empty content error
    - [ ] Empty tags handling
  - [ ] Implement TestToolRecall (7 cases):
    - [ ] Recall all
    - [ ] Filter by kind
    - [ ] With limit
    - [ ] Query filter
    - [ ] Tag filter
    - [ ] No memories
    - [ ] Scope filter
  - [ ] Implement TestToolReflect (3 cases)
  - [ ] Implement TestToolForget (3 cases)
  - [ ] Run tests: `go test -v -run TestTool.*Brain ./internal/butler`
  - [ ] Verify coverage on mcp_tools_brain.go: 15% → 85%
  - [ ] Commit: "test: add comprehensive brain MCP tool tests"

- [ ] **Task 2.2: Review & Refactor** (2 hours)
  - [ ] Extract common helpers to butler_test_helpers.go
  - [ ] Ensure test isolation (no side effects)
  - [ ] Add table-driven test improvements

**End of Day 3 Checkpoint:**

- [ ] 1 MCP tool test file complete
- [ ] ~500 test lines added
- [ ] Butler package coverage: 30% → 50%+
- [ ] All tests passing

---

### Day 4 - January 11, 2026

- [ ] **Task 2.3: Search MCP Tool Tests** (4 hours)

  - [ ] Create `apps/cli/internal/butler/mcp_tools_search_test.go`
  - [ ] Implement TestToolExplore (10 cases):
    - [ ] Basic query
    - [ ] With limit
    - [ ] Room filter
    - [ ] Fuzzy match enabled
    - [ ] Fuzzy match disabled
    - [ ] Missing query error
    - [ ] Empty query error
    - [ ] Limit capping
    - [ ] No results
    - [ ] Special characters
  - [ ] Implement additional tool tests:
    - [ ] TestToolExploreImpact
    - [ ] TestToolExploreSymbols
    - [ ] TestToolExploreFile
    - [ ] TestToolExploreDeps
    - [ ] TestToolExploreCallers
  - [ ] Run tests: `go test -v -run TestToolExplore ./internal/butler`
  - [ ] Verify coverage on mcp_tools_search.go: 25% → 90%
  - [ ] Commit: "test: add comprehensive search MCP tool tests"

- [ ] **Task 2.4: Briefing MCP Tool Tests** (3 hours)

  - [ ] Create `apps/cli/internal/butler/mcp_tools_briefing_test.go`
  - [ ] Create mock LLM client
  - [ ] Implement TestToolBriefingSmart (10 cases):
    - [ ] Workspace context
    - [ ] File context
    - [ ] Room context
    - [ ] Task context
    - [ ] Summary style
    - [ ] Detailed style
    - [ ] Actionable style
    - [ ] Default context
    - [ ] Invalid context
    - [ ] Missing context path
  - [ ] Run tests: `go test -v -run TestToolBriefing ./internal/butler`
  - [ ] Verify coverage on mcp_tools_briefing.go: 0% → 75%
  - [ ] Commit: "test: add briefing MCP tool tests with LLM mocks"

- [ ] **Task 2.5: Integration** (1 hour)
  - [ ] Run all butler tests: `go test ./internal/butler/...`
  - [ ] Check coverage: `go test -coverprofile=coverage.out ./internal/butler`
  - [ ] Fix any race conditions: `go test -race ./internal/butler/...`

**End of Day 4 Checkpoint:**

- [ ] 3 MCP tool test files complete
- [ ] ~1,250 total test lines
- [ ] Butler package coverage: 30% → 70%+
- [ ] All tests passing
- [ ] Commit: "test: complete MCP tool test suite"

---

## Phase 3: HTTP Handlers (Day 5)

### Day 5 - January 12, 2026

- [ ] **Task 3.1: Brain HTTP Handler Tests** (3 hours)

  - [ ] Create `apps/cli/internal/dashboard/handlers_brain_test.go`
  - [ ] Implement TestHandleRemember (15 cases):
    - [ ] Valid idea POST
    - [ ] Valid decision POST
    - [ ] Valid learning POST
    - [ ] Auto-classification
    - [ ] Room scope
    - [ ] File scope
    - [ ] Multiple tags
    - [ ] Missing content error
    - [ ] Empty content error
    - [ ] Invalid JSON error
    - [ ] GET method not allowed
    - [ ] PUT method not allowed
    - [ ] No memory (503)
    - [ ] Default scope
    - [ ] Response validation
  - [ ] Implement TestHandleRecall
  - [ ] Implement TestHandleReflect
  - [ ] Implement TestHandleForget
  - [ ] Run tests: `go test -v -run TestHandle.*Brain ./internal/dashboard`
  - [ ] Verify coverage on handlers_brain.go: 20% → 85%
  - [ ] Commit: "test: add comprehensive brain HTTP handler tests"

- [ ] **Task 3.2: Search HTTP Handler Tests** (2 hours)

  - [ ] Extend `apps/cli/internal/dashboard/handlers_test.go` or create new file
  - [ ] Implement TestHandleSearch comprehensive (12 cases):
    - [ ] Basic search
    - [ ] With limit
    - [ ] Empty query error
    - [ ] Missing query error
    - [ ] No results
    - [ ] Special characters
    - [ ] Unicode query
    - [ ] No butler
    - [ ] No memory
    - [ ] No corridor
    - [ ] POST not allowed
    - [ ] Large limit
  - [ ] Run tests: `go test -v -run TestHandleSearch ./internal/dashboard`
  - [ ] Verify coverage on handlers_search.go: 30% → 90%
  - [ ] Commit: "test: add comprehensive search HTTP handler tests"

- [ ] **Task 3.3: Context HTTP Handler Tests** (2 hours)

  - [ ] Create `apps/cli/internal/dashboard/handlers_context_test.go`
  - [ ] Implement TestHandleContextPreview (12 cases):
    - [ ] Valid file path
    - [ ] With max tokens
    - [ ] Include learnings
    - [ ] Include decisions
    - [ ] Include failures
    - [ ] Min confidence filter
    - [ ] All options combined
    - [ ] Missing file path
    - [ ] Empty file path
    - [ ] Invalid JSON
    - [ ] GET not allowed
    - [ ] No butler (503)
  - [ ] Implement TestHandleContextPack
  - [ ] Implement TestHandleContextValidate
  - [ ] Run tests: `go test -v -run TestHandleContext ./internal/dashboard`
  - [ ] Verify coverage on handlers_context.go: 0% → 80%
  - [ ] Commit: "test: add context HTTP handler tests"

- [ ] **Task 3.4: Final Integration & Validation** (1 hour)
  - [ ] Run ALL tests: `go test ./...`
  - [ ] Generate full coverage: `make test-coverage`
  - [ ] Verify overall coverage increased by 25-30%
  - [ ] Run race detector: `go test -race ./...`
  - [ ] Run CI/CD locally if possible
  - [ ] Final commit: "test: complete critical test suite implementation"

**End of Day 5 Checkpoint:**

- [ ] 3+ HTTP handler test files complete
- [ ] ~1,200 test lines added
- [ ] Dashboard package coverage: 25% → 70%+
- [ ] Overall project coverage: 55% → 80%+
- [ ] All tests passing
- [ ] No race conditions
- [ ] Ready for PR/review

---

## Final Deliverables

- [ ] **9 Test Files Created:**

  1. ✅ parser_python_test.go (~400 lines)
  2. ✅ parser_javascript_test.go (~350 lines)
  3. ✅ parser_typescript_test.go (~450 lines)
  4. ✅ mcp_tools_brain_test.go (~500 lines)
  5. ✅ mcp_tools_search_test.go (~400 lines)
  6. ✅ mcp_tools_briefing_test.go (~350 lines)
  7. ✅ handlers_brain_test.go (~450 lines)
  8. ✅ handlers_search_test.go (~350 lines)
  9. ✅ handlers_context_test.go (~400 lines)

- [ ] **Documentation:**

  - [x] CRITICAL-TEST-DESIGN.md (detailed specifications)
  - [x] CRITICAL-TEST-QUICK-REF.md (quick reference)
  - [x] SAMPLE-parser_python_test.go (template)
  - [x] SAMPLE-mcp_tools_brain_test.go (template)
  - [ ] Coverage report (HTML)
  - [ ] Test results summary

- [ ] **Metrics Achieved:**
  - [ ] Total test lines: ~3,650
  - [ ] Overall coverage: 55% → 80%+
  - [ ] Language parsers: 0% → 75%+
  - [ ] MCP tools: 20% → 85%+
  - [ ] HTTP handlers: 25% → 85%+
  - [ ] Test execution time: <30 seconds
  - [ ] Zero failing tests
  - [ ] Zero race conditions

---

## Quality Gates

Before marking complete:

- [ ] **Code Quality:**

  - [ ] All tests use table-driven approach
  - [ ] Helper functions extracted and reusable
  - [ ] Clear test names describing behavior
  - [ ] Comprehensive error case coverage
  - [ ] No hardcoded values (use constants)

- [ ] **Coverage:**

  - [ ] Coverage report generated
  - [ ] All target percentages met
  - [ ] No critical paths untested
  - [ ] Edge cases covered

- [ ] **Reliability:**

  - [ ] Tests pass consistently (3x runs)
  - [ ] No flaky tests
  - [ ] Tests run in isolation
  - [ ] No race conditions detected
  - [ ] Cleanup properly implemented

- [ ] **Documentation:**
  - [ ] Test purposes documented
  - [ ] Complex test logic explained
  - [ ] Mock/stub patterns documented
  - [ ] Coverage gaps identified

---

## Commands Reference

```bash
# Run specific package tests
go test -v ./internal/analysis/...
go test -v ./internal/butler/...
go test -v ./internal/dashboard/...

# Run specific test
go test -v -run TestPythonParser ./internal/analysis

# Coverage for package
go test -coverprofile=coverage.out ./internal/analysis
go tool cover -html=coverage.out

# Coverage summary
go tool cover -func=coverage.out

# Race detection
go test -race ./...

# Run all tests
go test ./...

# Full coverage (via Makefile)
make test-coverage

# Benchmark (if needed)
go test -bench=. -benchmem ./internal/analysis
```

---

## Troubleshooting

**Issue:** Tests fail due to missing dependencies

- **Fix:** Run `go mod tidy` and `go mod download`

**Issue:** Coverage not improving

- **Fix:** Ensure test cases actually exercise the code paths
- Use `go test -coverprofile` and check HTML report

**Issue:** Race conditions detected

- **Fix:** Use proper locking, avoid shared state in tests
- Use `t.Parallel()` carefully

**Issue:** Tests are flaky

- **Fix:** Ensure proper cleanup with `t.Cleanup()`
- Avoid time-dependent assertions
- Use test isolation

**Issue:** Slow tests

- **Fix:** Use `t.Parallel()` for independent tests
- Mock external dependencies
- Avoid unnecessary sleeps

---

## Success Criteria

✅ **Complete when:**

1. All 9 test files created and passing
2. Coverage increased by 25-30% overall
3. No failing tests in CI/CD
4. No race conditions
5. Test execution time <30 seconds
6. Code reviewed and approved
7. Documentation updated

---

## Next Steps (Post-Implementation)

After completing this sprint:

1. **Identify remaining gaps:**

   - Less common parsers (Rust, Scala, etc.)
   - Secondary MCP tools
   - Edge case handlers

2. **Integration tests:**

   - Multi-component workflows
   - Real workspace scenarios
   - End-to-end MCP flows

3. **Performance tests:**

   - Parser benchmarks
   - Search performance
   - Memory query optimization

4. **Continuous improvement:**
   - Regular coverage reviews
   - Test maintenance
   - Refactoring as needed

---

**Last Updated:** January 8, 2026  
**Owner:** Development Team  
**Status:** 🔴 Ready to Start
