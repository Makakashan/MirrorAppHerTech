import { Injectable, signal } from '@angular/core';
import { HttpClient } from '@angular/common/http';
import { Metrics } from '../models/metrics.model';

@Injectable({ providedIn: 'root' })
export class MetricsService {
  readonly metrics = signal<Metrics | null>(null);

  constructor(private http: HttpClient) {}

  load(): void {
    this.http.get<Metrics>('/api/metrics').subscribe({
      next: metrics => this.metrics.set(metrics),
      error: () => undefined,
    });
  }

  set(metrics: Metrics): void {
    this.metrics.set(metrics);
  }
}
