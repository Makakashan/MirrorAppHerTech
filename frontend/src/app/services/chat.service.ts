import { Injectable, signal } from '@angular/core';
import { HttpClient } from '@angular/common/http';
import { Message } from '../models/message.model';
import { Metrics } from '../models/metrics.model';
import { MetricsService } from './metrics.service';

const AI_RESPONSES: { keywords: string[]; text: string }[] = [
  {
    keywords: ['anxious', 'anxiety', 'panic', 'worried', 'stress', 'stressed'],
    text: 'It sounds like your nervous system is staying on high alert. Try naming the trigger, the body signal, and the thought that followed it.',
  },
  {
    keywords: ['sad', 'empty', 'lonely', 'tired', 'hopeless'],
    text: 'I hear a lower emotional tone in that. Separate what happened from what you concluded about yourself, then track whether this pattern repeats.',
  },
  {
    keywords: ['angry', 'irritated', 'frustrated', 'rage'],
    text: 'Anger often points to a crossed boundary, unmet need, or accumulated pressure. Write down what felt unfair and what you needed in that moment.',
  },
  {
    keywords: ['avoid', 'procrastinate', 'postpone', 'focus'],
    text: 'Avoidance usually protects you from discomfort in the short term. The useful question is what emotion appears right before you switch away from the task.',
  },
];

const DEFAULT_RESPONSE =
  'Describe the situation, emotion, body reaction, and what you did next. Over several entries, Mirror can make recurring triggers and coping strategies easier to see.';

@Injectable({ providedIn: 'root' })
export class ChatService {
  private readonly apiBase = '/api';

  readonly messages = signal<Message[]>([
    {
      id: '0',
      role: 'ai',
      text: "Hi! I'm Mirror. Tell me what has been on your mind, and I will help you reflect on the patterns behind it.",
      time: this.getTime(),
    },
  ]);

  readonly isTyping = signal(false);

  constructor(
    private http: HttpClient,
    private metricsService: MetricsService,
  ) {
    this.http.get<Message[]>(`${this.apiBase}/messages`).subscribe({
      next: messages => {
        if (messages.length) {
          this.messages.set(messages);
        }
      },
      error: () => undefined,
    });
  }

  sendUserMessage(text: string): void {
    this.messages.update(msgs => [
      ...msgs,
      { id: Date.now().toString(), role: 'user', text, time: this.getTime() },
    ]);
    this.isTyping.set(true);

    this.http.post<{ message: Message; metrics: Metrics }>(`${this.apiBase}/chat`, { message: text }).subscribe({
      next: response => {
        this.metricsService.set(response.metrics);
        this.addAiMessage(response.message);
      },
      error: () => this.addAiMessage({
        id: (Date.now() + 1).toString(),
        role: 'ai',
        text: this.getResponse(text),
        time: this.getTime(),
      }),
    });
  }

  private getTime(): string {
    return new Date().toLocaleTimeString([], { hour: '2-digit', minute: '2-digit' });
  }

  private getResponse(query: string): string {
    const q = query.toLowerCase();
    const match = AI_RESPONSES.find(r => r.keywords.some(k => q.includes(k)));
    return match ? match.text : DEFAULT_RESPONSE;
  }

  private addAiMessage(message: Message): void {
    this.isTyping.set(false);
    this.messages.update(msgs => [...msgs, message]);
  }
}
