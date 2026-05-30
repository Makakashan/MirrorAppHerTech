import { ChangeDetectionStrategy, Component } from '@angular/core';
import { NgtCanvas } from 'angular-three/dom';
import { MirrorSceneGraphComponent } from './mirror-scene-graph.component';

@Component({
  selector: 'app-mirror-three-scene',
  standalone: true,
  imports: [NgtCanvas, MirrorSceneGraphComponent],
  changeDetection: ChangeDetectionStrategy.OnPush,
  template: `
    <ngt-canvas
      [camera]="camera"
      [gl]="gl"
      [dpr]="[1, 1.75]"
      [linear]="true"
      [lookAt]="[0, 0, -1]"
    >
      <app-mirror-scene-graph *canvasContent />
    </ngt-canvas>
  `,
  styles: [`
    :host {
      display: block;
      position: absolute;
      inset: 0;
      pointer-events: none;
    }

    ngt-canvas {
      display: block;
      width: 100%;
      height: 100%;
    }
  `],
})
export class MirrorThreeSceneComponent {
  readonly camera = { position: [0, 0, 5.2] as [number, number, number], fov: 42, near: 0.1, far: 100 };
  readonly gl = {
    alpha: true,
    antialias: true,
    preserveDrawingBuffer: true,
    powerPreference: 'high-performance' as WebGLPowerPreference,
  };
}
