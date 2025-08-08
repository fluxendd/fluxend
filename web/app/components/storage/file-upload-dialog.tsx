import { useState, useCallback } from "react";
import { Upload, FileIcon } from "lucide-react";
import {
  Dialog,
  DialogContent,
  DialogDescription,
  DialogHeader,
  DialogTitle,
} from "~/components/ui/dialog";
import { Button } from "~/components/ui/button";
import type { StorageContainer } from "~/types/storage";
import { formatBytes } from "~/lib/utils";

interface FileUploadDialogProps {
  open: boolean;
  onOpenChange: (open: boolean) => void;
  container: StorageContainer;
  onUpload: (file: File) => Promise<void>;
}

export function FileUploadDialog({
  open,
  onOpenChange,
  container,
  onUpload,
}: FileUploadDialogProps) {
  const [selectedFile, setSelectedFile] = useState<File | null>(null);
  const [isUploading, setIsUploading] = useState(false);
  const [uploadProgress, setUploadProgress] = useState(0);
  const [dragActive, setDragActive] = useState(false);

  const handleDrag = useCallback((e: React.DragEvent) => {
    e.preventDefault();
    e.stopPropagation();
    if (e.type === "dragenter" || e.type === "dragover") {
      setDragActive(true);
    } else if (e.type === "dragleave") {
      setDragActive(false);
    }
  }, []);

  const handleDrop = useCallback((e: React.DragEvent) => {
    e.preventDefault();
    e.stopPropagation();
    setDragActive(false);

    if (e.dataTransfer.files && e.dataTransfer.files[0]) {
      const file = e.dataTransfer.files[0];
      if (file.size > container.maxFileSize) {
        alert(
          `File size exceeds maximum allowed size of ${formatBytes(
            container.maxFileSize
          )}`
        );
        return;
      }
      setSelectedFile(file);
    }
  }, [container.maxFileSize]);

  const handleFileSelect = (e: React.ChangeEvent<HTMLInputElement>) => {
    if (e.target.files && e.target.files[0]) {
      const file = e.target.files[0];
      if (file.size > container.maxFileSize) {
        alert(
          `File size exceeds maximum allowed size of ${formatBytes(
            container.maxFileSize
          )}`
        );
        return;
      }
      setSelectedFile(file);
    }
  };

  const handleUpload = async () => {
    if (!selectedFile) return;

    setIsUploading(true);
    setUploadProgress(0);

    try {
      // Simulate upload progress (in real implementation, this would track actual upload)
      const progressInterval = setInterval(() => {
        setUploadProgress((prev) => {
          if (prev >= 90) {
            clearInterval(progressInterval);
            return 90;
          }
          return prev + 10;
        });
      }, 200);

      await onUpload(selectedFile);
      
      clearInterval(progressInterval);
      setUploadProgress(100);
      
      // Reset state immediately after successful upload
      setSelectedFile(null);
      setIsUploading(false);
      setUploadProgress(0);
    } catch (error) {
      setIsUploading(false);
      setUploadProgress(0);
    }
  };

  const handleCancel = () => {
    if (!isUploading) {
      setSelectedFile(null);
      onOpenChange(false);
    }
  };

  // Reset state when dialog closes
  const handleOpenChange = (newOpen: boolean) => {
    if (!isUploading && !newOpen) {
      setSelectedFile(null);
      setUploadProgress(0);
    }
    onOpenChange(newOpen);
  };

  return (
    <Dialog open={open} onOpenChange={handleOpenChange}>
      <DialogContent className="sm:max-w-xl">
        <DialogHeader>
          <DialogTitle>Upload File</DialogTitle>
          <DialogDescription>
            Upload a file to <span className="font-semibold text-foreground">{container.name}</span> (Max: {formatBytes(container.maxFileSize)})
          </DialogDescription>
        </DialogHeader>

        <div className="space-y-4">
          {!selectedFile ? (
            <div
              className={`border-2 border-dashed rounded-lg p-8 text-center transition-colors ${
                dragActive
                  ? "border-primary bg-primary/5"
                  : "border-muted-foreground/25"
              }`}
              onDragEnter={handleDrag}
              onDragLeave={handleDrag}
              onDragOver={handleDrag}
              onDrop={handleDrop}
            >
              <Upload className="mx-auto h-12 w-12 text-muted-foreground" />
              <p className="mt-2 text-sm text-muted-foreground">
                Drag and drop a file here, or click to select
              </p>
              <input
                type="file"
                className="hidden"
                id="file-upload"
                onChange={handleFileSelect}
                disabled={isUploading}
              />
              <Button
                variant="outline"
                className="mt-4"
                onClick={() => document.getElementById("file-upload")?.click()}
                disabled={isUploading}
              >
                Select File
              </Button>
            </div>
          ) : (
            <div className="space-y-4 w-full">
              <div className="flex items-center gap-3 p-4 border rounded-lg overflow-hidden w-full">
                <FileIcon className="h-8 w-8 text-muted-foreground flex-shrink-0" />
                <div className="min-w-0 flex-1">
                  <p className="font-medium break-all" title={selectedFile.name}>{selectedFile.name}</p>
                  <p className="text-sm text-muted-foreground">
                    {formatBytes(selectedFile.size)}
                  </p>
                </div>
                {!isUploading && (
                  <Button
                    variant="ghost"
                    size="sm"
                    onClick={() => setSelectedFile(null)}
                    className="flex-shrink-0 ml-2"
                  >
                    Remove
                  </Button>
                )}
              </div>

              {isUploading && (
                <div className="space-y-2">
                  <div className="flex justify-between text-sm">
                    <span>Uploading...</span>
                    <span>{uploadProgress}%</span>
                  </div>
                  <div className="w-full bg-secondary rounded-full h-2">
                    <div
                      className="bg-primary h-2 rounded-full transition-all"
                      style={{ width: `${uploadProgress}%` }}
                    />
                  </div>
                </div>
              )}

              <div className="flex justify-end gap-2">
                <Button
                  variant="outline"
                  onClick={handleCancel}
                  disabled={isUploading}
                >
                  Cancel
                </Button>
                <Button onClick={handleUpload} disabled={isUploading}>
                  {isUploading ? "Uploading..." : "Upload"}
                </Button>
              </div>
            </div>
          )}
        </div>
      </DialogContent>
    </Dialog>
  );
}