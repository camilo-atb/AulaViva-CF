import { ChangeDetectionStrategy, Component, computed, inject, signal } from '@angular/core';
import { CommonModule } from '@angular/common';
import { RouterLink } from '@angular/router';
import { forkJoin } from 'rxjs';
import { ErrorMessageService } from '../../core/services/error-message.service';
import { EventService } from '../../core/services/event.service';
import { EventResponse } from '../../core/models/api.models';
import { LoadingStateComponent } from '../../shared/components/loading-state.component';
import { EmptyStateComponent } from '../../shared/components/empty-state.component';
import { formatDate } from '../../shared/utils/date.utils';

@Component({
  selector: 'app-dashboard',
  standalone: true,
  imports: [CommonModule, RouterLink, LoadingStateComponent, EmptyStateComponent],
  changeDetection: ChangeDetectionStrategy.OnPush,
  template: `
    <section class="page-section">
      <div class="page-heading">
        <div>
          <span class="eyebrow">Miércoles, 01 de noviembre de 2026</span>
          <h1>Todo en orden,<br><em>con calma.</em></h1>
          <p class="lede">Una mirada serena al estado de las autorizaciones de tu escuela.</p>
        </div>
        <a class="button button-primary" routerLink="/events/new"><span>+</span> Crear evento</a>
      </div>

      @if (loading()) { <app-loading-state /> }
      @else if (error()) { <div class="notice notice-error">{{ error() }}</div> }
      @else {
        <div class="dashboard-grid">
          <section class="welcome-panel">
            <div class="panel-kicker">Panorama actual</div>
            <h2>{{ totalAuthorizations() }} <small>solicitudes</small></h2>
            <p>El pulso de tus eventos escolares, actualizado al abrir esta vista.</p>
            <div class="mini-stats">
              <span><b>{{ pending() }}</b> por responder</span>
              <span><b>{{ authorized() }}</b> autorizadas</span>
            </div>
          </section>
          <div class="stat-block"><span class="stat-label">Autorizadas</span><strong>{{ authorized() }}</strong><span class="stat-note">respuestas positivas</span></div>
          <div class="stat-block accent"><span class="stat-label">Pendientes</span><strong>{{ pending() }}</strong><span class="stat-note">requieren atención</span></div>
          <div class="stat-block"><span class="stat-label">No autorizadas</span><strong>{{ notAuthorized() }}</strong><span class="stat-note">respuestas negativas</span></div>
        </div>
        <div class="section-header"><div><span class="eyebrow">Seguimiento</span><h2>Eventos recientes</h2></div><a class="text-link" routerLink="/events">Ver todos <span>→</span></a></div>
        @if (events().length) {
          <div class="event-list">
            @for (event of events().slice(0, 5); track event.id) {
              <a class="event-row" [routerLink]="['/events', event.id]">
                <span class="event-date"><b>{{ event.eventDate | slice: 8:10 }}</b><small>{{ event.eventDate | slice: 5:7 }}</small></span>
                <span class="event-copy"><strong>{{ event.name }}</strong><small>{{ event.courseName || 'Curso seleccionado' }} · {{ formatEventDate(event.eventDate) }}</small></span>
                <span class="event-arrow">↗</span>
              </a>
            }
          </div>
        } @else {
          <app-empty-state title="Aún no hay eventos" message="Crea el primer evento para comenzar a recibir respuestas." />
        }
      }
    </section>
  `
})
export class DashboardComponent {
  private readonly eventService = inject(EventService);
  private readonly errorMessage = inject(ErrorMessageService);
  readonly events = signal<EventResponse[]>([]);
  readonly loading = signal(true);
  readonly error = signal('');
  readonly totalAuthorizations = computed(() => this.events().reduce((sum, event) => sum + (event.summary?.pending ?? 0) + (event.summary?.authorized ?? 0) + (event.summary?.notAuthorized ?? 0), 0));
  readonly pending = computed(() => this.events().reduce((sum, event) => sum + (event.summary?.pending ?? 0), 0));
  readonly authorized = computed(() => this.events().reduce((sum, event) => sum + (event.summary?.authorized ?? 0), 0));
  readonly notAuthorized = computed(() => this.events().reduce((sum, event) => sum + (event.summary?.notAuthorized ?? 0), 0));

  constructor() {
    this.eventService.list().subscribe({
      next: (events) => {
        if (!events.length) {
          this.events.set([]);
          this.loading.set(false);
          return;
        }
        forkJoin(events.map((event) => this.eventService.get(event.id))).subscribe({
          next: (details) => { this.events.set(details); this.loading.set(false); },
          error: (error: unknown) => { this.error.set(this.errorMessage.getMessage(error)); this.loading.set(false); }
        });
      },
      error: (error: unknown) => { this.error.set(this.errorMessage.getMessage(error)); this.loading.set(false); }
    });
  }

  formatEventDate(date: string): string { return formatDate(date); }
}
