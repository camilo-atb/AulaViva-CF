import { ChangeDetectionStrategy, Component, inject, signal } from '@angular/core';
import { CommonModule } from '@angular/common';
import { ActivatedRoute } from '@angular/router';
import { AuthorizationService } from '../../core/services/authorization.service';
import { ErrorMessageService } from '../../core/services/error-message.service';
import { RequestResponse } from '../../core/models/api.models';
import { formatDate } from '../../shared/utils/date.utils';

@Component({
  selector: 'app-guardian-portal',
  standalone: true,
  imports: [CommonModule],
  changeDetection: ChangeDetectionStrategy.OnPush,
  template: `
    <main class="guardian-page">
      <div class="guardian-brand"><span class="brand-mark">av</span><span><strong>Aula</strong>Viva</span></div>
      @if (loading()) { <div class="guardian-card"><div class="guardian-loader"></div><p>Estamos preparando la invitación…</p></div> }
      @if (!loading() && error()) { <div class="guardian-card guardian-error"><span class="guardian-symbol">!</span><span class="eyebrow">Enlace no disponible</span><h1>No encontramos<br><em>esta solicitud.</em></h1><p>{{ error() }}</p></div> }
      @if (!loading() && !error() && confirmed()) { <div class="guardian-card confirmed"><span class="confirmation-mark">✓</span><span class="eyebrow">Respuesta registrada</span><h1>Gracias por<br><em>responder.</em></h1><p>La escuela ya recibió tu respuesta para <strong>{{ request()?.student?.fullName }}</strong>.</p><div class="confirmation-note">Puedes cerrar esta ventana con tranquilidad.</div></div> }
      @if (!loading() && !error() && !confirmed() && request()) {
        <div class="guardian-card">
          <span class="eyebrow">Una invitación de la Escuela Primaria Sofía</span>
          <h1>Hola, familia<br><em>{{ firstName(currentRequest.student.fullName) }}.</em></h1>
          <div class="guardian-divider"></div>
          <p class="guardian-question">Queremos contar contigo para una actividad especial.</p>
          <div class="guardian-event"><span class="event-label">EVENTO ESCOLAR</span><h2>{{ currentRequest.event.name }}</h2><p>{{ currentRequest.event.description || 'Pronto compartiremos más detalles de esta actividad.' }}</p><dl><div><dt>Fecha</dt><dd>{{ formatEventDate(currentRequest.event.eventDate) }}</dd></div><div><dt>Alumno</dt><dd>{{ currentRequest.student.fullName }}</dd></div></dl></div>
          @if (currentRequest.authorization.status === 'pending') {
            <p class="guardian-prompt">Por favor, indícanos si autorizas su participación.</p>
            <div class="guardian-actions"><button class="guardian-button approve" [disabled]="sending()" (click)="respond('authorized')">Sí, autorizo <span>→</span></button><button class="guardian-button decline" [disabled]="sending()" (click)="respond('not_authorized')">No autorizo</button></div>
          } @else {
            <div class="notice notice-warm">Esta solicitud ya fue respondida como <strong>{{ currentRequest.authorization.status === 'authorized' ? 'autorizada' : 'no autorizada' }}</strong>.</div>
          }
          @if (error()) { <div class="notice notice-error">{{ error() }}</div> }
          <p class="guardian-footnote">Tus datos se usan únicamente para gestionar esta autorización escolar.</p>
        </div>
      }
    </main>
  `
})
export class GuardianPortalComponent {
  private readonly route = inject(ActivatedRoute);
  private readonly authorizationService = inject(AuthorizationService);
  private readonly errorMessage = inject(ErrorMessageService);
  readonly request = signal<RequestResponse | null>(null);
  readonly loading = signal(true);
  readonly sending = signal(false);
  readonly confirmed = signal(false);
  readonly error = signal('');
  private readonly token = this.route.snapshot.paramMap.get('token') ?? '';

  get currentRequest(): RequestResponse {
    return this.request() as RequestResponse;
  }

  constructor() {
    this.authorizationService.get(this.token).subscribe({
      next: (request) => { this.request.set(request); this.loading.set(false); },
      error: (error: unknown) => { this.error.set(this.errorMessage.getMessage(error, 'El enlace puede estar vencido o ser incorrecto.')); this.loading.set(false); }
    });
  }

  respond(status: 'authorized' | 'not_authorized'): void {
    this.sending.set(true);
    this.error.set('');
    this.authorizationService.respond(this.token, status).subscribe({
      next: (request) => { this.request.set(request); this.confirmed.set(true); this.sending.set(false); },
      error: (error: unknown) => { this.error.set(this.errorMessage.getMessage(error, 'Esta solicitud ya fue respondida.')); this.sending.set(false); }
    });
  }

  firstName(name: string): string { return name.split(' ')[0]; }
  formatEventDate(date: string): string { return formatDate(date); }
}
