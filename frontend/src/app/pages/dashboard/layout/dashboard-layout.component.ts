import { Component } from "@angular/core";
import { Router, RouterOutlet, RouterLink, RouterLinkActive } from "@angular/router";
import {
	LucideCalendarDays,
	LucideChartColumn,
	LucideHouse,
	LucideMessageCircle,
	LucideSettings,
	LucideUser,
} from "@lucide/angular";
import { ThemeService } from "../../../services/theme.service";
import { ChatService } from "../../../services/chat.service";
import { ChatSession } from "../../../models/message.model";

@Component({
	selector: "app-dashboard-layout",
	templateUrl: "./dashboard-layout.component.html",
	styleUrl: "./dashboard-layout.component.scss",
	standalone: true,
	imports: [
		RouterOutlet,
		RouterLink,
		RouterLinkActive,
		LucideCalendarDays,
		LucideChartColumn,
		LucideHouse,
		LucideMessageCircle,
		LucideSettings,
		LucideUser,
	],
})
export class DashboardLayoutComponent {
	sidebarWidth = 240;
	sidebarCollapsed = true;

	constructor(
		readonly themeService: ThemeService,
		readonly chatService: ChatService,
		private router: Router,
	) {}

	toggleSidebar(): void {
		this.sidebarCollapsed = !this.sidebarCollapsed;
	}

	openNewChat(): void {
		this.chatService.startNewChat();
		this.router.navigate(['/chat']);
	}

	openSession(session: ChatSession): void {
		this.chatService.loadSession(session);
		this.router.navigate(['/chat']);
	}

	formatDate(iso: string): string {
		const d = new Date(iso);
		return d.toLocaleDateString([], { month: 'short', day: 'numeric' });
	}
}
