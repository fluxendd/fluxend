import { useState, useEffect } from "react";
import {
  useParams,
  useNavigate,
  useOutletContext,
} from "react-router";
import { AppHeader } from "~/components/shared/header";
import { TableForm } from "~/components/shared/table-form";
import { toast } from "sonner";
import { useQuery } from "@tanstack/react-query";
import { DataTableSkeleton } from "~/components/shared/data-table-skeleton";
import type { ProjectLayoutOutletContext } from "~/components/shared/project-layout";
import type { ColumnType } from "~/types/table";
import { queryClient } from "~/lib/query";
import type {Route} from "../../../.react-router/types/app/routes/tables/+types/page";

interface Column {
  name: string;
  position: number;
  notNull: boolean;
  type: string;
  defaultValue: string;
  primary: boolean;
  unique: boolean;
  foreign: boolean;
  referenceTable: string | null;
  referenceColumn: string | null;
}

// Map PostgreSQL data types to our column types
const mapPostgresToColumnType = (pgType: string): ColumnType => {
  // Extract base type (remove size specifications like (255))
  const baseType = pgType.toLowerCase().split('(')[0].trim();
  
  // Map PostgreSQL types to our simplified types
  const typeMap: Record<string, ColumnType> = {
    'integer': 'integer',
    'serial': 'serial',
    'character varying': 'varchar',
    'varchar': 'varchar',
    'text': 'text',
    'boolean': 'boolean',
    'date': 'date',
    'timestamp': 'timestamp',
    'timestamp without time zone': 'timestamp',
    'timestamp with time zone': 'timestamp',
    'float': 'float',
    'double precision': 'float',
    'real': 'float',
    'uuid': 'uuid',
    'json': 'json',
    'jsonb': 'json'
  };

  return typeMap[baseType] || 'text';
};

export function meta({}: Route.MetaArgs) {
  return [
    { title: "Edit Table - Fluxend" },
    { name: "description", content: "Edit your table structure" },
  ];
}

export default function EditTable() {
  const { tableId, projectId } = useParams();
  const navigate = useNavigate();
  const [isSubmitting, setIsSubmitting] = useState(false);
  const { services } = useOutletContext<ProjectLayoutOutletContext>();

  // Fetch existing table columns
  const { data: columnsData, isLoading, error } = useQuery({
    queryKey: ["columns", projectId, tableId],
    queryFn: async () => {
      if (!projectId || !tableId) return null;
      const response = await services.tables.getTableColumns(projectId, tableId);
      if (!response.ok) {
        throw new Error("Failed to fetch table columns");
      }
      const data = await response.json();
      return data.content;
    },
    enabled: !!projectId && !!tableId,
  });

  const handleSubmit = async (data: { tableName: string; columns: any[] }) => {
    if (!projectId || !tableId) {
      toast.error("Project ID and Table ID are required");
      return;
    }

    setIsSubmitting(true);
    try {
      // Get current columns to compare
      const currentResponse = await services.tables.getTableColumns(projectId, tableId);
      if (!currentResponse.ok) {
        throw new Error("Failed to fetch current table columns");
      }
      const currentColumnsData = await currentResponse.json();
      const currentColumns = currentColumnsData.content;

      // Extract current column names for comparison
      const currentColumnNames = new Set(currentColumns.map((col: Column) => col.name));

      // Separate new columns from existing ones
      const newColumns = data.columns.filter(col => !currentColumnNames.has(col.name));

      // Create new columns first if any exist
      if (newColumns.length > 0) {
        const createRequestBody = {
          columns: newColumns.map(col => ({
            name: col.name,
            type: col.type,
          })),
        };

        const createResponse = await services.tables.createTableColumns(
            projectId,
            tableId,
            createRequestBody
        );

        if (!createResponse.ok) {
          const errorData = await createResponse.json().catch(() => ({}));
          throw new Error(errorData.message || "Failed to create new columns");
        }
      }

      // Update with ALL columns (existing + new) - this will also handle deletions of columns not in the request
      const updateRequestBody = {
        columns: data.columns.map(col => ({
          name: col.name,
          type: col.type,
        })),
      };

      const updateResponse = await services.tables.updateTableColumns(
          projectId,
          tableId,
          updateRequestBody
      );

      if (!updateResponse.ok) {
        const errorData = await updateResponse.json().catch(() => ({}));
        throw new Error(errorData.message || "Failed to update table columns");
      }

      toast.success("Table updated successfully");

      // Invalidate queries to refresh data
      await queryClient.invalidateQueries({
        queryKey: ["columns", projectId, tableId],
      });
      await queryClient.invalidateQueries({
        queryKey: ["tables", projectId],
      });

      // Navigate back to the table view
      navigate(`/projects/${projectId}/tables/${tableId}`);

    } catch (error) {
      console.error("Error updating table:", error);
      toast.error("Failed to update table. Please try again.");
    } finally {
      setIsSubmitting(false);
    }
  };

  if (isLoading) {
    return (
      <div className="flex flex-col h-full">
        <AppHeader title={`Edit Table: ${tableId || "Loading..."}`} />
        <div className="flex-1 p-6">
          <div className="max-w-4xl mx-auto rounded-lg border p-4">
            <DataTableSkeleton columns={3} rows={5} />
          </div>
        </div>
      </div>
    );
  }

  if (error || !columnsData) {
    return (
      <div className="flex flex-col h-full">
        <AppHeader title={`Edit Table: ${tableId || "Unknown"}`} />
        <div className="flex-1 p-6">
          <div className="max-w-4xl mx-auto">
            <div className="text-center p-8 border rounded-lg bg-destructive/10">
              <p className="text-destructive">
                Failed to load table columns. Please try again.
              </p>
            </div>
          </div>
        </div>
      </div>
    );
  }

  // Transform columns data to match our form structure
  const transformedData = {
    tableName: tableId || "",
    columns: columnsData.map((col: Column) => ({
      name: col.name,
      type: mapPostgresToColumnType(col.type),
      primary: col.primary,
    })),
  };

  return (
    <div className="flex flex-col h-full">
      <AppHeader title={`Edit Table: ${tableId || "Unknown"}`} />
      <div className="flex-1 p-6">
        <div className="max-w-4xl mx-auto">
          <TableForm
            mode="edit"
            initialData={transformedData}
            onSubmit={handleSubmit}
            isSubmitting={isSubmitting}
          />
        </div>
      </div>
    </div>
  );
}