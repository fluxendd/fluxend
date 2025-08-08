import type { Route } from "./+types/sidebar";
import {
  Outlet,
  redirect,
  useNavigate,
  useOutletContext,
  useParams,
} from "react-router";
import { useState, useCallback, useEffect } from "react";
import { SidebarProvider, SidebarInset } from "~/components/ui/sidebar";
import { StorageSidebar } from "~/components/storage/sidebar";
import { useQuery, useQueryClient } from "@tanstack/react-query";
import type { ProjectLayoutOutletContext } from "~/components/shared/project-layout";
import { CreateContainerDialog } from "~/components/storage/create-container-dialog";
import { toast } from "sonner";
import { Button } from "~/components/ui/button";
import { Plus } from "lucide-react";

export default function StorageLayout() {
  const { projectDetails, services } =
    useOutletContext<ProjectLayoutOutletContext>();
  const projectId = projectDetails?.uuid;
  const { containerId } = useParams();
  const navigate = useNavigate();
  const queryClient = useQueryClient();

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

  const containers = containersData?.content || [];

  // Auto-navigate to the first container if none is selected
  useEffect(() => {
    if (
      !containerId &&
      containers.length > 0 &&
      projectId &&
      !isContainersLoading
    ) {
      navigate(`/projects/${projectId}/storage/${containers[0].uuid}`, {
        replace: true,
      });
    }
  }, [containerId, containers, projectId, navigate, isContainersLoading]);

  const handleCreateContainer = useCallback(
    async (container: {
      name: string;
      description?: string;
      is_public: boolean;
      max_file_size: number;
    }) => {
      if (!projectId) return;

      try {
        const response = await services.storage.createContainer(projectId, {
          projectUUID: projectId,
          name: container.name,
          description: container.description || "",
          is_public: container.is_public,
          max_file_size: container.max_file_size,
        });

        if (response.success) {
          await queryClient.invalidateQueries({
            queryKey: ["storage-containers", projectId],
          });
          toast.success("Container created successfully");
          setCreateContainerOpen(false);

          // Navigate to the new container
          if (response.content?.uuid) {
            navigate(`/projects/${projectId}/storage/${response.content.uuid}`);
          }
        } else {
          toast.error(response.errors?.[0] || "Failed to create container");
          // Throw error to prevent form reset in the dialog
          throw new Error(response.errors?.[0] || "Failed to create container");
        }
      } catch (error) {
        toast.error("Failed to create container");
        // Re-throw to prevent form reset in the dialog
        throw error;
      }
    },
    [projectId, services.storage, queryClient, navigate]
  );

  if (containersError) {
    return (
      <div className="flex flex-col h-full">
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
        <StorageSidebar
          containers={containers}
          activeContainerId={containerId}
          isLoading={isContainersLoading}
          projectId={projectId!}
          onCreateContainer={() => setCreateContainerOpen(true)}
        />
        <SidebarInset className="flex-1 overflow-hidden">
          {containers.length === 0 && !isContainersLoading ? (
            <div className="h-full flex items-center justify-center">
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
          ) : (
            <div className="h-full overflow-auto">
              <Outlet
                context={{
                  projectDetails,
                  services,
                  containers,
                  isContainersLoading,
                  setCreateContainerOpen,
                }}
              />
            </div>
          )}
        </SidebarInset>

        <CreateContainerDialog
          open={createContainerOpen}
          onOpenChange={setCreateContainerOpen}
          onSubmit={handleCreateContainer}
        />
      </div>
    </SidebarProvider>
  );
}

