import { useRef, useEffect, useCallback, useMemo, memo } from "react";
import { flexRender, getCoreRowModel, useReactTable } from "@tanstack/react-table";
import { useVirtualizer } from "@tanstack/react-virtual";
import type { ColumnDef } from "@tanstack/react-table";
import { Table, TableBody, TableCell, TableHead, TableHeader, TableRow } from "~/components/ui/table";
import { cn } from "~/lib/utils";
import { Loader2 } from "lucide-react";
import type { LogEntry } from "~/services/logs";

interface VirtualizedLogsTableProps {
  columns: ColumnDef<LogEntry>[];
  data: LogEntry[];
  onRowClick?: (row: LogEntry) => void;
  fetchNextPage: () => void;
  hasNextPage: boolean;
  isFetchingNextPage: boolean;
  isLoading: boolean;
  windowStartIndex?: number;
}

// Memoized table row component
const VirtualRow = memo(({ 
  row, 
  onRowClick 
}: { 
  row: any; 
  onRowClick?: (row: LogEntry) => void;
}) => {
  const handleClick = useCallback((e: React.MouseEvent) => {
    const target = e.target as HTMLElement;
    if (target.closest('[data-no-row-click]')) {
      return;
    }
    onRowClick?.(row.original);
  }, [onRowClick, row.original]);

  return (
    <TableRow
      className={cn(
        "cursor-pointer hover:bg-muted/50",
        "data-[state=selected]:bg-muted"
      )}
      onClick={handleClick}
    >
      {row.getVisibleCells().map((cell: any) => (
        <TableCell 
          key={cell.id}
          style={{ width: cell.column.getSize() }}
        >
          {flexRender(cell.column.columnDef.cell, cell.getContext())}
        </TableCell>
      ))}
    </TableRow>
  );
});

VirtualRow.displayName = "VirtualRow";

export function VirtualizedLogsTable({
  columns,
  data,
  onRowClick,
  fetchNextPage,
  hasNextPage,
  isFetchingNextPage,
  isLoading,
  windowStartIndex = 0,
}: VirtualizedLogsTableProps) {
  const scrollContainerRef = useRef<HTMLDivElement>(null);
  const tableContainerRef = useRef<HTMLDivElement>(null);

  const table = useReactTable({
    data,
    columns,
    getCoreRowModel: getCoreRowModel(),
  });

  const { rows } = table.getRowModel();

  // Row virtualizer with optimized settings
  const rowVirtualizer = useVirtualizer({
    count: rows.length,
    getScrollElement: () => scrollContainerRef.current,
    estimateSize: useCallback(() => 48, []), // Estimated row height
    overscan: 5, // Reduced overscan for better performance
    scrollMargin: scrollContainerRef.current?.offsetTop ?? 0,
  });

  const virtualRows = rowVirtualizer.getVirtualItems();
  const totalSize = rowVirtualizer.getTotalSize();

  // Track last visible index for infinite scroll
  const lastVisibleIndex = virtualRows[virtualRows.length - 1]?.index ?? -1;
  
  // Handle infinite scroll
  useEffect(() => {
    if (lastVisibleIndex < 0) return;
    
    // Check if we're near the end of the data
    const remainingItems = data.length - lastVisibleIndex;
    const shouldFetch = remainingItems <= 20 && hasNextPage && !isFetchingNextPage;
    
    console.log('Infinite scroll check:', {
      lastVisibleIndex,
      dataLength: data.length,
      remainingItems,
      shouldFetch,
      hasNextPage,
      isFetchingNextPage
    });
    
    if (shouldFetch) {
      console.log('Fetching next page...');
      fetchNextPage();
    }
  }, [lastVisibleIndex, data.length, hasNextPage, isFetchingNextPage, fetchNextPage]);

  const paddingTop = virtualRows.length > 0 ? virtualRows[0]?.start || 0 : 0;
  const paddingBottom =
    virtualRows.length > 0
      ? totalSize - (virtualRows[virtualRows.length - 1]?.end || 0)
      : 0;

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
    <div className="flex flex-col h-full overflow-hidden">
      {/* Scrollable container */}
      <div 
        ref={scrollContainerRef}
        className="flex-1 overflow-auto"
      >
        <div ref={tableContainerRef}>
          <Table>
            <TableHeader className="sticky top-0 z-10 bg-background">
              {table.getHeaderGroups().map((headerGroup) => (
                <TableRow key={headerGroup.id}>
                  {headerGroup.headers.map((header) => (
                    <TableHead 
                      key={header.id}
                      style={{ width: header.column.getSize() }}
                    >
                      {header.isPlaceholder
                        ? null
                        : flexRender(header.column.columnDef.header, header.getContext())}
                    </TableHead>
                  ))}
                </TableRow>
              ))}
            </TableHeader>
            <TableBody>
              {paddingTop > 0 && (
                <tr>
                  <td style={{ height: `${paddingTop}px` }} />
                </tr>
              )}
              {virtualRows.map((virtualRow) => {
                const row = rows[virtualRow.index];
                return (
                  <VirtualRow
                    key={row.id}
                    row={row}
                    onRowClick={onRowClick}
                  />
                );
              })}
              {paddingBottom > 0 && (
                <tr>
                  <td style={{ height: `${paddingBottom}px` }} />
                </tr>
              )}
            </TableBody>
          </Table>
          
        </div>
        
        {/* Loading indicator at the bottom */}
        <div className="border-t bg-background">
          {isFetchingNextPage && (
            <div className="py-4 flex items-center justify-center">
              <div className="flex items-center gap-2">
                <Loader2 className="h-4 w-4 animate-spin text-primary" />
                <span className="text-sm font-medium">Loading more logs...</span>
              </div>
            </div>
          )}
          {!hasNextPage && !isFetchingNextPage && data.length > 0 && (
            <div className="py-3 flex items-center justify-center">
              <span className="text-sm text-muted-foreground">No more logs to load</span>
            </div>
          )}
          {hasNextPage && !isFetchingNextPage && (
            <div className="py-3 flex items-center justify-center">
              <span className="text-sm text-muted-foreground">Scroll to load more</span>
            </div>
          )}
        </div>
      </div>
    </div>
  );
}