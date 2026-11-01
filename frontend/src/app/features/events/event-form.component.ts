import { ChangeDetectionStrategy, Component, inject, signal } from '@angular/core';
import { CommonModule } from '@angular/common';
import { FormBuilder, ReactiveFormsModule, Validators } from '@angular/forms';
import { Router, RouterLink } from '@angular/router';
import { CourseService } from '../../core/services/course.service';
import { EventService } from '../../core/services/event.service';
import { ErrorMessageService } from '../../core/services/error-message.service';
import { Course } from '../../core/models/api.models';
import { LoadingStateComponent } from '../../shared/components/loading-state.component';

@Component({
  selector: 'app-event-form',
  standalone: true,
  imports: [CommonModule, ReactiveFormsModule, RouterLink, LoadingStateComponent],
  changeDetection: ChangeDetectionStrategy.OnPush,
  template: `
    <section class="page-section narrow-section">
      <a class="back-link" routerLink="/events">← Volver a eventos</a>
      <div class="form-intro"><span class="eyebrow">Nuevo evento</span><h1>Una invitación<br><em>bien cuidada.</em></h1><p>Al crear el evento generaremos una solicitud individual para cada alumno del curso.</p></div>
      <form class="event-form" [formGroup]="form" (ngSubmit)="submit()">
        <label>Nombre del evento <input formControlName="name" placeholder="Ej. Salida al museo" /></label>
        @if (form.controls.name.touched && form.controls.name.invalid) { <span class="field-error">Escribe un nombre para el evento.</span> }
        <label>Curso
          <select formControlName="courseId">
            <option [ngValue]="null">Selecciona un curso</option>
            @for (course of courses(); track course.id) { <option [ngValue]="course.id">{{ course.name }}</option> }
          </select>
        </label>
        <label>Fecha del evento <input type="date" formControlName="eventDate" /></label>
        <label>Descripción <textarea formControlName="description" rows="4" placeholder="Cuéntales a las familias lo esencial…"></textarea></label>
        @if (error()) { <div class="notice notice-error">{{ error() }}</div> }
        <div class="form-actions"><a class="button button-quiet" routerLink="/events">Cancelar</a><button class="button button-primary" [disabled]="saving() || form.invalid">{{ saving() ? 'Creando…' : 'Crear y generar solicitudes' }}</button></div>
      </form>
    </section>
  `
})
export class EventFormComponent {
  private readonly formBuilder = inject(FormBuilder);
  private readonly courseService = inject(CourseService);
  private readonly eventService = inject(EventService);
  private readonly errorMessage = inject(ErrorMessageService);
  private readonly router = inject(Router);
  readonly courses = signal<Course[]>([]);
  readonly saving = signal(false);
  readonly error = signal('');
  readonly form = this.formBuilder.nonNullable.group({
    name: ['', [Validators.required, Validators.maxLength(120)]],
    courseId: [0, [Validators.required, Validators.min(1)]],
    eventDate: ['', Validators.required],
    description: ['', Validators.maxLength(500)]
  });

  constructor() {
    this.courseService.list().subscribe({
      next: (courses) => this.courses.set(courses),
      error: (error: unknown) => this.error.set(this.errorMessage.getMessage(error, 'No pudimos cargar los cursos.'))
    });
  }

  submit(): void {
    if (this.form.invalid) { this.form.markAllAsTouched(); return; }
    this.saving.set(true);
    this.error.set('');
    this.eventService.create(this.form.getRawValue()).subscribe({
      next: (event) => this.router.navigate(['/events', event.id]),
      error: (error: unknown) => { this.error.set(this.errorMessage.getMessage(error)); this.saving.set(false); }
    });
  }
}
