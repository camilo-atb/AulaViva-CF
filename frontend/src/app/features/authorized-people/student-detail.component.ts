import { ChangeDetectionStrategy, Component, inject, signal } from '@angular/core';
import { ActivatedRoute, RouterLink } from '@angular/router';
import { AuthorizedPerson, AuthorizedPersonInput, StudentDetail } from '../../core/models/api.models';
import { StudentService } from '../../core/services/student.service';
import { ErrorMessageService } from '../../core/services/error-message.service';
import { LoadingStateComponent } from '../../shared/components/loading-state.component';
import { AuthorizedPersonFormComponent } from './authorized-person-form.component';

@Component({
  selector: 'app-student-detail',
  standalone: true,
  imports: [RouterLink, LoadingStateComponent, AuthorizedPersonFormComponent],
  changeDetection: ChangeDetectionStrategy.OnPush,
  template: `
    <section class="page-section narrow-section">
      <a class="back-link" routerLink="/students">← Volver a estudiantes</a>
      @if (loading()) { <app-loading-state /> }
      @else if (error() && !student()) { <div class="notice notice-error">{{ error() }}</div> }
      @else if (student()) {
        <div class="student-detail-intro">
          <span class="student-avatar large">{{ initials(currentStudent.fullName) }}</span>
          <div><span class="eyebrow">Ficha del estudiante</span><h1>{{ currentStudent.fullName }}</h1><p>Curso {{ currentStudent.courseId }} · Registro operativo de recogida</p></div>
        </div>
        @if (error()) { <div class="notice notice-error">{{ error() }}</div> }
        @if (success()) { <div class="notice notice-warm">{{ success() }}</div> }
        <div class="section-header"><div><span class="eyebrow">Personas registradas</span><h2>¿Quién puede recogerle?</h2></div><button class="button button-primary" (click)="startCreate()">+ Añadir persona</button></div>
        @if (!currentStudent.authorizedPeople.length && !editing()) { <div class="empty-state person-empty"><div class="empty-mark">—</div><h3>Sin personas autorizadas</h3><p>Añade el primer registro para tener esta información a mano.</p><button class="button button-primary" (click)="startCreate()">Registrar persona</button></div> }
        @if (editing()) {
          <div class="person-form-panel"><span class="eyebrow">{{ editingPerson() ? 'Editar registro' : 'Nuevo registro' }}</span><h3>{{ editingPerson() ? 'Actualiza los datos' : 'Añade una persona de confianza' }}</h3><app-authorized-person-form [person]="editingPerson()" [saving]="saving()" (saved)="save($event)" (cancel)="cancelEdit()" /></div>
        }
        @if (currentStudent.authorizedPeople.length) {
          <div class="people-list">
            @for (person of currentStudent.authorizedPeople; track person.id) {
              <article class="person-row">
                <div class="person-badge">{{ initials(person.fullName) }}</div>
                <div class="person-copy"><strong>{{ person.fullName }}</strong><span>{{ person.relationship }} · {{ person.documentType }} {{ person.documentNumber }}</span><small>{{ person.phone }}</small></div>
                <div class="person-actions"><button class="copy-button" (click)="startEdit(person)">Editar</button><button class="copy-button remove-action" (click)="remove(person)">Eliminar</button></div>
              </article>
            }
          </div>
        }
      }
    </section>
  `
})
export class StudentDetailComponent {
  private readonly service = inject(StudentService);
  private readonly errorMessage = inject(ErrorMessageService);
  private readonly route = inject(ActivatedRoute);
  readonly student = signal<StudentDetail | null>(null);
  readonly loading = signal(true);
  readonly saving = signal(false);
  readonly error = signal('');
  readonly success = signal('');
  readonly editing = signal(false);
  readonly editingPerson = signal<AuthorizedPerson | null>(null);
  private readonly studentId = Number(this.route.snapshot.paramMap.get('id'));

  get currentStudent(): StudentDetail {
    return this.student() as StudentDetail;
  }

  constructor() { this.load(); }

  load(): void {
    this.service.get(this.studentId).subscribe({
      next: (student) => { this.student.set(student); this.loading.set(false); },
      error: (error: unknown) => { this.error.set(this.errorMessage.getMessage(error, 'No pudimos cargar la ficha.')); this.loading.set(false); }
    });
  }

  startCreate(): void { this.error.set(''); this.success.set(''); this.editingPerson.set(null); this.editing.set(true); }
  startEdit(person: AuthorizedPerson): void { this.error.set(''); this.success.set(''); this.editingPerson.set(person); this.editing.set(true); }
  cancelEdit(): void { this.editing.set(false); this.editingPerson.set(null); }

  save(input: AuthorizedPersonInput): void {
    this.saving.set(true);
    const person = this.editingPerson();
    const request = person ? this.service.updateAuthorizedPerson(person.id, input) : this.service.createAuthorizedPerson(this.studentId, input);
    request.subscribe({
      next: () => { this.saving.set(false); this.editing.set(false); this.editingPerson.set(null); this.success.set(person ? 'Los datos fueron actualizados.' : 'La persona fue registrada.'); this.load(); },
      error: (error: unknown) => { this.saving.set(false); this.error.set(this.errorMessage.getMessage(error, 'No pudimos guardar este registro.')); }
    });
  }

  remove(person: AuthorizedPerson): void {
    if (!window.confirm(`¿Eliminar a ${person.fullName} del registro?`)) return;
    this.service.deleteAuthorizedPerson(person.id).subscribe({
      next: () => { this.success.set('El registro fue eliminado.'); this.load(); },
      error: (error: unknown) => this.error.set(this.errorMessage.getMessage(error, 'No pudimos eliminar este registro.'))
    });
  }

  initials(name: string): string { return name.split(' ').slice(0, 2).map((part) => part[0]).join('').toUpperCase(); }
}
