import { Component, Input } from '@angular/core';
import { ChartPoint } from '../../models/metrics.model';

interface PlotPoint extends ChartPoint {
  chartX: number;
  chartY: number;
}

@Component({
  selector: 'app-mood-line-chart',
  standalone: true,
  imports: [],
  templateUrl: './mood-line-chart.component.html',
  styleUrl: './mood-line-chart.component.scss',
})
export class MoodLineChartComponent {
  @Input({ required: true }) data: ChartPoint[] = [];

  readonly width = 720;
  readonly height = 320;
  readonly plot = {
    top: 28,
    right: 22,
    bottom: 48,
    left: 44,
  };

  get points(): PlotPoint[] {
    const domain = this.domain;
    const plotWidth = this.width - this.plot.left - this.plot.right;
    const plotHeight = this.height - this.plot.top - this.plot.bottom;
    const denominator = Math.max(this.data.length - 1, 1);

    return this.data.map((item, index) => {
      const progress = index / denominator;
      const valueProgress = (item.value - domain.min) / (domain.max - domain.min);

      return {
        ...item,
        chartX: Math.round(this.plot.left + progress * plotWidth),
        chartY: Math.round(this.plot.top + (1 - valueProgress) * plotHeight),
      };
    });
  }

  get linePoints(): string {
    return this.points.map(point => `${point.chartX},${point.chartY}`).join(' ');
  }

  get areaPoints(): string {
    const points = this.points;
    if (!points.length) {
      return '';
    }

    const baseline = this.height - this.plot.bottom;
    const first = points[0];
    const last = points[points.length - 1];

    return `${first.chartX},${baseline} ${this.linePoints} ${last.chartX},${baseline}`;
  }

  get ticks(): Array<{ value: number; y: number }> {
    const domain = this.domain;
    const plotHeight = this.height - this.plot.top - this.plot.bottom;

    return Array.from({ length: 5 }, (_, index) => {
      const progress = index / 4;
      const value = Math.round(domain.max - progress * (domain.max - domain.min));

      return {
        value,
        y: Math.round(this.plot.top + progress * plotHeight),
      };
    });
  }

  get baselineY(): number {
    return this.height - this.plot.bottom;
  }

  get domain(): { min: number; max: number } {
    const values = this.data.map(item => item.value);
    const low = values.length ? Math.min(...values) : 0;
    const high = values.length ? Math.max(...values) : 100;
    const min = Math.max(0, Math.floor((low - 8) / 5) * 5);
    const max = Math.min(100, Math.ceil((high + 8) / 5) * 5);

    if (max - min < 20) {
      return {
        min: Math.max(0, min - 10),
        max: Math.min(100, max + 10),
      };
    }

    return { min, max };
  }
}
