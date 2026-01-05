import { describe, it, expect, beforeEach, afterEach } from 'vitest';
import { TestBed } from '@angular/core/testing';
import { provideHttpClient } from '@angular/common/http';
import { HttpTestingController, provideHttpClientTesting } from '@angular/common/http/testing';
import { ApiService, Stats } from './api.service';

describe('ApiService', () => {
  let service: ApiService;
  let httpMock: HttpTestingController;

  beforeEach(() => {
    TestBed.configureTestingModule({
      providers: [
        ApiService,
        provideHttpClient(),
        provideHttpClientTesting()
      ]
    });

    service = TestBed.inject(ApiService);
    httpMock = TestBed.inject(HttpTestingController);
  });

  afterEach(() => {
    httpMock.verify();
  });

  it('should be created', () => {
    expect(service).toBeTruthy();
  });

  describe('Health & Stats', () => {
    it('should fetch health status', () => {
      const mockHealth = { status: 'healthy', timestamp: '2025-01-01T00:00:00Z' };

      service.getHealth().subscribe(health => {
        expect(health).toEqual(mockHealth);
      });

      const req = httpMock.expectOne('/api/health');
      expect(req.request.method).toBe('GET');
      req.flush(mockHealth);
    });

    it('should fetch stats', () => {
      const mockStats: Stats = {
        sessions: { total: 10, active: 2 },
        learnings: 100,
        filesTracked: 50,
        rooms: 5
      };

      service.getStats().subscribe(stats => {
        expect(stats).toEqual(mockStats);
      });

      const req = httpMock.expectOne('/api/stats');
      expect(req.request.method).toBe('GET');
      req.flush(mockStats);
    });
  });

  describe('Sessions', () => {
    it('should fetch sessions with default params', () => {
      const mockResponse = { sessions: [], count: 0 };

      service.getSessions().subscribe(response => {
        expect(response).toEqual(mockResponse);
      });

      const req = httpMock.expectOne(r => r.url === '/api/sessions');
      expect(req.request.params.get('active')).toBe('false');
      expect(req.request.params.get('limit')).toBe('50');
      req.flush(mockResponse);
    });

    it('should fetch active sessions only', () => {
      service.getSessions(true, 10).subscribe(() => {
        // subscription
      });

      const req = httpMock.expectOne(r => r.url === '/api/sessions');
      expect(req.request.params.get('active')).toBe('true');
      expect(req.request.params.get('limit')).toBe('10');
      req.flush({ sessions: [], count: 0 });
    });

    it('should fetch single session by id', () => {
      const sessionId = 'test-123';

      service.getSession(sessionId).subscribe(response => {
        expect(response.session.id).toBe(sessionId);
      });

      const req = httpMock.expectOne(`/api/sessions/${sessionId}`);
      expect(req.request.method).toBe('GET');
      req.flush({ session: { id: sessionId }, activities: [] });
    });
  });

  describe('Learnings', () => {
    it('should fetch learnings with default params', () => {
      service.getLearnings().subscribe(response => {
        expect(response.count).toBe(0);
      });

      const req = httpMock.expectOne(r => r.url === '/api/learnings');
      expect(req.request.params.get('limit')).toBe('50');
      expect(req.request.params.has('scope')).toBe(false);
      expect(req.request.params.has('query')).toBe(false);
      req.flush({ learnings: [], count: 0 });
    });

    it('should fetch learnings with scope and query', () => {
      service.getLearnings('workspace', 'test', 10).subscribe(response => {
        expect(response.count).toBeGreaterThanOrEqual(0);
      });

      const req = httpMock.expectOne(r => r.url === '/api/learnings');
      expect(req.request.params.get('scope')).toBe('workspace');
      expect(req.request.params.get('query')).toBe('test');
      expect(req.request.params.get('limit')).toBe('10');
      req.flush({ learnings: [], count: 0 });
    });
  });

  describe('Search', () => {
    it('should perform search with query', () => {
      const query = 'test search';

      service.search(query).subscribe(response => {
        expect(response).toBeDefined();
      });

      const req = httpMock.expectOne(r => r.url === '/api/search');
      expect(req.request.params.get('q')).toBe(query);
      expect(req.request.params.get('limit')).toBe('20');
      req.flush({ results: [] });
    });
  });

  describe('Rooms', () => {
    it('should fetch rooms', () => {
      const mockResponse = { rooms: [], count: 0 };

      service.getRooms().subscribe(response => {
        expect(response).toEqual(mockResponse);
      });

      const req = httpMock.expectOne('/api/rooms');
      expect(req.request.method).toBe('GET');
      req.flush(mockResponse);
    });
  });

  describe('Error Handling', () => {
    it('should handle HTTP errors', () => {
      service.getHealth().subscribe({
        error: (error) => {
          expect(error.status).toBe(500);
        }
      });

      const req = httpMock.expectOne('/api/health');
      req.flush('Server error', { status: 500, statusText: 'Internal Server Error' });
    });
  });
});
