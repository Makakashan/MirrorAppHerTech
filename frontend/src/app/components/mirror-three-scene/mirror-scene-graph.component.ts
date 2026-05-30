import { CUSTOM_ELEMENTS_SCHEMA, ChangeDetectionStrategy, Component, OnDestroy } from '@angular/core';
import { beforeRender, NgtArgs } from 'angular-three';
import * as THREE from 'three';

@Component({
  selector: 'app-mirror-scene-graph',
  standalone: true,
  imports: [NgtArgs],
  schemas: [CUSTOM_ELEMENTS_SCHEMA],
  changeDetection: ChangeDetectionStrategy.OnPush,
  template: `<ngt-primitive *args="[sceneGroup]" />`,
})
export class MirrorSceneGraphComponent implements OnDestroy {
  readonly sceneGroup = new THREE.Group();
  private readonly particles: THREE.Points;
  private readonly line: THREE.Line;
  private readonly cards = new THREE.Group();
  private readonly destroyFrame: () => void;

  constructor() {
    this.particles = this.createParticles();
    this.line = this.createMoodLine();
    this.sceneGroup.add(this.particles, this.line, this.cards);
    this.addCardPlanes();

    this.destroyFrame = beforeRender(({ clock }) => {
      const time = clock.getElapsedTime();
      this.particles.rotation.y = time * 0.035;
      this.particles.rotation.x = Math.sin(time * 0.2) * 0.04;
      this.line.rotation.y = Math.sin(time * 0.35) * 0.08;
      this.cards.children.forEach((card: THREE.Object3D, index: number) => {
        card.position.y = 1.35 + Math.sin(time * 0.8 + index) * 0.08;
        card.rotation.y = Math.sin(time * 0.45 + index) * 0.12;
      });
    });
  }

  ngOnDestroy(): void {
    this.destroyFrame();
    this.disposeObject(this.sceneGroup);
  }

  private createParticles(): THREE.Points {
    const count = 900;
    const positions = new Float32Array(count * 3);
    const colors = new Float32Array(count * 3);
    const teal = new THREE.Color('#2dd4bf');
    const blue = new THREE.Color('#38bdf8');

    for (let i = 0; i < count; i++) {
      const index = i * 3;
      positions[index] = THREE.MathUtils.randFloatSpread(9);
      positions[index + 1] = THREE.MathUtils.randFloatSpread(5.6);
      positions[index + 2] = THREE.MathUtils.randFloat(-5.5, 1.8);

      const color = teal.clone().lerp(blue, Math.random() * 0.45);
      colors[index] = color.r;
      colors[index + 1] = color.g;
      colors[index + 2] = color.b;
    }

    const geometry = new THREE.BufferGeometry();
    geometry.setAttribute('position', new THREE.BufferAttribute(positions, 3));
    geometry.setAttribute('color', new THREE.BufferAttribute(colors, 3));
    geometry.computeBoundingSphere();

    const material = new THREE.PointsMaterial({
      size: 0.025,
      transparent: true,
      opacity: 0.62,
      vertexColors: true,
      depthWrite: false,
    });

    return new THREE.Points(geometry, material);
  }

  private createMoodLine(): THREE.Line {
    const scores = [68, 74, 71, 86, 82, 90, 78];
    const points = scores.map((score, index) => {
      const x = (index / (scores.length - 1)) * 4.8 - 2.4;
      const y = ((score - 50) / 50) * 1.7 - 0.75;
      const z = Math.sin(index * 0.8) * 0.22;
      return new THREE.Vector3(x, y, z);
    });

    const geometry = new THREE.BufferGeometry().setFromPoints(points);
    const material = new THREE.LineBasicMaterial({
      color: '#2dd4bf',
      transparent: true,
      opacity: 0.96,
    });
    const line = new THREE.Line(geometry, material);
    line.position.set(0.35, -0.85, -1.25);
    line.rotation.x = -0.18;

    points.forEach(point => {
      const dot = new THREE.Mesh(
        new THREE.SphereGeometry(0.055, 18, 18),
        new THREE.MeshBasicMaterial({ color: '#f8fafc' }),
      );
      dot.position.copy(point);
      line.add(dot);
    });

    return line;
  }

  private addCardPlanes(): void {
    const material = new THREE.MeshBasicMaterial({
      color: '#2dd4bf',
      transparent: true,
      opacity: 0.12,
      side: THREE.DoubleSide,
    });

    [
      [-2.35, 1.15, -1.6],
      [0, 1.32, -1.9],
      [2.35, 1.08, -1.6],
    ].forEach(([x, y, z], index) => {
      const mesh = new THREE.Mesh(new THREE.PlaneGeometry(1.35, 0.72), material.clone());
      mesh.position.set(x, y, z);
      mesh.rotation.set(-0.12, (index - 1) * 0.18, 0);
      this.cards.add(mesh);
    });

  }

  private disposeObject(object: THREE.Object3D): void {
    object.traverse((child: THREE.Object3D) => {
      const mesh = child as THREE.Mesh | THREE.Points | THREE.Line;
      mesh.geometry?.dispose();
      const material = mesh.material;
      if (Array.isArray(material)) {
        material.forEach(item => item.dispose());
      } else {
        material?.dispose();
      }
    });
  }
}
