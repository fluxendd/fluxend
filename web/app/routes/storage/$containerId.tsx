import { useCallback, useMemo, useState } from "react";
import { useOutletContext, useParams } from "react-router";
import { useQueryClient } from "@tanstack/react-query";
import type { Route } from "./+types/$containerId";
import type { ProjectLayoutOutletContext } from "~/components/shared/project-layout";
import type { StorageContainer } from "~/types/storage";
import { RefreshButton } from "~/components/shared/refresh-button";
import { DataTableSkeleton } from "~/components/shared/data-table-skeleton";
import { Button } from "~/components/ui/button";
import { Plus, Upload } from "lucide-react";
import { ContainerList } from "~/components/storage/container-list";
import { FileList } from "~/components/storage/file-list";
import { FileUploadDialog } from "~/components/storage/file-upload-dialog";

interface StorageLayoutContext extends ProjectLayoutOutletContext {
  containers: StorageContainer[];
}

export function meta({}: Route.MetaArgs) {
  return [
    { title: "Storage - Fluxend" },
    { name: "description", content: "Manage your storage containers and files" },
  ];
}

export default function StorageContainer() {
  const { projectDetails, services, containers } = useOutletContext<StorageLayoutContext>();
  const projectId = projectDetails?.uuid;
  const { containerId } = useParams();
  const queryClient = useQueryClient();
  const [uploadDialogOpen, setUploadDialogOpen] = useState(false);

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

  return (
    <div className="flex flex-col h-full">
      <div className="border-b px-4 py-2 flex-shrink-0">
        <div className="flex items-center justify-between">
          <div className="text-base font-bold text-foreground h-[32px] flex flex-col justify-center">
            Storage {activeContainer && `/ ${activeContainer.name}`}
          </div>
          <div className="flex items-center gap-2">
            {containerId && activeContainer && (
              <Button
                size="sm"
                onClick={() => setUploadDialogOpen(true)}
              >
                <Upload className="h-4 w-4 mr-1" />
                Upload File
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
        {!containers.length ? (
          <div className="rounded-lg border h-full overflow-hidden">
            <DataTableSkeleton columns={5} rows={8} />
          </div>
        ) : containerId && activeContainer ? (
          <FileList
            projectId={projectId!}
            container={activeContainer}
            services={services}
            uploadDialogOpen={uploadDialogOpen}
            setUploadDialogOpen={setUploadDialogOpen}
          />
        ) : !containerId && containers.length > 0 ? (
          <ContainerList
            containers={containers}
            projectId={projectId!}
            onContainerDeleted={handleRefresh}
            services={services}
          />
        ) : null}
      </div>
    </div>
  );
}