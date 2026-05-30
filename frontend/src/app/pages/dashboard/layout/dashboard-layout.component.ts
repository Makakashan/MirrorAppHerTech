import { Component } from '@angular/core';
import { RouterOutlet, RouterLink, RouterLinkActive } from '@angular/router';

@Component({
  selector: 'app-dashboard-layout',
  templateUrl: './dashboard-layout.component.html',
  styleUrl: './dashboard-layout.component.scss',
  standalone: true,
  imports: [RouterOutlet, RouterLink, RouterLinkActive],
})
export class DashboardLayoutComponent {}
