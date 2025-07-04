import { useState, useCallback, useMemo } from "react";
import { useQuery, useQueryClient } from "@tanstack/react-query";
import type { LogsFilters, LogEntry, LogsResponse } from "~/services/logs";

interface PageData {
  page: number;
  logs: LogEntry[];
  metadata: LogsResponse['metadata'];
}

interface UseBidirectionalLogsOptions {
  projectId: string | undefined;
  filters: LogsFilters;
  services: any;
  enabled: boolean;
  autoRefresh: boolean;
  refreshInterval: number;
  logsPerPage: number;
  maxPagesInMemory: number;
}

export function useBidirectionalLogs({
  projectId,
  filters,
  services,
  enabled,
  autoRefresh,
  refreshInterval,
  logsPerPage,
  maxPagesInMemory,
}: UseBidirectionalLogsOptions) {
  const queryClient = useQueryClient();
  
  // Track loaded pages and their data
  const [loadedPages, setLoadedPages] = useState<Map<number, PageData>>(new Map());
  const [pageWindow, setPageWindow] = useState({ start: 1, end: 1 });
  const [isLoadingPage, setIsLoadingPage] = useState(false);
  const [totalAvailable, setTotalAvailable] = useState(0);

  // Query key for caching
  const queryKey = ["logs", projectId, filters];

  // Function to fetch a specific page
  const fetchPage = useCallback(async (pageNum: number): Promise<PageData> => {
    if (!projectId) throw new Error("Project ID required");
    
    const response = await services.logs.getLogs(projectId, {
      ...filters,
      limit: logsPerPage,
      sort: "createdAt",
      order: "desc" as const,
      page: pageNum,
    });

    return {
      page: pageNum,
      logs: response.content,
      metadata: response.metadata,
    };
  }, [projectId, services.logs, filters, logsPerPage]);

  // Initial data fetch
  const { isLoading: isInitialLoading, error: initialError } = useQuery({
    queryKey: [...queryKey, "initial"],
    queryFn: async () => {
      const firstPage = await fetchPage(1);
      setLoadedPages(new Map([[1, firstPage]]));
      setPageWindow({ start: 1, end: 1 });
      setTotalAvailable(firstPage.metadata.total);
      return firstPage;
    },
    enabled: enabled && !loadedPages.has(1),
    refetchInterval: autoRefresh ? refreshInterval : false,
    staleTime: 5 * 60 * 1000,
    gcTime: 10 * 60 * 1000,
  });

  // Load next page
  const loadNextPage = useCallback(async () => {
    if (isLoadingPage) return;
    
    const nextPage = pageWindow.end + 1;
    const totalPages = Math.ceil(totalAvailable / logsPerPage);
    
    if (nextPage > totalPages) return;
    
    setIsLoadingPage(true);
    try {
      const pageData = await fetchPage(nextPage);
      
      setLoadedPages(prev => {
        const newPages = new Map(prev);
        newPages.set(nextPage, pageData);
        
        // Remove oldest page if we exceed the limit
        if (newPages.size > maxPagesInMemory) {
          const oldestPage = pageWindow.start;
          newPages.delete(oldestPage);
          setPageWindow(w => ({ start: w.start + 1, end: nextPage }));
        } else {
          setPageWindow(w => ({ ...w, end: nextPage }));
        }
        
        return newPages;
      });
      
      setTotalAvailable(pageData.metadata.total);
    } finally {
      setIsLoadingPage(false);
    }
  }, [isLoadingPage, pageWindow, totalAvailable, logsPerPage, fetchPage, maxPagesInMemory]);

  // Load previous page (when scrolling up after pages were removed)
  const loadPreviousPage = useCallback(async () => {
    if (isLoadingPage || pageWindow.start === 1) return;
    
    const previousPage = pageWindow.start - 1;
    
    setIsLoadingPage(true);
    try {
      const pageData = await fetchPage(previousPage);
      
      setLoadedPages(prev => {
        const newPages = new Map(prev);
        newPages.set(previousPage, pageData);
        
        // Remove newest page if we exceed the limit
        if (newPages.size > maxPagesInMemory) {
          const newestPage = pageWindow.end;
          newPages.delete(newestPage);
          setPageWindow(w => ({ start: previousPage, end: w.end - 1 }));
        } else {
          setPageWindow(w => ({ start: previousPage, ...w }));
        }
        
        return newPages;
      });
    } finally {
      setIsLoadingPage(false);
    }
  }, [isLoadingPage, pageWindow, fetchPage, maxPagesInMemory]);

  // Refresh data (reload current window)
  const refresh = useCallback(async () => {
    setIsLoadingPage(true);
    try {
      const newPages = new Map<number, PageData>();
      
      // Reload all pages in current window
      for (let page = pageWindow.start; page <= pageWindow.end; page++) {
        const pageData = await fetchPage(page);
        newPages.set(page, pageData);
        if (page === pageWindow.end) {
          setTotalAvailable(pageData.metadata.total);
        }
      }
      
      setLoadedPages(newPages);
    } finally {
      setIsLoadingPage(false);
    }
  }, [pageWindow, fetchPage]);

  // Computed values
  const { allLogs, hasNextPage, hasPreviousPage, paginationInfo } = useMemo(() => {
    // Sort pages and flatten logs
    const sortedPages = Array.from(loadedPages.values()).sort((a, b) => a.page - b.page);
    const logs = sortedPages.flatMap(page => page.logs);
    
    const totalPages = Math.ceil(totalAvailable / logsPerPage);
    const hasNext = pageWindow.end < totalPages;
    const hasPrevious = pageWindow.start > 1;
    
    return {
      allLogs: logs,
      hasNextPage: hasNext,
      hasPreviousPage: hasPrevious,
      paginationInfo: {
        totalDisplayed: logs.length,
        totalAvailable,
        currentPage: pageWindow.end,
        totalPages,
        firstPageInMemory: pageWindow.start,
        lastPageInMemory: pageWindow.end,
        hasRemovedPages: pageWindow.start > 1,
        hasMorePages: hasNext,
      },
    };
  }, [loadedPages, pageWindow, totalAvailable, logsPerPage]);

  return {
    data: { pages: Array.from(loadedPages.values()) },
    allLogs,
    isLoading: isInitialLoading,
    isFetchingNextPage: isLoadingPage && pageWindow.end < Math.ceil(totalAvailable / logsPerPage),
    isFetchingPreviousPage: isLoadingPage && pageWindow.start > 1,
    error: initialError,
    refetch: refresh,
    fetchNextPage: loadNextPage,
    fetchPreviousPage: loadPreviousPage,
    hasNextPage,
    hasPreviousPage,
    paginationInfo,
  };
}