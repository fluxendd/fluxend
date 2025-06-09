import { Outlet, useNavigation } from "react-router";
import { SidebarInset, SidebarProvider } from "~/components/ui/sidebar";
import { AppSidebar } from "~/components/shared/app-sidebar";
import { QueryClientProvider } from "@tanstack/react-query";
import { queryClient } from "~/lib/query";
import { LoaderCircle } from "lucide-react";
import { TooltipProvider } from "~/components/ui/tooltip";

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

export default function AppLayout() {
  return (
    <QueryClientProvider client={queryClient}>
      <TooltipProvider>
        <SidebarProvider open={false}>
          <AppSidebar />
          <SidebarInset>
            <Outlet />
            <FloatingLoadingIcon />
          </SidebarInset>
        </SidebarProvider>
      </TooltipProvider>
    </QueryClientProvider>
  );
}
