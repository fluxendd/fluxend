import { useCallback, useState, useEffect } from "react";
import { useOutletContext, useParams } from "react-router";
import { useQueryClient } from "@tanstack/react-query";
import type { Route } from "./+types/$containerId";
import type { ProjectLayoutOutletContext } from "~/components/shared/project-layout";
import type { StorageContainer } from "~/types/storage";
import { RefreshButton } from "~/components/shared/refresh-button";
import { DataTableSkeleton } from "~/components/shared/data-table-skeleton";
import { Button } from "~/components/ui/button";
import { Upload, Grid3X3, List } from "lucide-react";
import { FileList } from "~/components/storage/file-list";
import { cn } from "~/lib/utils";

interface StorageLayoutContext extends ProjectLayoutOutletContext {
  containers: StorageContainer[];
  isContainersLoading: boolean;
}

export function meta({}: Route.MetaArgs) {
  return [
    { title: "Storage - Fluxend" },
    {
      name: "description",
      content: "Manage your storage containers and files",
    },
  ];
}

export default function StorageContainer() {
  const { projectDetails, services, containers, isContainersLoading } =
    useOutletContext<StorageLayoutContext>();
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

  // Get the current container from the containers list
  const currentContainer = containers.find((c) => c.uuid === containerId);
  const containerName = currentContainer?.name || "";

  // Show loading skeleton while containers are being loaded
  if (isContainersLoading) {
    return (
      <div className="flex flex-col h-full">
        <div className="border-b px-4 py-2 flex-shrink-0">
          <div className="flex items-center justify-between">
            <div className="h-8 w-48 bg-muted animate-pulse rounded" />
            <div className="flex items-center gap-2">
              <div className="h-8 w-20 bg-muted animate-pulse rounded" />
              <div className="h-8 w-24 bg-muted animate-pulse rounded" />
              <div className="h-8 w-8 bg-muted animate-pulse rounded" />
            </div>
          </div>
        </div>
        <div className="flex-1 p-4">
          <DataTableSkeleton columns={4} rows={5} />
        </div>
      </div>
    );
  }

  // Only show "Container not found" after loading is complete
  if (!containerId || !currentContainer) {
    return (
      <div className="flex flex-col h-full">
        <div className="flex-1 flex items-center justify-center">
          <div className="text-muted-foreground">Container not found</div>
        </div>
      </div>
    );
  }

  return (
    <div className="flex flex-col h-full">
      <div className="border-b px-4 py-2 flex-shrink-0">
        <div className="flex items-center justify-between">
          <div className="text-base font-bold text-foreground h-[32px] flex flex-col justify-center">
            {containerName}
          </div>
          <div className="flex items-center gap-2">
            <div className="inline-flex h-8 items-center justify-center rounded-md bg-muted p-1 text-muted-foreground">
              <Button
                variant="ghost"
                size="sm"
                onClick={() => handleViewModeChange("grid")}
                className={cn(
                  "h-6 px-2 py-1 text-xs",
                  viewMode === "grid" &&
                    "bg-background text-foreground shadow-sm"
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
                  viewMode === "table" &&
                    "bg-background text-foreground shadow-sm"
                )}
              >
                <List className="h-3 w-3" />
              </Button>
            </div>
            <Button size="sm" onClick={() => setUploadDialogOpen(true)}>
              <Upload className="h-4 w-4 mr-1" />
              Upload File
            </Button>
            <RefreshButton
              onRefresh={handleRefresh}
              size="sm"
              title="Refresh Files"
            />
          </div>
        </div>
      </div>

      {/* Content */}
      <div className="flex-1 overflow-auto p-4 min-h-0">
        <FileList
          projectId={projectId!}
          container={currentContainer}
          services={services}
          uploadDialogOpen={uploadDialogOpen}
          setUploadDialogOpen={setUploadDialogOpen}
          viewMode={viewMode}
        />
      </div>
    </div>
  );
}

