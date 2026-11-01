import { Component, input } from '@angular/core';
import { AuthorizationStatus } from '../../core/models/api.models';

@Component({
  selector: 'app-status-pill',
  standalone: true,
  template: `
    <span class="status-pill" [class]="status()">
      <span class="status-indicator"></span>{{ label }}
    </span>
  `
})
export class StatusPillComponent {
  readonly status = input.required<AuthorizationStatus>();

  get label(): string {
    return {
      pending: 'Pendiente',
      authorized: 'Autorizada',
      not_authorized: 'No autorizada'
    }[this.status()];
  }
}
