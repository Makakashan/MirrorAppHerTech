import { Component } from '@angular/core';
import { ChatComponent } from '../../components/chat/chat.component';

@Component({
  selector: 'app-chat-page',
  standalone: true,
  imports: [ChatComponent],
  template: `<app-chat />`,
  styles: [`:host { display: flex; flex-direction: column; height: 100%; }`],
})
export class ChatPageComponent {}
