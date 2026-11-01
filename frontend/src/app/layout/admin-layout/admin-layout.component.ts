import { Component } from '@angular/core';
import { RouterLink, RouterLinkActive, RouterOutlet } from '@angular/router';

@Component({
  selector: 'app-admin-layout',
  standalone: true,
  imports: [RouterOutlet, RouterLink, RouterLinkActive],
  template: `
    <div class="admin-shell">
      <aside class="sidebar">
        <a class="brand" routerLink="/dashboard" aria-label="Ir al inicio de AulaViva">
          <span class="brand-mark">av</span>
          <span><strong>Aula</strong>Viva</span>
        </a>
        <div class="sidebar-intro">
          <span class="eyebrow">Escuela Primaria Sofía</span>
          <p>Un lugar claro para cada permiso.</p>
        </div>
        <nav aria-label="Navegación principal">
          <a routerLink="/dashboard" routerLinkActive="active" [routerLinkActiveOptions]="{ exact: true }">
            <span class="nav-glyph">01</span> Resumen
          </a>
          <a routerLink="/events" routerLinkActive="active">
            <span class="nav-glyph">02</span> Eventos
          </a>
          <a routerLink="/students" routerLinkActive="active">
            <span class="nav-glyph">03</span> Recogida
          </a>
        </nav>
        <div class="sidebar-footer">
          <span class="status-dot"></span>
          <span>Entorno local</span>
        </div>
      </aside>
      <main class="main-content">
        <header class="mobile-header">
          <a class="brand" routerLink="/dashboard"><span class="brand-mark">av</span><span><strong>Aula</strong>Viva</span></a>
          <a class="mobile-action" routerLink="/students">Recogida</a>
        </header>
        <router-outlet />
      </main>
    </div>
  `
})
export class AdminLayoutComponent {}
