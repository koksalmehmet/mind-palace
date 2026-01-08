# Critical Test Design - Index

**Research Date:** January 8, 2026  
**Status:** ✅ Research Complete, Ready for Implementation  
**Goal:** +25-30% test coverage via 9 critical test files

---

## 📋 Document Index

All documents are in `sprint-tasks/` directory:

### 1. Executive Summary

**[TEST-RESEARCH-SUMMARY.md](TEST-RESEARCH-SUMMARY.md)**

- Research findings and analysis
- Coverage gap identification
- ROI analysis and risk assessment
- **Read this FIRST** for overview

### 2. Detailed Test Specifications

**[CRITICAL-TEST-DESIGN.md](CRITICAL-TEST-DESIGN.md)**

- Complete test case specifications (2,300 lines)
- All 119 test cases detailed
- Mock/stub requirements
- Expected coverage per file
- **Use this** for detailed implementation guidance

### 3. Quick Reference Guide

**[CRITICAL-TEST-QUICK-REF.md](CRITICAL-TEST-QUICK-REF.md)**

- Quick lookup for test patterns
- Common commands
- Helper function templates
- Priority ordering
- **Use this** during implementation for quick reference

### 4. Implementation Checklist

**[TEST-IMPLEMENTATION-CHECKLIST.md](TEST-IMPLEMENTATION-CHECKLIST.md)**

- Day-by-day task breakdown (5 days)
- Quality gates and checkpoints
- Success criteria
- Troubleshooting guide
- **Follow this** step-by-step during implementation

### 5. Code Templates

**[SAMPLE-parser_python_test.go](SAMPLE-parser_python_test.go)**

- Complete Python parser test example (350 lines)
- 18 test cases with patterns
- Helper functions
- **Copy this** as template for parser tests

**[SAMPLE-mcp_tools_brain_test.go](SAMPLE-mcp_tools_brain_test.go)**

- Complete MCP tool test example (300 lines)
- Table-driven test patterns
- Mock strategies
- **Copy this** as template for MCP tool tests

---

## 🎯 Quick Start

### For Implementation Team

1. **Start here:** Read [TEST-RESEARCH-SUMMARY.md](TEST-RESEARCH-SUMMARY.md)
2. **Plan work:** Review [TEST-IMPLEMENTATION-CHECKLIST.md](TEST-IMPLEMENTATION-CHECKLIST.md)
3. **Implement:** Follow day-by-day tasks using templates
4. **Reference:** Keep [CRITICAL-TEST-QUICK-REF.md](CRITICAL-TEST-QUICK-REF.md) handy
5. **Details:** Check [CRITICAL-TEST-DESIGN.md](CRITICAL-TEST-DESIGN.md) for specifics

### For Review/Approval

1. Read [TEST-RESEARCH-SUMMARY.md](TEST-RESEARCH-SUMMARY.md) - Executive overview
2. Review [CRITICAL-TEST-DESIGN.md](CRITICAL-TEST-DESIGN.md) - Verify approach
3. Check [TEST-IMPLEMENTATION-CHECKLIST.md](TEST-IMPLEMENTATION-CHECKLIST.md) - Approve timeline

### For Ongoing Reference

- **Commands:** [CRITICAL-TEST-QUICK-REF.md](CRITICAL-TEST-QUICK-REF.md)
- **Patterns:** Sample files (SAMPLE-\*.go)
- **Coverage targets:** [CRITICAL-TEST-DESIGN.md](CRITICAL-TEST-DESIGN.md)

---

## 📊 Test File Summary

| #         | File                       | Lines     | Tests   | Coverage Δ  | Priority |
| --------- | -------------------------- | --------- | ------- | ----------- | -------- |
| 1         | parser_python_test.go      | 400       | 15      | +75%        | P0       |
| 2         | parser_javascript_test.go  | 350       | 15      | +75%        | P0       |
| 3         | parser_typescript_test.go  | 450       | 18      | +80%        | P0       |
| 4         | mcp_tools_brain_test.go    | 500       | 12+     | +70%        | P0       |
| 5         | mcp_tools_search_test.go   | 400       | 10+     | +65%        | P1       |
| 6         | mcp_tools_briefing_test.go | 350       | 10      | +75%        | P1       |
| 7         | handlers_brain_test.go     | 450       | 15      | +65%        | P1       |
| 8         | handlers_search_test.go    | 350       | 12      | +60%        | P1       |
| 9         | handlers_context_test.go   | 400       | 12      | +80%        | P2       |
| **Total** | **9 files**                | **3,650** | **119** | **+25-30%** | -        |

---

## 📅 Timeline

| Day | Date   | Focus                 | Deliverable                           |
| --- | ------ | --------------------- | ------------------------------------- |
| 1   | Jan 8  | Python & JS parsers   | 2 test files, +20% analysis coverage  |
| 2   | Jan 9  | TypeScript parser     | 1 test file, +10% analysis coverage   |
| 3   | Jan 10 | Brain MCP tool        | 1 test file, +15% butler coverage     |
| 4   | Jan 11 | Search & Briefing MCP | 2 test files, +25% butler coverage    |
| 5   | Jan 12 | HTTP handlers         | 3 test files, +20% dashboard coverage |

**Total:** 5 days, 9 files, +25-30% overall coverage

---

## ✅ Deliverables Checklist

### Research & Design (COMPLETE ✅)

- [x] Coverage gap analysis
- [x] Test prioritization
- [x] Detailed test specifications
- [x] Implementation timeline
- [x] Code templates
- [x] Quick reference guides
- [x] Implementation checklist

### Implementation (PENDING 🔴)

- [ ] 3 parser test files
- [ ] 3 MCP tool test files
- [ ] 3 HTTP handler test files
- [ ] Coverage reports
- [ ] All tests passing
- [ ] No race conditions
- [ ] Documentation updates

---

## 🔑 Key Files by Purpose

### For Understanding the Approach

- [TEST-RESEARCH-SUMMARY.md](TEST-RESEARCH-SUMMARY.md) - Why these tests?
- [CRITICAL-TEST-DESIGN.md](CRITICAL-TEST-DESIGN.md) - What exactly to test?

### For Implementation

- [TEST-IMPLEMENTATION-CHECKLIST.md](TEST-IMPLEMENTATION-CHECKLIST.md) - Day-by-day tasks
- [SAMPLE-parser_python_test.go](SAMPLE-parser_python_test.go) - Parser template
- [SAMPLE-mcp_tools_brain_test.go](SAMPLE-mcp_tools_brain_test.go) - MCP template

### For Reference During Work

- [CRITICAL-TEST-QUICK-REF.md](CRITICAL-TEST-QUICK-REF.md) - Commands & patterns
- [CRITICAL-TEST-DESIGN.md](CRITICAL-TEST-DESIGN.md) - Test case details

---

## 🎓 Learning Resources

### Go Testing Best Practices

- Table-driven tests: https://go.dev/wiki/TableDrivenTests
- Testing package: https://pkg.go.dev/testing
- httptest: https://pkg.go.dev/net/http/httptest

### Project-Specific Patterns

- Existing parser tests: `apps/cli/internal/analysis/parser_cpp_test.go`
- Existing MCP tests: `apps/cli/internal/butler/mcp_test.go`
- Existing handler tests: `apps/cli/internal/dashboard/handlers_test.go`

### Templates Provided

- Parser pattern: [SAMPLE-parser_python_test.go](SAMPLE-parser_python_test.go)
- MCP tool pattern: [SAMPLE-mcp_tools_brain_test.go](SAMPLE-mcp_tools_brain_test.go)

---

## 🚀 Getting Started Commands

```bash
# Navigate to project
cd C:\git\mind-palace\apps\cli

# Start with Day 1, Task 1.1
# Create Python parser tests
touch internal/analysis/parser_python_test.go

# Copy template as starting point
# (See SAMPLE-parser_python_test.go)

# Run tests
go test -v -run TestPythonParser ./internal/analysis

# Check coverage
go test -coverprofile=coverage.out ./internal/analysis
go tool cover -html=coverage.out

# Continue with checklist...
```

---

## 📈 Success Metrics

### Quantitative

- [ ] Overall coverage: 55% → 80%+ ✅
- [ ] Parser coverage: 0% → 75%+ ✅
- [ ] MCP coverage: 20% → 85%+ ✅
- [ ] Handler coverage: 25% → 85%+ ✅
- [ ] Test count: 119 tests ✅
- [ ] Test LOC: ~3,650 lines ✅

### Qualitative

- [ ] All tests passing ✅
- [ ] No flaky tests ✅
- [ ] Execution time <30s ✅
- [ ] Table-driven patterns ✅
- [ ] Clear error messages ✅
- [ ] Good test coverage ✅

---

## 💬 Questions?

### Research Questions

- See [TEST-RESEARCH-SUMMARY.md](TEST-RESEARCH-SUMMARY.md) - Research methodology section

### Design Questions

- See [CRITICAL-TEST-DESIGN.md](CRITICAL-TEST-DESIGN.md) - Specific test cases

### Implementation Questions

- See [TEST-IMPLEMENTATION-CHECKLIST.md](TEST-IMPLEMENTATION-CHECKLIST.md) - Troubleshooting section
- See [CRITICAL-TEST-QUICK-REF.md](CRITICAL-TEST-QUICK-REF.md) - Common patterns

### Code Questions

- See template files (SAMPLE-\*.go)
- Check existing tests in codebase

---

## 📝 Updates & Changes

| Date        | Change                      | Updated By   |
| ----------- | --------------------------- | ------------ |
| Jan 8, 2026 | Initial research and design | AI Assistant |
| -           | -                           | -            |

---

## 🎯 Status Overview

| Component      | Status         | Progress |
| -------------- | -------------- | -------- |
| Research       | ✅ Complete    | 100%     |
| Design         | ✅ Complete    | 100%     |
| Documentation  | ✅ Complete    | 100%     |
| Templates      | ✅ Complete    | 100%     |
| Implementation | 🔴 Not Started | 0%       |
| Validation     | ⚪ Pending     | 0%       |

---

## 📞 Next Actions

1. **Review** this index and research summary
2. **Approve** the implementation plan
3. **Start** Day 1 tasks from checklist
4. **Follow** templates for each test file
5. **Track** progress against checkpoints
6. **Report** completion with coverage metrics

---

**Ready to implement. All design work complete. Let's boost that coverage! 🚀**

---

_This index was generated as part of the critical test design research for Mind Palace project._
