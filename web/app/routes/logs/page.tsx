import { useState, useCallback, useMemo, useEffect } from "react";
import { useOutletContext, useSearchParams } from "react-router";
import { useBidirectionalLogs } from "./use-bidirectional-logs";
import { DataTableSkeleton } from "~/components/shared/data-table-skeleton";
import { RefreshButton } from "~/components/shared/refresh-button";
import { Button } from "~/components/ui/button";
import { RefreshCw, Pause, Play } from "lucide-react";
import type { ProjectLayoutOutletContext } from "~/components/shared/project-layout";
import { LogsTable } from "./logs-table";
import { createLogsColumns } from "./columns";
import { LogFilters } from "./log-filters";
import { LogDetailSheet } from "./log-detail-sheet";
import type { LogsFilters, LogEntry } from "~/services/logs";
import type { Route } from "./+types/page";

export function meta({}: Route.MetaArgs) {
  return [
    { title: "Logs - Fluxend" },
    { name: "description", content: "View and manage your application logs" },
  ];
}

const LOGS_PER_PAGE = 100;
const MAX_PAGES_IN_MEMORY = 5; // Keep maximum 5 pages (500 logs) in memory
const MAX_LOGS_IN_MEMORY = LOGS_PER_PAGE * MAX_PAGES_IN_MEMORY;

export default function Logs() {
  const { projectDetails, services } =
    useOutletContext<ProjectLayoutOutletContext>();
  const projectId = projectDetails?.uuid;
  const [searchParams, setSearchParams] = useSearchParams();

  const columns = useMemo(() => createLogsColumns(), []);

  // Initialize filters from URL search params
  const getFiltersFromSearchParams = useCallback(() => {
    const params: LogsFilters = {};
    
    // Get all filter parameters from URL
    const method = searchParams.get('method');
    const status = searchParams.get('status');
    const ipAddress = searchParams.get('ipAddress');
    const endpoint = searchParams.get('endpoint');
    const startTime = searchParams.get('startTime');
    const endTime = searchParams.get('endTime');
    
    if (method) params.method = method;
    if (status) params.status = status;
    if (ipAddress) params.ipAddress = ipAddress;
    if (endpoint) params.endpoint = endpoint;
    if (startTime) params.startTime = parseInt(startTime, 10);
    if (endTime) params.endTime = parseInt(endTime, 10);
    
    return params;
  }, [searchParams]);

  const [filters, setFilters] = useState<LogsFilters>(getFiltersFromSearchParams());
  const [autoRefresh, setAutoRefresh] = useState(false);
  const [refreshInterval] = useState(5000); // 5 seconds
  const [selectedLog, setSelectedLog] = useState<LogEntry | null>(null);
  const [sheetOpen, setSheetOpen] = useState(false);

  // Build query filters
  const queryFilters = useMemo(
    () => ({
      ...filters,
      limit: LOGS_PER_PAGE,
      sort: "createdAt",
      order: "desc" as const,
    }),
    [filters]
  );

  const {
    allLogs,
    isLoading,
    isFetchingNextPage,
    isFetchingPreviousPage,
    error,
    refetch,
    fetchNextPage,
    fetchPreviousPage,
    hasNextPage,
    hasPreviousPage,
    paginationInfo,
  } = useBidirectionalLogs({
    projectId,
    filters: queryFilters,
    services,
    enabled: !!projectId,
    autoRefresh,
    refreshInterval,
    logsPerPage: LOGS_PER_PAGE,
    maxPagesInMemory: MAX_PAGES_IN_MEMORY,
  });

  const handleFilterChange = useCallback((newFilters: LogsFilters) => {
    setFilters(newFilters);
    
    // Update URL search params
    const newSearchParams = new URLSearchParams();
    
    // Add all non-empty filters to URL
    if (newFilters.method) newSearchParams.set('method', newFilters.method);
    if (newFilters.status) newSearchParams.set('status', newFilters.status);
    if (newFilters.ipAddress) newSearchParams.set('ipAddress', newFilters.ipAddress);
    if (newFilters.endpoint) newSearchParams.set('endpoint', newFilters.endpoint);
    if (newFilters.startTime !== undefined) newSearchParams.set('startTime', newFilters.startTime.toString());
    if (newFilters.endTime !== undefined) newSearchParams.set('endTime', newFilters.endTime.toString());
    
    setSearchParams(newSearchParams, { replace: true });
  }, [setSearchParams]);

  const handleRefresh = useCallback(() => {
    refetch();
  }, [refetch]);

  const handleRowClick = useCallback((row: LogEntry) => {
    setSelectedLog(row);
    setSheetOpen(true);
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
    <div className="absolute inset-0 flex flex-col">
      <div className="border-b px-4 py-2 flex-shrink-0">
        <div className="flex items-center justify-between">
          <div className="flex items-center gap-4">
            <div className="text-base font-bold text-foreground h-[32px] flex flex-col justify-center">
              Logs
            </div>
            {paginationInfo && paginationInfo.totalDisplayed > 0 && (
              <div className="text-sm text-muted-foreground">
                Showing {paginationInfo.firstLogNumber}-{paginationInfo.lastLogNumber} of {paginationInfo.totalAvailable} logs
              </div>
            )}
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

      <LogFilters onFiltersChange={handleFilterChange} initialFilters={getFiltersFromSearchParams()} />

      <div className="flex-1 overflow-hidden p-4">
        {isLoading && allLogs.length === 0 ? (
          <div className="rounded-lg border h-full overflow-hidden">
            <DataTableSkeleton columns={7} rows={10} />
          </div>
        ) : (
          <LogsTable
            columns={columns}
            data={allLogs}
            onRowClick={handleRowClick}
            fetchNextPage={fetchNextPage}
            hasNextPage={hasNextPage}
            isFetchingNextPage={isFetchingNextPage}
            isLoading={isLoading}
            error={error}
            hasRemovedPages={paginationInfo?.hasRemovedPages ?? false}
            onLoadPreviousPage={fetchPreviousPage}
            isFetchingPreviousPage={isFetchingPreviousPage}
          />
        )}
      </div>

      <LogDetailSheet
        log={selectedLog}
        open={sheetOpen}
        onOpenChange={setSheetOpen}
      />
    </div>
  );
}