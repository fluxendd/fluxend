import { useState, useCallback, useMemo, useEffect } from "react";
import { useQuery } from "@tanstack/react-query";
import { useOutletContext } from "react-router";
import { DataTableSkeleton } from "~/components/shared/data-table-skeleton";
import { RefreshButton } from "~/components/shared/refresh-button";
import { Button } from "~/components/ui/button";
import { RefreshCw, Pause, Play } from "lucide-react";
import type { ProjectLayoutOutletContext } from "~/components/shared/project-layout";
import { DataTable } from "~/routes/tables/data-table";
import { createLogsColumns } from "./columns";
import { LogFilters } from "./log-filters";
import { LogDetailSheet } from "./log-detail-sheet";
import type { LogsFilters, LogEntry } from "~/services/logs";

const DEFAULT_PAGE_SIZE = 50;
const DEFAULT_PAGE_INDEX = 0;

type PaginationType = {
  pageIndex: number;
  pageSize: number;
};

export default function Logs() {
  const { projectDetails, services } = useOutletContext<ProjectLayoutOutletContext>();
  const projectId = projectDetails?.uuid;
  
  const columns = useMemo(() => createLogsColumns(), []);

  const [pagination, setPagination] = useState<PaginationType>({
    pageIndex: DEFAULT_PAGE_INDEX,
    pageSize: DEFAULT_PAGE_SIZE,
  });

  const [filters, setFilters] = useState<LogsFilters>({});
  const [autoRefresh, setAutoRefresh] = useState(false);
  const [refreshInterval] = useState(5000); // 5 seconds
  const [sheetState, setSheetState] = useState<{ log: LogEntry | null; open: boolean }>({
    log: null,
    open: false
  });

  // Build query filters
  const queryFilters = useMemo(() => ({
    ...filters,
    page: pagination.pageIndex + 1, // API uses 1-based pagination
    limit: pagination.pageSize,
    order: 'desc' as const,
    sort: 'createdAt',
  }), [filters, pagination]);

  const {
    data: logsData,
    isLoading,
    isFetching,
    error,
    refetch,
  } = useQuery({
    queryKey: ["logs", projectId, queryFilters],
    queryFn: async () => {
      if (!projectId) throw new Error("Project ID required");
      return services.logs.getLogs(projectId, queryFilters);
    },
    enabled: !!projectId,
    refetchInterval: autoRefresh ? refreshInterval : false,
  });

  const handlePaginationChange = useCallback((updaterOrValue: PaginationType | ((old: PaginationType) => PaginationType)) => {
    setPagination(updaterOrValue);
  }, []);

  const handleFilterChange = useCallback((newFilters: LogsFilters) => {
    setFilters(newFilters);
    // Reset to first page when filters change
    setPagination(prev => ({ ...prev, pageIndex: 0 }));
  }, []);

  const handleRefresh = useCallback(() => {
    refetch();
  }, [refetch]);

  const handleRowClick = useCallback((row: LogEntry) => {
    setSheetState({ log: row, open: true });
  }, []);

  if (error) {
    return (
      <div className="flex flex-col h-full">
        <div className="border-b px-4 py-2 flex-shrink-0">
          <div className="flex items-center justify-between">
            <div className="text-base font-bold text-foreground h-[32px] flex flex-col justify-center">
              Logs
            </div>
          </div>
        </div>
        <div className="flex-1 flex items-center justify-center">
          <div className="text-destructive">
            Error loading logs: {error.message}
          </div>
        </div>
      </div>
    );
  }

  return (
    <div className="flex flex-col h-full overflow-hidden">
      <div className="border-b px-4 py-2 flex-shrink-0">
        <div className="flex items-center justify-between">
          <div className="text-base font-bold text-foreground h-[32px] flex flex-col justify-center">
            Logs
          </div>
          <div className="flex items-center gap-2">
            <Button
              size="sm"
              variant={autoRefresh ? "secondary" : "outline"}
              onClick={() => setAutoRefresh(!autoRefresh)}
              title={autoRefresh ? "Stop auto-refresh" : "Start auto-refresh"}
            >
              {autoRefresh ? (
                <>
                  <Pause className="h-4 w-4 mr-1" />
                  Auto-refresh ON
                </>
              ) : (
                <>
                  <Play className="h-4 w-4 mr-1" />
                  Auto-refresh OFF
                </>
              )}
            </Button>
            <RefreshButton
              onRefresh={handleRefresh}
              size="sm"
              title="Refresh Logs"
            />
          </div>
        </div>
      </div>

      <LogFilters onFiltersChange={handleFilterChange} />

      {isLoading && (
        <div className="rounded-lg border mx-4 py-4 mt-4">
          <DataTableSkeleton columns={7} rows={8} />
        </div>
      )}

      {!isLoading && logsData && (
        <div className="flex-1 min-h-0 p-4 overflow-hidden">
          <div className="rounded-lg border h-full">
            <DataTable
              columns={columns}
              data={logsData.content || []}
              pagination={pagination}
              onPaginationChange={handlePaginationChange}
              totalRows={logsData.totalCount}
              emptyMessage="No logs found."
              onRowClick={handleRowClick}
            />
          </div>
        </div>
      )}

      <LogDetailSheet
        log={sheetState.log}
        open={sheetState.open}
        onOpenChange={(open) => setSheetState(prev => ({ ...prev, open }))}
      />
    </div>
  );
}
