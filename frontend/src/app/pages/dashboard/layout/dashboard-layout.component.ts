import { Component } from "@angular/core";
import { RouterOutlet, RouterLink, RouterLinkActive } from "@angular/router";
import {
	LucideCalendarDays,
	LucideChartColumn,
	LucideHouse,
	LucideMessageCircle,
	LucideSettings,
	LucideUser,
} from "@lucide/angular";
import { ThemeService } from "../../../services/theme.service";

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
	readonly recentChats = ['Morning anxiety check-in', 'Work boundary reflection', 'Evening gratitude note'];

	constructor(readonly themeService: ThemeService) {}

	toggleSidebar(): void {
		this.sidebarCollapsed = !this.sidebarCollapsed;
	}
}
