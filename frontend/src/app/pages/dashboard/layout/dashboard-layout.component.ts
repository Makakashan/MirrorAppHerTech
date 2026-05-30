import { Component } from '@angular/core';
import { RouterOutlet, RouterLink, RouterLinkActive } from '@angular/router';
import {
  LucideBell,
  LucideCalendarDays,
  LucideChartColumn,
  LucideHouse,
  LucideSearch,
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
    LucideBell,
    LucideCalendarDays,
    LucideChartColumn,
    LucideHouse,
    LucideSearch,
    LucideSettings,
    LucideUser,
  ],
})
export class DashboardLayoutComponent {
  constructor(readonly themeService: ThemeService) {}
}
