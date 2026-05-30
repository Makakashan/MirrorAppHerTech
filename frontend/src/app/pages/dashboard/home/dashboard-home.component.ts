import { Component } from '@angular/core';
import { ChatComponent } from '../../../components/chat/chat.component';
import { MetricsPanelComponent } from '../../../components/metrics-panel/metrics-panel.component';

@Component({
  selector: 'app-dashboard-home',
  templateUrl: './dashboard-home.component.html',
  styleUrl: './dashboard-home.component.scss',
  standalone: true,
  imports: [ChatComponent, MetricsPanelComponent],
})
export class DashboardHomeComponent {}
