import { useCallback, useMemo, useState, useEffect } from "react";
import { useOutletContext, useParams } from "react-router";
import { useQueryClient } from "@tanstack/react-query";
import type { Route } from "./+types/$containerId";
import type { ProjectLayoutOutletContext } from "~/components/shared/project-layout";
import type { StorageContainer } from "~/types/storage";
import { RefreshButton } from "~/components/shared/refresh-button";
import { DataTableSkeleton } from "~/components/shared/data-table-skeleton";
import { Button } from "~/components/ui/button";
import { Plus, Upload, Grid3X3, List } from "lucide-react";
import { ContainerList } from "~/components/storage/container-list";
import { FileList } from "~/components/storage/file-list";
import { FileUploadDialog } from "~/components/storage/file-upload-dialog";
import { cn } from "~/lib/utils";

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
  const [viewMode, setViewMode] = useState<"grid" | "table">("grid");

  // Load view preference from localStorage
  useEffect(() => {
    const savedView = localStorage.getItem("storage-view-mode");
    if (savedView === "table" || savedView === "grid") {
      setViewMode(savedView);
    }
  }, []);

  // Save view preference to localStorage
  const handleViewModeChange = useCallback((value: string) => {
    if (value === "grid" || value === "table") {
      setViewMode(value);
      localStorage.setItem("storage-view-mode", value);
    }
  }, []);

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
              <>
                <div className="inline-flex h-8 items-center justify-center rounded-md bg-muted p-1 text-muted-foreground">
                  <Button
                    variant="ghost"
                    size="sm"
                    onClick={() => handleViewModeChange("grid")}
                    className={cn(
                      "h-6 px-2 py-1 text-xs",
                      viewMode === "grid" && "bg-background text-foreground shadow-sm"
                    )}
                  >
                    <Grid3X3 className="h-3 w-3" />
                  </Button>
                  <Button
                    variant="ghost"
                    size="sm"
                    onClick={() => handleViewModeChange("table")}
                    className={cn(
                      "h-6 px-2 py-1 text-xs",
                      viewMode === "table" && "bg-background text-foreground shadow-sm"
                    )}
                  >
                    <List className="h-3 w-3" />
                  </Button>
                </div>
                <Button
                  size="sm"
                  onClick={() => setUploadDialogOpen(true)}
                >
                  <Upload className="h-4 w-4 mr-1" />
                  Upload File
                </Button>
              </>
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
      <div className="flex-1 overflow-auto p-4 min-h-0">
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
            viewMode={viewMode}
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