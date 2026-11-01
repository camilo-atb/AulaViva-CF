import { Component, input } from '@angular/core';

@Component({
  selector: 'app-empty-state',
  standalone: true,
  template: `
    <div class="empty-state">
      <span class="empty-mark">—</span>
      <h3>{{ title() }}</h3>
      <p>{{ message() }}</p>
    </div>
  `
})
export class EmptyStateComponent {
  readonly title = input.required<string>();
  readonly message = input.required<string>();
}
