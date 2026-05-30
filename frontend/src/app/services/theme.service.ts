import { DOCUMENT, isPlatformBrowser } from '@angular/common';
import { Inject, Injectable, PLATFORM_ID, signal } from '@angular/core';

export type ThemeMode = 'light' | 'dark' | 'system';

const storageKey = 'mirror-theme';

@Injectable({ providedIn: 'root' })
export class ThemeService {
  readonly mode = signal<ThemeMode>('system');
  readonly reducedMotion = signal(false);
  readonly privateMode = signal(true);
  readonly compactMode = signal(false);

  constructor(
    @Inject(DOCUMENT) private document: Document,
    @Inject(PLATFORM_ID) private platformId: object,
  ) {
    if (!this.isBrowser) {
      return;
    }

    const savedTheme = window.localStorage.getItem(storageKey) as ThemeMode | null;
    const savedMotion = window.localStorage.getItem('mirror-reduced-motion');
    const savedPrivacy = window.localStorage.getItem('mirror-private-mode');
    const savedDensity = window.localStorage.getItem('mirror-compact-mode');

    if (savedTheme === 'light' || savedTheme === 'dark' || savedTheme === 'system') {
      this.mode.set(savedTheme);
    }

    this.reducedMotion.set(savedMotion === 'true');
    this.privateMode.set(savedPrivacy !== 'false');
    this.compactMode.set(savedDensity === 'true');
    this.applyTheme(this.mode());
    this.applyPreferences();
  }

  setTheme(mode: ThemeMode): void {
    this.mode.set(mode);

    if (this.isBrowser) {
      window.localStorage.setItem(storageKey, mode);
    }

    this.applyTheme(mode);
  }

  setReducedMotion(value: boolean): void {
    this.reducedMotion.set(value);
    this.persistFlag('mirror-reduced-motion', value);
    this.applyPreferences();
  }

  setPrivateMode(value: boolean): void {
    this.privateMode.set(value);
    this.persistFlag('mirror-private-mode', value);
    this.applyPreferences();
  }

  setCompactMode(value: boolean): void {
    this.compactMode.set(value);
    this.persistFlag('mirror-compact-mode', value);
    this.applyPreferences();
  }

  reset(): void {
    this.setTheme('system');
    this.setReducedMotion(false);
    this.setPrivateMode(true);
    this.setCompactMode(false);
  }

  private applyTheme(mode: ThemeMode): void {
    if (!this.isBrowser) {
      return;
    }

    const prefersDark = window.matchMedia('(prefers-color-scheme: dark)').matches;
    const resolved = mode === 'system' ? (prefersDark ? 'dark' : 'light') : mode;
    this.document.documentElement.dataset['theme'] = resolved;
    this.document.documentElement.dataset['themeMode'] = mode;
  }

  private applyPreferences(): void {
    this.document.documentElement.dataset['motion'] = this.reducedMotion() ? 'reduced' : 'default';
    this.document.documentElement.dataset['privacy'] = this.privateMode() ? 'private' : 'standard';
    this.document.documentElement.dataset['density'] = this.compactMode() ? 'compact' : 'comfortable';
  }

  private persistFlag(key: string, value: boolean): void {
    if (this.isBrowser) {
      window.localStorage.setItem(key, String(value));
    }
  }

  private get isBrowser(): boolean {
    return isPlatformBrowser(this.platformId);
  }
}
