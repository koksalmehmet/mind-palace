# Angular 21 Quick Reference Card

**Print this and keep at your desk during implementation**

---

## 🎯 Top 5 Changes

```
1. SWITCH TO VITEST (5x faster tests)
   npm install -D vitest @vitest/ui @vitest/coverage-v8 happy-dom
   npm run test (instead of ng test)

2. USE SIGNALS (70% less change detection)
   const rooms = signal<Room[]>([]);
   readonly getRooms = this.rooms.asReadonly();

3. ADD LAZY LOADING (40% smaller initial bundle)
   loadComponent: () => import('./room-detail')

4. ALWAYS USE OnPush (prevents change detection thrashing)
   changeDetection: ChangeDetectionStrategy.OnPush

5. ENHANCE TYPESCRIPT (100% strict mode)
   useUnknownInCatchVariables: true
```

---

## 🏗️ Architecture Pattern

```
┌────────────────────────────────────┐
│    Component (OnPush)              │
│    - input signals                 │
│    - output signals                │
│    - NO subscriptions              │
└──────────────┬──────────────────────┘
               │ inject()
               ▼
┌────────────────────────────────────┐
│    Signal-Based Service            │
│    - signal (mutable state)         │
│    - computed (derived)             │
│    - effect (side effects)          │
│    - .asReadonly() exports          │
└──────────────┬──────────────────────┘
               │ internal .pipe()
               ▼
         HTTP Observable
         (converts to signal)
```

---

## 📝 Component Template

```typescript
// ✅ DO THIS
@Component({
  selector: 'app-room-list',
  standalone: true,
  imports: [CommonModule],
  changeDetection: ChangeDetectionStrategy.OnPush  // ← REQUIRED
})
export class RoomListComponent {
  // Input signals
  readonly rooms = input.required<Room[]>();
  readonly isLoading = input(false);
  
  // Output signals
  readonly roomSelected = output<string>();
  
  // Local state (computed derived values)
  readonly filteredRooms = computed(() => {
    return this.rooms().filter(r => !r.archived);
  });
  
  // Methods
  onSelect(id: string) {
    this.roomSelected.emit(id);
  }
}
```

**Template:**
```html
<div class="room-list">
  @for (room of filteredRooms(); track room.id) {
    <app-room-card 
      [room]="room"
      (select)="onSelect($event)"
    />
  }
</div>
```

---

## 📦 Service Template

```typescript
// ✅ DO THIS
@Injectable({ providedIn: 'root' })
export class RoomService {
  private readonly http = inject(HttpClient);
  
  // Private mutable signals
  private readonly rooms = signal<Room[]>([]);
  private readonly selectedId = signal<string | null>(null);
  
  // Public read-only signals
  readonly getRooms = this.rooms.asReadonly();
  readonly getSelectedRoom = computed(() => {
    const id = this.selectedId();
    return this.rooms().find(r => r.id === id);
  });
  
  // Mutations
  loadRooms() {
    this.http.get<Room[]>('/api/rooms')
      .pipe(
        tap(rooms => this.rooms.set(rooms))
      )
      .subscribe();
  }
  
  selectRoom(id: string) {
    this.selectedId.set(id);
  }
}
```

---

## 🧪 Test Template

```typescript
// ✅ DO THIS WITH VITEST
import { describe, it, expect, beforeEach } from 'vitest';
import { TestBed } from '@angular/core/testing';

describe('RoomService', () => {
  let service: RoomService;

  beforeEach(() => {
    TestBed.configureTestingModule({
      providers: [RoomService]
    });
    service = TestBed.inject(RoomService);
  });

  it('should load rooms', () => {
    service.loadRooms();
    expect(service.getRooms().length).toBeGreaterThan(0);
  });

  it('should select room', () => {
    service['rooms'].set([{ id: '1', name: 'Test' }]);
    service.selectRoom('1');
    expect(service.getSelectedRoom()?.id).toBe('1');
  });
});
```

---

## 🔄 Signal Patterns

### Read-Only Signal Export
```typescript
private readonly data = signal<T[]>([]);
readonly getData = this.data.asReadonly();
```

### Computed Derived Values
```typescript
readonly filtered = computed(() => 
  this.getData().filter(item => !item.archived)
);
```

### Convert Observable to Signal
```typescript
readonly data = toSignal(
  this.http.get('/api/data'),
  { initialValue: [] }
);
```

### Watch Signal Changes
```typescript
constructor() {
  effect(() => {
    console.log('Data changed:', this.getData());
  });
}
```

### Update Signal
```typescript
// Direct replacement
this.data.set([...]);

// Update existing
this.data.update(items => [
  ...items,
  newItem
]);
```

---

## 🚫 Common Mistakes

### ❌ WRONG: Subscriptions in Components
```typescript
constructor(service: RoomService) {
  service.rooms$.subscribe(rooms => {  // ← Memory leak!
    this.rooms = rooms;
  });
}
```

### ✅ RIGHT: Read Signals Directly
```typescript
readonly rooms = inject(RoomService).getRooms;
```

---

### ❌ WRONG: Modifying Service State from Component
```typescript
rooms.push(newRoom);  // ← Race conditions!
```

### ✅ RIGHT: Call Service Method
```typescript
this.roomService.addRoom(newRoom);  // ← Service handles mutation
```

---

### ❌ WRONG: No OnPush Change Detection
```typescript
@Component({
  selector: 'app-room'
  // changeDetection omitted = Default strategy
})
```

### ✅ RIGHT: Always Use OnPush with Signals
```typescript
@Component({
  selector: 'app-room',
  changeDetection: ChangeDetectionStrategy.OnPush
})
```

---

## 📊 Performance Checklist

- [ ] Component has `changeDetection: ChangeDetectionStrategy.OnPush`
- [ ] No subscriptions in components (use signals instead)
- [ ] Service uses private `signal()` + public `asReadonly()`
- [ ] Computed values used for derived state
- [ ] Routes use `loadComponent` for lazy loading
- [ ] D3 bundled locally (not from CDN)
- [ ] Tests run with Vitest (not Karma)
- [ ] Bundle size < 250KB gzipped (initial)
- [ ] Change detection < 16ms per frame
- [ ] Zero `any` type usage

---

## 🎯 TypeScript Strict Mode Checklist

```json
{
  "compilerOptions": {
    "strict": true,                          ✅
    "noImplicitOverride": true,              ✅
    "noPropertyAccessFromIndexSignature": true, ✅
    "noImplicitReturns": true,               ✅
    "noFallthroughCasesInSwitch": true,      ✅
    "useUnknownInCatchVariables": true,      ← ADD
    "noUncheckedIndexedAccess": true,        ← ADD
    "noImplicitThis": true,                  ← ADD
    "exactOptionalPropertyTypes": true       ← ADD
  }
}
```

---

## 🧠 Error Handling (Strict)

### ❌ WRONG: Any Catch
```typescript
try {
  // ...
} catch (e) {  // ← 'e' is any
  console.log(e.message);
}
```

### ✅ RIGHT: Type Unknown
```typescript
try {
  // ...
} catch (error: unknown) {
  if (error instanceof Error) {
    console.log(error.message);
  }
}
```

---

## 📦 Build & Bundle Commands

```bash
# Development
npm run start              # Start dev server

# Testing
npm run test              # Watch mode (Vitest)
npm run test:ui          # Browser UI
npm run test:coverage    # Coverage report

# Production
npm run build            # Build for prod
ng build --stats-json   # Analyze bundle
npm run lighthouse       # Performance audit

# Type Checking
npx tsc --noEmit        # Type check only
```

---

## 🔗 Lazy Loading Routes

```typescript
const routes: Routes = [
  {
    path: 'dashboard',
    loadComponent: () =>
      import('./dashboard/dashboard.component')
        .then(m => m.DashboardComponent)
  },
  {
    path: 'analysis',
    loadComponent: () =>
      import('./analysis/analysis.component')
        .then(m => m.AnalysisComponent)
  }
];
```

**Result:** Separate chunks for each route, loaded on-demand.

---

## 📏 Bundle Size Targets

| Type | Limit | Check Command |
|------|-------|---|
| Initial bundle | < 250 KB | `ng build --stats-json` |
| Route chunk | < 100 KB | See dist/ files |
| CSS | < 50 KB | Check styles.css |
| D3 library | ~150 KB | Normal size |

---

## ⚡ Performance Metrics

| Metric | Target | Acceptable |
|--------|--------|-----------|
| FCP | < 1.8s | < 2.5s |
| LCP | < 2.5s | < 4s |
| TTI | < 3.5s | < 5s |
| CLS | < 0.1 | < 0.25 |
| Perf Score | > 90 | > 80 |

---

## 🛠️ Vitest Setup

```bash
# Install
npm install -D vitest @vitest/ui @vitest/coverage-v8

# Create vitest.config.ts
import { defineConfig } from 'vitest/config'

export default defineConfig({
  test: {
    globals: true,
    environment: 'happy-dom',
    setupFiles: ['src/test-setup.ts']
  }
})

# Update package.json
"test": "vitest"
"test:ui": "vitest --ui"
"test:coverage": "vitest run --coverage"

# Run
npm run test
```

---

## 💡 Pro Tips

1. **Use `track` in @for loops** - Improves performance
   ```html
   @for (room of rooms(); track room.id) { ... }
   ```

2. **Keep computed() simple** - Don't do heavy work
   ```typescript
   readonly filtered = computed(() => 
     this.rooms().filter(r => !r.archived)  // ← Simple
   );
   ```

3. **Effect cleanup automatic** - No manual unsubscribe needed
   ```typescript
   effect(() => {
     // This automatically cleans up when component destroys
     this.logger.log(this.signal());
   });
   ```

4. **Test signals synchronously** - No async needed
   ```typescript
   it('should update signal', () => {
     service.selectRoom('1');
     expect(service.getSelectedRoom()?.id).toBe('1');  // Sync!
   });
   ```

5. **Use `asReadonly()` for public signals** - Prevents accidental mutation
   ```typescript
   readonly data = this.internalData.asReadonly();
   ```

---

## 🔗 Quick Links

**Documents:**
- [Full Best Practices](./ANGULAR-21-BEST-PRACTICES.md)
- [Code Templates](./ANGULAR-21-CODE-TEMPLATES.md)
- [Implementation Checklist](./ANGULAR-21-IMPLEMENTATION-CHECKLIST.md)
- [Executive Summary](./ANGULAR-21-EXECUTIVE-SUMMARY.md)

**Official Docs:**
- [Angular 21](https://angular.io)
- [Vitest](https://vitest.dev)
- [TypeScript 5.x](https://www.typescriptlang.org/docs/handbook/)

---

## 📞 Help Commands

```bash
# Check Vitest works
npm run test -- --version

# Check TypeScript version
npx tsc --version

# Check Angular version
ng version

# Analyze bundle
ng build --stats-json && \
  npx webpack-bundle-analyzer dist/*/stats.json

# TypeScript strict check
npx tsc --noEmit --strict

# Run single test
npm run test -- room.spec

# Watch specific file
npm run test -- room.spec --watch
```

---

## ✅ Pre-Implementation Sign-Off

Before starting, verify:
- [ ] Read this quick reference card
- [ ] Read Executive Summary (5 min)
- [ ] Understood the 5 main changes
- [ ] Vitest/TypeScript tools installed
- [ ] Current tests passing
- [ ] Current build succeeding
- [ ] Have questions? Asked tech lead

---

**Print this → Keep at desk → Reference while coding**

**Version:** 1.0  
**Last Updated:** January 5, 2025

