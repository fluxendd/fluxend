import { useState, useCallback, useMemo, useEffect } from "react";
import { useInfiniteQuery } from "@tanstack/react-query";
import { useOutletContext } from "react-router";
import { DataTableSkeleton } from "~/components/shared/data-table-skeleton";
import { RefreshButton } from "~/components/shared/refresh-button";
import { Button } from "~/components/ui/button";
import { RefreshCw, Pause, Play } from "lucide-react";
import type { ProjectLayoutOutletContext } from "~/components/shared/project-layout";
import { VirtualizedLogsTable } from "./virtualized-logs-table";
import { createLogsColumnsVirtualized } from "./columns";
import { LogFilters } from "./log-filters";
import { LogDetailSheet } from "./log-detail-sheet";
import type { LogsFilters, LogEntry } from "~/services/logs";

const LOGS_PER_PAGE = 100;
const MAX_PAGES_IN_MEMORY = 5; // Keep maximum 5 pages (500 logs) in memory

export default function Logs() {
  const { projectDetails, services } = useOutletContext<ProjectLayoutOutletContext>();
  const projectId = projectDetails?.uuid;
  
  const columns = useMemo(() => createLogsColumnsVirtualized(), []);

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
    limit: LOGS_PER_PAGE,
    sort: 'createdAt',
    order: 'desc',
  }), [filters]);

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
    getNextPageParam: (lastPage, pages) => {
      console.log('getNextPageParam:', {
        lastPageLength: lastPage.content.length,
        LOGS_PER_PAGE,
        totalPages: pages.length,
        willFetchMore: lastPage.content.length >= LOGS_PER_PAGE
      });
      // If we got a full page or close to it, there might be more
      if (lastPage.content.length >= LOGS_PER_PAGE - 10) {
        return pages.length + 1;
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
  });

  // Implement windowed pagination - keep only MAX_PAGES_IN_MEMORY pages
  const windowedData = useMemo(() => {
    if (!data?.pages) return { logs: [], startIndex: 0 };
    
    const totalPages = data.pages.length;
    
    // If we have more than MAX_PAGES_IN_MEMORY, create a window
    if (totalPages > MAX_PAGES_IN_MEMORY) {
      // Keep the most recent MAX_PAGES_IN_MEMORY pages
      const startPageIndex = totalPages - MAX_PAGES_IN_MEMORY;
      const windowedPages = data.pages.slice(startPageIndex);
      const logs = windowedPages.flatMap(page => page.content);
      const startIndex = startPageIndex * LOGS_PER_PAGE;
      
      return { logs, startIndex };
    }
    
    // Otherwise, return all logs
    return { 
      logs: data.pages.flatMap(page => page.content),
      startIndex: 0
    };
  }, [data]);

  const handleFilterChange = useCallback((newFilters: LogsFilters) => {
    setFilters(newFilters);
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

      <div className="flex-1 min-h-0 p-4 overflow-hidden">
        {isLoading && windowedData.logs.length === 0 ? (
          <div className="rounded-lg border">
            <div className="p-4">
              <DataTableSkeleton columns={7} rows={8} />
            </div>
          </div>
        ) : (
          <div className="rounded-lg border h-full">
            <VirtualizedLogsTable
              columns={columns}
              data={windowedData.logs}
              onRowClick={handleRowClick}
              fetchNextPage={fetchNextPage}
              hasNextPage={hasNextPage ?? false}
              isFetchingNextPage={isFetchingNextPage}
              isLoading={isLoading}
              windowStartIndex={windowedData.startIndex}
            />
          </div>
        )}
      </div>

      <LogDetailSheet
        log={sheetState.log}
        open={sheetState.open}
        onOpenChange={(open) => setSheetState(prev => ({ ...prev, open }))}
      />
    </div>
  );
}