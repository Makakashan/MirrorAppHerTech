import { ChangeDetectionStrategy, Component } from '@angular/core';
import { RouterLink } from '@angular/router';
import { LucideCalendarDays, LucideChartColumn, LucideMessageCircle, LucideSend } from '@lucide/angular';
import { MirrorThreeSceneComponent } from '../../components/mirror-three-scene/mirror-three-scene.component';

@Component({
  selector: 'app-mirror-main-page',
  standalone: true,
  imports: [
    RouterLink,
    LucideCalendarDays,
    LucideChartColumn,
    LucideMessageCircle,
    LucideSend,
    MirrorThreeSceneComponent,
  ],
  templateUrl: './mirror-main-page.component.html',
  styleUrl: './mirror-main-page.component.scss',
  changeDetection: ChangeDetectionStrategy.OnPush,
})
export class MirrorMainPageComponent {
  readonly stats = [
    { label: 'Mood score', value: '8.1', change: '+0.9', type: 'score' },
    { label: 'Streak', value: '12', change: 'days', type: 'streak' },
    { label: 'Stability', value: '76%', change: '+6%', type: 'stability' },
  ];

  readonly days = ['Mon', 'Tue', 'Wed', 'Thu', 'Fri', 'Sat', 'Sun'];
}
