import { ComponentFixture, TestBed } from '@angular/core/testing';
import { provideHttpClient } from '@angular/common/http';
import { provideHttpClientTesting } from '@angular/common/http/testing';
import { ActivatedRoute } from '@angular/router';
import { of } from 'rxjs';
import { GuardianPortalComponent } from './guardian-portal.component';
import { AuthorizationService } from '../../core/services/authorization.service';

describe('GuardianPortalComponent', () => {
  let fixture: ComponentFixture<GuardianPortalComponent>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      imports: [GuardianPortalComponent],
      providers: [
        provideHttpClient(),
        provideHttpClientTesting(),
        { provide: ActivatedRoute, useValue: { snapshot: { paramMap: { get: () => 'token-demo' } } } },
        {
          provide: AuthorizationService,
          useValue: {
            get: () => of({
              authorization: { id: 1, eventId: 1, studentId: 1, requestToken: 'token-demo', status: 'pending' },
              event: { id: 1, courseId: 1, name: 'Salida escolar', description: 'Visita', eventDate: '2026-11-20' },
              student: { id: 1, courseId: 1, fullName: 'Sofía Martínez' }
            })
          }
        }
      ]
    }).compileComponents();
    fixture = TestBed.createComponent(GuardianPortalComponent);
    fixture.detectChanges();
  });

  it('shows the event and student from the API response', () => {
    expect(fixture.nativeElement.textContent).toContain('Salida escolar');
    expect(fixture.nativeElement.textContent).toContain('Sofía');
  });
});
