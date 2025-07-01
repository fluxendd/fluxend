import { useRef, useEffect, useCallback, memo, useMemo } from "react";
import { useVirtualizer } from "@tanstack/react-virtual";
import { cn } from "~/lib/utils";
import { Loader2, Clock, FileText, Globe } from "lucide-react";
import type { LogEntry } from "~/services/logs";
import { formatTimestamp } from "~/lib/utils";
import { Badge } from "~/components/ui/badge";

interface SimpleVirtualLogsTableProps {
  data: LogEntry[];
  onRowClick?: (row: LogEntry) => void;
  fetchNextPage: () => void;
  hasNextPage: boolean;
  isFetchingNextPage: boolean;
  isLoading: boolean;
}

// Simple memoized row without complex components
const LogRow = memo(({ 
  log, 
  onRowClick,
  style
}: { 
  log: LogEntry; 
  onRowClick?: (row: LogEntry) => void;
  style: React.CSSProperties;
}) => {
  const { date, time } = useMemo(() => formatTimestamp(log.createdAt), [log.createdAt]);
  
  const handleClick = useCallback((e: React.MouseEvent) => {
    e.preventDefault();
    onRowClick?.(log);
  }, [onRowClick, log]);

  return (
    <div
      style={style}
      className="flex items-center gap-4 px-4 py-2 border-b border-border cursor-pointer"
      onClick={handleClick}
      onMouseEnter={undefined} // Disable hover effects during scroll
    >
      {/* Timestamp */}
      <div className="w-[150px] flex-shrink-0">
        <div className="text-sm font-medium">{date}</div>
        <div className="text-xs text-muted-foreground">{time}</div>
      </div>
      
      {/* Method */}
      <div className="w-[80px] flex-shrink-0">
        <Badge variant="outline" className="font-mono">
          {log.method}
        </Badge>
      </div>
      
      {/* Endpoint */}
      <div className="flex-1 min-w-0">
        <span className="font-mono text-sm truncate block">
          {log.endpoint}
        </span>
      </div>
      
      {/* Status */}
      <div className="w-[80px] flex-shrink-0">
        <Badge 
          variant="outline" 
          className={cn(
            "font-mono",
            log.status >= 200 && log.status < 300 && "border-green-600 text-green-600",
            log.status >= 400 && log.status < 500 && "border-yellow-600 text-yellow-600",
            log.status >= 500 && "border-red-600 text-red-600"
          )}
        >
          {log.status}
        </Badge>
      </div>
      
      {/* IP Address */}
      <div className="w-[120px] flex-shrink-0">
        <Badge variant="outline" className="font-mono">
          {log.ipAddress}
        </Badge>
      </div>
    </div>
  );
});

LogRow.displayName = "LogRow";

export function SimpleVirtualLogsTable({
  data,
  onRowClick,
  fetchNextPage,
  hasNextPage,
  isFetchingNextPage,
  isLoading,
}: SimpleVirtualLogsTableProps) {
  const parentRef = useRef<HTMLDivElement>(null);

  // Simple virtualizer without complex options
  const virtualizer = useVirtualizer({
    count: data.length,
    getScrollElement: () => parentRef.current,
    estimateSize: () => 60,
    overscan: 2,
    getItemKey: useCallback((index: number) => data[index]?.uuid || index, [data]),
  });

  const items = virtualizer.getVirtualItems();

  // Handle infinite scroll
  useEffect(() => {
    const lastItem = items[items.length - 1];
    if (!lastItem) return;

    if (
      lastItem.index >= data.length - 1 &&
      hasNextPage &&
      !isFetchingNextPage
    ) {
      fetchNextPage();
    }
  }, [hasNextPage, fetchNextPage, isFetchingNextPage, items, data.length]);

  if (isLoading && data.length === 0) {
    return (
      <div className="flex items-center justify-center h-full">
        <Loader2 className="h-8 w-8 animate-spin text-muted-foreground" />
      </div>
    );
  }

  if (data.length === 0 && !isLoading) {
    return (
      <div className="flex items-center justify-center h-full">
        <span className="text-muted-foreground">No logs found</span>
      </div>
    );
  }

  return (
    <div className="flex flex-col h-full">
      {/* Header */}
      <div className="flex items-center gap-4 px-4 py-2 border-b bg-background font-medium text-sm">
        <div className="w-[150px] flex-shrink-0 flex items-center gap-2">
          <Clock className="h-3 w-3" />
          Timestamp
        </div>
        <div className="w-[80px] flex-shrink-0">Method</div>
        <div className="flex-1 flex items-center gap-2">
          <FileText className="h-3 w-3" />
          Endpoint
        </div>
        <div className="w-[80px] flex-shrink-0">Status</div>
        <div className="w-[120px] flex-shrink-0 flex items-center gap-2">
          <Globe className="h-3 w-3" />
          IP Address
        </div>
      </div>

      {/* Scrollable content */}
      <div
        ref={parentRef}
        className="flex-1 overflow-auto"
      >
        <div
          style={{
            height: `${virtualizer.getTotalSize()}px`,
            width: '100%',
            position: 'relative',
          }}
        >
          {items.map((virtualItem) => {
            const log = data[virtualItem.index];
            return (
              <LogRow
                key={virtualItem.key}
                log={log}
                onRowClick={onRowClick}
                style={{
                  position: 'absolute',
                  top: 0,
                  left: 0,
                  width: '100%',
                  height: `${virtualItem.size}px`,
                  transform: `translateY(${virtualItem.start}px)`,
                }}
              />
            );
          })}
        </div>
        
        {/* Loading indicator */}
        {isFetchingNextPage && (
          <div className="flex items-center justify-center py-4">
            <div className="flex items-center gap-2 text-muted-foreground">
              <Loader2 className="h-4 w-4 animate-spin" />
              <span className="text-sm">Loading more logs...</span>
            </div>
          </div>
        )}
        
        {!hasNextPage && data.length > 0 && (
          <div className="flex items-center justify-center py-4">
            <span className="text-sm text-muted-foreground">No more logs to load</span>
          </div>
        )}
      </div>
    </div>
  );
}