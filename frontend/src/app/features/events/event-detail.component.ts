import { ChangeDetectionStrategy, Component, inject, signal } from '@angular/core';
import { CommonModule } from '@angular/common';
import { ActivatedRoute, RouterLink } from '@angular/router';
import { EventResponse } from '../../core/models/api.models';
import { EventService } from '../../core/services/event.service';
import { ErrorMessageService } from '../../core/services/error-message.service';
import { StatusPillComponent } from '../../shared/components/status-pill.component';
import { LoadingStateComponent } from '../../shared/components/loading-state.component';
import { formatDate } from '../../shared/utils/date.utils';

@Component({
  selector: 'app-event-detail',
  standalone: true,
  imports: [CommonModule, RouterLink, StatusPillComponent, LoadingStateComponent],
  changeDetection: ChangeDetectionStrategy.OnPush,
  template: `
    <section class="page-section">
      <a class="back-link" routerLink="/events">← Volver a eventos</a>
      @if (loading()) { <app-loading-state /> }
      @if (!loading() && error()) { <div class="notice notice-error">{{ error() }}</div> }
      @if (!loading() && !error() && event()) {
        <div class="detail-heading"><div><span class="eyebrow">{{ currentEvent.courseName || 'Evento escolar' }} · {{ formatEventDate(currentEvent.eventDate) }}</span><h1>{{ currentEvent.name }}</h1><p>{{ currentEvent.description || 'Una actividad para compartir en comunidad.' }}</p></div><span class="detail-stamp">AulaViva<br><small>autorizaciones</small></span></div>
        <div class="summary-strip">
          <div><span>Pendientes</span><strong>{{ currentEvent.summary?.pending ?? 0 }}</strong></div>
          <div><span>Autorizadas</span><strong>{{ currentEvent.summary?.authorized ?? 0 }}</strong></div>
          <div><span>No autorizadas</span><strong>{{ currentEvent.summary?.notAuthorized ?? 0 }}</strong></div>
        </div>
        <div class="section-header"><div><span class="eyebrow">Seguimiento individual</span><h2>Respuestas de las familias</h2></div></div>
        <div class="authorization-table">
          <div class="table-head"><span>Alumno</span><span>Estado</span><span>Enlace individual</span></div>
          @for (item of currentEvent.authorizations ?? []; track item.id) {
            <div class="table-row auth-row"><span><strong>{{ item.studentName || 'Alumno' }}</strong><small>Solicitud #{{ item.id }}</small></span><app-status-pill [status]="item.status" /><span class="auth-actions"><button class="copy-button" (click)="copyLink(item.requestToken)">{{ copiedToken() === item.requestToken ? 'Copiado' : 'Copiar enlace' }}</button><button class="copy-button whatsapp-button" (click)="shareOnWhatsApp(item.requestToken)">WhatsApp</button></span></div>
          }
        </div>
      }
    </section>
  `
})
export class EventDetailComponent {
  private readonly eventService = inject(EventService);
  private readonly errorMessage = inject(ErrorMessageService);
  private readonly route = inject(ActivatedRoute);
  readonly event = signal<EventResponse | null>(null);
  readonly loading = signal(true);
  readonly error = signal('');
  readonly copiedToken = signal('');

  get currentEvent(): EventResponse {
    return this.event() as EventResponse;
  }

  constructor() {
    const id = Number(this.route.snapshot.paramMap.get('id'));
    this.eventService.get(id).subscribe({
      next: (event) => { this.event.set(event); this.loading.set(false); },
      error: (error: unknown) => { this.error.set(this.errorMessage.getMessage(error, 'No pudimos cargar este evento.')); this.loading.set(false); }
    });
  }

  formatEventDate(date: string): string { return formatDate(date); }

  copyLink(token: string): void {
    const link = `${window.location.origin}/request/${token}`;
    navigator.clipboard.writeText(link).then(() => {
      this.copiedToken.set(token);
      window.setTimeout(() => this.copiedToken.set(''), 1800);
    }).catch(() => this.error.set('No pudimos copiar el enlace. Cópialo desde la barra del navegador.'));
  }

  shareOnWhatsApp(token: string): void {
    const link = `${window.location.origin}/request/${token}`;
    window.open(`https://wa.me/?text=${encodeURIComponent(`AulaViva: revisa esta autorización escolar ${link}`)}`, '_blank', 'noopener');
  }
}
