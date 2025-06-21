import { useState, useCallback, useMemo, useEffect } from "react";
import { useQueryClient } from "@tanstack/react-query";
import type { PaginationState } from "@tanstack/react-table";

/**
 * Custom hook for managing filter state and handling PostgREST filters
 */
export function useDataFilters(
  projectId: string,
  tableId: string,
  pagination: PaginationState,
  onPaginationChange: (newPagination: PaginationState) => void
) {
  const queryClient = useQueryClient();
  const [filterParams, setFilterParams] = useState<Record<string, string>>({});

  // Handle filter changes
  const handleFilterChange = useCallback(
    (params: Record<string, string>) => {
      setFilterParams(params);

      // Reset to first page when filters change
      if (pagination.pageIndex !== 0) {
        onPaginationChange({
          ...pagination,
          pageIndex: 0,
        });
      } else {
        // If we're already on first page, manually invalidate the query
        queryClient.invalidateQueries({
          queryKey: [
            "rows",
            projectId,
            tableId,
            pagination.pageSize,
            0,
            filterParams,
          ],
        });
      }
    },
    [
      pagination,
      onPaginationChange,
      queryClient,
      projectId,
      tableId,
      filterParams,
    ]
  );

  // Get query keys for the current data
  const queryKeys = useMemo(() => {
    return [
      "rows",
      projectId,
      tableId,
      pagination.pageSize,
      pagination.pageIndex,
      filterParams,
    ];
  }, [
    projectId,
    tableId,
    pagination.pageSize,
    pagination.pageIndex,
    filterParams,
  ]);

  // Reset filters
  const resetFilters = useCallback(() => {
    setFilterParams({});

    // Invalidate query to refresh data without filters
    queryClient.invalidateQueries({
      queryKey: [
        "rows",
        projectId,
        tableId,
        pagination.pageSize,
        pagination.pageIndex,
      ],
    });
  }, [queryClient, projectId, tableId, pagination]);

  // Check if any filters are applied
  const hasActiveFilters = useMemo(() => {
    return Object.keys(filterParams).length > 0;
  }, [filterParams]);

  return {
    filterParams,
    handleFilterChange,
    resetFilters,
    hasActiveFilters,
    queryKeys,
  };
}
