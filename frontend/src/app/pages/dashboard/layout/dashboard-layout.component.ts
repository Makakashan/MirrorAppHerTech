import { Component } from '@angular/core';
import { RouterOutlet, RouterLink, RouterLinkActive } from '@angular/router';
import {
  LucideCalendarDays,
  LucideChartColumn,
  LucideHouse,
  LucideSettings,
  LucideUser,
} from '@lucide/angular';
import { ThemeService } from '../../../services/theme.service';

@Component({
  selector: 'app-dashboard-layout',
  templateUrl: './dashboard-layout.component.html',
  styleUrl: './dashboard-layout.component.scss',
  standalone: true,
  imports: [
    RouterOutlet,
    RouterLink,
    RouterLinkActive,
    LucideCalendarDays,
    LucideChartColumn,
    LucideHouse,
    LucideSettings,
    LucideUser,
  ],
})
export class DashboardLayoutComponent {
  constructor(readonly themeService: ThemeService) {}
}
