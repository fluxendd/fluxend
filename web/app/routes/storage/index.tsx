import { useCallback } from "react";
import { useOutletContext } from "react-router";
import { useQueryClient } from "@tanstack/react-query";
import type { Route } from "./+types/index";
import type { ProjectLayoutOutletContext } from "~/components/shared/project-layout";
import type { StorageContainer } from "~/types/storage";
import { RefreshButton } from "~/components/shared/refresh-button";
import { DataTableSkeleton } from "~/components/shared/data-table-skeleton";
import { Button } from "~/components/ui/button";
import { Plus } from "lucide-react";
import { ContainerList } from "~/components/storage/container-list";

interface StorageLayoutContext extends ProjectLayoutOutletContext {
  containers: StorageContainer[];
  isContainersLoading?: boolean;
  setCreateContainerOpen?: (open: boolean) => void;
}

export function meta({}: Route.MetaArgs) {
  return [
    { title: "Storage - Fluxend" },
    { name: "description", content: "Manage your storage containers" },
  ];
}

export default function StorageIndex() {
  const { projectDetails, services, containers, isContainersLoading, setCreateContainerOpen } = useOutletContext<StorageLayoutContext>();
  const projectId = projectDetails?.uuid;
  const queryClient = useQueryClient();

  const handleRefresh = useCallback(async () => {
    await queryClient.invalidateQueries({
      queryKey: ["storage-containers", projectId],
    });
  }, [queryClient, projectId]);

  return (
    <div className="flex flex-col h-full">
      <div className="border-b px-4 py-2 flex-shrink-0">
        <div className="flex items-center justify-between">
          <div className="text-base font-bold text-foreground h-[32px] flex flex-col justify-center">
            Storage
          </div>
          <div className="flex items-center gap-2">
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
        {isContainersLoading ? (
          <div className="rounded-lg border h-full overflow-hidden">
            <DataTableSkeleton columns={5} rows={8} />
          </div>
        ) : containers.length > 0 ? (
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
              <Button onClick={() => setCreateContainerOpen?.(true)}>
                <Plus className="h-4 w-4 mr-2" />
                Create Container
              </Button>
            </div>
          </div>
        )}
      </div>
    </div>
  );
}