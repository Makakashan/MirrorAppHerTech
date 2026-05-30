import { Routes } from '@angular/router';

export const routes: Routes = [
  {
    path: '',
    loadComponent: () =>
      import('./pages/dashboard/layout/dashboard-layout.component').then(m => m.DashboardLayoutComponent),
    children: [
      {
        path: '',
        loadComponent: () =>
          import('./pages/chat/chat-page.component').then(m => m.ChatPageComponent),
      },
      {
        path: 'analytics',
        loadComponent: () =>
          import('./pages/analytics/analytics-page.component').then(m => m.AnalyticsPageComponent),
      },
      {
        path: 'calendar',
        loadComponent: () =>
          import('./pages/calendar/calendar-page.component').then(m => m.CalendarPageComponent),
      },
      {
        path: 'settings',
        loadComponent: () =>
          import('./pages/analytics/analytics-page.component').then(m => m.AnalyticsPageComponent),
      },
    ],
  },
];
