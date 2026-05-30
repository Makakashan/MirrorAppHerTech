import { Component } from '@angular/core';

interface DayMood {
  day: number;
  mood: 'great' | 'good' | 'neutral' | 'bad';
  isFuture: boolean;
}

@Component({
  selector: 'app-metrics-panel',
  templateUrl: './metrics-panel.component.html',
  styleUrl: './metrics-panel.component.scss',
  standalone: true,
  imports: [],
})
export class MetricsPanelComponent {
  readonly today = new Date();
  readonly currentMonth = this.today.toLocaleDateString('en-US', { month: 'long' });

  readonly chartData = [
    { label: 'Mon', value: 72, x: 10 },
    { label: 'Tue', value: 85, x: 107 },
    { label: 'Wed', value: 78, x: 203 },
    { label: 'Thu', value: 91, x: 300 },
    { label: 'Fri', value: 88, x: 397 },
    { label: 'Sat', value: 76, x: 493 },
    { label: 'Sun', value: 84, x: 590 },
  ];

  readonly chartPolyline = this.chartData
    .map(d => `${d.x},${this.toY(d.value)}`)
    .join(' ');

  readonly chartArea =
    `10,130 ${this.chartPolyline} 590,130`;

  readonly calendarDays: DayMood[] = this.buildCalendar();

  private toY(v: number): number {
    return Math.round(130 - ((v - 50) / 50) * 120);
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
