import { useState, useCallback, useMemo, useEffect } from "react";
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
  const [loadedPageNumbers, setLoadedPageNumbers] = useState<number[]>([]);
  const [isLoadingPage, setIsLoadingPage] = useState(false);
  const [totalAvailable, setTotalAvailable] = useState(0);

  // Query key for caching
  const queryKey = ["logs", projectId, filters];

  // Reset state when filters change
  useEffect(() => {
    setLoadedPages(new Map());
    setLoadedPageNumbers([]);
  }, [filters]);

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
      setLoadedPageNumbers([1]);
      setTotalAvailable(firstPage.metadata.total);
      return firstPage;
    },
    enabled: enabled && loadedPageNumbers.length === 0,
    refetchInterval: autoRefresh ? refreshInterval : false,
    staleTime: 5 * 60 * 1000,
    gcTime: 10 * 60 * 1000,
  });

  // Load next page
  const loadNextPage = useCallback(async () => {
    if (isLoadingPage || loadedPageNumbers.length === 0) return;
    
    const lastPage = Math.max(...loadedPageNumbers);
    const nextPage = lastPage + 1;
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
          const oldestPage = Math.min(...loadedPageNumbers);
          newPages.delete(oldestPage);
          setLoadedPageNumbers(prev => [...prev.filter(p => p !== oldestPage), nextPage].sort((a, b) => a - b));
        } else {
          setLoadedPageNumbers(prev => [...prev, nextPage].sort((a, b) => a - b));
        }
        
        return newPages;
      });
      
      setTotalAvailable(pageData.metadata.total);
    } finally {
      setIsLoadingPage(false);
    }
  }, [isLoadingPage, loadedPageNumbers, totalAvailable, logsPerPage, fetchPage, maxPagesInMemory]);

  // Load previous page (when scrolling up after pages were removed)
  const loadPreviousPage = useCallback(async () => {
    if (isLoadingPage || loadedPageNumbers.length === 0) return;
    
    const firstPage = Math.min(...loadedPageNumbers);
    if (firstPage === 1) return; // Already at the beginning
    
    const previousPage = firstPage - 1;
    
    console.log('Loading previous page:', {
      previousPage,
      currentPages: loadedPageNumbers,
      firstPage,
      maxPages: maxPagesInMemory
    });
    
    setIsLoadingPage(true);
    try {
      const pageData = await fetchPage(previousPage);
      
      setLoadedPages(prev => {
        const newPages = new Map(prev);
        newPages.set(previousPage, pageData);
        
        // Remove newest page if we exceed the limit
        if (newPages.size > maxPagesInMemory) {
          const newestPage = Math.max(...loadedPageNumbers);
          newPages.delete(newestPage);
          setLoadedPageNumbers(prev => [previousPage, ...prev.filter(p => p !== newestPage)].sort((a, b) => a - b));
        } else {
          setLoadedPageNumbers(prev => [previousPage, ...prev].sort((a, b) => a - b));
        }
        
        return newPages;
      });
    } finally {
      setIsLoadingPage(false);
    }
  }, [isLoadingPage, loadedPageNumbers, fetchPage, maxPagesInMemory]);

  // Refresh data (reload current window)
  const refresh = useCallback(async () => {
    if (loadedPageNumbers.length === 0) return;
    
    setIsLoadingPage(true);
    try {
      const newPages = new Map<number, PageData>();
      
      // Reload all currently loaded pages
      for (const page of loadedPageNumbers) {
        const pageData = await fetchPage(page);
        newPages.set(page, pageData);
        setTotalAvailable(pageData.metadata.total);
      }
      
      setLoadedPages(newPages);
    } finally {
      setIsLoadingPage(false);
    }
  }, [loadedPageNumbers, fetchPage]);

  // Computed values
  const { allLogs, hasNextPage, hasPreviousPage, paginationInfo } = useMemo(() => {
    // Sort pages and flatten logs
    const sortedPages = Array.from(loadedPages.values()).sort((a, b) => a.page - b.page);
    const logs = sortedPages.flatMap(page => page.logs);
    
    const totalPages = Math.ceil(totalAvailable / logsPerPage);
    const firstPage = loadedPageNumbers.length > 0 ? Math.min(...loadedPageNumbers) : 1;
    const lastPage = loadedPageNumbers.length > 0 ? Math.max(...loadedPageNumbers) : 1;
    const hasNext = lastPage < totalPages;
    const hasPrevious = firstPage > 1;
    
    return {
      allLogs: logs,
      hasNextPage: hasNext,
      hasPreviousPage: hasPrevious,
      paginationInfo: {
        totalDisplayed: logs.length,
        totalAvailable,
        currentPage: lastPage,
        totalPages,
        firstPageInMemory: firstPage,
        lastPageInMemory: lastPage,
        hasRemovedPages: firstPage > 1,
        hasMorePages: hasNext,
      },
    };
  }, [loadedPages, loadedPageNumbers, totalAvailable, logsPerPage]);

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