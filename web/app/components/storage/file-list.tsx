import { useState, useCallback, useMemo, useEffect } from "react";
import { useQuery, useQueryClient } from "@tanstack/react-query";
import {
  FileIcon,
  Download,
  Trash2,
  Edit,
  MoreVertical,
  Upload,
  FileText,
  FileImage,
  FileVideo,
  FileAudio,
  FileArchive,
  FileCode,
  Grid3X3,
  List,
} from "lucide-react";
import { Button } from "~/components/ui/button";
import {
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
} from "~/components/ui/table";
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
  AlertDialogTrigger,
} from "~/components/ui/alert-dialog";
import { DataTableSkeleton } from "~/components/shared/data-table-skeleton";
import { toast } from "sonner";
import type { StorageContainer, StorageFile } from "~/types/storage";
import type { Services } from "~/services";
import { formatBytes, cn } from "~/lib/utils";
import { RenameFileDialog } from "~/components/storage/rename-file-dialog";
import { FileUploadDialog } from "~/components/storage/file-upload-dialog";
import { FileGrid } from "~/components/storage/file-grid";

interface FileListProps {
  projectId: string;
  container: StorageContainer;
  services: Services;
  uploadDialogOpen: boolean;
  setUploadDialogOpen: (open: boolean) => void;
  viewMode: "grid" | "table";
}

const getFileIcon = (mimeType: string) => {
  if (mimeType.startsWith("image/")) return FileImage;
  if (mimeType.startsWith("video/")) return FileVideo;
  if (mimeType.startsWith("audio/")) return FileAudio;
  if (mimeType.includes("zip") || mimeType.includes("tar")) return FileArchive;
  if (mimeType.includes("text") || mimeType.includes("document")) return FileText;
  if (
    mimeType.includes("javascript") ||
    mimeType.includes("json") ||
    mimeType.includes("xml") ||
    mimeType.includes("html")
  )
    return FileCode;
  return FileIcon;
};

export function FileList({ projectId, container, services, uploadDialogOpen, setUploadDialogOpen, viewMode }: FileListProps) {
  const queryClient = useQueryClient();
  const [renameFile, setRenameFile] = useState<StorageFile | null>(null);

  // Fetch files
  const {
    isLoading,
    data: filesData,
    error,
  } = useQuery({
    queryKey: ["storage-files", projectId, container.uuid],
    queryFn: () => services.storage.listFiles(projectId, container.uuid),
    enabled: !!projectId && !!container.uuid,
  });

  const files = useMemo(() => {
    return filesData?.content || [];
  }, [filesData]);

  const handleDelete = useCallback(async (fileId: string) => {
    try {
      const response = await services.storage.deleteFile(
        projectId,
        container.uuid,
        fileId
      );

      if (response.ok) {
        toast.success("File deleted successfully");
        // Invalidate both files and containers to update file count
        await Promise.all([
          queryClient.invalidateQueries({
            queryKey: ["storage-files", projectId, container.uuid],
          }),
          queryClient.invalidateQueries({
            queryKey: ["storage-containers", projectId],
          }),
        ]);
      } else {
        toast.error("Failed to delete file");
      }
    } catch (error) {
      toast.error("Failed to delete file");
    }
  }, [projectId, container.uuid, services.storage, queryClient]);

  const handleRename = useCallback(
    async (newName: string) => {
      if (!renameFile) return;

      try {
        const response = await services.storage.renameFile(
          projectId,
          container.uuid,
          renameFile.uuid,
          {
            full_file_name: newName,
            projectUUID: projectId,
          }
        );

        if (response.success) {
          toast.success("File renamed successfully");
          await queryClient.invalidateQueries({
            queryKey: ["storage-files", projectId, container.uuid],
          });
          setRenameFile(null);
        } else {
          toast.error(response.errors?.[0] || "Failed to rename file");
        }
      } catch (error) {
        toast.error("Failed to rename file");
      }
    },
    [renameFile, projectId, container.uuid, services.storage, queryClient]
  );

  const handleUpload = useCallback(
    async (file: File) => {
      try {
        const response = await services.storage.uploadFile(
          projectId,
          container.uuid,
          file
        );

        if (response.success) {
          toast.success("File uploaded successfully");
          // Invalidate both files and containers to update file count
          await Promise.all([
            queryClient.invalidateQueries({
              queryKey: ["storage-files", projectId, container.uuid],
            }),
            queryClient.invalidateQueries({
              queryKey: ["storage-containers", projectId],
            }),
          ]);
          setUploadDialogOpen(false);
        } else {
          toast.error(response.errors?.[0] || "Failed to upload file");
        }
      } catch (error) {
        toast.error("Failed to upload file");
      }
    },
    [projectId, container.uuid, services.storage, queryClient, setUploadDialogOpen]
  );

  const handleDownload = useCallback(
    async (file: StorageFile) => {
      try {
        const response = await services.storage.downloadFile(
          projectId,
          container.uuid,
          file.uuid,
          file.fullFileName
        );

        if (!response.success) {
          toast.error(response.errors?.[0] || "Failed to download file");
        }
      } catch (error) {
        toast.error("Failed to download file");
      }
    },
    [projectId, container.uuid, services.storage]
  );

  if (error) {
    return (
      <div className="flex items-center justify-center h-full">
        <div className="text-destructive">
          Error loading files: {error.message}
        </div>
      </div>
    );
  }

  return (
    <>
      <div className="h-full flex flex-col">
        {isLoading ? (
          <div className="rounded-lg border overflow-hidden">
            <DataTableSkeleton columns={5} rows={10} />
          </div>
        ) : files.length === 0 ? (
          <div className="rounded-lg border p-8 overflow-hidden">
            <div className="text-center">
              <FileIcon className="mx-auto h-12 w-12 text-muted-foreground" />
              <h3 className="mt-2 text-sm font-semibold">No files</h3>
              <p className="mt-1 text-sm text-muted-foreground">
                Get started by uploading a file to this container.
              </p>
              <div className="mt-6">
                <Button onClick={() => setUploadDialogOpen(true)}>
                  <Upload className="h-4 w-4 mr-2" />
                  Upload File
                </Button>
              </div>
            </div>
          </div>
        ) : viewMode === "table" ? (
          <div className="rounded-lg border overflow-hidden">
            <div className="overflow-auto">
              <Table>
              <TableHeader>
                <TableRow>
                  <TableHead>Name</TableHead>
                  <TableHead>Type</TableHead>
                  <TableHead>Size</TableHead>
                  <TableHead>Modified</TableHead>
                  <TableHead className="w-[50px]"></TableHead>
                </TableRow>
              </TableHeader>
              <TableBody>
                {files.map((file) => {
                  const FileIconComponent = getFileIcon(file.mimeType);
                  return (
                    <TableRow key={file.uuid}>
                      <TableCell>
                        <div className="flex items-center gap-2">
                          <FileIconComponent className="h-4 w-4 text-muted-foreground" />
                          <span className="font-medium">
                            {file.fullFileName}
                          </span>
                        </div>
                      </TableCell>
                      <TableCell className="text-muted-foreground">
                        {file.mimeType}
                      </TableCell>
                      <TableCell className="text-muted-foreground">
                        {formatBytes(file.size)}
                      </TableCell>
                      <TableCell className="text-muted-foreground">
                        {new Date(file.updatedAt).toLocaleString()}
                      </TableCell>
                      <TableCell>
                        <DropdownMenu>
                          <DropdownMenuTrigger asChild>
                            <Button variant="ghost" size="icon" className="h-8 w-8">
                              <MoreVertical className="h-4 w-4" />
                            </Button>
                          </DropdownMenuTrigger>
                          <DropdownMenuContent align="end">
                            <DropdownMenuItem
                              onClick={() => handleDownload(file)}
                            >
                              <Download className="h-4 w-4 mr-2" />
                              Download
                            </DropdownMenuItem>
                            <DropdownMenuItem
                              onClick={() => setRenameFile(file)}
                            >
                              <Edit className="h-4 w-4 mr-2" />
                              Rename
                            </DropdownMenuItem>
                            <DropdownMenuSeparator />
                            <AlertDialog>
                              <AlertDialogTrigger asChild>
                                <DropdownMenuItem
                                  className="text-destructive"
                                  onSelect={(e) => e.preventDefault()}
                                >
                                  <Trash2 className="h-4 w-4 mr-2" />
                                  Delete
                                </DropdownMenuItem>
                              </AlertDialogTrigger>
                              <AlertDialogContent>
                                <AlertDialogHeader>
                                  <AlertDialogTitle>Delete File</AlertDialogTitle>
                                  <AlertDialogDescription>
                                    Are you sure you want to delete "{file.fullFileName}"? This action cannot be undone.
                                  </AlertDialogDescription>
                                </AlertDialogHeader>
                                <AlertDialogFooter>
                                  <AlertDialogCancel>Cancel</AlertDialogCancel>
                                  <AlertDialogAction
                                    onClick={() => handleDelete(file.uuid)}
                                    className="bg-destructive text-destructive-foreground hover:bg-destructive/90"
                                  >
                                    Delete
                                  </AlertDialogAction>
                                </AlertDialogFooter>
                              </AlertDialogContent>
                            </AlertDialog>
                          </DropdownMenuContent>
                        </DropdownMenu>
                      </TableCell>
                    </TableRow>
                  );
                })}
              </TableBody>
              </Table>
            </div>
          </div>
        ) : (
          <div className="overflow-auto">
            <FileGrid
              files={files}
              onRename={setRenameFile}
              onDelete={handleDelete}
              onDownload={handleDownload}
            />
          </div>
        )}
      </div>


      {/* Rename File Dialog */}
      {renameFile && (
        <RenameFileDialog
          open={true}
          onOpenChange={(open) => !open && setRenameFile(null)}
          file={renameFile}
          onSubmit={handleRename}
        />
      )}

      {/* Upload File Dialog */}
      {uploadDialogOpen && (
        <FileUploadDialog
          open={true}
          onOpenChange={setUploadDialogOpen}
          container={container}
          onUpload={handleUpload}
        />
      )}
    </>
  );
}