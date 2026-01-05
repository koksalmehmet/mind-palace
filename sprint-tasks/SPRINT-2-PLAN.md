# Sprint 2 Implementation Plan: Testing & Stability

**Date:** January 5, 2026  
**Version:** 1.0-RC1  
**Status:** 🟢 Ready for Execution  
**Prepared by:** AI Engineering Team with Sub-Agent Research  

---

## Executive Summary

Sprint 1 successfully completed **6 critical fixes** (version sync, CORS security, MIT license, documentation links, deprecated code removal, TypeScript strict mode). Sprint 2 focuses on **closing the frontend testing gap** and **stabilizing the system** for Beta release.

**Key Metrics:**
- **Current Test Coverage:** CLI 95% | Dashboard 0% | VS Code 0%
- **Target Coverage:** CLI 95% | Dashboard 70%+ | VS Code 70%+
- **Estimated Timeline:** 3-4 weeks
- **Risk Level:** LOW (all changes are additive, no breaking changes)

---

## What Was Accomplished (Sprint 1)

✅ Version synchronization across all apps  
✅ WebSocket CORS security implementation  
✅ Documentation link validation  
✅ MIT license modernization  
✅ Deprecated code removal  
✅ TypeScript strict mode enablement  

**Result:** 38 files modified, 3,757 insertions, production security baseline established.

---

## Sprint 2: Phase-by-Phase Roadmap

### Phase 1: Frontend Test Infrastructure (Week 1)

**Goal:** Set up testing frameworks and write foundational test cases.

#### 1.1 Dashboard Testing with Vitest

**Why Vitest?**
- 5x faster than Jest
- Angular 21 native support
- Built-in browser mode for D3.js/Cytoscape testing
- Lower configuration overhead

**Implementation:**

```bash
# Install Vitest
npm install -D vitest @vitest/ui @vitest/browser
npm install -D @testing-library/angular @testing-library/dom

# Create vitest config
echo "import { defineConfig } from 'vitest/config';
export default defineConfig({
  test: {
    globals: true,
    environment: 'jsdom',
    coverage: { provider: 'v8', reporter: ['text', 'json', 'html'] }
  }
});" > apps/dashboard/vitest.config.ts
```

**Test Structure:**

```
apps/dashboard/src/
├── app/
│   ├── core/
│   │   ├── services/
│   │   │   ├── websocket.service.ts
│   │   │   └── websocket.service.spec.ts          [NEW]
│   │   └── guards/
│   ├── features/
│   │   ├── overview/
│   │   │   ├── neural-map/
│   │   │   │   ├── neural-map.component.ts
│   │   │   │   └── neural-map.component.spec.ts   [NEW]
│   │   │   └── ...
│   ├── layouts/
│   └── shared/
└── test/
    └── setup.ts                                    [NEW]
```

**Priority Test Files (Week 1):**

1. `websocket.service.spec.ts` - WebSocket connection, reconnection, message handling
2. `api.service.spec.ts` - HTTP interceptors, error handling
3. `neural-map.component.spec.ts` - D3.js visualization rendering
4. `contradictions.component.spec.ts` - Complex UI logic

**Target:** 25 test files, 15+ hours effort

#### 1.2 VS Code Extension Testing

**Framework:** @vscode/test-electron with Mocha/Sinon

**Installation:**

```bash
npm install -D @vscode/test-electron mocha chai sinon
```

**Test Structure:**

```
apps/vscode/src/
├── extension.ts
├── extension.test.ts                              [NEW]
├── commands/
│   ├── explore.ts
│   └── explore.test.ts                            [NEW]
├── providers/
│   └── hover-provider.test.ts                     [NEW]
└── test/
    ├── e2e/
    │   ├── extension.e2e.test.ts                  [NEW]
    │   └── mcp-bridge.e2e.test.ts                 [NEW]
    └── integration/
        ├── sidebar-sync.test.ts                   [NEW]
        └── status-bar.test.ts                     [NEW]
```

**Priority Tests (Week 1):**

1. Extension activation and command registration
2. MCP bridge communication (mock MCP server)
3. Configuration loading and updates
4. Event emitter patterns

**Target:** 12 test files, 10+ hours effort

**Update package.json:**

```json
{
  "scripts": {
    "test": "vitest run",
    "test:watch": "vitest",
    "test:ui": "vitest --ui",
    "test:coverage": "vitest --coverage",
    "test:e2e": "node ./test/e2e/runTests.js"
  }
}
```

---

### Phase 2: Bundle CDN Dependencies Locally (Week 1-2)

**Security & Performance Impact:** Eliminates supply chain risk, reduces external dependencies.

#### 2.1 Dashboard D3.js & Cytoscape Bundling

**Current Issue:**

```html
<!-- apps/dashboard/src/index.html -->
<script src="https://cdnjs.cloudflare.com/ajax/libs/d3/7.8.5/d3.min.js"></script>
<script src="https://cdnjs.cloudflare.com/ajax/libs/cytoscape/3.28.1/cytoscape.min.js"></script>
```

**Solution: Add to npm dependencies**

```bash
npm install d3 cytoscape
```

**Update Angular configuration:**

```json
// angular.json
{
  "projects": {
    "dashboard": {
      "architect": {
        "build": {
          "options": {
            "scripts": [
              "node_modules/d3/dist/d3.min.js",
              "node_modules/cytoscape/dist/cytoscape.min.js"
            ]
          }
        }
      }
    }
  }
}
```

**Update components:**

```typescript
// neural-map.component.ts
import * as d3 from 'd3';
import cytoscape from 'cytoscape';
```

**Build optimization:**

```json
// angular.json optimization section
{
  "optimization": true,
  "sourceMap": false,
  "namedChunks": false,
  "aot": true,
  "extractCss": true,
  "buildOptimizer": true
}
```

**Expected Impact:**
- Eliminate 2 external requests
- Add ~500KB to bundle (gzipped: ~150KB)
- Zero supply chain risk
- Bundle version control

#### 2.2 VS Code Extension Webview Bundling

**Current Issue:**

```typescript
// apps/vscode/src/webviews/neural-map.ts
panel.webview.html = `
  <script src="https://cdnjs.cloudflare.com/ajax/libs/d3/7.8.5/d3.min.js"></script>
  <script src="https://cdnjs.cloudflare.com/ajax/libs/cytoscape/3.28.1/cytoscape.min.js"></script>
`;
```

**Solution: Use VS Code Resource API**

```bash
npm install -D esbuild
```

**Build script:**

```javascript
// scripts/bundle-webview-deps.js
const esbuild = require('esbuild');

esbuild.buildSync({
  entryPoints: ['node_modules/d3/dist/d3.js', 'node_modules/cytoscape/dist/cytoscape.js'],
  bundle: true,
  minify: true,
  outdir: 'dist/webview-deps',
  format: 'iife',
  globalName: ['d3', 'cytoscape']
});
```

**Update webview loading:**

```typescript
const d3Path = panel.webview.asWebviewUri(vscode.Uri.joinPath(extensionUri, 'dist/webview-deps/d3.js'));
const cytoscapePath = panel.webview.asWebviewUri(vscode.Uri.joinPath(extensionUri, 'dist/webview-deps/cytoscape.js'));

panel.webview.html = `
  <script src="${d3Path}"></script>
  <script src="${cytoscapePath}"></script>
`;
```

**Expected Impact:**
- Zero external dependencies
- Better security (CSP compliance)
- Offline functionality
- Deterministic versioning

---

### Phase 3: Structured Logging Framework (Week 2)

**Goal:** Replace 20+ console.log statements with production-grade logging.

#### 3.1 Dashboard Logger Service

**Create logger service:**

```typescript
// apps/dashboard/src/app/core/services/logger.service.ts
import { Injectable } from '@angular/core';

export interface LogEntry {
  timestamp: Date;
  level: 'debug' | 'info' | 'warn' | 'error';
  message: string;
  data?: any;
  stack?: string;
}

@Injectable({ providedIn: 'root' })
export class LoggerService {
  private logs: LogEntry[] = [];
  private maxLogs = 1000;

  debug(message: string, data?: any) {
    this.log('debug', message, data);
  }

  info(message: string, data?: any) {
    this.log('info', message, data);
  }

  warn(message: string, data?: any) {
    this.log('warn', message, data);
  }

  error(message: string, data?: any) {
    this.log('error', message, data);
  }

  private log(level: string, message: string, data?: any) {
    const entry: LogEntry = {
      timestamp: new Date(),
      level: level as any,
      message,
      data
    };

    this.logs.push(entry);
    if (this.logs.length > this.maxLogs) {
      this.logs.shift();
    }

    // Production: Only log errors, warnings
    if (level === 'error' || level === 'warn') {
      console[level](`[${level.toUpperCase()}] ${message}`, data);
    }
  }

  getLogs(): LogEntry[] {
    return [...this.logs];
  }
}
```

**Audit file replacements:**

```bash
# Find all console.log in dashboard
grep -r "console\." apps/dashboard/src --include="*.ts"

# Results to replace:
# apps/dashboard/src/app/core/services/websocket.service.ts (6 instances)
# apps/dashboard/src/app/features/overview/neural-map/neural-map.service.ts (2 instances)
# apps/dashboard/src/app/features/analysis/analysis.service.ts (3 instances)
# ... (8-12 more instances)
```

**Update websocket.service.ts:**

```typescript
// BEFORE
console.log('WebSocket connected');

// AFTER
this.logger.info('WebSocket connected');
```

**Target:** Remove all debug console.log, keep error/warning console calls.

#### 3.2 VS Code Extension Logger

**Create logging utility:**

```typescript
// apps/vscode/src/services/logger.ts
import * as vscode from 'vscode';

const outputChannel = vscode.window.createOutputChannel('Mind Palace');

export const logger = {
  debug: (message: string, data?: any) => {
    if (process.env.DEBUG) {
      outputChannel.appendLine(`[DEBUG] ${message}${data ? ' ' + JSON.stringify(data) : ''}`);
    }
  },
  
  info: (message: string, data?: any) => {
    outputChannel.appendLine(`[INFO] ${message}${data ? ' ' + JSON.stringify(data) : ''}`);
  },
  
  warn: (message: string, error?: Error) => {
    outputChannel.appendLine(`[WARN] ${message}${error ? '\n' + error.stack : ''}`);
  },
  
  error: (message: string, error?: Error) => {
    outputChannel.appendLine(`[ERROR] ${message}${error ? '\n' + error.stack : ''}`);
    vscode.window.showErrorMessage(message);
  }
};
```

**Update files:**

```bash
grep -r "console\." apps/vscode/src --include="*.ts"

# Replace in:
# apps/vscode/src/bridge.ts (1 instance)
# apps/vscode/src/sidebar.ts (2-3 instances)
# apps/vscode/src/commands/*.ts (5+ instances)
```

---

### Phase 4: Advanced Features & Stabilization (Week 3-4)

#### 4.1 Performance Benchmarks

**Create benchmark suite:**

```typescript
// apps/cli/tests/benchmarks_test.go
func BenchmarkIndexing(b *testing.B) {
  // Test workspace with 1000 files
  b.Run("FileIndexing", func(b *testing.B) {
    // Benchmark file parsing
  })
  
  b.Run("SemanticSearch", func(b *testing.B) {
    // Benchmark FTS queries
  })
  
  b.Run("NeuralMapGeneration", func(b *testing.B) {
    // Benchmark graph building
  })
}
```

**Document scaling:**

```markdown
# Mind Palace Scaling Characteristics

## Tested Configurations

| Metric | Max Tested | Recommended Limit |
|--------|-----------|------------------|
| Workspace Size | 50GB | 100GB |
| File Count | 50,000 | 100,000 |
| Symbols | 500,000 | 1M |
| Knowledge Items | 10,000 | 50,000 |

## Performance Profile

| Operation | Time (1000 files) | Time (10,000 files) |
|-----------|------------------|-------------------|
| Initial Scan | 2.5s | 25s |
| Incremental Update | 100ms | 500ms |
| Semantic Search | 50ms | 200ms |
| Neural Map (100 nodes) | 150ms | 200ms |
```

#### 4.2 LLM Integration Tests

**Current gap:** Ollama and OpenAI clients untested

```go
// apps/cli/internal/llm/ollama_test.go
func TestOllamaEmbedding(t *testing.T) {
  server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
    json.NewEncoder(w).Encode(map[string]interface{}{
      "embedding": []float32{0.1, 0.2, 0.3},
    })
  }))
  defer server.Close()

  client := NewOllamaClient(server.URL)
  embedding, err := client.Embed(context.Background(), "test string")
  
  require.NoError(t, err)
  require.NotNil(t, embedding)
}
```

#### 4.3 Cache Management Layer (VS Code)

**Current issue:** Unbounded cache growth

```typescript
// apps/vscode/src/services/cache.ts
export class LRUCache<K, V> {
  private maxSize: number;
  private cache: Map<K, V>;
  private order: K[] = [];

  constructor(maxSize: number = 100) {
    this.maxSize = maxSize;
    this.cache = new Map();
  }

  set(key: K, value: V) {
    if (this.cache.has(key)) {
      this.order = this.order.filter(k => k !== key);
    } else if (this.cache.size >= this.maxSize) {
      const oldest = this.order.shift();
      if (oldest) this.cache.delete(oldest);
    }
    
    this.cache.set(key, value);
    this.order.push(key);
  }

  get(key: K): V | undefined {
    return this.cache.get(key);
  }

  clear() {
    this.cache.clear();
    this.order = [];
  }
}
```

#### 4.4 Postmortem Feature for VS Code

**Currently missing from extension.**

```typescript
// apps/vscode/src/commands/postmortem.ts
export async function addPostmortem() {
  const title = await vscode.window.showInputBox({
    prompt: 'Postmortem title'
  });
  
  const notes = await vscode.window.showInputBox({
    prompt: 'What did you learn?'
  });

  // Call MCP tool: toolAddPostmortem
  const mcp = PalaceBridge.getInstance();
  await mcp.call('toolAddPostmortem', {
    title,
    notes,
    timestamp: new Date().toISOString(),
    context: await getCurrentContext()
  });
}
```

#### 4.5 Butler Refactoring (Reduce God Object)

**Current state:** 850+ lines in extension.ts

**Target structure:**

```
apps/vscode/src/
├── extension.ts (150 lines - activation/registration)
├── command-registry.ts (200 lines - command handling)
├── sidebar-manager.ts (200 lines - sidebar views)
├── status-bar-manager.ts (100 lines - status bar)
├── mcp-bridge.ts (300 lines - MCP communication)
└── event-bus.ts (100 lines - event coordination)
```

#### 4.6 Interactive Onboarding

**Dashboard welcome flow:**

```typescript
// apps/dashboard/src/app/features/onboarding/onboarding.component.ts
@Component({
  selector: 'app-onboarding',
  template: `
    <div class="onboarding-container">
      <div *ngIf="step === 'welcome'" class="welcome-screen">
        <h1>Welcome to Mind Palace</h1>
        <p>Your AI-powered second brain for code</p>
        <button (click)="nextStep()">Get Started</button>
      </div>
      
      <div *ngIf="step === 'init'" class="init-screen">
        <h2>Initialize Your Palace</h2>
        <input [(ngModel)]="projectName" placeholder="Project name">
        <button (click)="initProject()">Create</button>
      </div>
      
      <div *ngIf="step === 'sample'" class="sample-screen">
        <h2>Add Sample Knowledge</h2>
        <p>Let's capture your first room...</p>
        <button (click)="createSampleRoom()">Create Sample</button>
      </div>
    </div>
  `
})
export class OnboardingComponent {
  step = 'welcome';
  projectName = '';
  
  nextStep() { this.step = 'init'; }
  
  async initProject() {
    await this.api.post('/api/init', { name: this.projectName });
    this.step = 'sample';
  }
}
```

---

## Implementation Priority Matrix

| Task | Impact | Effort | Blockers | Priority |
|------|--------|--------|----------|----------|
| Vitest setup (Dashboard) | 🔴 Critical | 1d | None | P0 |
| Extension tests setup | 🔴 Critical | 1d | None | P0 |
| Bundle D3/Cytoscape | 🟡 High | 2d | Tests | P1 |
| Logger service | 🟡 High | 2d | None | P1 |
| Benchmarks | 🟢 Medium | 1d | None | P2 |
| LLM tests | 🟢 Medium | 1d | None | P2 |
| Postmortem VS Code | 🟢 Medium | 2d | Logger | P2 |
| Cache management | 🟢 Medium | 1d | None | P2 |
| Butler refactoring | 🟢 Medium | 2-3d | Tests | P3 |
| Onboarding flow | 🟢 Medium | 2d | Tests | P3 |

---

## Success Metrics (Sprint 2)

| Metric | Current | Target | Success Criteria |
|--------|---------|--------|------------------|
| Dashboard Test Coverage | 0% | 70%+ | ✅ |
| VS Code Test Coverage | 0% | 70%+ | ✅ |
| External Dependencies | 2 CDN | 0 | ✅ |
| Console.log Statements | 20+ | 0 | ✅ |
| Scaling Documentation | Missing | Complete | ✅ |
| LLM Coverage | 0% | 90%+ | ✅ |
| Version | 0.0.2-alpha | 0.1.0-beta | Ready |

---

## Risk Assessment & Mitigation

| Risk | Probability | Impact | Mitigation |
|------|-------------|--------|-----------|
| Tests break on upgrade | Medium | High | Use LTS versions, extensive testing |
| Bundling increases size | Low | Medium | Monitor bundle metrics, optimize |
| Refactoring bugs | Low | High | Maintain test coverage >80% |
| Performance regression | Low | Medium | Add performance benchmarks |

---

## Deliverables Checklist

- [ ] Vitest configuration (dashboard)
- [ ] @vscode/test-electron setup (extension)
- [ ] 25+ Dashboard test files
- [ ] 12+ Extension test files
- [ ] D3.js/Cytoscape bundled (dashboard)
- [ ] WebView bundle setup (extension)
- [ ] Logger service deployed
- [ ] All console.log audit completed
- [ ] Performance benchmarks documented
- [ ] LLM integration tests added
- [ ] Cache management layer implemented
- [ ] Postmortem feature in VS Code
- [ ] Butler refactored into modules
- [ ] Onboarding tutorial created
- [ ] Version bumped to 0.1.0-beta
- [ ] All tests passing (CI green)

---

## Timeline & Capacity

**Team:** 1 Full-time engineer + AI assistance

**Weeks 1-2:** Phase 1 & 2 (Testing + Bundling)
- 40 hours (tests) + 16 hours (bundling) = 56 hours

**Weeks 3-4:** Phase 3 & 4 (Logging + Advanced Features)
- 24 hours (logging) + 40 hours (features) = 64 hours

**Total:** 120 hours = 3-4 weeks at 1 FTE

---

## Next Steps

1. **This week:** Review this plan with team
2. **Monday:** Begin Phase 1 (Vitest setup)
3. **Daily:** 15-min stand-up on progress
4. **Weekly:** Demo of working features
5. **Sprint end:** Beta release readiness assessment

---

**Prepared by:** AI Engineering Team  
**Date:** January 5, 2026  
**Status:** 🟢 Ready for Execution
