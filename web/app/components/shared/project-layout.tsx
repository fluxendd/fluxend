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
import type { User } from "~/services/user";

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

  const projectPromise = services.projects.getProjectDetails(projectId);
  const userPromise = services.user.getCurrentUser();

  const [project, user] = await Promise.all([projectPromise, userPromise]);

  return data({ projectDetails: project, user, authToken }, { status: 200 });
}

export function shouldRevalidate(arg: ShouldRevalidateFunctionArgs) {
  return false;
}

export type ProjectLayoutOutletContext = {
  userDetails: User;
  projectDetails: Project;
  services: Services;
};

export default function ProjectLayout({ loaderData }: Route.ComponentProps) {
  const { projectId } = useParams();
  const { user } = loaderData;
  const services = initializeServices(loaderData.authToken);

  if (!projectId || !user.content) {
    throw new Error("Unauthorized");
  }

  return (
    <QueryClientProvider client={queryClient}>
      <TooltipProvider>
        <SidebarProvider open={false}>
          <ProjectSidebar projectId={projectId} userDetails={user.content} />
          <SidebarInset>
            <Outlet
              context={{
                projectDetails: loaderData.projectDetails,
                userDetails: loaderData.user,
                services,
              }}
            />
            <FloatingLoadingIcon />
          </SidebarInset>
        </SidebarProvider>
      </TooltipProvider>
    </QueryClientProvider>
  );
}
