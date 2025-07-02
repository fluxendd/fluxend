import { useState, useCallback, useMemo } from "react";
import { useInfiniteQuery } from "@tanstack/react-query";
import { useOutletContext } from "react-router";
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

const LOGS_PER_PAGE = 100;
const MAX_PAGES_IN_MEMORY = 5; // Keep maximum 5 pages (500 logs) in memory

export default function Logs() {
  const { projectDetails, services } =
    useOutletContext<ProjectLayoutOutletContext>();
  const projectId = projectDetails?.uuid;

  const columns = useMemo(() => createLogsColumns(), []);

  const [filters, setFilters] = useState<LogsFilters>({});
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
    data,
    isLoading,
    isFetchingNextPage,
    error,
    refetch,
    fetchNextPage,
    hasNextPage,
  } = useInfiniteQuery({
    queryKey: ["logs", projectId, queryFilters],
    queryFn: async ({ pageParam = 1 }) => {
      if (!projectId) throw new Error("Project ID required");
      return services.logs.getLogs(projectId, {
        ...queryFilters,
        page: pageParam,
      });
    },
    enabled: !!projectId,
    refetchInterval: autoRefresh ? refreshInterval : false,
    retry: 1, // Only retry once on failure
    getNextPageParam: (lastPage) => {
      // Use the hasMore flag from the API response
      if (lastPage.hasMore) {
        // API returns current page number, so next page is current + 1
        return lastPage.metadata.page + 1;
      }
      return undefined;
    },
    initialPageParam: 1,
    // Prevent unnecessary refetches
    staleTime: 5 * 60 * 1000, // 5 minutes
    gcTime: 10 * 60 * 1000, // 10 minutes (was cacheTime)
    refetchOnMount: false,
    refetchOnWindowFocus: false,
    refetchOnReconnect: false,
    // Keep only MAX_PAGES_IN_MEMORY pages
    maxPages: MAX_PAGES_IN_MEMORY,
  });

  // Flatten pages data and get pagination info
  const { allLogs, paginationInfo } = useMemo(() => {
    if (!data?.pages || data.pages.length === 0) {
      return { allLogs: [], paginationInfo: null };
    }
    
    const logs = data.pages.flatMap((page) => page.content);
    const lastPage = data.pages[data.pages.length - 1];
    
    // Calculate total logs displayed vs total available
    const totalDisplayed = logs.length;
    const totalAvailable = lastPage.metadata.total;
    const currentPage = lastPage.metadata.page;
    const totalPages = Math.ceil(totalAvailable / lastPage.metadata.limit);
    
    return {
      allLogs: logs,
      paginationInfo: {
        totalDisplayed,
        totalAvailable,
        currentPage,
        totalPages,
      }
    };
  }, [data]);

  const handleFilterChange = useCallback((newFilters: LogsFilters) => {
    setFilters(newFilters);
  }, []);

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
            {paginationInfo && (
              <div className="text-sm text-muted-foreground">
                Showing {paginationInfo.totalDisplayed} of {paginationInfo.totalAvailable} logs
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

      <LogFilters onFiltersChange={handleFilterChange} />

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
            hasNextPage={hasNextPage ?? false}
            isFetchingNextPage={isFetchingNextPage}
            isLoading={isLoading}
            error={error}
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