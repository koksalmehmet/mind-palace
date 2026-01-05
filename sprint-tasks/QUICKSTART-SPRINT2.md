# 🚀 Quick Start: Sprint 2 Execution Guide

**For:** Development Team  
**Date:** January 5, 2026  
**Timeline:** Start Monday, 3-4 weeks, 120 hours effort  

---

## TL;DR - What You Need to Know

✅ **Sprint 1 is complete** - All 6 critical fixes deployed  
✅ **Three comprehensive plans created** - Ready to execute  
✅ **Research completed** - Latest best practices from official sources  
✅ **You're here:** Picking up to continue with Sprint 2  

**Your job:** Follow the phase-by-phase plan in SPRINT-2-PLAN.md and pick tests to add first.

---

## 📚 Required Reading (In Order)

### For Everyone (30 min)
1. This document (5 min)
2. PROJECT-STATUS.md Executive Summary (10 min)  
3. SPRINT-2-PLAN.md Overview section (15 min)

### For Developers (2 hours)
1. SPRINT-2-PLAN.md Phases 1-4 (1 hour)
2. TECHNOLOGY-RESEARCH-2025.md Testing sections (45 min)
3. Start setting up Vitest (15 min)

### For Architecture Review
1. TECHNOLOGY-RESEARCH-2025.md all sections (1 hour)
2. SPRINT-2-PLAN.md Phases 3-4 (30 min)
3. Make approval decisions

---

## 🎯 Phase 1: This Week (Vitest + Extension Tests)

### Monday Morning: Setup

```bash
# 1. Update to latest branch
cd /Users/mehmetkoksal/Documents/Projects/Personal/mind-palace
git checkout plc-001
git pull origin plc-001

# 2. Create Sprint 2 branch
git checkout -b plc-002

# 3. Install testing dependencies
cd apps/dashboard
npm install -D vitest @vitest/ui @vitest/browser @testing-library/angular

cd ../vscode
npm install -D @vscode/test-electron mocha chai sinon

cd ../..
```

### Dashboard Testing (Tuesday-Wednesday)

**Setup Vitest:**

```bash
cd apps/dashboard

# Create config
cat > vitest.config.ts << 'EOF'
import { defineConfig } from 'vitest/config';
import angular from 'vite-plugin-angular';

export default defineConfig({
  plugins: [angular()],
  test: {
    globals: true,
    environment: 'jsdom',
    coverage: {
      provider: 'v8',
      reporter: ['text', 'json', 'html']
    }
  }
});
EOF

# Update package.json scripts
npm set-script test "vitest run"
npm set-script test:watch "vitest"
npm set-script test:coverage "vitest --coverage"
```

**Write first 3 critical tests:**

1. `src/app/core/services/websocket.service.spec.ts` (2 hours)
   - Connection/reconnection logic
   - Message handling
   - Error recovery

2. `src/app/core/services/api.service.spec.ts` (2 hours)
   - HTTP request/response
   - Error handling
   - Interceptors

3. `src/app/features/overview/neural-map/neural-map.component.spec.ts` (3 hours)
   - Component initialization
   - D3 rendering
   - Data updates

**Target:** 3 test files, 50+ test cases, ~6% coverage by Wednesday EOD

### VS Code Extension Testing (Wednesday-Friday)

**Setup @vscode/test-electron:**

```bash
cd apps/vscode

# Create test structure
mkdir -p src/test/{e2e,integration}

# Create test runner
cat > src/test/runTests.ts << 'EOF'
import * as path from 'path';
import { runTests } from '@vscode/test-electron';

async function main() {
  try {
    const extensionDevelopmentPath = path.resolve(__dirname, '../../');
    const extensionTestsPath = path.resolve(__dirname, 'suite');
    
    await runTests({
      extensionDevelopmentPath,
      extensionTestsPath,
      launchArgs: ['--disable-extensions']
    });
  } catch (err) {
    console.error('Failed to run tests', err);
    process.exit(1);
  }
}

main();
EOF
```

**Write first 3 critical tests:**

1. `src/extension.test.ts` (1.5 hours)
   - Extension activation
   - Command registration
   - Initial state setup

2. `src/commands/explore.test.ts` (1.5 hours)
   - Command execution
   - Error handling
   - Output validation

3. `src/mcp-bridge.test.ts` (2 hours)
   - MCP connection
   - Tool invocation (mocked)
   - Response handling

**Target:** 3 test files, 30+ test cases, ~5% coverage by Friday EOD

### Week 1 Deliverables
- ✅ Vitest configured and running
- ✅ @vscode/test-electron working
- ✅ 6 test files written
- ✅ CI pipeline updated with test runs
- ✅ Test coverage reports generated

---

## 🔒 Phase 2: Next Week (Bundling)

### Monday: D3.js & Cytoscape Migration

```bash
# Add to npm
cd apps/dashboard
npm install d3 cytoscape

# Update angular.json
# See SPRINT-2-PLAN.md Phase 2 for exact JSON changes

# Update components
# Replace CDN loads with npm imports
# Test builds locally
npm run build

# Check bundle size
ls -lh dist/dashboard
```

### Tuesday: VS Code WebView Bundling

```bash
# Install esbuild
npm install -D esbuild

# Create bundle script (see SPRINT-2-PLAN.md)
cat > scripts/bundle-webview.mjs

# Run bundling
node scripts/bundle-webview.mjs

# Update webview loading code
# See SPRINT-2-PLAN.md Phase 2 for implementation
```

### Week 2 Deliverables
- ✅ D3.js locally bundled
- ✅ Cytoscape locally bundled
- ✅ VS Code CDN dependencies eliminated
- ✅ Security verified (no external scripts)
- ✅ Build times documented

---

## 📝 Phase 3: Week 3 (Structured Logging)

### Create Logger Service

```typescript
// Dashboard
src/app/core/services/logger.service.ts (150 lines)

// VS Code
src/services/logger.ts (100 lines)
```

### Replace console.log

```bash
# Find all console statements
grep -r "console\." apps/dashboard/src --include="*.ts" | wc -l
grep -r "console\." apps/vscode/src --include="*.ts" | wc -l

# Replace each with logger.* calls
# Use find/replace in VS Code
```

### Week 3 Deliverables
- ✅ Logger service implemented
- ✅ 20+ console.log calls removed
- ✅ Production logging configured
- ✅ Test coverage maintained

---

## ⚡ Phase 4: Week 4 (Advanced Features)

### Pick 2-3 from menu:
1. **Benchmarks** (2 hours) - Go benchmark suite
2. **LLM Tests** (2 hours) - Ollama/OpenAI client tests
3. **Postmortem** (4 hours) - VS Code feature
4. **Cache Layer** (2 hours) - LRU cache impl
5. **Onboarding** (4 hours) - Welcome flow

**Recommended:** Start with Benchmarks + LLM Tests (quick wins)

### Week 4 Deliverables
- ✅ 2-3 advanced features implemented
- ✅ Performance characteristics documented
- ✅ Test coverage >70% on all components

---

## ✅ Success Metrics (End of Sprint 2)

| Metric | Target | Acceptance |
|--------|--------|-----------|
| Dashboard coverage | 70%+ | ≥70 ✅ |
| VS Code coverage | 70%+ | ≥70 ✅ |
| CLI coverage | 95%+ | ≥95 ✅ |
| CDN dependencies | 0 | 0 external CDN ✅ |
| console.log in prod | 0 | No debug logs ✅ |
| Test execution time | <30s | CI passes quickly ✅ |
| Version | 0.1.0-beta | Released ✅ |

---

## 🚨 Critical Milestones

- **Monday EOD:** Vitest + @vscode/test-electron set up, first tests passing
- **Wednesday EOD:** 6 test files written, coverage reports showing
- **Friday EOD:** Phase 1 complete, PR submitted for review
- **Week 2 EOD:** Phase 2 complete, all CDN deps bundled locally
- **Week 3 EOD:** Phase 3 complete, logging service deployed
- **Week 4 EOD:** Phase 4 complete, 0.1.0-beta ready

---

## 📞 If You Get Stuck

### For Code Examples
→ See SPRINT-2-PLAN.md (has full code snippets)

### For Best Practices
→ See TECHNOLOGY-RESEARCH-2025.md (official sources)

### For Architecture Questions
→ See ANALYSIS.md component analysis

### For Sprint 1 Verification
→ See IMPLEMENTATION-LOG.md + TEST-RESULTS.md

---

## 🎬 First Task (DO THIS NOW)

1. Open SPRINT-2-PLAN.md
2. Read "Phase 1: Frontend Test Infrastructure"
3. Follow the Vitest setup section exactly
4. Run: `npm install -D vitest @vitest/ui @vitest/browser @testing-library/angular`
5. Create `vitest.config.ts` as shown
6. Run: `npm test` (should show 0 tests found - that's OK)
7. Slack: "✅ Vitest set up and ready, writing first test now"

---

## 💪 You Got This!

Sprint 1 laid the foundation. Sprint 2 is about building the quality and stability needed for beta release. Follow the plan, write the tests, and we'll have a production-ready system in 4 weeks.

**Questions?** See the documents. They have answers.

**Ready?** Go to SPRINT-2-PLAN.md and start Phase 1.

---

**Prepared by:** AI Engineering Team  
**Status:** 🟢 READY TO EXECUTE  
**Start:** Monday, January 6, 2026
