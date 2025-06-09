import { createCollectionsService } from "./collections";
import { createDashboardService } from "./dashboard";
import { createProjectsService } from "./projects";
import { createUserService } from "./user";

export function initializeServices(authToken: string) {
  return {
    collections: createCollectionsService(authToken),
    dashboard: createDashboardService(authToken),
    user: createUserService(authToken),
    projects: createProjectsService(authToken),
  };
}

export type Services = ReturnType<typeof initializeServices>;
