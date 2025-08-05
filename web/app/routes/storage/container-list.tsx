import { useState, useCallback } from "react";
import { useNavigate } from "react-router";
import { Package2, Globe, Lock, Trash2, Edit, MoreVertical } from "lucide-react";
import { Button } from "~/components/ui/button";
import {
  Card,
  CardContent,
  CardDescription,
  CardFooter,
  CardHeader,
  CardTitle,
} from "~/components/ui/card";
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuSeparator,
  DropdownMenuTrigger,
} from "~/components/ui/dropdown-menu";
import {
  AlertDialog,
  AlertDialogAction,
  AlertDialogCancel,
  AlertDialogContent,
  AlertDialogDescription,
  AlertDialogFooter,
  AlertDialogHeader,
  AlertDialogTitle,
} from "~/components/ui/alert-dialog";
import { Badge } from "~/components/ui/badge";
import { toast } from "sonner";
import type { StorageContainer } from "~/types/storage";
import type { Services } from "~/services";
import { formatBytes } from "~/lib/utils";
import { UpdateContainerDialog } from "./update-container-dialog";

interface ContainerListProps {
  containers: StorageContainer[];
  projectId: string;
  onContainerDeleted: () => void;
  services: Services;
}

export function ContainerList({
  containers,
  projectId,
  onContainerDeleted,
  services,
}: ContainerListProps) {
  const navigate = useNavigate();
  const [deleteContainerId, setDeleteContainerId] = useState<string | null>(null);
  const [editContainer, setEditContainer] = useState<StorageContainer | null>(null);

  const handleDelete = useCallback(async () => {
    if (!deleteContainerId) return;

    try {
      const response = await services.storage.deleteContainer(projectId, deleteContainerId);
      
      if (response.ok) {
        toast.success("Container deleted successfully");
        onContainerDeleted();
      } else {
        toast.error("Failed to delete container");
      }
    } catch (error) {
      toast.error("Failed to delete container");
    } finally {
      setDeleteContainerId(null);
    }
  }, [deleteContainerId, projectId, services.storage, onContainerDeleted]);

  const handleUpdate = useCallback(async (updates: {
    name: string;
    description?: string;
    is_public: boolean;
    max_file_size: number;
  }) => {
    if (!editContainer) return;

    try {
      const response = await services.storage.updateContainer(
        projectId,
        editContainer.uuid,
        {
          ...updates,
          projectUUID: projectId,
        }
      );

      if (response.success) {
        toast.success("Container updated successfully");
        onContainerDeleted(); // This will refresh the list
        setEditContainer(null);
      } else {
        toast.error(response.errors?.[0] || "Failed to update container");
      }
    } catch (error) {
      toast.error("Failed to update container");
    }
  }, [editContainer, projectId, services.storage, onContainerDeleted]);

  return (
    <>
      <div className="grid gap-4 md:grid-cols-2 lg:grid-cols-3">
        {containers.map((container) => (
          <Card
            key={container.uuid}
            className="cursor-pointer hover:shadow-md transition-shadow"
            onClick={() => navigate(`/projects/${projectId}/storage/${container.uuid}`)}
          >
            <CardHeader className="pb-3">
              <div className="flex items-start justify-between">
                <div className="flex items-center gap-2">
                  <Package2 className="h-5 w-5 text-muted-foreground" />
                  <CardTitle className="text-lg">{container.name}</CardTitle>
                </div>
                <DropdownMenu>
                  <DropdownMenuTrigger asChild>
                    <Button
                      variant="ghost"
                      size="icon"
                      className="h-8 w-8"
                      onClick={(e) => e.stopPropagation()}
                    >
                      <MoreVertical className="h-4 w-4" />
                    </Button>
                  </DropdownMenuTrigger>
                  <DropdownMenuContent align="end">
                    <DropdownMenuItem
                      onClick={(e) => {
                        e.stopPropagation();
                        setEditContainer(container);
                      }}
                    >
                      <Edit className="h-4 w-4 mr-2" />
                      Edit
                    </DropdownMenuItem>
                    <DropdownMenuSeparator />
                    <DropdownMenuItem
                      className="text-destructive"
                      onClick={(e) => {
                        e.stopPropagation();
                        setDeleteContainerId(container.uuid);
                      }}
                    >
                      <Trash2 className="h-4 w-4 mr-2" />
                      Delete
                    </DropdownMenuItem>
                  </DropdownMenuContent>
                </DropdownMenu>
              </div>
              {container.description && (
                <CardDescription className="mt-1">
                  {container.description}
                </CardDescription>
              )}
            </CardHeader>
            <CardContent className="pb-3">
              <div className="flex items-center gap-4 text-sm text-muted-foreground">
                <div className="flex items-center gap-1">
                  {container.isPublic ? (
                    <>
                      <Globe className="h-3 w-3" />
                      <span>Public</span>
                    </>
                  ) : (
                    <>
                      <Lock className="h-3 w-3" />
                      <span>Private</span>
                    </>
                  )}
                </div>
                <div>
                  {container.totalFiles} file{container.totalFiles !== 1 ? 's' : ''}
                </div>
              </div>
            </CardContent>
            <CardFooter className="pt-3 border-t">
              <div className="flex items-center justify-between w-full text-xs text-muted-foreground">
                <span>Max file size: {formatBytes(container.maxFileSize)}</span>
                <Badge variant="secondary" className="text-xs">
                  {new Date(container.updatedAt).toLocaleDateString()}
                </Badge>
              </div>
            </CardFooter>
          </Card>
        ))}
      </div>

      {/* Delete Confirmation Dialog */}
      <AlertDialog
        open={!!deleteContainerId}
        onOpenChange={(open) => !open && setDeleteContainerId(null)}
      >
        <AlertDialogContent>
          <AlertDialogHeader>
            <AlertDialogTitle>Delete Container</AlertDialogTitle>
            <AlertDialogDescription>
              Are you sure you want to delete this container? This action cannot be
              undone and all files within the container will be permanently deleted.
            </AlertDialogDescription>
          </AlertDialogHeader>
          <AlertDialogFooter>
            <AlertDialogCancel>Cancel</AlertDialogCancel>
            <AlertDialogAction
              onClick={handleDelete}
              className="bg-destructive text-destructive-foreground hover:bg-destructive/90"
            >
              Delete
            </AlertDialogAction>
          </AlertDialogFooter>
        </AlertDialogContent>
      </AlertDialog>

      {/* Update Container Dialog */}
      {editContainer && (
        <UpdateContainerDialog
          open={!!editContainer}
          onOpenChange={(open) => !open && setEditContainer(null)}
          container={editContainer}
          onSubmit={handleUpdate}
        />
      )}
    </>
  );
}