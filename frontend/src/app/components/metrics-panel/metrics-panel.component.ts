import { Component, effect } from '@angular/core';
import { ChartPoint, DayMood, Metrics } from '../../models/metrics.model';
import { MetricsService } from '../../services/metrics.service';

@Component({
  selector: 'app-metrics-panel',
  templateUrl: './metrics-panel.component.html',
  styleUrl: './metrics-panel.component.scss',
  standalone: true,
  imports: [],
})
export class MetricsPanelComponent {
  readonly today = new Date();
  currentMonth = this.today.toLocaleDateString('en-US', { month: 'long' });
  averageScore = 8.4;
  averageDelta = 1.2;
  streakDays = 14;

  chartData: ChartPoint[] = [
    { label: 'Mon', value: 72, x: 10 },
    { label: 'Tue', value: 85, x: 107 },
    { label: 'Wed', value: 78, x: 203 },
    { label: 'Thu', value: 91, x: 300 },
    { label: 'Fri', value: 88, x: 397 },
    { label: 'Sat', value: 76, x: 493 },
    { label: 'Sun', value: 84, x: 590 },
  ];

  chartPolyline = this.buildPolyline();

  chartArea = `10,130 ${this.chartPolyline} 590,130`;

  calendarDays: DayMood[] = this.buildCalendar();

  constructor(private metricsService: MetricsService) {
    effect(() => {
      const metrics = this.metricsService.metrics();
      if (metrics) {
        this.applyMetrics(metrics);
      }
    });
    this.metricsService.load();
  }

  private toY(v: number): number {
    return Math.round(130 - ((v - 50) / 50) * 120);
  }

  private buildPolyline(): string {
    return this.chartData
      .map(d => `${d.x},${this.toY(d.value)}`)
      .join(' ');
  }

  private applyMetrics(metrics: Metrics): void {
    this.averageScore = metrics.averageScore;
    this.averageDelta = metrics.averageDelta;
    this.streakDays = metrics.streakDays;
    this.currentMonth = metrics.currentMonth;
    this.chartData = metrics.chartData;
    this.calendarDays = metrics.calendarDays;
    this.chartPolyline = this.buildPolyline();
    this.chartArea = `10,130 ${this.chartPolyline} 590,130`;
  }

  private buildCalendar(): DayMood[] {
    const today = new Date();
    const dow = today.getDay() || 7;
    const monday = new Date(today);
    monday.setDate(today.getDate() - dow + 1 - 21);

    const moods: DayMood['mood'][] = ['great', 'good', 'neutral', 'bad', 'great'];
    return Array.from({ length: 28 }, (_, i) => {
      const d = new Date(monday);
      d.setDate(monday.getDate() + i);
      return { day: d.getDate(), mood: moods[d.getDate() % 5], isFuture: d > today };
    });
  }
}
