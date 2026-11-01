import { ChangeDetectionStrategy, Component, inject, signal } from '@angular/core';
import { RouterLink } from '@angular/router';
import { StudentDetail } from '../../core/models/api.models';
import { StudentService } from '../../core/services/student.service';
import { ErrorMessageService } from '../../core/services/error-message.service';
import { EmptyStateComponent } from '../../shared/components/empty-state.component';
import { LoadingStateComponent } from '../../shared/components/loading-state.component';

@Component({
  selector: 'app-student-list',
  standalone: true,
  imports: [RouterLink, EmptyStateComponent, LoadingStateComponent],
  changeDetection: ChangeDetectionStrategy.OnPush,
  template: `
    <section class="page-section">
      <div class="page-heading compact">
        <div><span class="eyebrow">Personas autorizadas</span><h1>La salida,<br><em>con claridad.</em></h1></div>
      </div>
      <div class="student-search">
        <label for="student-search">Buscar estudiante</label>
        <div class="search-row">
          <input id="student-search" [value]="query()" (input)="query.set(inputValue($event))" (keyup.enter)="search()" placeholder="Escribe un nombre" />
          <button class="button button-primary" (click)="search()">Buscar</button>
        </div>
      </div>
      @if (loading()) { <app-loading-state /> }
      @else if (error()) { <div class="notice notice-error">{{ error() }} <button class="text-link inline-action" (click)="search()">Reintentar</button></div> }
      @else if (!students().length) { <app-empty-state title="No encontramos estudiantes" message="Prueba con otro nombre o revisa el listado completo." /> }
      @else {
        <div class="student-list">
          @for (student of students(); track student.id) {
            <a class="student-row" [routerLink]="['/students', student.id]">
              <span class="student-avatar">{{ initials(student.fullName) }}</span>
              <span class="student-copy"><strong>{{ student.fullName }}</strong><small>Curso {{ student.courseId }}</small></span>
              <span class="student-status" [class.has-people]="student.authorizedPeople.length > 0">
                {{ student.authorizedPeople.length ? student.authorizedPeople.length + ' persona(s)' : 'Sin personas registradas' }}
              </span>
              <span class="event-arrow">↗</span>
            </a>
          }
        </div>
      }
    </section>
  `
})
export class StudentListComponent {
  private readonly service = inject(StudentService);
  private readonly errorMessage = inject(ErrorMessageService);
  readonly students = signal<StudentDetail[]>([]);
  readonly query = signal('');
  readonly loading = signal(true);
  readonly error = signal('');

  constructor() { this.search(); }

  inputValue(event: Event): string { return (event.target as HTMLInputElement).value; }

  search(): void {
    this.loading.set(true);
    this.error.set('');
    this.service.list(this.query()).subscribe({
      next: (students) => { this.students.set(students); this.loading.set(false); },
      error: (error: unknown) => { this.error.set(this.errorMessage.getMessage(error, 'No pudimos cargar los estudiantes.')); this.loading.set(false); }
    });
  }

  initials(name: string): string {
    return name.split(' ').slice(0, 2).map((part) => part[0]).join('').toUpperCase();
  }
}
