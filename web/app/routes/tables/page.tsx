import { useQuery, useQueryClient } from "@tanstack/react-query";
import type { Route } from "./+types/page";
import { columnsQuery, rowsQuery, prepareColumns } from "./columns";
import { useState, useMemo, useCallback } from "react";
import { RefreshButton } from "~/components/shared/refresh-button";
import { SearchDataTableWrapper } from "~/components/shared/search-data-table-wrapper";
import { DataTableSkeleton } from "~/components/shared/data-table-skeleton";
import { useNavigate, useOutletContext } from "react-router";
import type { ProjectLayoutOutletContext } from "~/components/shared/project-layout";
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
import { Button } from "~/components/ui/button";
import { Trash2, Pencil } from "lucide-react";
import { toast } from "sonner";

const DEFAULT_PAGE_SIZE = 50;
const DEFAULT_PAGE_INDEX = 0;

type PaginationType = {
  pageIndex: number;
  pageSize: number;
};

export default function TablePageContent({ params }: Route.ComponentProps) {
  const { projectId, tableId } = params;
  const { projectDetails, services } =
    useOutletContext<ProjectLayoutOutletContext>();

  const queryClient = useQueryClient();
  const navigate = useNavigate();
  const [pagination, setPagination] = useState<PaginationType>({
    pageIndex: DEFAULT_PAGE_INDEX,
    pageSize: DEFAULT_PAGE_SIZE,
  });

  const {
    isLoading: isColumnsLoading,
    data: columnsData = [],
    isFetching: isColumnsFetching,
    error: columnsError,
  } = useQuery(columnsQuery(projectId, tableId)) || { data: [] };

  const columns = useMemo(() => {
    if (!columnsData || !Array.isArray(columnsData)) {
      return [];
    }

    return prepareColumns(columnsData, tableId);
  }, [columnsData, tableId]);

  const [filterParams, setFilterParams] = useState<Record<string, string>>({});

  // Handle filter changes with pagination reset
  const handleFilterChange = useCallback(
    (params: Record<string, string>) => {
      setFilterParams(params);
      // Reset to first page when filters change
      if (pagination.pageIndex !== 0) {
        setPagination({
          ...pagination,
          pageIndex: 0,
        });
      }
    },
    [pagination, setPagination]
  );

  const resetFilters = useCallback(() => {
    setFilterParams({});
  }, []);

  const {
    isLoading: isRowsLoading,
    data: rowsData = { totalCount: 0, rows: [] },
    isFetching: isRowsFetching,
    error: rowsError,
  } = useQuery({
    ...rowsQuery(
      projectId,
      projectDetails?.dbName as string,
      tableId,
      pagination,
      filterParams
    ),
  });

  // Safely destructure to handle undefined
  const { rows = [], totalCount = 0 } = rowsData;

  // Track initial loading vs pagination loading separately
  const isInitialLoading = isColumnsLoading || isRowsLoading;
  const isFetching = isColumnsFetching || isRowsFetching;

  const handleRefresh = async () => {
    if (tableId) {
      // Invalidate queries
      await Promise.all([
        queryClient.invalidateQueries({
          queryKey: ["columns", projectId, tableId],
        }),
        queryClient.invalidateQueries({
          queryKey: [
            "rows",
            projectId,
            tableId,
            pagination.pageSize,
            pagination.pageIndex,
            filterParams,
          ],
        }),
      ]);
    }
  };

  const onPaginationChange = useCallback(
    (
      updaterOrValue: PaginationType | ((old: PaginationType) => PaginationType)
    ) => {
      const newPagination =
        typeof updaterOrValue === "function"
          ? updaterOrValue(pagination)
          : updaterOrValue;
      setPagination(newPagination);
      queryClient.invalidateQueries({
        queryKey: [
          "rows",
          projectId,
          tableId,
          newPagination.pageSize,
          newPagination.pageIndex,
          filterParams,
        ],
      });
    },
    [pagination, projectId, tableId, filterParams, queryClient]
  );

  const handleDeleteTable = useCallback(async () => {
    if (!tableId || !projectId) return;

    const response = await services.tables.deleteTable(projectId, tableId);

    if (response.ok) {
      // Invalidate collections query to refresh the sidebar
      await queryClient.invalidateQueries({
        queryKey: ["tables", projectId],
      });

      navigate(`/projects/${projectId}/tables`);
    } else if (response?.errors) {
      toast.error(response?.errors[0]);
    } else {
      throw new Error("Unknown error deleting collection");
    }
  }, [tableId, projectId, queryClient, navigate]);

  const noTableSelected = !tableId;

  return (
    <div className="flex flex-col h-full overflow-hidden">
      <div className="border-b px-4 py-2 mb-2 flex-shrink-0">
        <div className="flex items-center justify-between">
          <div className="text-base font-bold text-foreground h-[32px] flex flex-col justify-center">
            Tables / {tableId && `${tableId}`}
          </div>
          <div className="flex items-center gap-2">
            {!noTableSelected && (
              <>
                <Button
                  type="button"
                  variant="ghost"
                  size="sm"
                  className="cursor-pointer"
                  onClick={() => navigate(`/projects/${projectId}/tables/${tableId}/edit`)}
                  title="Edit Table"
                >
                  <Pencil />
                </Button>
                <AlertDialog>
                  <AlertDialogTrigger asChild>
                    <Button
                      type="button"
                      variant="ghost"
                      size="sm"
                      className="cursor-pointer"
                    >
                      <Trash2 />
                    </Button>
                  </AlertDialogTrigger>
                  <AlertDialogContent>
                    <AlertDialogHeader>
                      <AlertDialogTitle>
                        Are you absolutely sure?
                      </AlertDialogTitle>
                      <AlertDialogDescription>
                        This action cannot be undone. This will permanently delete
                        table{" "}
                        <strong className="text-destructive">{tableId}</strong>{" "}
                        from our servers.
                      </AlertDialogDescription>
                    </AlertDialogHeader>
                    <AlertDialogFooter>
                      <AlertDialogCancel>Cancel</AlertDialogCancel>
                      <AlertDialogAction
                        onClick={handleDeleteTable}
                        className="cursor-pointer"
                      >
                        Continue
                      </AlertDialogAction>
                    </AlertDialogFooter>
                  </AlertDialogContent>
                </AlertDialog>
              </>
            )}
            <RefreshButton
              onRefresh={useCallback(handleRefresh, [
                tableId,
                projectId,
                queryClient,
              ])}
              size="sm"
              title="Refresh Tables and Table Data"
            />
          </div>
        </div>
      </div>

      {isInitialLoading && (
        <div className="rounded-md border mx-4 py-4">
          <DataTableSkeleton columns={5} rows={8} />
        </div>
      )}

      {!isInitialLoading &&
        Array.isArray(columns) &&
        columns.length > 0 &&
        tableId && (
          <div className="flex-1 min-h-0 py-2 pb-3 overflow-hidden">
            <SearchDataTableWrapper
              columns={columns}
              rawColumns={columnsData}
              data={Array.isArray(rows) ? rows : []}
              isFetching={isFetching}
              emptyMessage="No table data found."
              className="w-full h-full"
              pagination={pagination}
              totalRows={totalCount}
              projectId={projectId}
              tableId={tableId}
              onFilterChange={handleFilterChange}
              onPaginationChange={onPaginationChange}
            />
          </div>
        )}

      {!isInitialLoading &&
        (!Array.isArray(columns) || columns.length === 0) &&
        tableId && (
          <div className="flex-1 min-h-0 flex items-center justify-center mx-4">
            <div className="text-md text-muted-foreground border rounded-md p-8 bg-muted/10">
              No Table Data Found
            </div>
          </div>
        )}

      {noTableSelected && !isInitialLoading && (
        <div className="flex-1 min-h-0 flex items-center justify-center mx-4">
          <div className="text-md text-muted-foreground border rounded-md p-8 bg-muted/10">
            Please select a collection from the sidebar
          </div>
        </div>
      )}
    </div>
  );
}
