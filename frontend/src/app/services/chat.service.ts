import { Injectable, signal } from '@angular/core';
import { Message } from '../models/message.model';

const AI_RESPONSES: { keywords: string[]; text: string }[] = [
  {
    keywords: ['gym', 'sport', 'workout', 'sweat'],
    text: 'For gym workouts I recommend the Jabra Elite 8 Active — IP57 water resistance, secure-fit wing tips, and 8 hours of battery life. They stay in place even during intense sessions.',
  },
  {
    keywords: ['budget', 'cheap', '$100', 'under', 'affordable'],
    text: "Under $100, the Sony WH-CH720N is a great choice — noise cancellation, 35-hour battery, and Sony's signature sound tuning. Outstanding value for the price.",
  },
  {
    keywords: ['ear', 'in-ear', 'earphone'],
    text: 'For all-day in-ear comfort, the Shure AONIC 215 are hard to beat — lightweight over-ear cable, memory foam tips, and detailed, accurate sound reproduction.',
  },
  {
    keywords: ['noise', 'cancel', 'wireless'],
    text: 'The Sony WH-1000XM5 leads the market in active noise cancellation. It offers 30-hour battery life, adaptive sound control, and exceptional audio quality — ideal for both commuting and focused work.',
  },
];

const DEFAULT_RESPONSE =
  "The Sony WH-1000XM5 is an excellent all-around choice — industry-leading noise cancellation, 30-hour battery, and premium sound. It adapts to your environment automatically and works great for both travel and work.";

@Injectable({ providedIn: 'root' })
export class ChatService {
  readonly messages = signal<Message[]>([
    {
      id: '0',
      role: 'ai',
      text: "Hi! I'm Mirror, your AI assistant. Tell me what you're looking for and I'll find the perfect headphones for you.",
      time: this.getTime(),
    },
  ]);

  readonly isTyping = signal(false);

  sendUserMessage(text: string): void {
    this.messages.update(msgs => [
      ...msgs,
      { id: Date.now().toString(), role: 'user', text, time: this.getTime() },
    ]);
    this.isTyping.set(true);
    setTimeout(() => {
      this.isTyping.set(false);
      this.messages.update(msgs => [
        ...msgs,
        { id: (Date.now() + 1).toString(), role: 'ai', text: this.getResponse(text), time: this.getTime() },
      ]);
    }, 1400);
  }

  private getTime(): string {
    return new Date().toLocaleTimeString([], { hour: '2-digit', minute: '2-digit' });
  }

  private getResponse(query: string): string {
    const q = query.toLowerCase();
    const match = AI_RESPONSES.find(r => r.keywords.some(k => q.includes(k)));
    return match ? match.text : DEFAULT_RESPONSE;
  }
}
