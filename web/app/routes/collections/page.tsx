import { useQuery, useQueryClient } from "@tanstack/react-query";
import type { Route } from "./+types/page";
import { columnsQuery, rowsQuery, prepareColumns } from "./columns";
import { LoaderCircle } from "lucide-react";
import { useState, useMemo, useCallback } from "react";
import { useNavigate, href } from "react-router";
import { RefreshButton } from "~/components/shared/refresh-button";
import { DataTableWrapper } from "~/components/shared/data-table-wrapper";
import { SearchDataTableWrapper } from "~/components/shared/search-data-table-wrapper";
import { DataTableSkeleton } from "~/components/shared/data-table-skeleton";
import type { ColumnDef, PaginationState } from "@tanstack/react-table";
import { dbCookie } from "~/lib/cookies";
import { getDbIdFromCookies } from "~/lib/utils";
import { Button } from "~/components/ui/button";

export default function CollectionPageContent({
  params,
}: Route.ComponentProps) {
  const { projectId, collectionId } = params;
  const queryClient = useQueryClient();
  const [pagination, setPagination] = useState({
    pageIndex: 0, //initial page index
    pageSize: 50, //default page size
  });

  const {
    isLoading,
    data: rawColumns = [],
    isFetching,
    error: columnsError,
  } = useQuery(columnsQuery(projectId, collectionId)) || { data: [] };

  const columns = useMemo(() => {
    if (!rawColumns || !Array.isArray(rawColumns) || rawColumns.length === 0) {
      return [];
    }
    return prepareColumns(rawColumns, collectionId);
  }, [rawColumns, collectionId]);

  // Simple filter state management
  const [filterParams, setFilterParams] = useState<Record<string, string>>({});
  const hasActiveFilters = Object.keys(filterParams).length > 0;
  
  // Handle filter changes with pagination reset
  const handleFilterChange = useCallback((params: Record<string, string>) => {
    setFilterParams(params);
    // Reset to first page when filters change
    if (pagination.pageIndex !== 0) {
      setPagination({
        ...pagination,
        pageIndex: 0,
      });
    }
  }, [pagination, setPagination]);
  
  const resetFilters = useCallback(() => {
    setFilterParams({});
  }, []);
  
  const {
    isLoading: isRowsLoading,
    data: rowData = { totalCount: 0, rows: [] },
    isFetching: isRowsFetching,
    error: rowsError,
  } = useQuery({
    ...rowsQuery(projectId, collectionId, pagination, filterParams),
    enabled: !!collectionId,
  });

  // Safely destructure to handle undefined
  const { rows = [], totalCount = 0 } = rowData || {};

  // Track initial loading vs pagination loading separately
  const isInitialDataLoading =
    isLoading || !rawColumns || rawColumns.length === 0;
  const isPaginationLoading = !isInitialDataLoading && isRowsFetching;

  const handleRefresh = async () => {
    if (collectionId) {
      // Invalidate queries
      await Promise.all([
        queryClient.invalidateQueries({
          queryKey: ["columns", projectId, collectionId],
        }),
        queryClient.invalidateQueries({
          queryKey: [
            "rows",
            projectId,
            collectionId,
            pagination.pageSize,
            pagination.pageIndex,
            filterParams,
          ],
        }),
      ]);
    }
  };

  const noCollectionSelected = !collectionId;

  // Calculate loading states for UI
  const isInitialLoading = isInitialDataLoading;
  const isRefetching = isFetching || isRowsFetching;
  const showLoadingState = isInitialLoading;

  return (
    <div className="flex flex-col h-full overflow-hidden">
      <div className="border-b px-4 py-2 mb-2 flex-shrink-0">
        <div className="flex items-center justify-between">
          <div className="text-base font-bold text-foreground h-[32px] flex flex-col justify-center">
            Collections / {collectionId && `(${collectionId})`}
          </div>
          <div className="flex items-center gap-2">
            <RefreshButton
              onRefresh={useCallback(handleRefresh, [
                collectionId,
                projectId,
                queryClient,
              ])}
              title="Refresh Collections and Collection Data"
            />
          </div>
        </div>
      </div>

      {showLoadingState && (
        <div className="rounded-md border mx-4">
          <DataTableSkeleton columns={5} rows={8} />
        </div>
      )}

      {!showLoadingState &&
        Array.isArray(columns) &&
        columns.length > 0 &&
        collectionId && (
          <div className="flex-1 min-h-0 overflow-hidden">
            <SearchDataTableWrapper
              columns={columns}
              rawColumns={rawColumns}
              data={Array.isArray(rows) ? rows : []}
              isLoading={isPaginationLoading || isRowsLoading}
              emptyMessage="No table data found."
              className="w-full h-full"
              pagination={pagination}
              totalRows={totalCount}
              projectId={projectId}
              collectionId={collectionId}
              onFilterChange={handleFilterChange}
              onPaginationChange={(newPagination) => {
                setPagination(newPagination);
                queryClient.invalidateQueries({
                  queryKey: [
                    "rows",
                    projectId,
                    collectionId,
                    newPagination.pageSize,
                    newPagination.pageIndex,
                    filterParams,
                  ],
                });
              }}
            />
          </div>
        )}

      {!isLoading &&
        !isRowsLoading &&
        (!Array.isArray(columns) || columns.length === 0) &&
        collectionId && (
          <div className="flex-1 min-h-0 flex items-center justify-center mx-4">
            <div className="text-md text-muted-foreground border rounded-md p-8 bg-muted/10">
              No Table Data Found
            </div>
          </div>
        )}

      {noCollectionSelected && !isLoading && (
        <div className="flex-1 min-h-0 flex items-center justify-center mx-4">
          <div className="text-md text-muted-foreground border rounded-md p-8 bg-muted/10">
            Please select a collection from the sidebar
          </div>
        </div>
      )}
    </div>
  );
}
