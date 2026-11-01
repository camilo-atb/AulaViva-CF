import { Routes } from '@angular/router';
import { AdminLayoutComponent } from './layout/admin-layout/admin-layout.component';

export const appRoutes: Routes = [
  {
    path: '',
    component: AdminLayoutComponent,
    children: [
      { path: '', pathMatch: 'full', redirectTo: 'dashboard' },
      {
        path: 'dashboard',
        loadComponent: () => import('./features/dashboard/dashboard.component').then((m) => m.DashboardComponent)
      },
      {
        path: 'events/new',
        loadComponent: () => import('./features/events/event-form.component').then((m) => m.EventFormComponent)
      },
      {
        path: 'events/:id',
        loadComponent: () => import('./features/events/event-detail.component').then((m) => m.EventDetailComponent)
      },
      {
        path: 'events',
        loadComponent: () => import('./features/events/events-list.component').then((m) => m.EventsListComponent)
      }
    ]
  },
  {
    path: 'request/:token',
    loadComponent: () => import('./features/authorizations/guardian-portal.component').then((m) => m.GuardianPortalComponent)
  },
  { path: '**', redirectTo: 'dashboard' }
];
