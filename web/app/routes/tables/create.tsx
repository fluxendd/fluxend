import { useState } from "react";
import { AppHeader } from "~/components/shared/header";
import { useParams, useNavigate, useOutletContext } from "react-router";
import { TableForm } from "~/components/shared/table-form";
import type { CreateTableRequest } from "~/types/table";
import { queryClient } from "~/lib/query";
import type { ProjectLayoutOutletContext } from "~/components/shared/project-layout";

export default function CreateTable() {
  const navigate = useNavigate();
  const { projectId } = useParams();
  const [isSubmitting, setIsSubmitting] = useState(false);
  const { services } = useOutletContext<ProjectLayoutOutletContext>();

  const onSubmit = async (data: { tableName: string; columns: any[] }) => {
    if (!projectId) {
      alert("Project ID is required");
      return;
    }

    setIsSubmitting(true);
    try {
      const requestBody: CreateTableRequest = {
        name: data.tableName,
        columns: data.columns,
      };

      const response = await services.tables.createTable(
        projectId,
        requestBody
      );

      if (response.ok) {
        const responseData = await response.json();
        const newTableName = responseData.content?.name || data.tableName;

        // Invalidate tables query to refresh the sidebar
        await queryClient.invalidateQueries({
          queryKey: ["tables", projectId],
        });

        // Redirect to the new table
        navigate(`/projects/${projectId}/tables/${newTableName}`);
      } else {
        const errorData = await response.json().catch(() => ({}));
        throw new Error(errorData.message || "Failed to create table");
      }
    } catch (error) {
      console.error("Error creating table:", error);
      alert("Failed to create table. Please try again.");
    } finally {
      setIsSubmitting(false);
    }
  };

  return (
    <div className="flex flex-col h-full">
      <AppHeader title="Create Table" />
      <div className="flex-1 p-6">
        <div className="max-w-4xl mx-auto">
          <TableForm
            mode="create"
            onSubmit={onSubmit}
            isSubmitting={isSubmitting}
          />
        </div>
      </div>
    </div>
  );
}