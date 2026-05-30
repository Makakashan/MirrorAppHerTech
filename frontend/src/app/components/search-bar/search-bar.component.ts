import { Component, EventEmitter, Input, Output } from '@angular/core';

@Component({
  selector: 'app-search-bar',
  templateUrl: './search-bar.component.html',
  styleUrl: './search-bar.component.scss',
  standalone: true,
})
export class SearchBarComponent {
  @Input() value = '';
  @Output() valueChange = new EventEmitter<string>();
  @Input() micIcon = '';
  @Output() submit = new EventEmitter<string>();

  onInput(event: Event): void {
    this.valueChange.emit((event.target as HTMLInputElement).value);
  }

  onSubmit(): void {
    if (this.value.trim()) {
      this.submit.emit(this.value.trim());
    }
  }
}
