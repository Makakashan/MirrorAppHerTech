export interface Message {
  id: string;
  role: 'user' | 'ai';
  text: string;
  time: string;
}

export interface ChatSession {
  id: string;
  title: string;
  startedAt: string;
  messages: Message[];
}
