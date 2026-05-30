import { Component, inject } from '@angular/core';
import { HttpClient } from '@angular/common/http';
import { Router } from '@angular/router';
import { NavbarComponent } from '../../components/navbar/navbar.component';
import { SearchBarComponent } from '../../components/search-bar/search-bar.component';
import { SuggestionCardComponent } from '../../components/suggestion-card/suggestion-card.component';
import { FooterComponent } from '../../components/footer/footer.component';
import { ChatService } from '../../services/chat.service';
import { SuggestionCard } from '../../models/suggestion.model';

@Component({
  selector: 'app-home',
  templateUrl: './home.component.html',
  styleUrl: './home.component.scss',
  standalone: true,
  imports: [NavbarComponent, SearchBarComponent, SuggestionCardComponent, FooterComponent],
})
export class HomeComponent {
  private router = inject(Router);
  private chatService = inject(ChatService);
  private http = inject(HttpClient);

  searchQuery = '';

  readonly menuIcon = 'https://www.figma.com/api/mcp/asset/bf38406a-41b7-469f-9494-bde7dcfaf424';
  readonly micIcon = 'https://www.figma.com/api/mcp/asset/812806d4-700d-427d-b54a-cd059ff6b726';
  readonly flashIcon = 'https://www.figma.com/api/mcp/asset/7f94b787-c3c1-4de7-a3ed-120ecd0013f9';

  suggestions: SuggestionCard[] = [
    {
      icon: 'https://www.figma.com/api/mcp/asset/929399b3-69ce-47ea-ac8a-5ca9f3d4a8bf',
      textMain: 'I felt anxious today and could not focus',
      textSecondary: 'help me understand what may have triggered it',
    },
  ];

  constructor() {
    this.http.get<SuggestionCard[]>('/api/suggestions').subscribe({
      next: suggestions => {
        if (suggestions.length) {
          this.suggestions = suggestions;
        }
      },
      error: () => undefined,
    });
  }

  sendMessage(text: string): void {
    if (!text.trim()) return;
    this.chatService.sendUserMessage(text);
    this.router.navigate(['/dashboard']);
  }

  onAskThis(card: SuggestionCard): void {
    this.sendMessage(card.textMain + ' ' + card.textSecondary);
  }
}
