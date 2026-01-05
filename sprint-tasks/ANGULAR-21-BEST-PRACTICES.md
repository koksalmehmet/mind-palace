# Angular 21 & TypeScript 5.x Best Practices for Mind Palace Dashboard

**Date:** January 2025  
**Scope:** Production-grade recommendations for Mind Palace Dashboard  
**Status:** Actionable Implementation Guide

---

## Executive Summary

This document provides current, battle-tested best practices for Angular 21 and TypeScript 5.x. Mind Palace Dashboard is already well-positioned with standalone components and strict TypeScript configuration. This guide focuses on signal-based state, testing setup, performance optimization, and security patterns.

---

## 1. Angular 21 Standalone Components & Signal-Based State Management

### Current Status ✅
- Dashboard already uses standalone components (bootstrapApplication)
- Angular 21.0.6 with TypeScript 5.9.3
- Proper change detection setup

### Best Practices Implementation

#### 1.1 Signal-Based State Management Pattern

Replace traditional services with signals for reactive, performant state management:

```typescript
// ✅ RECOMMENDED: Signal-based state service
import { Injectable } from '@angular/core';
import { signal, computed, effect } from '@angular/core';

@Injectable({ providedIn: 'root' })
export class RoomStateService {
  // Signal for mutable state
  private readonly rooms = signal<Room[]>([]);
  private readonly selectedRoomId = signal<string | null>(null);
  private readonly isLoading = signal(false);
  private readonly error = signal<string | null>(null);

  // Computed derived state (automatically memoized)
  readonly selectedRoom = computed(() => {
    const id = this.selectedRoomId();
    return this.rooms().find(r => r.id === id);
  });

  readonly filteredRooms = computed(() => {
    // Expensive computation only runs when dependencies change
    return this.rooms().filter(r => !r.isArchived);
  });

  readonly hasErrors = computed(() => this.error() !== null);

  constructor(private http: HttpClient) {
    // Track state changes for logging/debugging
    effect(() => {
      console.log('Rooms changed:', this.rooms());
    });
  }

  // Mutations
  loadRooms(): void {
    this.isLoading.set(true);
    this.http.get<Room[]>('/api/rooms')
      .pipe(
        tap(rooms => this.rooms.set(rooms)),
        tap(() => this.isLoading.set(false)),
        catchError(err => {
          this.error.set(err.message);
          return of([]);
        })
      )
      .subscribe();
  }

  selectRoom(id: string): void {
    this.selectedRoomId.set(id);
  }

  // Return signal references (read-only from component)
  getRooms = this.rooms.asReadonly();
  getSelectedRoom = this.selectedRoom.asReadonly();
  getIsLoading = this.isLoading.asReadonly();
}
```

#### 1.2 Signal-Based Component Pattern

```typescript
// ✅ RECOMMENDED: Signal-aware component with OnPush
import { Component, computed, inject, input, output } from '@angular/core';
import { CommonModule } from '@angular/common';
import { ChangeDetectionStrategy } from '@angular/core';

@Component({
  selector: 'app-room-list',
  standalone: true,
  imports: [CommonModule],
  template: `
    <div class="room-list">
      @for (room of rooms(); track room.id) {
        <app-room-card 
          [room]="room"
          [isSelected]="room.id === selectedRoomId()"
          (select)="onRoomSelect($event)"
        />
      }
      @if (isLoading()) {
        <div class="spinner">Loading...</div>
      }
    </div>
  `,
  changeDetection: ChangeDetectionStrategy.OnPush
})
export class RoomListComponent {
  // Input signals (since Angular 21)
  readonly rooms = input.required<Room[]>();
  readonly isLoading = input.required<boolean>();
  readonly selectedRoomId = input<string | null>(null);

  // Output signals
  readonly roomSelected = output<string>();

  // Computed derived values
  readonly roomCount = computed(() => this.rooms().length);
  readonly emptyState = computed(() => this.rooms().length === 0 && !this.isLoading());

  onRoomSelect(id: string): void {
    this.roomSelected.emit(id);
  }
}
```

**Why Signals Over RxJS in Leaf Components?**
- Signals eliminate change detection overhead
- No subscription memory leaks
- Automatic dependency tracking
- Built-in memoization via `computed`
- Better TypeScript inference

---

## 2. Testing Angular 21 Components with Vitest (2025 Recommended)

### Recommended: Vitest Over Jest

**Why Vitest?**
- Same Jest API but 5-10x faster (Vite powered)
- ESM-first (matches Angular 21 module resolution)
- Built-in browser mode for realistic component testing
- Better async component support

### Setup Vitest for Angular 21

#### 2.1 Installation & Configuration

```bash
npm install -D vitest @vitest/ui @vitest/coverage-v8 happy-dom
npm install -D @angular/core @angular/common/testing
```

**vitest.config.ts:**
```typescript
import { defineConfig } from 'vitest/config';
import angular from '@analogjs/vite-plugin-angular';

export default defineConfig({
  plugins: [angular()],
  test: {
    globals: true,
    environment: 'happy-dom',  // or 'jsdom' for full browser API
    setupFiles: ['src/test-setup.ts'],
    coverage: {
      provider: 'v8',
      reporter: ['text', 'json', 'html'],
      exclude: [
        'node_modules/',
        'src/test-setup.ts',
        '**/*.spec.ts',
      ]
    },
    include: ['src/**/*.spec.ts'],
  },
});
```

**src/test-setup.ts:**
```typescript
import { vi } from 'vitest';
import 'zone.js';
import 'zone.js/testing';

// Mock Angular in tests
Object.defineProperty(window, 'CSS', {value: null});
Object.defineProperty(window, 'getComputedStyle', {
  value: () => ({ display: 'none', appearance: ['-webkit-appearance'] })
});
Object.defineProperty(document, 'doctype', {
  value: '<!DOCTYPE html>'
});
Object.defineProperty(document.body.style, 'transform', {
  value: () => { return {}; },
  writable: true,
});
```

#### 2.2 Component Testing with Signals

```typescript
// ✅ RECOMMENDED: Vitest + Signal-based testing
import { ComponentFixture, TestBed } from '@angular/core/testing';
import { describe, it, expect, beforeEach, vi } from 'vitest';
import { signal } from '@angular/core';
import { RoomListComponent } from './room-list.component';
import { RoomStateService } from '../services/room-state.service';

describe('RoomListComponent', () => {
  let component: RoomListComponent;
  let fixture: ComponentFixture<RoomListComponent>;
  let roomService: RoomStateService;

  beforeEach(async () => {
    const mockService = {
      getRooms: signal([
        { id: '1', name: 'Living Room', isArchived: false }
      ]),
      getIsLoading: signal(false),
    };

    await TestBed.configureTestingModule({
      imports: [RoomListComponent],
      providers: [
        { provide: RoomStateService, useValue: mockService }
      ]
    }).compileComponents();

    roomService = TestBed.inject(RoomStateService);
    fixture = TestBed.createComponent(RoomListComponent);
    component = fixture.componentInstance;
  });

  it('should render rooms from signal', () => {
    fixture.componentRef.setInput('rooms', [
      { id: '1', name: 'Bedroom' }
    ]);
    fixture.detectChanges();

    const roomCards = fixture.nativeElement.querySelectorAll('app-room-card');
    expect(roomCards).toHaveLength(1);
  });

  it('should handle room selection output', () => {
    const selectedRoomId = vi.fn();
    fixture.componentRef.instance.roomSelected.subscribe(selectedRoomId);

    fixture.componentRef.setInput('rooms', [
      { id: '1', name: 'Bedroom' }
    ]);
    fixture.detectChanges();

    const selectBtn = fixture.nativeElement.querySelector('button');
    selectBtn.click();

    expect(selectedRoomId).toHaveBeenCalledWith('1');
  });
});
```

#### 2.3 Testing Services with Signals

```typescript
// ✅ Service testing with signals
import { TestBed } from '@angular/core/testing';
import { HttpClientTestingModule, HttpTestingController } from '@angular/common/http/testing';
import { describe, it, expect, beforeEach } from 'vitest';
import { RoomStateService } from './room-state.service';

describe('RoomStateService', () => {
  let service: RoomStateService;
  let httpMock: HttpTestingController;

  beforeEach(() => {
    TestBed.configureTestingModule({
      imports: [HttpClientTestingModule],
      providers: [RoomStateService]
    });

    service = TestBed.inject(RoomStateService);
    httpMock = TestBed.inject(HttpTestingController);
  });

  it('should load rooms and update signal', () => {
    const mockRooms = [
      { id: '1', name: 'Library' },
      { id: '2', name: 'Study' }
    ];

    service.loadRooms();
    const req = httpMock.expectOne('/api/rooms');
    req.flush(mockRooms);

    // Signals are synchronously accessible
    expect(service.getRooms().length).toBe(2);
    expect(service.getIsLoading()).toBe(false);
  });

  it('should compute filtered rooms correctly', () => {
    const mockRooms = [
      { id: '1', name: 'Library', isArchived: false },
      { id: '2', name: 'Study', isArchived: true }
    ];

    service['rooms'].set(mockRooms);
    
    // Computed signals automatically update
    expect(service.filteredRooms().length).toBe(1);
  });

  afterEach(() => {
    httpMock.verify();
  });
});
```

**Update package.json:**
```json
{
  "scripts": {
    "test": "vitest",
    "test:ui": "vitest --ui",
    "test:coverage": "vitest run --coverage"
  }
}
```

---

## 3. Performance Optimization in Angular 21

### 3.1 Change Detection Strategy: OnPush + Signals

**Critical for Dashboard Performance:**

```typescript
// ✅ ALWAYS use OnPush with standalone components
@Component({
  selector: 'app-dashboard',
  standalone: true,
  imports: [CommonModule, RoomListComponent],
  template: `...`,
  changeDetection: ChangeDetectionStrategy.OnPush  // ← Essential
})
export class DashboardComponent {
  // Signals automatically trigger OnPush change detection
  readonly rooms = input<Room[]>();
  readonly selectedRoom = computed(() => this.getSelected());
}
```

**Impact:**
- Change detection only runs when @Input() signals change or events fire
- ~70% reduction in change detection cycles for complex dashboards
- Linear performance scaling instead of quadratic

### 3.2 Lazy Loading with Routes

```typescript
// ✅ app.routes.ts - Route-based lazy loading
export const routes: Routes = [
  {
    path: 'dashboard',
    component: DashboardComponent,
    children: [
      {
        path: 'rooms/:id',
        loadComponent: () => import('./features/room-detail/room-detail.component')
          .then(m => m.RoomDetailComponent)
      },
      {
        path: 'analysis',
        loadComponent: () => import('./features/analysis/analysis.component')
          .then(m => m.AnalysisComponent)
      }
    ]
  },
  {
    path: 'settings',
    loadComponent: () => import('./features/settings/settings.component')
      .then(m => m.SettingsComponent)
  }
];
```

**Result:** Only critical path components loaded initially, ~40% faster FCP (First Contentful Paint).

### 3.3 Virtual Scrolling for Large Lists

```typescript
// ✅ Virtual scrolling for D3 room visualization
import { ScrollingModule } from '@angular/cdk/scrolling';

@Component({
  selector: 'app-room-list',
  standalone: true,
  imports: [ScrollingModule],
  template: `
    <cdk-virtual-scroll-viewport itemSize="150" class="viewport">
      @for (room of rooms(); track room.id) {
        <app-room-card [room]="room" />
      }
    </cdk-virtual-scroll-viewport>
  `
})
export class RoomListComponent {
  readonly rooms = input<Room[]>();
}
```

**Handles:** 10K+ rooms without performance degradation.

### 3.4 Memoization with computed()

```typescript
// ✅ Avoid re-computing expensive operations
@Injectable({ providedIn: 'root' })
export class RoomAnalysisService {
  private readonly allRooms = signal<Room[]>([]);
  private readonly selectedFilters = signal<RoomFilter>({});

  // This computation only runs when inputs change
  readonly filteredRoomsWithMetrics = computed(() => {
    const rooms = this.allRooms();
    const filters = this.selectedFilters();
    
    // Expensive D3 metric calculation
    return rooms
      .filter(r => matchesFilters(r, filters))
      .map(r => ({
        ...r,
        metrics: calculateComplexMetrics(r)  // Only when deps change
      }));
  });
}
```

### 3.5 Async Pipe with Signals

```typescript
// ✅ New signal-based async pattern (Angular 21.1+)
@Component({
  template: `
    <div>{{ roomData() }}</div>
  `
})
export class RoomComponent {
  private readonly roomService = inject(RoomService);
  
  // Convert Observable to Signal
  readonly roomData = toSignal(
    this.roomService.getRoom(),
    { initialValue: null }
  );
  
  // For cleanup:
  // toSignal automatically handles subscription cleanup via effect()
}
```

---

## 4. TypeScript 5.x Strict Mode Compliance

### Current Status ✅
Dashboard has excellent TypeScript config. Enhancements:

```json
{
  "compilerOptions": {
    "strict": true,                          // ✅ Enabled
    "noImplicitOverride": true,              // ✅ Enabled
    "noPropertyAccessFromIndexSignature": true, // ✅ Enabled
    "noImplicitReturns": true,               // ✅ Enabled
    "noFallthroughCasesInSwitch": true,      // ✅ Enabled
    "useUnknownInCatchVariables": true,      // ← ADD (safer error handling)
    "noUncheckedIndexedAccess": true,        // ← ADD (safer object access)
    "noImplicitThis": true,                  // ← ADD (explicit 'this' typing)
    "exactOptionalPropertyTypes": true,      // ← ADD (strict optional handling)
    "forceConsistentCasingInFileNames": true // ✅ Enabled
  }
}
```

### Best Practices

#### 4.1 Avoid `any`, Use Unknown First

```typescript
// ❌ WRONG
function processData(data: any) {
  return data.transform(); // No type safety
}

// ✅ CORRECT
function processData(data: unknown) {
  if (typeof data === 'object' && data !== null && 'transform' in data) {
    return (data as { transform: () => void }).transform();
  }
  throw new Error('Invalid data');
}
```

#### 4.2 Type Guards with Signals

```typescript
// ✅ Type-safe signal discrimination
type RoomState = 
  | { status: 'idle' }
  | { status: 'loading' }
  | { status: 'success'; data: Room[] }
  | { status: 'error'; error: Error };

const roomState = signal<RoomState>({ status: 'idle' });

const displayRooms = computed(() => {
  const state = roomState();
  switch (state.status) {
    case 'success':
      return state.data;  // TypeScript knows data exists
    case 'loading':
      return [];
    case 'error':
      console.error(state.error);
      return [];
    case 'idle':
      return [];
  }
});
```

#### 4.3 Strict Constructor Initialization

```typescript
// ✅ All properties initialized
@Injectable()
class RoomService {
  private readonly http = inject(HttpClient);
  private readonly config = inject(APP_CONFIG);
  
  readonly rooms = signal<Room[]>([]);
  
  // No uninitialized properties
  constructor() {
    this.loadInitialData();
  }
}
```

---

## 5. Build Optimization & Code Splitting

### 5.1 Angular CLI Build Configuration

**angular.json updates:**
```json
{
  "projects": {
    "dashboard": {
      "architect": {
        "build": {
          "options": {
            "outputHashing": "all",           // Hash all outputs
            "vendorChunk": false,             // Don't separate vendor
            "commonChunk": true,              // Merge common deps
            "namedChunks": false,             // Minimize chunk names
            "optimization": true,             // Minify/uglify
            "sourceMap": false,               // No source maps in prod
            "buildOptimizer": true,           // Remove unused code
            "fileReplacements": [
              {
                "replace": "src/environments/environment.ts",
                "with": "src/environments/environment.prod.ts"
              }
            ],
            "scripts": [],
            "styles": ["src/styles.scss"],
            "inlineStyleLanguage": "scss"
          },
          "configurations": {
            "production": {
              "budgets": [
                {
                  "type": "initial",
                  "maximumWarning": "250kb",
                  "maximumError": "350kb"
                },
                {
                  "type": "anyComponentStyle",
                  "maximumWarning": "2kb",
                  "maximumError": "4kb"
                }
              ]
            }
          }
        }
      }
    }
  }
}
```

### 5.2 Route-Based Code Splitting

Already covered in §3.2 - use `loadComponent` for automatic chunk generation.

**Expected bundle sizes:**
- Initial: ~150KB (gzipped)
- Room detail chunk: ~25KB
- Analysis chunk: ~40KB

### 5.3 Esbuild Configuration

Angular 21 uses Esbuild by default (faster than Webpack):

```typescript
// No special config needed - already optimized
// Verify in build output:
// $ ng build --configuration production
// 
// Should show:
// ✔ Optimization: 15.234s
// ✔ Initial chunk bundle: 147KB (gzipped)
```

---

## 6. CDN vs Local Bundling for Third-Party Libraries

### Current Stack Analysis

**Mind Palace uses:**
- `d3@^7.9.0` (2.2MB minified)
- `rxjs@7.8.0` (included in Angular bundles)
- `tslib@^2.6.0` (small, always bundle)

### Recommendation: BUNDLE D3 (Don't use CDN)

**Why NOT use D3 from CDN:**

| Factor | CDN | Bundled |
|--------|-----|---------|
| Initial request | +1 extra request | Bundled, parallel |
| Latency | 50-200ms additional | 0ms (http2 multiplexed) |
| Caching | Shared, but often busted | App-specific, stable |
| Bundle size control | No treeshaking | Full treeshaking benefit |
| Security | Dependency on external service | Supply chain integrated |
| TypeScript types | Separate @types pkg | Inline types |
| **Recommendation** | ❌ | ✅ |

**Exception:** Only use CDN if D3 bundle exceeds 400KB gzipped (it won't).

### Security-First D3 Integration

```typescript
// ✅ src/environments/environment.prod.ts
export const environment = {
  production: true,
  d3: {
    // Never load from CDN - security risk with data
    bundled: true,
    version: '7.9.0'
  },
  csp: {
    scriptSrc: ["'self'"],  // No unsafe-inline, no external CDN
    styleSrc: ["'self'"],
  }
};

// ✅ service.ts - Safe D3 usage
import * as d3 from 'd3';

@Injectable({ providedIn: 'root' })
export class RoomVisualizationService {
  renderRoomMap(element: SVGElement, rooms: Room[]): void {
    // D3 is bundled, type-safe, no network dependency
    d3.select(element).selectAll('circle')
      .data(rooms, d => d.id)
      .enter()
      .append('circle');
  }
}
```

### CSP Headers (Helmet/Server-side)

```typescript
// production-server.ts or middleware
res.setHeader('Content-Security-Policy', [
  "default-src 'self'",
  "script-src 'self'",     // No CDN, no inline
  "style-src 'self'",      // No CDN, no inline
  "img-src 'self' data:",  // Allow data URIs for D3
  "connect-src 'self' /api", // API only
].join(';'));
```

---

## 7. Async Component Composition in Angular 21

### 7.1 Deferrable Views Pattern

**Latest Angular 21 feature for async composition:**

```typescript
// ✅ Template-level async handling
@Component({
  selector: 'app-dashboard',
  standalone: true,
  imports: [CommonModule, RoomDetailComponent],
  template: `
    <div class="dashboard">
      @defer (on interaction) {
        <app-room-detail [roomId]="selectedRoomId()" />
      } @placeholder {
        <div class="placeholder">Click to load room details</div>
      } @error {
        <div class="error">Failed to load room details</div>
      } @loading {
        <app-spinner />
      }
    </div>
  `
})
export class DashboardComponent {
  readonly selectedRoomId = signal<string | null>(null);
}
```

### 7.2 Async Component with Dependency Injection

```typescript
// ✅ Async initialization pattern
import { Component, computed } from '@angular/core';

@Component({
  selector: 'app-room-analysis',
  standalone: true,
  imports: [CommonModule],
  template: `
    @if (analyticsReady(); track analyticsReady()) {
      <app-d3-visualization [data]="analyticsData()" />
    }
  `
})
export class RoomAnalysisComponent {
  private readonly analyticsService = inject(AnalyticsService);
  
  // Async initialization with signal
  private readonly analyticsInitialized = signal(false);
  
  readonly analyticsData = computed(() => {
    if (this.analyticsInitialized()) {
      return this.analyticsService.getRoomMetrics();
    }
    return null;
  });
  
  readonly analyticsReady = computed(() => this.analyticsInitialized());
  
  constructor() {
    // Async initialization in constructor with proper cleanup
    this.analyticsService.initialize()
      .pipe(
        tap(() => this.analyticsInitialized.set(true)),
        catchError(() => {
          this.analyticsInitialized.set(false);
          return of(null);
        })
      )
      .subscribe();
  }
}
```

### 7.3 Parallel Component Loading

```typescript
// ✅ Load multiple async components in parallel
const routes: Routes = [
  {
    path: 'dashboard',
    component: DashboardComponent,
    resolve: {
      rooms: RoomResolver,
      analytics: AnalyticsResolver,  // Runs in parallel
      settings: SettingsResolver     // All load together
    }
  }
];

@Injectable({ providedIn: 'root' })
export class RoomResolver implements Resolve<Room[]> {
  constructor(private service: RoomService) {}

  resolve(): Observable<Room[]> {
    return this.service.getAllRooms();
  }
}
```

### 7.4 Injectable Factory Pattern

```typescript
// ✅ Create components with async dependencies
@Injectable()
export class ComponentFactory {
  async createRoomVisualization(element: HTMLElement): Promise<void> {
    // Load D3 and data in parallel
    const [d3Module, roomData] = await Promise.all([
      import('d3'),
      this.roomService.getRooms()
    ]);
    
    this.renderVisualization(element, d3Module, roomData);
  }
}
```

---

## Implementation Roadmap for Mind Palace Dashboard

### Phase 1: Testing Setup (Week 1)
- [ ] Install Vitest + @analogjs/vite-plugin-angular
- [ ] Create vitest.config.ts
- [ ] Convert 1 component spec to Vitest pattern
- [ ] Setup coverage reporting

### Phase 2: Signal Migration (Week 2-3)
- [ ] Create signal-based state services
- [ ] Migrate room list to input/output signals
- [ ] Add computed derived state
- [ ] Profile performance improvements

### Phase 3: Build Optimization (Week 3)
- [ ] Update angular.json with split strategies
- [ ] Configure route lazy loading for /analysis
- [ ] Set bundle size budgets
- [ ] Measure bundle size reduction

### Phase 4: TypeScript Strictness (Week 4)
- [ ] Add new strict compiler options
- [ ] Fix strict errors found
- [ ] Add type guards to d3 integration
- [ ] 100% strict mode compliance

### Phase 5: Performance Benchmarking (Week 4-5)
- [ ] Lighthouse audit before/after
- [ ] Change detection profiling
- [ ] Memory leak detection
- [ ] Target: LCP < 2.5s, TTI < 3.5s

---

## Recommended Package Updates

```json
{
  "dependencies": {
    "@angular/animations": "^21.0.6",
    "@angular/common": "^21.0.6",
    "@angular/compiler": "^21.0.6",
    "@angular/core": "^21.0.6",
    "@angular/forms": "^21.0.6",
    "@angular/platform-browser": "^21.0.6",
    "@angular/platform-browser-dynamic": "^21.0.6",
    "@angular/router": "^21.0.6",
    "@angular/cdk": "^21.0.0",
    "d3": "^7.9.0",
    "rxjs": "~7.8.0",
    "tslib": "^2.6.0",
    "zone.js": "~0.15.1"
  },
  "devDependencies": {
    "@angular-devkit/build-angular": "^21.0.4",
    "@angular/cli": "^21.0.4",
    "@angular/compiler-cli": "^21.0.6",
    "@types/d3": "^7.4.0",
    "@vitest/ui": "^1.1.0",
    "@vitest/coverage-v8": "^1.1.0",
    "vitest": "^1.1.0",
    "happy-dom": "^12.10.3",
    "typescript": "~5.9.3"
  }
}
```

---

## Common Pitfalls & Solutions

| Pitfall | Impact | Solution |
|---------|--------|----------|
| Not using OnPush with signals | 70% slower change detection | Always set `changeDetection: OnPush` |
| Subscribing in templates | Memory leaks | Use `toSignal()` or `async` pipe |
| Large components | Hard to test | Use input/output signals, 1 responsibility |
| D3 from CDN | Network + security risk | Bundle locally, control update cycles |
| No lazy loading | Large initial bundle | Use `loadComponent` in routes |
| Missing type guards | Runtime errors | Use discriminated unions with switch |
| Change detection debugging | Silent performance issues | Use Angular DevTools profiler |

---

## Validation Checklist

- [ ] All components use `ChangeDetectionStrategy.OnPush`
- [ ] No raw Observable subscriptions in components (use signals)
- [ ] All route-able features use `loadComponent`
- [ ] D3 bundled (not from CDN)
- [ ] TypeScript with `useUnknownInCatchVariables: true`
- [ ] 100% test coverage for state services
- [ ] Vitest running with coverage reporting
- [ ] Bundle size < 350KB initial (gzipped)
- [ ] Lighthouse score >= 90 (performance)
- [ ] Zero `console.error` or `console.warn` in tests

---

## References

- **Angular 21 Docs:** https://angular.io (v21)
- **Vitest Guide:** https://vitest.dev/guide/
- **TypeScript Handbook:** https://www.typescriptlang.org/docs/handbook/
- **Web Vitals:** https://web.dev/vitals/
- **Angular Performance:** https://angular.io/guide/change-detection

---

**Document Version:** 1.0  
**Last Updated:** January 2025  
**Next Review:** Q2 2025 (post-implementation)
