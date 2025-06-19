import { data, Outlet, useNavigation } from "react-router";
import { SidebarInset, SidebarProvider } from "~/components/ui/sidebar";
import { AppSidebar } from "~/components/shared/app-sidebar";
import { QueryClientProvider } from "@tanstack/react-query";
import { queryClient } from "~/lib/query";
import { LoaderCircle } from "lucide-react";
import { TooltipProvider } from "~/components/ui/tooltip";
import { getServerAuthToken } from "~/lib/auth";
import type { Route } from "./+types/app-layout";
import { initializeServices } from "~/services";

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
  const authToken = await getServerAuthToken(request.headers);

  if (!authToken) {
    throw new Error("Unauthorized");
  }

  const services = initializeServices(authToken);

  const user = await services.user.getCurrentUser();

  return data({ user, authToken }, { status: 200 });
}

export default function AppLayout({ loaderData }: Route.ComponentProps) {
  const { user, authToken } = loaderData;

  if (!user.content) {
    return null;
  }

  return (
    <QueryClientProvider client={queryClient}>
      <TooltipProvider>
        <SidebarProvider open={false}>
          <AppSidebar authToken={authToken} userDetails={user.content} />
          <SidebarInset>
            <Outlet />
            <FloatingLoadingIcon />
          </SidebarInset>
        </SidebarProvider>
      </TooltipProvider>
    </QueryClientProvider>
  );
}
