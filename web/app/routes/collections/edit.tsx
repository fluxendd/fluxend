import { useState } from "react";
import {
  useParams,
  useNavigate,
  useRouteLoaderData,
  useOutletContext,
} from "react-router";
import { AppHeader } from "~/components/shared/header";
import { Button } from "~/components/ui/button";
import { Input } from "~/components/ui/input";
import { Label } from "~/components/ui/label";
import {
  Card,
  CardContent,
  CardDescription,
  CardHeader,
  CardTitle,
} from "~/components/ui/card";
import type { ProjectLayoutOutletContext } from "~/components/shared/project-layout";

export default function EditCollection() {
  const { tableName, projectId } = useParams();
  const navigate = useNavigate();
  const [isDeleting, setIsDeleting] = useState(false);
  const { services } = useOutletContext<ProjectLayoutOutletContext>();

  const handleDeleteTable = async () => {
    if (!tableName || !projectId) return;

    const confirmDelete = window.confirm(
      `Are you sure you want to delete the table "${tableName}"? This action cannot be undone.`
    );

    if (!confirmDelete) return;

    setIsDeleting(true);
    try {
      const response = await services.collections.deleteCollection(
        projectId,
        tableName
      );

      if (response.ok) {
        // Navigate back to collections list after successful deletion
        navigate("/collections");
      } else {
        const errorData = await response.json().catch(() => ({}));
        throw new Error(errorData.message || "Failed to delete table");
      }
    } catch (error) {
      console.error("Error deleting table:", error);
      // You might want to show a toast notification here
      alert("Failed to delete table. Please try again.");
    } finally {
      setIsDeleting(false);
    }
  };

  return (
    <div className="flex flex-col h-full">
      <AppHeader
        title={`Edit Collection: ${tableName || "Unknown"}`}
        showDelete={true}
        onDelete={handleDeleteTable}
        deleteLabel={isDeleting ? "Deleting..." : "Delete Table"}
      />

      <div className="flex-1 p-6">
        <div className="max-w-4xl mx-auto">
          <Card>
            <CardHeader>
              <CardTitle>Edit Collection</CardTitle>
              <CardDescription>
                Modify your collection settings and structure
              </CardDescription>
            </CardHeader>
            <CardContent className="space-y-4">
              <div className="space-y-2">
                <Label htmlFor="tableName">Table Name</Label>
                <Input
                  id="tableName"
                  value={tableName || ""}
                  disabled
                  placeholder="Table name"
                />
              </div>

              <div className="pt-4">
                <p className="text-sm text-muted-foreground">
                  Collection editing features are coming soon. For now, you can
                  delete this collection using the delete button in the header.
                </p>
              </div>

              <div className="flex justify-end space-x-2">
                <Button
                  variant="outline"
                  onClick={() => navigate("/collections")}
                >
                  Back to Collections
                </Button>
                <Button disabled>Save Changes</Button>
              </div>
            </CardContent>
          </Card>
        </div>
      </div>
    </div>
  );
}
