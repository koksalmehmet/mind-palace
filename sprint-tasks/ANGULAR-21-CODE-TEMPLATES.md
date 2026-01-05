# Angular 21 Implementation Guide - Code Templates

Quick-start code templates for Mind Palace Dashboard improvements.

---

## Template 1: Signal-Based State Service

```typescript
// src/app/core/services/room-state.service.ts
import { Injectable, inject, signal, computed, effect } from '@angular/core';
import { HttpClient } from '@angular/common/http';
import { of } from 'rxjs';
import { tap, catchError } from 'rxjs/operators';

export interface Room {
  id: string;
  name: string;
  description: string;
  isArchived: boolean;
  createdAt: Date;
  metadata?: Record<string, unknown>;
}

export interface RoomFilter {
  search?: string;
  archivedOnly?: boolean;
  dateRange?: { start: Date; end: Date };
}

@Injectable({ providedIn: 'root' })
export class RoomStateService {
  private readonly http = inject(HttpClient);

  // ===== State Signals =====
  private readonly rooms = signal<Room[]>([]);
  private readonly selectedRoomId = signal<string | null>(null);
  private readonly isLoading = signal(false);
  private readonly error = signal<string | null>(null);
  private readonly filters = signal<RoomFilter>({});

  // ===== Computed Derived State =====
  readonly selectedRoom = computed(() => {
    const id = this.selectedRoomId();
    return this.rooms().find(r => r.id === id) || null;
  });

  readonly filteredRooms = computed(() => {
    const rooms = this.rooms();
    const filters = this.filters();
    
    return rooms.filter(room => {
      if (filters.search) {
        const search = filters.search.toLowerCase();
        if (!room.name.toLowerCase().includes(search) &&
            !room.description.toLowerCase().includes(search)) {
          return false;
        }
      }
      
      if (filters.archivedOnly && !room.isArchived) {
        return false;
      }
      
      return true;
    });
  });

  readonly roomCount = computed(() => this.rooms().length);
  readonly filteredCount = computed(() => this.filteredRooms().length);
  readonly hasError = computed(() => this.error() !== null);
  readonly isEmpty = computed(() => this.rooms().length === 0 && !this.isLoading());

  // ===== Public Read-Only Signals =====
  readonly getRooms = this.rooms.asReadonly();
  readonly getFilteredRooms = this.filteredRooms.asReadonly();
  readonly getSelectedRoom = this.selectedRoom.asReadonly();
  readonly getIsLoading = this.isLoading.asReadonly();
  readonly getError = this.error.asReadonly();
  readonly getRoomCount = this.roomCount.asReadonly();

  constructor() {
    // Debug effect - logs state changes in development
    if (!this.isProduction()) {
      effect(() => {
        console.log('Rooms state changed:', {
          rooms: this.rooms().length,
          selected: this.selectedRoomId(),
          loading: this.isLoading(),
          error: this.error()
        });
      });
    }
  }

  // ===== Mutations =====
  loadRooms(): void {
    this.isLoading.set(true);
    this.error.set(null);

    this.http.get<Room[]>('/api/rooms')
      .pipe(
        tap(rooms => {
          this.rooms.set(rooms);
          this.isLoading.set(false);
        }),
        catchError(err => {
          this.error.set(err.message || 'Failed to load rooms');
          this.isLoading.set(false);
          return of([]);
        })
      )
      .subscribe();
  }

  addRoom(room: Omit<Room, 'id' | 'createdAt'>): void {
    this.isLoading.set(true);
    this.http.post<Room>('/api/rooms', room)
      .pipe(
        tap(newRoom => {
          this.rooms.update(rooms => [...rooms, newRoom]);
          this.isLoading.set(false);
        }),
        catchError(err => {
          this.error.set(err.message);
          this.isLoading.set(false);
          return of(null);
        })
      )
      .subscribe();
  }

  updateRoom(id: string, updates: Partial<Room>): void {
    this.http.patch<Room>(`/api/rooms/${id}`, updates)
      .pipe(
        tap(updatedRoom => {
          this.rooms.update(rooms =>
            rooms.map(r => r.id === id ? updatedRoom : r)
          );
        }),
        catchError(err => {
          this.error.set(err.message);
          return of(null);
        })
      )
      .subscribe();
  }

  deleteRoom(id: string): void {
    this.http.delete(`/api/rooms/${id}`)
      .pipe(
        tap(() => {
          this.rooms.update(rooms => rooms.filter(r => r.id !== id));
          if (this.selectedRoomId() === id) {
            this.selectedRoomId.set(null);
          }
        }),
        catchError(err => {
          this.error.set(err.message);
          return of(null);
        })
      )
      .subscribe();
  }

  selectRoom(id: string | null): void {
    this.selectedRoomId.set(id);
  }

  setFilter(filter: RoomFilter): void {
    this.filters.set(filter);
  }

  clearError(): void {
    this.error.set(null);
  }

  private isProduction(): boolean {
    return typeof ngDevMode === 'undefined' || !ngDevMode;
  }
}
```

---

## Template 2: Signal-Based Component with OnPush

```typescript
// src/app/features/rooms/room-list.component.ts
import { Component, ChangeDetectionStrategy, inject, output, input, computed } from '@angular/core';
import { CommonModule } from '@angular/common';
import { FormsModule } from '@angular/forms';
import { RoomStateService, Room } from '../../core/services/room-state.service';
import { RoomCardComponent } from './room-card.component';

@Component({
  selector: 'app-room-list',
  standalone: true,
  imports: [CommonModule, FormsModule, RoomCardComponent],
  template: `
    <div class="room-list">
      <div class="toolbar">
        <input
          type="text"
          placeholder="Search rooms..."
          [value]="searchQuery()"
          (change)="onSearchChange($event)"
          class="search-input"
        />
        <button (click)="onNewRoom()" class="btn btn-primary">
          New Room
        </button>
      </div>

      <div class="stats">
        <span>Total: {{ roomCount() }} | Filtered: {{ filteredCount() }}</span>
      </div>

      @if (isLoading()) {
        <div class="spinner">Loading rooms...</div>
      } @else if (isEmpty()) {
        <div class="empty-state">
          <p>No rooms found. Create your first room to get started.</p>
        </div>
      } @else {
        <div class="room-grid">
          @for (room of filteredRooms(); track room.id) {
            <app-room-card
              [room]="room"
              [isSelected]="selectedRoomId() === room.id"
              (select)="roomSelected.emit($event)"
              (delete)="roomDeleted.emit($event)"
              (archive)="roomArchived.emit($event)"
            />
          }
        </div>
      }

      @if (hasError()) {
        <div class="error-banner">
          <p>{{ error() }}</p>
          <button (click)="onDismissError()">Dismiss</button>
        </div>
      }
    </div>
  `,
  styles: [`
    .room-list {
      display: flex;
      flex-direction: column;
      gap: 1rem;
      padding: 1rem;
    }

    .toolbar {
      display: flex;
      gap: 1rem;
      align-items: center;
    }

    .search-input {
      flex: 1;
      padding: 0.5rem;
      border: 1px solid #ccc;
      border-radius: 4px;
    }

    .room-grid {
      display: grid;
      grid-template-columns: repeat(auto-fill, minmax(300px, 1fr));
      gap: 1rem;
    }

    .empty-state, .spinner, .error-banner {
      padding: 2rem;
      text-align: center;
      border: 1px solid #ddd;
      border-radius: 4px;
    }

    .stats {
      font-size: 0.9rem;
      color: #666;
    }
  `],
  changeDetection: ChangeDetectionStrategy.OnPush
})
export class RoomListComponent {
  private readonly roomService = inject(RoomStateService);

  // ===== Outputs =====
  readonly roomSelected = output<string>();
  readonly roomDeleted = output<string>();
  readonly roomArchived = output<string>();
  readonly newRoomRequested = output<void>();

  // ===== Local State =====
  private readonly searchQuery = signal('');

  // ===== Service Inputs =====
  readonly rooms = this.roomService.getFilteredRooms;
  readonly isLoading = this.roomService.getIsLoading;
  readonly error = this.roomService.getError;
  readonly hasError = this.roomService.getError; // Computed in service

  // ===== Derived Local State =====
  readonly filteredRooms = computed(() => {
    const query = this.searchQuery().toLowerCase();
    const rooms = this.roomService.getFilteredRooms();

    if (!query) return rooms;

    return rooms.filter(r =>
      r.name.toLowerCase().includes(query) ||
      r.description.toLowerCase().includes(query)
    );
  });

  readonly isEmpty = computed(() =>
    this.filteredRooms().length === 0 && !this.isLoading()
  );

  readonly roomCount = computed(() => this.roomService.getRoomCount());
  readonly filteredCount = computed(() => this.filteredRooms().length);
  readonly selectedRoomId = this.roomService.getSelectedRoom; // For display

  onSearchChange(event: Event): void {
    const target = event.target as HTMLInputElement;
    this.searchQuery.set(target.value);
  }

  onNewRoom(): void {
    this.newRoomRequested.emit();
  }

  onDismissError(): void {
    this.roomService.clearError();
  }
}
```

---

## Template 3: Vitest Component Test

```typescript
// src/app/features/rooms/room-list.component.spec.ts
import { ComponentFixture, TestBed } from '@angular/core/testing';
import { HttpClientTestingModule } from '@angular/common/http/testing';
import { describe, it, expect, beforeEach, vi } from 'vitest';
import { signal } from '@angular/core';
import { RoomListComponent } from './room-list.component';
import { RoomStateService } from '../../core/services/room-state.service';

describe('RoomListComponent', () => {
  let component: RoomListComponent;
  let fixture: ComponentFixture<RoomListComponent>;
  let roomService: RoomStateService;

  const mockRooms = [
    { id: '1', name: 'Living Room', description: 'Main gathering space', isArchived: false, createdAt: new Date() },
    { id: '2', name: 'Bedroom', description: 'Private space', isArchived: false, createdAt: new Date() },
    { id: '3', name: 'Archive Room', description: 'Old data', isArchived: true, createdAt: new Date() }
  ];

  beforeEach(async () => {
    const mockService = {
      getFilteredRooms: signal(mockRooms),
      getIsLoading: signal(false),
      getError: signal(null),
      getRoomCount: signal(3),
      getSelectedRoom: signal(null),
      clearError: vi.fn(),
      selectRoom: vi.fn()
    };

    await TestBed.configureTestingModule({
      imports: [RoomListComponent, HttpClientTestingModule],
      providers: [
        { provide: RoomStateService, useValue: mockService }
      ]
    }).compileComponents();

    roomService = TestBed.inject(RoomStateService);
    fixture = TestBed.createComponent(RoomListComponent);
    component = fixture.componentInstance;
  });

  describe('rendering', () => {
    it('should render room list', () => {
      fixture.detectChanges();

      const cards = fixture.nativeElement.querySelectorAll('app-room-card');
      expect(cards).toHaveLength(3);
    });

    it('should show empty state when no rooms', () => {
      (roomService.getFilteredRooms as any).set([]);
      fixture.detectChanges();

      const emptyState = fixture.nativeElement.querySelector('.empty-state');
      expect(emptyState).toBeTruthy();
    });

    it('should show loading state', () => {
      (roomService.getIsLoading as any).set(true);
      fixture.detectChanges();

      const spinner = fixture.nativeElement.querySelector('.spinner');
      expect(spinner).toBeTruthy();
    });

    it('should show error banner', () => {
      (roomService.getError as any).set('Failed to load rooms');
      fixture.detectChanges();

      const errorBanner = fixture.nativeElement.querySelector('.error-banner');
      expect(errorBanner?.textContent).toContain('Failed to load rooms');
    });
  });

  describe('search', () => {
    it('should filter rooms by search query', () => {
      fixture.detectChanges();

      const input = fixture.nativeElement.querySelector('.search-input');
      input.value = 'Living';
      input.dispatchEvent(new Event('change'));
      fixture.detectChanges();

      expect(component.filteredRooms().length).toBe(1);
      expect(component.filteredRooms()[0].name).toBe('Living Room');
    });

    it('should be case-insensitive', () => {
      fixture.detectChanges();

      const input = fixture.nativeElement.querySelector('.search-input');
      input.value = 'BEDROOM';
      input.dispatchEvent(new Event('change'));
      fixture.detectChanges();

      expect(component.filteredRooms().length).toBe(1);
      expect(component.filteredRooms()[0].name).toBe('Bedroom');
    });
  });

  describe('outputs', () => {
    it('should emit roomSelected when room card selected', () => {
      const emitSpy = vi.spyOn(component.roomSelected, 'emit');
      fixture.detectChanges();

      expect(emitSpy).toHaveBeenCalledWith('1');
    });

    it('should emit newRoomRequested when new button clicked', () => {
      const emitSpy = vi.spyOn(component.newRoomRequested, 'emit');
      fixture.detectChanges();

      const newBtn = fixture.nativeElement.querySelector('.btn-primary');
      newBtn.click();

      expect(emitSpy).toHaveBeenCalled();
    });
  });

  describe('error handling', () => {
    it('should call clearError when dismiss clicked', () => {
      (roomService.getError as any).set('Test error');
      fixture.detectChanges();

      const dismissBtn = fixture.nativeElement.querySelector('.error-banner button');
      dismissBtn.click();

      expect(roomService.clearError).toHaveBeenCalled();
    });
  });
});
```

---

## Template 4: Vitest Service Test

```typescript
// src/app/core/services/room-state.service.spec.ts
import { TestBed } from '@angular/core/testing';
import { HttpClientTestingModule, HttpTestingController } from '@angular/common/http/testing';
import { describe, it, expect, beforeEach, afterEach } from 'vitest';
import { RoomStateService, Room } from './room-state.service';

describe('RoomStateService', () => {
  let service: RoomStateService;
  let httpMock: HttpTestingController;

  const mockRooms: Room[] = [
    {
      id: '1',
      name: 'Library',
      description: 'Knowledge repository',
      isArchived: false,
      createdAt: new Date()
    },
    {
      id: '2',
      name: 'Archive',
      description: 'Archived items',
      isArchived: true,
      createdAt: new Date()
    }
  ];

  beforeEach(() => {
    TestBed.configureTestingModule({
      imports: [HttpClientTestingModule],
      providers: [RoomStateService]
    });

    service = TestBed.inject(RoomStateService);
    httpMock = TestBed.inject(HttpTestingController);
  });

  afterEach(() => {
    httpMock.verify();
  });

  describe('loadRooms', () => {
    it('should fetch rooms and update signal', () => {
      service.loadRooms();

      const req = httpMock.expectOne('/api/rooms');
      expect(req.request.method).toBe('GET');
      req.flush(mockRooms);

      expect(service.getRooms().length).toBe(2);
      expect(service.getIsLoading()).toBe(false);
    });

    it('should set loading state during fetch', () => {
      service.loadRooms();

      expect(service.getIsLoading()).toBe(true);

      const req = httpMock.expectOne('/api/rooms');
      req.flush(mockRooms);

      expect(service.getIsLoading()).toBe(false);
    });

    it('should handle errors', () => {
      service.loadRooms();

      const req = httpMock.expectOne('/api/rooms');
      req.error(new ErrorEvent('Network error'), { status: 500 });

      expect(service.getError()).toBeTruthy();
      expect(service.getIsLoading()).toBe(false);
    });
  });

  describe('computed state', () => {
    beforeEach(() => {
      service['rooms'].set(mockRooms);
    });

    it('should compute filtered rooms excluding archived', () => {
      expect(service.filteredRooms().length).toBe(1);
      expect(service.filteredRooms()[0].name).toBe('Library');
    });

    it('should compute room count', () => {
      expect(service.roomCount()).toBe(2);
    });

    it('should select room by id', () => {
      service.selectRoom('1');

      expect(service.selectedRoom()).toEqual(mockRooms[0]);
    });
  });

  describe('mutations', () => {
    it('should add new room', () => {
      service['rooms'].set(mockRooms);

      const newRoom: Omit<Room, 'id' | 'createdAt'> = {
        name: 'New Room',
        description: 'Test',
        isArchived: false
      };

      service.addRoom(newRoom);

      const req = httpMock.expectOne('/api/rooms');
      const createdRoom: Room = { ...newRoom, id: '3', createdAt: new Date() };
      req.flush(createdRoom);

      expect(service.getRooms().length).toBe(3);
    });

    it('should update room', () => {
      service['rooms'].set(mockRooms);

      service.updateRoom('1', { name: 'Updated Library' });

      const req = httpMock.expectOne('/api/rooms/1');
      const updated: Room = { ...mockRooms[0], name: 'Updated Library' };
      req.flush(updated);

      const room = service.getRooms().find(r => r.id === '1');
      expect(room?.name).toBe('Updated Library');
    });

    it('should delete room', () => {
      service['rooms'].set(mockRooms);

      service.deleteRoom('1');

      const req = httpMock.expectOne('/api/rooms/1');
      req.flush(null);

      expect(service.getRooms().length).toBe(1);
      expect(service.getRooms().some(r => r.id === '1')).toBe(false);
    });
  });

  describe('filtering', () => {
    beforeEach(() => {
      service['rooms'].set(mockRooms);
    });

    it('should filter by search term', () => {
      service.setFilter({ search: 'Library' });

      expect(service.filteredRooms().length).toBe(1);
    });

    it('should handle empty search', () => {
      service.setFilter({ search: '' });

      expect(service.filteredRooms().length).toBe(1); // Still excludes archived
    });
  });
});
```

---

## Template 5: d3 Service with Type Safety

```typescript
// src/app/core/services/d3-visualization.service.ts
import { Injectable, inject } from '@angular/core';
import * as d3 from 'd3';
import { Room } from './room-state.service';

export interface D3Node {
  id: string;
  name: string;
  value: number;
  archived: boolean;
}

export interface D3Link {
  source: string;
  target: string;
  strength: number;
}

@Injectable({ providedIn: 'root' })
export class D3VisualizationService {
  /**
   * Render D3 force-directed graph for room connections
   */
  renderRoomGraph(
    selector: string,
    rooms: Room[],
    links: D3Link[]
  ): void {
    const element = document.querySelector(selector);
    if (!element || !(element instanceof SVGElement)) {
      console.error(`Invalid selector or element: ${selector}`);
      return;
    }

    const width = element.clientWidth || 800;
    const height = element.clientHeight || 600;

    // Convert rooms to D3 nodes
    const nodes: D3Node[] = rooms.map(room => ({
      id: room.id,
      name: room.name,
      value: 10,
      archived: room.isArchived
    }));

    // Create simulation
    const simulation = d3.forceSimulation<D3Node>(nodes)
      .force('link', d3.forceLink<D3Node, D3Link>()
        .id(d => d.id)
        .links(links)
      )
      .force('charge', d3.forceManyBody().strength(-300))
      .force('center', d3.forceCenter(width / 2, height / 2));

    // Clear previous
    d3.select(selector).selectAll('*').remove();

    const svg = d3.select(selector)
      .attr('width', width)
      .attr('height', height);

    // Draw links
    const link = svg.selectAll('line')
      .data(links)
      .enter()
      .append('line')
      .attr('stroke', '#999')
      .attr('stroke-opacity', 0.6)
      .attr('stroke-width', d => Math.sqrt(d.strength) * 2);

    // Draw nodes
    const node = svg.selectAll('circle')
      .data(nodes)
      .enter()
      .append('circle')
      .attr('r', d => d.value)
      .attr('fill', d => d.archived ? '#ccc' : '#4285f4')
      .call(drag(simulation) as any);

    // Draw labels
    const label = svg.selectAll('text')
      .data(nodes)
      .enter()
      .append('text')
      .text(d => d.name)
      .attr('font-size', 10)
      .attr('dx', 15)
      .attr('dy', 4);

    // Update positions on tick
    simulation.on('tick', () => {
      link
        .attr('x1', d => d.source.x ?? 0)
        .attr('y1', d => d.source.y ?? 0)
        .attr('x2', d => d.target.x ?? 0)
        .attr('y2', d => d.target.y ?? 0);

      node
        .attr('cx', d => d.x ?? 0)
        .attr('cy', d => d.y ?? 0);

      label
        .attr('x', d => d.x ?? 0)
        .attr('y', d => d.y ?? 0);
    });
  }

  private drag(simulation: d3.Simulation<D3Node, D3Link>) {
    return d3.drag<any, D3Node>()
      .on('start', (event: d3.D3DragEvent<any, D3Node, any>, d) => {
        if (!event.active) simulation.alphaTarget(0.3).restart();
        d.fx = d.x;
        d.fy = d.y;
      })
      .on('drag', (event: d3.D3DragEvent<any, D3Node, any>, d) => {
        d.fx = event.x;
        d.fy = event.y;
      })
      .on('end', (event: d3.D3DragEvent<any, D3Node, any>, d) => {
        if (!event.active) simulation.alphaTarget(0);
        d.fx = null;
        d.fy = null;
      });
  }
}
```

---

## Template 6: Lazy Loaded Feature Route

```typescript
// src/app/app.routes.ts
import { Routes } from '@angular/router';
import { DashboardComponent } from './features/dashboard/dashboard.component';

export const routes: Routes = [
  {
    path: '',
    redirectTo: 'dashboard',
    pathMatch: 'full'
  },
  {
    path: 'dashboard',
    component: DashboardComponent,
    children: [
      {
        path: 'rooms/:id',
        loadComponent: () =>
          import('./features/room-detail/room-detail.component')
            .then(m => m.RoomDetailComponent),
        data: { title: 'Room Details' }
      },
      {
        path: 'analysis',
        loadComponent: () =>
          import('./features/analysis/analysis.component')
            .then(m => m.AnalysisComponent),
        data: { title: 'Analysis' }
      },
      {
        path: 'connections',
        loadComponent: () =>
          import('./features/connections/connections.component')
            .then(m => m.ConnectionsComponent),
        data: { title: 'Room Connections' }
      }
    ]
  },
  {
    path: 'settings',
    loadComponent: () =>
      import('./features/settings/settings.component')
        .then(m => m.SettingsComponent),
    data: { title: 'Settings' }
  },
  {
    path: '**',
    loadComponent: () =>
      import('./features/not-found/not-found.component')
        .then(m => m.NotFoundComponent)
  }
];
```

---

## Template 7: Strict TypeScript Service

```typescript
// src/app/core/services/analytics.service.ts
import { Injectable, inject } from '@angular/core';
import { HttpClient } from '@angular/common/http';
import { Observable, throwError } from 'rxjs';
import { catchError, map } from 'rxjs/operators';

export interface RoomMetric {
  roomId: string;
  views: number;
  engagementScore: number;
  lastAccessed: Date;
}

export interface AnalyticsError {
  code: string;
  message: string;
  timestamp: Date;
}

// Type-safe error handling
type AnalyticsResult<T> =
  | { ok: true; data: T }
  | { ok: false; error: AnalyticsError };

@Injectable({ providedIn: 'root' })
export class AnalyticsService {
  private readonly http = inject(HttpClient);

  /**
   * Fetch analytics safely with proper error typing
   */
  getMetrics(roomId: string): Observable<AnalyticsResult<RoomMetric>> {
    return this.http.get<RoomMetric>(`/api/analytics/${roomId}`)
      .pipe(
        map(data => ({ ok: true as const, data })),
        catchError(error => {
          const analyticsError: AnalyticsError = this.parseError(error);
          return throwError(() => ({ ok: false as const, error: analyticsError }));
        })
      );
  }

  /**
   * Type-safe error parsing
   */
  private parseError(error: unknown): AnalyticsError {
    if (error instanceof Error) {
      return {
        code: 'UNKNOWN_ERROR',
        message: error.message,
        timestamp: new Date()
      };
    }

    if (typeof error === 'object' && error !== null && 'error' in error) {
      const httpError = error as { error?: { message?: string; code?: string } };
      return {
        code: httpError.error?.code ?? 'HTTP_ERROR',
        message: httpError.error?.message ?? 'An error occurred',
        timestamp: new Date()
      };
    }

    return {
      code: 'UNKNOWN_ERROR',
      message: 'An unknown error occurred',
      timestamp: new Date()
    };
  }
}
```

---

**Using these templates:**

1. Copy-paste into your component/service files
2. Update import paths as needed
3. Adapt to your actual Room interface
4. Run `npm run test` with Vitest
5. Test coverage will highlight any issues

**Next Steps:**
- Implement Template 1 (RoomStateService)
- Create specs using Template 3-4
- Update components to use input/output signals (Template 2)
- Integrate D3 service (Template 5)
- Add route lazy loading (Template 6)
