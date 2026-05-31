import { Component, ElementRef, ViewChild, computed, effect, inject } from '@angular/core';
import { FormsModule } from '@angular/forms';
import { LucideArrowUp, LucideFileText, LucideSparkles, LucideTrendingUp } from '@lucide/angular';
import { ChatService } from '../../services/chat.service';

@Component({
  selector: 'app-chat',
  templateUrl: './chat.component.html',
  styleUrl: './chat.component.scss',
  standalone: true,
  imports: [
    FormsModule,
    LucideArrowUp,
    LucideFileText,
    LucideSparkles,
    LucideTrendingUp,
  ],
})
export class ChatComponent {
  private chatService = inject(ChatService);

  messages = this.chatService.messages;
  visibleMessages = computed(() => this.messages().filter(message => message.id !== '0'));
  isTyping = this.chatService.isTyping;
  inputText = '';
  focused = false;
  readonly promptCards = [
    {
      title: 'What changed my mood today?',
      text: 'Find the strongest emotional shift in the last 24 hours.',
      icon: 'sparkles',
    },
    {
      title: 'Weekly mood pattern',
      text: 'Summarize repeated triggers and steady moments.',
      icon: 'trend',
    },
    {
      title: 'Deep reflection',
      text: 'Turn a messy thought into a clear journal note.',
      icon: 'file',
    },
  ];

  @ViewChild('messagesEl') messagesEl!: ElementRef<HTMLElement>;

  constructor() {
    effect(() => {
      this.messages();
      this.isTyping();
      setTimeout(() => {
        const el = this.messagesEl?.nativeElement;
        if (el) el.scrollTop = el.scrollHeight;
      }, 0);
    });
  }

  send(): void {
    const text = this.inputText.trim();
    if (!text) return;
    this.inputText = '';
    this.chatService.sendUserMessage(text);
  }

  askSuggestion(text: string): void {
    this.inputText = text;
    this.send();
  }

  onKeydown(e: KeyboardEvent): void {
    if (e.key === 'Enter' && !e.shiftKey) {
      e.preventDefault();
      this.send();
    }
  }
}
