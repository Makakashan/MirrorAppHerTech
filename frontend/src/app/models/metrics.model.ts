export interface DayMood {
  day: number;
  mood: 'great' | 'good' | 'neutral' | 'bad';
  isFuture: boolean;
}

export interface ChartPoint {
  label: string;
  value: number;
  x: number;
}

export interface Metrics {
  averageScore: number;
  averageDelta: number;
  streakDays: number;
  currentMonth: string;
  chartData: ChartPoint[];
  calendarDays: DayMood[];
}
