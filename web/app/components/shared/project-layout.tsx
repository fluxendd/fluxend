import {
  data,
  Outlet,
  useNavigation,
  useParams,
  type ShouldRevalidateFunctionArgs,
} from "react-router";
import { SidebarInset, SidebarProvider } from "~/components/ui/sidebar";
import { ProjectSidebar } from "~/components/shared/project-sidebar";
import { QueryClientProvider } from "@tanstack/react-query";
import { queryClient } from "~/lib/query";
import { LoaderCircle } from "lucide-react";
import { TooltipProvider } from "~/components/ui/tooltip";
import { getServerAuthToken } from "~/lib/auth";
import type { Route } from "./+types/project-layout";
import { initializeServices, type Services } from "~/services";
import type { Project } from "~/services/projects";

const FloatingLoadingIcon = () => {
  const navigation = useNavigation();
  if (navigation.state === "idle") {
    return null;
  }

  return (
    <div className="absolute bottom-6 right-6 bg-muted rounded-lg p-2">
      <LoaderCircle className="loading-icon" />
    </div>
  );
};

export async function loader({ request, params }: Route.LoaderArgs) {
  const { projectId } = params;
  const authToken = await getServerAuthToken(request.headers);
  if (!authToken || !projectId) {
    throw new Error("Unauthorized");
  }

  const services = initializeServices(authToken);

  const project = await services.projects.getProjectDetails(projectId);
  return data({ projectDetails: project, authToken }, { status: 200 });
}

export function shouldRevalidate(arg: ShouldRevalidateFunctionArgs) {
  return false;
}

export type ProjectLayoutOutletContext = {
  projectDetails: Project;
  services: Services;
};

export default function ProjectLayout({ loaderData }: Route.ComponentProps) {
  const { projectId } = useParams();
  const services = initializeServices(loaderData.authToken);

  return (
    <QueryClientProvider client={queryClient}>
      <TooltipProvider>
        <SidebarProvider open={false}>
          <ProjectSidebar projectId={projectId} />
          <SidebarInset>
            <Outlet
              context={{ projectDetails: loaderData.projectDetails, services }}
            />
            <FloatingLoadingIcon />
          </SidebarInset>
        </SidebarProvider>
      </TooltipProvider>
    </QueryClientProvider>
  );
}
