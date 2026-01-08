# Mind Palace - Integration Testing Strategy

**Date:** January 8, 2026  
**Status:** RESEARCH & DESIGN DOCUMENT  
**Author:** GitHub Copilot

---

## Executive Summary

This document outlines a comprehensive integration testing strategy for Mind Palace, covering cross-component communication between the Go CLI, VS Code extension, Angular dashboard, and MCP server. The strategy focuses on critical user flows and data contract validation across the ecosystem.

### Architecture Overview

Mind Palace consists of four interconnected components:

```
┌─────────────────┐         ┌──────────────────┐
│  VS Code Ext    │◄────────│   Palace CLI     │
│  (TypeScript)   │  spawn  │   (Go Binary)    │
└────────┬────────┘         └────────┬─────────┘
         │                           │
         │ HTTP API                  │ JSON-RPC (stdio)
         ▼                           ▼
┌─────────────────┐         ┌──────────────────┐
│  Dashboard      │◄────────│   MCP Server     │
│  (Angular 21)   │  HTTP   │   (Go Butler)    │
└─────────────────┘         └──────────────────┘
```

**Key Integration Points:**

1. **VS Code ↔ CLI**: Child process spawning, stdout/stderr parsing, JSON-RPC (MCP)
2. **Dashboard ↔ CLI**: REST API over HTTP (port 3001)
3. **MCP ↔ Agents**: JSON-RPC over stdio (AI agent communication)
4. **Schemas**: Shared data contracts (JSON schemas in `apps/cli/schemas/`)
5. **Config Files**: `.palace/config.jsonc`, `palace.jsonc`, room definitions

---

## Current Test Coverage Analysis

### Existing Test Infrastructure

#### ✅ CLI (Go) - Well Covered

- **Unit Tests**: ~74% coverage (78 test files for 105 source files)
- **Integration Tests**:
  - [apps/cli/tests/integration_test.go](apps/cli/tests/integration_test.go) - Basic lifecycle
  - [apps/cli/tests/e2e_flows_test.go](apps/cli/tests/e2e_flows_test.go) - E2E workflows
  - [apps/cli/tests/uat_test.go](apps/cli/tests/uat_test.go) - User acceptance tests
- **Test Framework**: Go standard library `testing` package
- **Coverage Tool**: `go test -cover`

#### ✅ VS Code Extension - Well Covered

- **Unit Tests**: 100% (49/49 passing)
- **Test Files**: `apps/vscode/src/test/suite/*.test.ts`
- **Test Framework**: Mocha + Chai + Sinon
- **Mocking**: VS Code API mocking, `child_process` stubbing
- **Coverage**: Tests for MCP communication, bridge, config, commands

#### ✅ Dashboard (Angular) - Well Covered

- **Unit Tests**: 100% (211/211 passing)
- **Test Framework**: Vitest + Testing Library
- **Coverage**: Component tests, service tests, HTTP mocking

#### ❌ Cross-Component Integration - **MISSING**

- No tests validating VS Code ↔ CLI communication
- No tests validating Dashboard ↔ API communication
- No end-to-end flows across multiple components
- No schema validation tests across boundaries
- No MCP client integration tests

### Current Test Gaps

1. **VS Code Extension → CLI Communication**

   - Spawning `palace serve` and validating MCP responses
   - Handling CLI version mismatches
   - Error handling when CLI is not installed
   - Config file changes triggering extension updates

2. **Dashboard → API Integration**

   - Real HTTP requests to dashboard server
   - WebSocket real-time updates
   - CORS configuration validation
   - Workspace switching

3. **End-to-End User Flows**

   - Store learning via CLI → View in VS Code extension
   - Start session in VS Code → View in Dashboard
   - Config change → Auto-reload in all components

4. **MCP Protocol Validation**

   - JSON-RPC request/response conformance
   - Tool calling with real CLI backend
   - Resource reading via MCP

5. **Schema Contract Validation**
   - JSON schema validation across all components
   - Breaking change detection in schemas

---

## Recommended Integration Test Framework

### Primary Framework: **Playwright + Custom Scripts**

**Rationale:**

- **Playwright** for browser automation (Dashboard testing)
- **Custom Go/TypeScript scripts** for CLI and extension integration
- **Testcontainers** (optional) for isolated test environments
- **JSON Schema validation** libraries for contract testing

### Alternative Considered: Pure E2E with Puppeteer

- **Rejected**: Less suited for multi-process orchestration (CLI + extension + dashboard)
- Playwright has better TypeScript support and VS Code integration

---

## Integration Test Scenarios (Critical Paths)

### Scenario 1: VS Code Extension ↔ CLI (MCP Communication)

**Priority**: 🔴 Critical  
**Description**: Validate that VS Code extension can spawn CLI MCP server and communicate correctly

**Test Steps:**

1. Build `palace` CLI binary
2. Initialize a test workspace with `palace init` and `palace scan`
3. Mock VS Code workspace pointing to test workspace
4. Start extension, verify it spawns `palace serve`
5. Call `bridge.getBrief()` → Verify JSON-RPC request/response
6. Call `bridge.store()` with learning → Verify storage via `palace recall`
7. Terminate extension → Verify MCP server cleanup

**Expected Outcomes:**

- MCP server starts successfully
- JSON-RPC requests follow protocol (see `apps/cli/internal/butler/mcp.go`)
- Tools return valid responses matching TypeScript types
- No zombie processes after cleanup

**Test Files:**

- `tests/integration/vscode-cli/mcp-communication.test.ts`

**Dependencies:**

- Compiled `palace` binary in test environment
- Mock VS Code API (already available in extension tests)
- JSON-RPC validation library

**Estimated Complexity**: Medium (3-5 days)

---

### Scenario 2: Dashboard ↔ API (HTTP Communication)

**Priority**: 🔴 Critical  
**Description**: Validate dashboard communicates correctly with CLI dashboard server

**Test Steps:**

1. Build `palace` CLI binary
2. Initialize test workspace and run `palace scan`
3. Start `palace dashboard --port 3001 --no-browser` in background
4. Use Playwright to navigate to `http://localhost:3001`
5. Verify dashboard loads and displays project data
6. Test API endpoints:
   - `GET /api/health` → Verify status
   - `GET /api/rooms` → Verify room data
   - `POST /api/remember` → Store learning
   - `GET /api/learnings` → Retrieve learning
7. Test WebSocket connection for real-time updates
8. Shutdown server gracefully

**Expected Outcomes:**

- Dashboard server starts successfully
- All API endpoints return valid JSON
- WebSocket connection established
- CORS headers present for allowed origins
- UI displays data correctly

**Test Files:**

- `tests/integration/dashboard-api/api-endpoints.spec.ts`
- `tests/integration/dashboard-api/websocket.spec.ts`

**Dependencies:**

- Playwright
- Compiled `palace` binary
- HTTP client for API testing (Playwright's built-in `request`)

**Estimated Complexity**: Medium (4-6 days)

---

### Scenario 3: End-to-End - Store Learning via CLI, Retrieve via Extension

**Priority**: 🟡 High  
**Description**: Validate data flows across CLI and VS Code extension

**Test Steps:**

1. Build `palace` CLI binary
2. Initialize test workspace
3. **Via CLI**: Run `palace store --as learning "Test learning content"`
4. Capture learning ID from CLI output
5. **Via Extension**: Mock extension environment, call `bridge.recallLearnings()`
6. Verify learning appears in extension response
7. **Via CLI**: Run `palace recall --type learning` → Verify same learning

**Expected Outcomes:**

- Learning stored via CLI is immediately available
- Extension retrieves learning with correct ID, content, and metadata
- Data consistency across CLI and extension

**Test Files:**

- `tests/integration/e2e/cli-to-extension.test.ts`

**Dependencies:**

- Compiled `palace` binary
- Mock VS Code environment
- Shared test workspace

**Estimated Complexity**: Low-Medium (2-3 days)

---

### Scenario 4: Config File Changes → Extension Auto-Reload

**Priority**: 🟡 High  
**Description**: Validate VS Code extension reacts to config file changes

**Test Steps:**

1. Build `palace` CLI binary and initialize workspace
2. Start VS Code extension with mocked file watcher
3. Modify `palace.jsonc` config file (e.g., change project name)
4. Trigger file watcher event
5. Verify extension reloads config via `watchProjectConfig()`
6. Verify extension calls `palace check` to verify workspace state

**Expected Outcomes:**

- File watcher detects config changes
- Extension reloads configuration
- Status check triggered automatically
- HUD updates with new config values

**Test Files:**

- `tests/integration/vscode-cli/config-watcher.test.ts`

**Dependencies:**

- Mock VS Code file system watcher
- Test workspace with modifiable config

**Estimated Complexity**: Low (1-2 days)

---

### Scenario 5: MCP Client Integration Flow (AI Agent Simulation)

**Priority**: 🟡 High  
**Description**: Validate MCP server handles tool calls from simulated AI agent

**Test Steps:**

1. Build `palace` CLI binary and initialize workspace
2. Start `palace serve` in test mode (stdio)
3. Send JSON-RPC `initialize` request
4. Verify `initialize` response with capabilities
5. Send `tools/list` request → Verify all tools listed (see [bridge.ts#L12-L63](apps/vscode/src/bridge.ts#L12-L63))
6. Call `tools/call` with `explore_rooms` → Verify room data returned
7. Call `tools/call` with `store` → Store learning
8. Call `tools/call` with `recall` → Verify learning retrieved
9. Send `shutdown` notification
10. Verify clean shutdown

**Expected Outcomes:**

- MCP server responds to all JSON-RPC methods
- Tool schemas match expected format
- Tool execution returns valid results
- No protocol violations

**Test Files:**

- `tests/integration/mcp-server/agent-simulation.test.go`

**Dependencies:**

- JSON-RPC client library (Go)
- MCP protocol spec (see `apps/cli/internal/butler/mcp.go`)

**Estimated Complexity**: Medium-High (5-7 days)

---

### Scenario 6: Schema Contract Validation

**Priority**: 🟢 Medium  
**Description**: Validate all components respect shared JSON schemas

**Test Steps:**

1. Load all JSON schemas from `apps/cli/schemas/`
2. **CLI Output Validation**:
   - Run `palace scan` → Validate output against `scan.schema.json`
   - Run `palace check --collect` → Validate `context-pack.schema.json`
3. **MCP Response Validation**:
   - Call MCP tools → Validate responses against expected schemas
4. **Dashboard API Validation**:
   - Call API endpoints → Validate responses against schemas
5. **Extension Bridge Validation**:
   - Call bridge methods → Validate TypeScript types match schemas

**Expected Outcomes:**

- All JSON outputs conform to schemas
- No breaking changes in data contracts
- TypeScript types align with JSON schemas

**Test Files:**

- `tests/integration/schema-validation/cli-output.test.go`
- `tests/integration/schema-validation/mcp-responses.test.ts`
- `tests/integration/schema-validation/api-responses.test.ts`

**Dependencies:**

- JSON Schema validator (Go: `github.com/xeipuuv/gojsonschema`)
- JSON Schema validator (TS: `ajv`)

**Estimated Complexity**: Medium (3-4 days)

---

### Scenario 7: Multi-Workspace Corridor Flow

**Priority**: 🟢 Medium  
**Description**: Validate corridor features work across workspaces

**Test Steps:**

1. Create workspace A and B
2. Initialize both with `palace init` and `palace scan`
3. **In Workspace A**: Store learning with `--corridor` flag
4. **In Workspace B**: Retrieve learning via `palace corridor learnings`
5. **Via Dashboard**: Switch between workspaces, verify corridor data syncs
6. **Via Extension**: Verify corridor view shows cross-workspace learnings

**Expected Outcomes:**

- Learnings promoted to corridor are accessible across workspaces
- Corridor database (`~/.palace/corridor/global.db`) updated correctly
- Dashboard and extension show corridor data

**Test Files:**

- `tests/integration/e2e/corridor-cross-workspace.test.go`

**Dependencies:**

- Multiple test workspaces
- Global corridor database access

**Estimated Complexity**: Medium-High (4-5 days)

---

## Test Placement Strategy

### Recommended Directory Structure

```
mind-palace/
├── apps/
│   ├── cli/
│   │   └── tests/           # Existing Go tests (keep as-is)
│   ├── vscode/
│   │   └── src/test/suite/  # Existing extension tests (keep as-is)
│   └── dashboard/
│       └── src/test/        # Existing Angular tests (keep as-is)
├── tests/                    # 🆕 NEW: Cross-component integration tests
│   ├── integration/
│   │   ├── vscode-cli/      # Extension ↔ CLI tests
│   │   │   ├── mcp-communication.test.ts
│   │   │   └── config-watcher.test.ts
│   │   ├── dashboard-api/   # Dashboard ↔ API tests
│   │   │   ├── api-endpoints.spec.ts
│   │   │   └── websocket.spec.ts
│   │   ├── mcp-server/      # MCP protocol tests
│   │   │   └── agent-simulation.test.go
│   │   ├── e2e/             # Full end-to-end flows
│   │   │   ├── cli-to-extension.test.ts
│   │   │   └── corridor-cross-workspace.test.go
│   │   └── schema-validation/ # Contract tests
│   │       ├── cli-output.test.go
│   │       ├── mcp-responses.test.ts
│   │       └── api-responses.test.ts
│   ├── fixtures/            # Shared test data
│   │   ├── workspaces/      # Sample workspaces
│   │   ├── configs/         # Test configs
│   │   └── schemas/         # Schema copies for validation
│   ├── utils/               # Shared test utilities
│   │   ├── cli-wrapper.ts   # Helper to run CLI commands
│   │   ├── workspace-setup.ts # Workspace initialization
│   │   └── mcp-client.go    # MCP client helper
│   ├── package.json         # TypeScript test dependencies
│   ├── tsconfig.json        # TypeScript config for tests
│   ├── go.mod               # Go test dependencies
│   └── README.md            # Integration test documentation
└── scripts/
    └── run-integration-tests.sh # Test runner script
```

### Rationale for Separate `tests/` Directory

1. **Cross-Component**: Integration tests span multiple apps
2. **Shared Fixtures**: Reusable test workspaces and data
3. **Independent Lifecycle**: Can run independently of component builds
4. **CI/CD Friendly**: Single entry point for integration tests

---

## Dependencies & Setup Requirements

### Language-Specific Dependencies

#### TypeScript (VS Code + Dashboard Integration Tests)

```json
// tests/package.json
{
  "devDependencies": {
    "@playwright/test": "^1.48.0",
    "@types/node": "^20.10.0",
    "typescript": "^5.9.3",
    "vitest": "^2.1.8",
    "ajv": "^8.12.0", // JSON Schema validation
    "ts-node": "^10.9.2",
    "execa": "^8.0.1" // CLI command execution
  }
}
```

#### Go (CLI + MCP Integration Tests)

```go
// tests/go.mod
require (
    github.com/stretchr/testify v1.9.0  // Assertions
    github.com/xeipuuv/gojsonschema v1.2.0 // JSON Schema validation
)
```

### Build Dependencies

1. **Compiled CLI Binary**: Tests require `palace` binary

   - Built via `make build-palace` or `go build ./apps/cli`
   - Place in `tests/bin/palace` or use `$PATH`

2. **VS Code Extension**: Compiled extension code

   - Built via `npm run compile` in `apps/vscode/`

3. **Dashboard Assets** (for Playwright tests):
   - Built via `npm run build` in `apps/dashboard/`
   - Or run dev server during tests

### Test Environment Setup

1. **Isolated Test Workspaces**: Each test creates temporary workspace
2. **Clean State**: Tests should not interfere with each other
3. **Port Management**: Dashboard server should use random available ports
4. **Process Cleanup**: Ensure CLI processes are terminated after tests

---

## Continuous Integration Strategy

### GitHub Actions Workflow

```yaml
# .github/workflows/integration-tests.yml
name: Integration Tests

on:
  pull_request:
    branches: [main]
  push:
    branches: [main]

jobs:
  integration-tests:
    runs-on: ubuntu-latest
    timeout-minutes: 20

    steps:
      - uses: actions/checkout@v4

      - name: Setup Go
        uses: actions/setup-go@v5
        with:
          go-version: "1.23"

      - name: Setup Node.js
        uses: actions/setup-node@v4
        with:
          node-version: "20"

      - name: Build Palace CLI
        run: make build-palace

      - name: Install Integration Test Dependencies
        run: |
          cd tests
          npm install
          go mod download

      - name: Run Integration Tests
        run: |
          cd tests
          npm run test:integration
          go test -v ./integration/...

      - name: Upload Test Results
        if: always()
        uses: actions/upload-artifact@v4
        with:
          name: integration-test-results
          path: tests/results/
```

### Test Execution Order

1. **Unit Tests** (existing) - Fast, no dependencies
2. **Integration Tests** (new) - Medium speed, requires binaries
3. **E2E Tests** (new) - Slowest, full system tests

### Parallel Execution

- Group independent integration tests for parallel execution
- Use test matrix for different OS environments (Linux, macOS, Windows)

---

## Testing Tools & Libraries Summary

| Tool/Library     | Purpose                | Language   | Used In                      |
| ---------------- | ---------------------- | ---------- | ---------------------------- |
| **Playwright**   | Browser automation     | TypeScript | Dashboard integration tests  |
| **Vitest**       | Test runner            | TypeScript | General TS integration tests |
| **Mocha/Chai**   | Assertion library      | TypeScript | VS Code extension tests      |
| **Ajv**          | JSON Schema validation | TypeScript | Schema contract tests        |
| **Execa**        | Process execution      | TypeScript | CLI wrapper in tests         |
| **Go testing**   | Test framework         | Go         | CLI and MCP server tests     |
| **gojsonschema** | JSON Schema validation | Go         | Schema validation in Go      |
| **testify**      | Assertions             | Go         | Go integration tests         |

---

## Implementation Roadmap

### Phase 1: Foundation (Week 1-2)

- [ ] Create `tests/` directory structure
- [ ] Set up TypeScript and Go test configurations
- [ ] Implement shared test utilities (CLI wrapper, workspace setup)
- [ ] Add CI workflow for integration tests

### Phase 2: Critical Path Tests (Week 3-4)

- [ ] **Scenario 1**: VS Code Extension ↔ CLI (MCP Communication)
- [ ] **Scenario 2**: Dashboard ↔ API (HTTP Communication)
- [ ] **Scenario 3**: End-to-End - Store Learning via CLI, Retrieve via Extension

### Phase 3: Additional Coverage (Week 5-6)

- [ ] **Scenario 4**: Config File Changes → Extension Auto-Reload
- [ ] **Scenario 5**: MCP Client Integration Flow
- [ ] **Scenario 6**: Schema Contract Validation

### Phase 4: Advanced Flows (Week 7-8)

- [ ] **Scenario 7**: Multi-Workspace Corridor Flow
- [ ] Performance benchmarks for integration tests
- [ ] Flakiness detection and stabilization

### Phase 5: Documentation & Maintenance (Week 9)

- [ ] Integration test documentation
- [ ] Developer onboarding guide
- [ ] Test maintenance runbook

---

## Success Metrics

### Coverage Targets

- **Critical Paths**: 100% (all 7 scenarios implemented)
- **API Endpoints**: 80% (key endpoints tested)
- **MCP Tools**: 75% (most commonly used tools)
- **Schema Validation**: 100% (all schemas validated)

### Quality Metrics

- **Test Flakiness**: < 2% (tests should be deterministic)
- **Execution Time**: < 10 minutes for full suite
- **False Positive Rate**: < 1%

### Maintenance Metrics

- **Documentation**: All tests documented with purpose and steps
- **CI Integration**: Tests run on every PR
- **Developer Experience**: Tests can run locally without complex setup

---

## Risks & Mitigation

### Risk 1: Test Flakiness (High Probability, High Impact)

**Mitigation:**

- Use explicit waits instead of implicit timeouts
- Implement retry logic for network operations
- Clean up processes and files deterministically
- Use unique ports for each test run

### Risk 2: Long Execution Time (Medium Probability, Medium Impact)

**Mitigation:**

- Run integration tests in parallel where possible
- Use test matrix for different scenarios
- Cache built binaries in CI
- Implement smart test selection (only run affected tests)

### Risk 3: Environment-Specific Failures (Medium Probability, Medium Impact)

**Mitigation:**

- Test on multiple OS platforms (Linux, macOS, Windows)
- Use containerization for consistent environments
- Document platform-specific quirks
- Use cross-platform path handling (`path.join`)

### Risk 4: Breaking Changes in Dependencies (Low Probability, High Impact)

**Mitigation:**

- Pin exact versions in `package.json` and `go.mod`
- Test dependency upgrades in isolated branches
- Monitor VS Code API changes
- Subscribe to Playwright and testing library changelogs

---

## Maintenance & Evolution

### Regular Maintenance Tasks

1. **Weekly**: Review test failures and flakiness reports
2. **Monthly**: Update dependencies and test frameworks
3. **Quarterly**: Review and refactor slow or brittle tests
4. **Per Release**: Add tests for new features before release

### Test Expansion Areas

1. **Performance Tests**: Add load testing for MCP server and dashboard
2. **Security Tests**: Validate input sanitization across components
3. **Accessibility Tests**: Dashboard WCAG compliance
4. **Browser Matrix**: Test dashboard on Chrome, Firefox, Safari

---

## Conclusion

This integration testing strategy provides:

1. ✅ **7 Critical Test Scenarios** covering all major integration points
2. ✅ **Clear Test Placement** in dedicated `tests/` directory
3. ✅ **Recommended Frameworks** (Playwright + Custom Scripts)
4. ✅ **Detailed Dependencies** (TypeScript and Go packages)
5. ✅ **Estimated Complexity** (23-32 days total for all scenarios)
6. ✅ **CI/CD Integration** (GitHub Actions workflow)
7. ✅ **Risk Mitigation** strategies for common failure modes

### Next Steps (Implementation Phase)

1. **Review & Approve**: Review this document with team
2. **Setup Foundation**: Create `tests/` directory and base configuration
3. **Prioritize Scenarios**: Start with Scenario 1 and 2 (critical paths)
4. **Implement Iteratively**: One scenario at a time, validate before moving forward
5. **Document Learnings**: Update this document as tests are implemented

### Total Estimated Effort

| Phase                                        | Duration    | Complexity      |
| -------------------------------------------- | ----------- | --------------- |
| Phase 1: Foundation                          | 2 weeks     | Low             |
| Phase 2: Critical Paths (Scenarios 1-3)      | 2 weeks     | Medium-High     |
| Phase 3: Additional Coverage (Scenarios 4-6) | 2 weeks     | Medium          |
| Phase 4: Advanced Flows (Scenario 7)         | 1 week      | Medium-High     |
| Phase 5: Documentation                       | 1 week      | Low             |
| **Total**                                    | **8 weeks** | **Medium-High** |

**Recommended Team Size**: 1-2 engineers dedicated to integration testing

---

**Document Version**: 1.0  
**Last Updated**: January 8, 2026  
**Status**: Ready for Review ✅
