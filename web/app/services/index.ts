import { createTablesService } from "./tables";
import { createDashboardService } from "./dashboard";
import { createProjectsService } from "./projects";
import { createUserService } from "./user";
import { createLogsService } from "./logs";
import { createSettingsService } from "./settings";
import { createStorageService } from "./storage";

export function initializeServices(authToken: string) {
  return {
    tables: createTablesService(authToken),
    dashboard: createDashboardService(authToken),
    user: createUserService(authToken),
    projects: createProjectsService(authToken),
    logs: createLogsService(authToken),
    settings: createSettingsService(authToken),
    storage: createStorageService(authToken),
  };
}

export type Services = ReturnType<typeof initializeServices>;
