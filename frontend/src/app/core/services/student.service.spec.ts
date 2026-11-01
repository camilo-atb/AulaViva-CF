import { TestBed } from '@angular/core/testing';
import { provideHttpClient } from '@angular/common/http';
import { HttpTestingController, provideHttpClientTesting } from '@angular/common/http/testing';
import { StudentService } from './student.service';

describe('StudentService', () => {
  let service: StudentService;
  let http: HttpTestingController;

  beforeEach(() => {
    TestBed.configureTestingModule({
      providers: [StudentService, provideHttpClient(), provideHttpClientTesting()]
    });
    service = TestBed.inject(StudentService);
    http = TestBed.inject(HttpTestingController);
  });

  afterEach(() => http.verify());

  it('searches students using the API contract', () => {
    service.list('Sofía').subscribe();
    const request = http.expectOne('http://localhost:8080/api/students?search=Sof%C3%ADa');
    expect(request.request.method).toBe('GET');
    request.flush([]);
  });

  it('creates an authorized person for a student', () => {
    const input = {
      fullName: 'Laura Gómez',
      documentType: 'CC' as const,
      documentNumber: '123456789',
      relationship: 'Tía',
      phone: '+57 300 555 0101'
    };
    service.createAuthorizedPerson(1, input).subscribe();
    const request = http.expectOne('http://localhost:8080/api/students/1/authorized-people');
    expect(request.request.method).toBe('POST');
    expect(request.request.body).toEqual(input);
    request.flush({ id: 1, studentId: 1, ...input });
  });
});
