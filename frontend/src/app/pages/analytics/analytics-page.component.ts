import { Component, effect } from '@angular/core';
import { LucideDownload } from '@lucide/angular';
import { MoodLineChartComponent } from '../../components/mood-line-chart/mood-line-chart.component';
import { ChartPoint, Metrics } from '../../models/metrics.model';
import { MetricsService } from '../../services/metrics.service';

type AnalyticsPeriod = '7d' | '30d';

@Component({
  selector: 'app-analytics-page',
  standalone: true,
  imports: [MoodLineChartComponent, LucideDownload],
  templateUrl: './analytics-page.component.html',
  styleUrl: './analytics-page.component.scss',
})
export class AnalyticsPageComponent {
  period: AnalyticsPeriod = '7d';
  averageScore = 8.4;
  averageDelta = 1.2;
  streakDays = 14;
  stability = 72;
  chartData: ChartPoint[] = [
    { label: 'Mon', value: 72, x: 10 },
    { label: 'Tue', value: 85, x: 107 },
    { label: 'Wed', value: 78, x: 203 },
    { label: 'Thu', value: 91, x: 300 },
    { label: 'Fri', value: 88, x: 397 },
    { label: 'Sat', value: 76, x: 493 },
    { label: 'Sun', value: 84, x: 590 },
  ];

  insights = [
    'Mood improves after moments of movement or routine change.',
    'Avoidance appears most often around pressure and evaluation.',
    'Recent entries show recovery after short breaks.',
  ];

  readonly sources = [
    { label: 'Journal entries', value: 18 },
    { label: 'Check-ins', value: 12 },
    { label: 'Tagged triggers', value: 9 },
  ];

  constructor(private metricsService: MetricsService) {
    effect(() => {
      const metrics = this.metricsService.metrics();
      if (metrics) {
        this.applyMetrics(metrics);
      }
    });
    this.metricsService.load();
  }

  setPeriod(period: AnalyticsPeriod): void {
    this.period = period;
  }

  get bestDay(): ChartPoint {
    return this.chartData.reduce((best, item) => item.value > best.value ? item : best, this.chartData[0]);
  }

  get lowDay(): ChartPoint {
    return this.chartData.reduce((low, item) => item.value < low.value ? item : low, this.chartData[0]);
  }

  get trendRows(): Array<ChartPoint & { change: number }> {
    return this.chartData.map((item, index, items) => ({
      ...item,
      change: index === 0 ? 0 : item.value - items[index - 1].value,
    }));
  }

  get deltaLabel(): string {
    return this.averageDelta >= 0 ? `+${this.averageDelta}` : `${this.averageDelta}`;
  }

  private applyMetrics(metrics: Metrics): void {
    this.averageScore = metrics.averageScore;
    this.averageDelta = metrics.averageDelta;
    this.streakDays = metrics.streakDays;
    this.chartData = metrics.chartData;
    this.stability = this.calculateStability(metrics.chartData);
  }

  private calculateStability(data: ChartPoint[]): number {
    const values = data.map(item => item.value);
    const average = values.reduce((sum, value) => sum + value, 0) / values.length;
    const variance = values.reduce((sum, value) => sum + Math.abs(value - average), 0) / values.length;
    return Math.max(0, Math.round(100 - variance * 2));
  }
}
