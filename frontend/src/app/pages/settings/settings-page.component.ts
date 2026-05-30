import { Component } from '@angular/core';
import {
  LucideMonitor,
  LucideMoon,
  LucidePalette,
  LucideSun,
} from '@lucide/angular';
import { ThemeMode, ThemeService } from '../../services/theme.service';

@Component({
  selector: 'app-settings-page',
  standalone: true,
  imports: [
    LucideMonitor,
    LucideMoon,
    LucidePalette,
    LucideSun,
  ],
  templateUrl: './settings-page.component.html',
  styleUrl: './settings-page.component.scss',
})
export class SettingsPageComponent {
  readonly themeOptions: Array<{ mode: ThemeMode; label: string; description: string }> = [
    { mode: 'light', label: 'Light', description: 'Default light interface' },
    { mode: 'dark', label: 'Dark', description: 'Use dark colors across the app' },
    { mode: 'system', label: 'System', description: 'Match your device setting' },
  ];

  constructor(readonly themeService: ThemeService) {}

  setTheme(mode: ThemeMode): void {
    this.themeService.setTheme(mode);
  }
}
