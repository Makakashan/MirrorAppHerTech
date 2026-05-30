import { Component, EventEmitter, Input, Output } from '@angular/core';
import { SuggestionCard } from '../../models/suggestion.model';

@Component({
  selector: 'app-suggestion-card',
  templateUrl: './suggestion-card.component.html',
  styleUrl: './suggestion-card.component.scss',
  standalone: true,
})
export class SuggestionCardComponent {
  @Input() card!: SuggestionCard;
  @Input() flashIcon = '';
  @Output() askThis = new EventEmitter<SuggestionCard>();
}
