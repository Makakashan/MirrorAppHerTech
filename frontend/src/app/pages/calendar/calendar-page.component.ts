import { Component, effect } from '@angular/core';
import { LucideChevronDown, LucideChevronLeft, LucideChevronRight } from '@lucide/angular';
import { DayMood } from '../../models/metrics.model';
import { MetricsService } from '../../services/metrics.service';

interface CalendarCell extends DayMood {
  id: string;
  isToday: boolean;
  note: string;
}

type CalendarView = 'month' | 'week';

@Component({
  selector: 'app-calendar-page',
  standalone: true,
  imports: [LucideChevronDown, LucideChevronLeft, LucideChevronRight],
  templateUrl: './calendar-page.component.html',
  styleUrl: './calendar-page.component.scss',
})
export class CalendarPageComponent {
  readonly weekdays = ['Mon', 'Tue', 'Wed', 'Thu', 'Fri', 'Sat', 'Sun'];
  currentMonth = new Date().toLocaleDateString('en-US', { month: 'long', year: 'numeric' });
  view: CalendarView = 'month';
  days: CalendarCell[] = this.buildFallbackCalendar();
  selectedDay = this.days.find(day => day.isToday) ?? this.days[0];

  constructor(private metricsService: MetricsService) {
    effect(() => {
      const metrics = this.metricsService.metrics();
      if (metrics) {
        this.currentMonth = `${metrics.currentMonth} ${new Date().getFullYear()}`;
        this.days = this.decorateDays(metrics.calendarDays);
        this.selectedDay = this.days.find(day => day.isToday) ?? this.days.find(day => !day.isFuture) ?? this.days[0];
      }
    });
    this.metricsService.load();
  }

  selectDay(day: CalendarCell): void {
    this.selectedDay = day;
  }

  setView(view: CalendarView): void {
    this.view = view;
  }

  get visibleDays(): CalendarCell[] {
    if (this.view === 'week') {
      const selectedIndex = this.days.findIndex(day => day.id === this.selectedDay.id);
      const start = Math.max(0, selectedIndex - (selectedIndex % 7));
      return this.days.slice(start, start + 7);
    }

    return this.days;
  }

  get completedDays(): number {
    return this.days.filter(day => !day.isFuture).length;
  }

  get difficultDays(): number {
    return this.days.filter(day => !day.isFuture && day.mood === 'bad').length;
  }

  get steadyDays(): number {
    return this.days.filter(day => !day.isFuture && (day.mood === 'good' || day.mood === 'great')).length;
  }

  get selectedWeekday(): string {
    const index = this.visibleDays.findIndex(day => day.id === this.selectedDay.id);
    return this.weekdays[Math.max(index, 0) % 7];
  }

  moodLabel(mood: DayMood['mood']): string {
    return {
      great: 'Great',
      good: 'Good',
      neutral: 'Neutral',
      bad: 'Difficult',
    }[mood];
  }

  moodTone(mood: DayMood['mood']): string {
    return {
      great: 'High',
      good: 'Steady',
      neutral: 'Mixed',
      bad: 'Low',
    }[mood];
  }

  private decorateDays(days: DayMood[]): CalendarCell[] {
    const today = new Date();
    return days.map((day, index) => ({
      ...day,
      id: `${day.day}-${index}`,
      isToday: day.day === today.getDate() && !day.isFuture,
      note: this.noteFor(day.mood),
    }));
  }

  private buildFallbackCalendar(): CalendarCell[] {
    const today = new Date();
    const dow = today.getDay() || 7;
    const monday = new Date(today);
    monday.setDate(today.getDate() - dow + 1 - 21);
    const moods: DayMood['mood'][] = ['neutral', 'good', 'bad', 'neutral', 'good', 'great', 'good'];

    return Array.from({ length: 28 }, (_, index) => {
      const date = new Date(monday);
      date.setDate(monday.getDate() + index);
      const mood = moods[index % moods.length];
      return {
        id: `${date.getDate()}-${index}`,
        day: date.getDate(),
        mood,
        isFuture: date > today,
        isToday: date.getDate() === today.getDate(),
        note: this.noteFor(mood),
      };
    });
  }

  private noteFor(mood: DayMood['mood']): string {
    return {
      great: 'Energy was high and routines seemed to support recovery.',
      good: 'Stable day with enough space to return to focus.',
      neutral: 'Mixed signals without a strong emotional swing.',
      bad: 'Pressure or avoidance likely shaped the day.',
    }[mood];
  }
}
