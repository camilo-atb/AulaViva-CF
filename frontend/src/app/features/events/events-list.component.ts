import { ChangeDetectionStrategy, Component, inject, signal } from '@angular/core';
import { CommonModule } from '@angular/common';
import { RouterLink } from '@angular/router';
import { EventResponse } from '../../core/models/api.models';
import { EventService } from '../../core/services/event.service';
import { ErrorMessageService } from '../../core/services/error-message.service';
import { LoadingStateComponent } from '../../shared/components/loading-state.component';
import { EmptyStateComponent } from '../../shared/components/empty-state.component';
import { formatDate } from '../../shared/utils/date.utils';

@Component({
  selector: 'app-events-list',
  standalone: true,
  imports: [CommonModule, RouterLink, LoadingStateComponent, EmptyStateComponent],
  changeDetection: ChangeDetectionStrategy.OnPush,
  template: `
    <section class="page-section">
      <div class="page-heading compact">
        <div><span class="eyebrow">Gestión de eventos</span><h1>La agenda<br><em>que nos reúne.</em></h1></div>
        <a class="button button-primary" routerLink="/events/new"><span>+</span> Nuevo evento</a>
      </div>
      @if (loading()) { <app-loading-state /> }
      @else if (error()) { <div class="notice notice-error">{{ error() }}</div> }
      @else if (!events().length) { <app-empty-state title="Aún no hay eventos" message="Crea el primer evento de tu escuela." /> }
      @else {
        <div class="event-table">
          <div class="table-head"><span>Evento</span><span>Fecha</span><span>Respuestas</span><span></span></div>
          @for (event of events(); track event.id) {
            <a class="table-row" [routerLink]="['/events', event.id]">
              <span><strong>{{ event.name }}</strong><small>{{ event.courseName || 'Curso seleccionado' }}</small></span>
              <span>{{ formatEventDate(event.eventDate) }}</span>
              <span class="response-count"><b>{{ event.summary?.authorized ?? 0 }}</b> / {{ total(event) }} respondidas</span>
              <span class="event-arrow">↗</span>
            </a>
          }
        </div>
      }
    </section>
  `
})
export class EventsListComponent {
  private readonly eventService = inject(EventService);
  private readonly errorMessage = inject(ErrorMessageService);
  readonly events = signal<EventResponse[]>([]);
  readonly loading = signal(true);
  readonly error = signal('');

  constructor() {
    this.eventService.list().subscribe({
      next: (events) => {
        if (!events.length) { this.events.set([]); this.loading.set(false); return; }
        Promise.all(events.map((event) => this.eventService.get(event.id).toPromise())).then((details) => {
          this.events.set(details.filter((event): event is EventResponse => !!event));
          this.loading.set(false);
        }).catch((error: unknown) => { this.error.set(this.errorMessage.getMessage(error)); this.loading.set(false); });
      },
      error: (error: unknown) => { this.error.set(this.errorMessage.getMessage(error)); this.loading.set(false); }
    });
  }

  total(event: EventResponse): number {
    const summary = event.summary;
    return (summary?.pending ?? 0) + (summary?.authorized ?? 0) + (summary?.notAuthorized ?? 0);
  }

  formatEventDate(date: string): string { return formatDate(date); }
}
