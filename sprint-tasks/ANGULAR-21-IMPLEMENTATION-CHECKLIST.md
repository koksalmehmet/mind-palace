# Angular 21 Implementation Checklist

**Status Tracking:** Mark items as `[x]` when completed

---

## Phase 1: Foundation Setup (Week 1) ⚡

### Testing Infrastructure
- [ ] Install Vitest: `npm install -D vitest @vitest/ui @vitest/coverage-v8 happy-dom`
- [ ] Install Angular testing utilities: `npm install -D @angular/core/testing`
- [ ] Create `vitest.config.ts` (see ANGULAR-21-CODE-TEMPLATES.md)
- [ ] Create `src/test-setup.ts` with zone.js mocks
- [ ] Add npm scripts:
  ```json
  "test": "vitest",
  "test:ui": "vitest --ui",
  "test:coverage": "vitest run --coverage"
  ```
- [ ] Run first test: `npm run test`
- [ ] Verify coverage reporting works

### TypeScript Strictness Enhancement
- [ ] Update `tsconfig.json` with new strict options:
  ```json
  "useUnknownInCatchVariables": true,
  "noUncheckedIndexedAccess": true,
  "noImplicitThis": true,
  "exactOptionalPropertyTypes": true
  ```
- [ ] Run `ng build` and fix any new strict mode errors
- [ ] Document all `as unknown` type assertions (should be minimal)

### Documentation Review
- [ ] Read [ANGULAR-21-BEST-PRACTICES.md](./ANGULAR-21-BEST-PRACTICES.md)
- [ ] Review [ANGULAR-21-CODE-TEMPLATES.md](./ANGULAR-21-CODE-TEMPLATES.md)
- [ ] Create team meeting to discuss patterns
- [ ] Share Signal advantages with team

---

## Phase 2: State Management Refactoring (Week 2-3) 📦

### Room State Service (High Priority)
- [ ] Create `src/app/core/services/room-state.service.ts`
  - Copy from Template 1 in CODE-TEMPLATES.md
  - Adapt Room interface to match your actual model
  - Test with `src/app/core/services/room-state.service.spec.ts`
- [ ] Create comprehensive spec file (use Template 4)
- [ ] Run tests: `npm run test -- room-state.service.spec`
- [ ] Achieve 100% code coverage
- [ ] Document in-code with JSDoc comments

### Component Refactoring
- [ ] Update `room-list.component.ts` to use signals
  - Replace `@Input() rooms` with `input.required<Room[]>()`
  - Replace `@Output() select` with `output<string>()`
  - Add `changeDetection: ChangeDetectionStrategy.OnPush`
  - Use computed() for derived state (see Template 2)
- [ ] Create spec: `room-list.component.spec.ts` (use Template 3)
- [ ] Update dashboard component similarly
- [ ] Verify all components have `changeDetection: OnPush`
- [ ] Run full test suite: `npm run test`

### Service Integration
- [ ] Inject RoomStateService in components
- [ ] Replace subscriptions with signal reads
- [ ] Remove manual change detection calls (not needed with OnPush)
- [ ] Test memory leaks with Angular DevTools

### Progress Check
- [ ] All room-related components use signal inputs
- [ ] State mutations only happen in services
- [ ] No Observable subscriptions in templates
- [ ] Tests pass with > 80% coverage

---

## Phase 3: Build Optimization (Week 3) 🚀

### Angular CLI Configuration
- [ ] Update `angular.json`:
  ```json
  "optimization": true,
  "buildOptimizer": true,
  "namedChunks": false,
  "outputHashing": "all",
  "vendorChunk": false,
  "commonChunk": true,
  "sourceMap": false
  ```
- [ ] Set bundle budgets:
  ```json
  "budgets": [
    {
      "type": "initial",
      "maximumWarning": "250kb",
      "maximumError": "350kb"
    }
  ]
  ```
- [ ] Build and verify: `ng build --configuration production`
- [ ] Check output: should be < 250KB gzipped

### Route Lazy Loading
- [ ] Update `app.routes.ts` with `loadComponent` (see Template 6)
  - Split features into separate files
  - Implement for: rooms/:id, analysis, connections, settings
- [ ] Verify bundle chunks created:
  ```bash
  ng build --configuration production
  # Look for: room-detail.chunk.js, analysis.chunk.js, etc.
  ```
- [ ] Test route navigation is smooth
- [ ] Measure chunk sizes:
  ```bash
  npm run build -- --stats-json
  webpack-bundle-analyzer dist/*/stats.json
  ```

### D3 Library Bundling
- [ ] Verify D3 is bundled locally (not CDN)
- [ ] Create `d3-visualization.service.ts` (see Template 5)
- [ ] Type all D3 interactions strictly
- [ ] Test D3 rendering with component tests
- [ ] Benchmark: D3 rendering should take < 100ms

### Production Build Verification
- [ ] Build command: `npm run build`
- [ ] Initial bundle: target < 250KB gzipped
- [ ] Lazy chunks: each < 100KB gzipped
- [ ] No console warnings about large chunks
- [ ] Source maps disabled in production build

---

## Phase 4: Performance Validation (Week 4) 📊

### Local Performance Testing
- [ ] Install Lighthouse: `npm install -D @lhci/cli@0.11.x`
- [ ] Create `.lhcirc.json`:
  ```json
  {
    "ci": {
      "collect": {
        "url": ["http://localhost:4200/dashboard"],
        "numberOfRuns": 3,
        "settings": {
          "precomputedLanternDataUrl": null
        }
      },
      "assert": {
        "preset": "lighthouse:recommended",
        "assertions": {
          "categories:performance": ["error", { "minScore": 0.9 }]
        }
      }
    }
  }
  ```
- [ ] Run Lighthouse audit:
  ```bash
  npm run build
  # In another terminal:
  npm run start
  # Then:
  lhci autorun --config=.lhcirc.json
  ```
- [ ] Document baseline metrics:
  - FCP: ___ ms
  - LCP: ___ ms (target: < 2500ms)
  - TTI: ___ ms (target: < 3500ms)
  - CLS: ___ (target: < 0.1)
  - Perf Score: ___ (target: > 90)

### Memory & Change Detection Profiling
- [ ] Install Angular DevTools extension
- [ ] Profile change detection:
  - Open DevTools > Profiler
  - Record while using dashboard
  - Check: no > 16ms frames (60fps)
- [ ] Memory leak detection:
  - DevTools Memory tab
  - Take heap snapshot
  - Navigate heavily
  - Take another snapshot
  - Compare: should not grow unbounded

### Metrics Before/After
Create file: `PERFORMANCE-METRICS.md` with:
```markdown
## Performance Improvements

### Bundle Size
- Before: XXX KB (initial)
- After: XXX KB (initial)
- Savings: XXX KB (XX%)

### Lighthouse Scores
| Metric | Before | After | Target |
|--------|--------|-------|--------|
| Performance | XX | XX | >90 |
| LCP | XXms | XXms | <2.5s |
| TTI | XXms | XXms | <3.5s |
| CLS | X.XX | X.XX | <0.1 |

### Code Coverage
- Services: XX%
- Components: XX%
- Total: XX%
```

---

## Phase 5: Type Safety & Compliance (Week 4-5) 🔒

### TypeScript Strict Mode Audit
- [ ] Run full build with strict options enabled
- [ ] List all `any` usages:
  ```bash
  grep -r "any" src/app --include="*.ts" | grep -v node_modules
  ```
- [ ] Replace `any` with `unknown` or proper types
- [ ] Target: 0 `any` usages (except legacy code with `// @ts-ignore`)

### Error Handling Refactoring
- [ ] Audit all try-catch blocks
- [ ] Use type-safe error handling (see Template 7):
  ```typescript
  catch (error: unknown) {
    if (error instanceof Error) {
      // Handle properly
    }
  }
  ```
- [ ] Create typed error responses for API calls
- [ ] Handle all error cases in services

### Type Guard Implementation
- [ ] Create discriminated union types for state:
  - RoomState: `{ status: 'loading' } | { status: 'success'; data: Room[] }`
  - AnalyticsState: similar pattern
- [ ] Use switch statements with exhaustive checking
- [ ] Compile with `noImplicitReturns` and `noFallthroughCasesInSwitch`

### Strict Compilation
- [ ] `ng build` should complete with 0 errors
- [ ] `npm run test` should pass all tests
- [ ] `npx tsc --noEmit` should have 0 errors

---

## Phase 6: Testing Coverage (Week 5) ✅

### Service Test Coverage
- [ ] RoomStateService: > 95% coverage
  - [ ] All mutations tested
  - [ ] All computed signals tested
  - [ ] Error scenarios tested
  - [ ] Filtering logic tested
- [ ] D3VisualizationService: > 85% coverage
  - [ ] Render methods tested
  - [ ] Error handling tested
  - [ ] Type conversions tested

### Component Test Coverage
- [ ] RoomListComponent: > 85% coverage
  - [ ] Rendering with data tested
  - [ ] Empty state tested
  - [ ] Loading state tested
  - [ ] Search filtering tested
  - [ ] Output events tested
  - [ ] Error handling tested
- [ ] DashboardComponent: > 80% coverage
- [ ] All feature components: > 80% coverage

### Test Quality Metrics
- [ ] Run: `npm run test:coverage`
- [ ] Generate HTML report: `open coverage/index.html`
- [ ] Document coverage baselines in PERFORMANCE-METRICS.md
- [ ] Set up pre-commit hook to prevent coverage regression

### Snapshot Testing
- [ ] For large components, consider snapshot tests
- [ ] Create snapshots for rendered DOM
- [ ] Document snapshot purpose with comments
- [ ] Review snapshots in git (don't blindly update)

---

## Phase 7: Documentation & Knowledge Share (Week 5) 📚

### Create Developer Guide
- [ ] Document signal patterns used in Mind Palace
- [ ] Create style guide for new components
- [ ] Document testing patterns with Vitest
- [ ] Add section on performance best practices

### Code Examples
- [ ] Maintain [ANGULAR-21-CODE-TEMPLATES.md](./ANGULAR-21-CODE-TEMPLATES.md)
- [ ] Add real examples from implemented code
- [ ] Create 2-3 tutorials for new developers

### Team Training
- [ ] Hold: "Signals & Standalone Components" workshop
- [ ] Hold: "Testing with Vitest" workshop
- [ ] Hold: "Performance Optimization" session
- [ ] Record sessions for async team members

### CI/CD Updates
- [ ] Update GitHub Actions to run tests
- [ ] Set up bundle size tracking
- [ ] Add Lighthouse CI (optional)
- [ ] Block PRs if tests fail or coverage drops

---

## Final Validation Checklist

Before marking complete, verify:

### Functionality ✅
- [ ] All dashboard features work unchanged
- [ ] No console errors or warnings
- [ ] All API calls successful
- [ ] D3 visualizations render correctly
- [ ] Room navigation smooth
- [ ] No memory leaks detected

### Performance ✅
- [ ] Initial bundle < 250KB gzipped
- [ ] LCP < 2.5s
- [ ] TTI < 3.5s
- [ ] Lighthouse Performance > 90
- [ ] Change detection < 16ms frame time

### Code Quality ✅
- [ ] 0 `any` type usages
- [ ] TypeScript strict mode: 0 errors
- [ ] Test coverage > 85% (services) & > 80% (components)
- [ ] All tests passing
- [ ] ESLint: 0 errors

### Documentation ✅
- [ ] ANGULAR-21-BEST-PRACTICES.md reviewed
- [ ] ANGULAR-21-CODE-TEMPLATES.md up-to-date
- [ ] PERFORMANCE-METRICS.md documented
- [ ] Team trained and equipped
- [ ] New component templates established

### Release Readiness ✅
- [ ] Version bumped appropriately
- [ ] CHANGELOG.md updated
- [ ] Tested in production-like environment
- [ ] Rollback plan documented
- [ ] Monitoring/alerts set up

---

## Quick Reference

**Current Status:** Phase 1 - Setup  
**Start Date:** January 2025  
**Target Completion:** End of January 2025  
**Team Lead:** _____  
**Last Updated:** 2025-01-05

### Key Contacts
- Angular/TypeScript Questions: _____
- Performance Issues: _____
- Testing Help: _____
- d3 Integration: _____

### Resources
- [ANGULAR-21-BEST-PRACTICES.md](./ANGULAR-21-BEST-PRACTICES.md) - Complete guide
- [ANGULAR-21-CODE-TEMPLATES.md](./ANGULAR-21-CODE-TEMPLATES.md) - Ready-to-use code
- Angular 21 Docs: https://angular.io
- Vitest Guide: https://vitest.dev/guide/
- TypeScript Handbook: https://www.typescriptlang.org/docs/handbook/

### Common Issues & Solutions

| Issue | Solution |
|-------|----------|
| Vitest not finding modules | Check moduleResolution in tsconfig |
| Signal change not triggering | Ensure component has OnPush strategy |
| Tests timing out | Increase timeout in vitest.config.ts |
| Coverage not running | Run `npm run test:coverage` separately |
| d3 types missing | Install `@types/d3` and check tsconfig |
| Build exceeding budget | Check for large unused dependencies |

---

**Questions?** Reference the implementation guide or ask the team lead listed above.
