import { useState, useCallback, useMemo } from "react";
import { useOutletContext, useParams } from "react-router";
import { useQuery, useQueryClient } from "@tanstack/react-query";
import type { Route } from "./+types/page";
import type { ProjectLayoutOutletContext } from "~/components/shared/project-layout";
import { RefreshButton } from "~/components/shared/refresh-button";
import { DataTableSkeleton } from "~/components/shared/data-table-skeleton";
import { Button } from "~/components/ui/button";
import { Plus } from "lucide-react";
import { StorageSidebar } from "./sidebar";
import { ContainerList } from "./container-list";
import { FileList } from "./file-list";
import { CreateContainerDialog } from "./create-container-dialog";
import { toast } from "sonner";
import type { StorageContainer } from "~/types/storage";
import { SidebarProvider, SidebarInset } from "~/components/ui/sidebar";

export function meta({}: Route.MetaArgs) {
  return [
    { title: "Storage - Fluxend" },
    { name: "description", content: "Manage your storage containers and files" },
  ];
}

export default function Storage() {
  const { projectDetails, services } = useOutletContext<ProjectLayoutOutletContext>();
  const projectId = projectDetails?.uuid;
  const { containerId } = useParams();
  const queryClient = useQueryClient();

  const [selectedContainer, setSelectedContainer] = useState<StorageContainer | null>(null);
  const [createContainerOpen, setCreateContainerOpen] = useState(false);

  // Fetch containers
  const {
    isLoading: isContainersLoading,
    data: containersData,
    error: containersError,
  } = useQuery({
    queryKey: ["storage-containers", projectId],
    queryFn: () => services.storage.listContainers(projectId!),
    enabled: !!projectId,
  });

  const containers = useMemo(() => {
    return containersData?.content || [];
  }, [containersData]);

  // Find selected container from containerId param
  const activeContainer = useMemo(() => {
    if (!containerId || !containers.length) return null;
    return containers.find(c => c.uuid === containerId) || null;
  }, [containerId, containers]);

  const handleRefresh = useCallback(async () => {
    await queryClient.invalidateQueries({
      queryKey: ["storage-containers", projectId],
    });
    if (containerId) {
      await queryClient.invalidateQueries({
        queryKey: ["storage-files", projectId, containerId],
      });
    }
  }, [queryClient, projectId, containerId]);

  const handleCreateContainer = useCallback(async (container: {
    name: string;
    description?: string;
    is_public: boolean;
    max_file_size: number;
  }) => {
    if (!projectId) return;

    try {
      const response = await services.storage.createContainer(projectId, {
        ...container,
        projectUUID: projectId,
      });

      if (response.success) {
        await queryClient.invalidateQueries({
          queryKey: ["storage-containers", projectId],
        });
        toast.success("Container created successfully");
        setCreateContainerOpen(false);
      } else {
        toast.error(response.errors?.[0] || "Failed to create container");
      }
    } catch (error) {
      toast.error("Failed to create container");
    }
  }, [projectId, services.storage, queryClient]);

  if (containersError) {
    return (
      <div className="flex flex-col h-full">
        <div className="border-b px-4 py-2 flex-shrink-0">
          <div className="flex items-center justify-between">
            <div className="text-base font-bold text-foreground h-[32px] flex flex-col justify-center">
              Storage
            </div>
          </div>
        </div>
        <div className="flex-1 flex items-center justify-center">
          <div className="text-destructive">
            Error loading containers: {containersError.message}
          </div>
        </div>
      </div>
    );
  }

  return (
    <SidebarProvider>
      <div className="flex h-screen w-full overflow-hidden">
        {/* Sidebar with containers list */}
        <StorageSidebar
          containers={containers}
          activeContainerId={containerId}
          isLoading={isContainersLoading}
          projectId={projectId!}
          onCreateContainer={() => setCreateContainerOpen(true)}
        />

        {/* Main content area */}
        <SidebarInset className="flex-1 overflow-hidden">
          <div className="h-full overflow-auto">
            <div className="flex flex-col h-full">
              <div className="border-b px-4 py-2 flex-shrink-0">
                <div className="flex items-center justify-between">
                  <div className="text-base font-bold text-foreground h-[32px] flex flex-col justify-center">
                    Storage {activeContainer && `/ ${activeContainer.name}`}
                  </div>
                  <div className="flex items-center gap-2">
                    {!containerId && (
                      <Button
                        size="sm"
                        onClick={() => setCreateContainerOpen(true)}
                      >
                        <Plus className="h-4 w-4 mr-1" />
                        New Container
                      </Button>
                    )}
                    <RefreshButton
                      onRefresh={handleRefresh}
                      size="sm"
                      title="Refresh Storage"
                    />
                  </div>
                </div>
              </div>

              {/* Content */}
              <div className="flex-1 overflow-hidden p-4">
                {isContainersLoading && !containers.length ? (
                  <div className="rounded-lg border h-full overflow-hidden">
                    <DataTableSkeleton columns={5} rows={8} />
                  </div>
                ) : containerId && activeContainer ? (
                  <FileList
                    projectId={projectId!}
                    container={activeContainer}
                    services={services}
                  />
                ) : !containerId && containers.length > 0 ? (
                  <ContainerList
                    containers={containers}
                    projectId={projectId!}
                    onContainerDeleted={handleRefresh}
                    services={services}
                  />
                ) : (
                  <div className="flex-1 flex items-center justify-center">
                    <div className="text-center">
                      <p className="text-muted-foreground mb-4">
                        No containers yet. Create your first container to get started.
                      </p>
                      <Button onClick={() => setCreateContainerOpen(true)}>
                        <Plus className="h-4 w-4 mr-2" />
                        Create Container
                      </Button>
                    </div>
                  </div>
                )}
              </div>
            </div>
          </div>
        </SidebarInset>

        {/* Dialogs */}
        <CreateContainerDialog
          open={createContainerOpen}
          onOpenChange={setCreateContainerOpen}
          onSubmit={handleCreateContainer}
        />
      </div>
    </SidebarProvider>
  );
}