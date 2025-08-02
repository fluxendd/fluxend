import type { Route } from "./+types/sidebar";
import {
  data,
  Outlet,
  useOutletContext,
} from "react-router";
import { useState } from "react";
import {
  Sidebar,
  SidebarContent,
  SidebarGroup,
  SidebarGroupContent,
  SidebarHeader,
  SidebarInput,
  SidebarInset,
  SidebarProvider,
} from "~/components/ui/sidebar";
import { DocsList } from "./docs-list";
import { FileText } from "lucide-react";
import type { ProjectLayoutOutletContext } from "~/components/shared/project-layout";
import { getServerAuthToken } from "~/lib/auth";
import { initializeServices } from "~/services";

export function HydrateFallback() {
  return (
    <SidebarProvider>
      <Sidebar
        collapsible="none"
        className="hidden md:flex h-screen"
        variant="inset"
      >
        <SidebarHeader className="border-b px-2">
          <div className="flex items-center gap-2">
            <SidebarInput
              placeholder="Search documentation..."
              disabled
              className="flex-1 rounded-lg"
            />
          </div>
        </SidebarHeader>
        <SidebarContent>
          <SidebarGroup className="px-0">
            <SidebarGroupContent>
              <div className="animate-pulse">
                <div className="h-8 bg-gray-200 rounded m-2 mb-2"></div>
                <div className="h-8 bg-gray-200 rounded m-2 mb-2"></div>
                <div className="h-8 bg-gray-200 rounded m-2 mb-2"></div>
              </div>
            </SidebarGroupContent>
          </SidebarGroup>
        </SidebarContent>
      </Sidebar>
      <SidebarInset className="overflow-hidden flex flex-col">
        <div className="p-4 flex-1 overflow-auto">
          <Outlet />
        </div>
      </SidebarInset>
    </SidebarProvider>
  );
}

export async function loader({ request, params }: Route.LoaderArgs) {
  const authToken = await getServerAuthToken(request.headers);
  const { projectId } = params;

  if (!authToken) {
    throw new Error("Unauthorized");
  }

  const services = initializeServices(authToken);

  // Get the OpenAPI spec to extract available tables
  const response = await services.openapi.getProjectOpenAPI(projectId);
  
  if (!response.success) {
    throw new Error(response.errors?.[0] || "Failed to load API documentation");
  }

  // Parse the OpenAPI spec to extract tables
  let openApiSpec;
  try {
    if (typeof response.content === 'object') {
      openApiSpec = response.content;
    } else {
      openApiSpec = JSON.parse(response.content);
    }
  } catch {
    throw new Error("Invalid API documentation format");
  }

  // Extract available tables from paths
  const tables = new Set<string>();
  if (openApiSpec?.paths) {
    Object.keys(openApiSpec.paths).forEach(path => {
      const match = path.match(/^\/([^\/]+)(?:\/|$)/);
      if (match && match[1] !== 'rpc') {
        tables.add(match[1]);
      }
    });
  }

  return data({
    tables: Array.from(tables).sort(),
    openApiSpec
  }, { status: 200 });
}

export default function DocsSidebar({
  loaderData,
  params,
}: Route.ComponentProps) {
  const { projectId } = params;
  const [searchTerm, setSearchTerm] = useState("");
  const { projectDetails, services } =
    useOutletContext<ProjectLayoutOutletContext>();

  const handleSearch = (e: React.ChangeEvent<HTMLInputElement>) => {
    setSearchTerm(e.target.value);
  };

  return (
    <SidebarProvider>
      <div className="flex h-screen w-full overflow-hidden">
        <Sidebar
          collapsible="none"
          className="hidden md:flex h-full flex-shrink-0 border-r"
          variant="inset"
        >
          <SidebarHeader className="gap-3 border-b p-2 mb-2 flex-shrink-0">
            <div className="flex items-center gap-2">
              <FileText className="h-5 w-5 text-muted-foreground" />
              <h2 className="text-lg font-semibold">API Documentation</h2>
            </div>
            <div className="flex items-center gap-2">
              <SidebarInput
                placeholder="Search documentation..."
                value={searchTerm}
                onChange={handleSearch}
                className="flex-1 rounded-lg"
              />
            </div>
          </SidebarHeader>
          <SidebarContent className="flex-1 min-h-0 flex flex-col">
            <SidebarGroup className="p-0 flex-1 overflow-hidden">
              <SidebarGroupContent className="h-full overflow-y-auto">
                <DocsList
                  tables={loaderData.tables}
                  projectId={projectId}
                  searchTerm={searchTerm}
                />
              </SidebarGroupContent>
            </SidebarGroup>
          </SidebarContent>
        </Sidebar>
        <SidebarInset className="flex-1 overflow-hidden">
          <div className="h-full overflow-auto">
            <Outlet context={{ projectDetails, services, openApiSpec: loaderData.openApiSpec }} />
          </div>
        </SidebarInset>
      </div>
    </SidebarProvider>
  );
}