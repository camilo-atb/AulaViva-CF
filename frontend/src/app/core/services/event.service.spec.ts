import { TestBed } from '@angular/core/testing';
import { provideHttpClient } from '@angular/common/http';
import { HttpTestingController, provideHttpClientTesting } from '@angular/common/http/testing';
import { EventService } from './event.service';

describe('EventService', () => {
  let service: EventService;
  let http: HttpTestingController;

  beforeEach(() => {
    TestBed.configureTestingModule({
      providers: [EventService, provideHttpClient(), provideHttpClientTesting()]
    });
    service = TestBed.inject(EventService);
    http = TestBed.inject(HttpTestingController);
  });

  afterEach(() => http.verify());

  it('requests events from the backend contract', () => {
    service.list().subscribe();
    const request = http.expectOne('http://localhost:8080/api/events');
    expect(request.request.method).toBe('GET');
    request.flush([]);
  });
});
