import { Component, ElementRef, ViewChild, inject, effect } from '@angular/core';
import { FormsModule } from '@angular/forms';
import { ChatService } from '../../services/chat.service';

@Component({
  selector: 'app-chat',
  templateUrl: './chat.component.html',
  styleUrl: './chat.component.scss',
  standalone: true,
  imports: [FormsModule],
})
export class ChatComponent {
  private chatService = inject(ChatService);

  messages = this.chatService.messages;
  isTyping = this.chatService.isTyping;
  inputText = '';
  focused = false;

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

  onKeydown(e: KeyboardEvent): void {
    if (e.key === 'Enter' && !e.shiftKey) {
      e.preventDefault();
      this.send();
    }
  }
}
