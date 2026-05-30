export interface Message {
  id: string;
  role: 'user' | 'ai';
  text: string;
  time: string;
}
