import { Component, input } from '@angular/core';

@Component({
  selector: 'app-loading-state',
  standalone: true,
  template: `<div class="loading-state" role="status"><span class="loader"></span>{{ message() }}</div>`
})
export class LoadingStateComponent {
  readonly message = input('Cargando información…');
}
