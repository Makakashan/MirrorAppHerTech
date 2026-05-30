import { Component } from '@angular/core';
import { MetricsPanelComponent } from '../../components/metrics-panel/metrics-panel.component';

@Component({
  selector: 'app-analytics-page',
  standalone: true,
  imports: [MetricsPanelComponent],
  templateUrl: './analytics-page.component.html',
  styleUrl: './analytics-page.component.scss',
})
export class AnalyticsPageComponent {}
