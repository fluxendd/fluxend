import { useRef, useEffect, useCallback, useMemo } from "react";
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
}

export function VirtualizedLogsTable({
  columns,
  data,
  onRowClick,
  fetchNextPage,
  hasNextPage,
  isFetchingNextPage,
  isLoading,
}: VirtualizedLogsTableProps) {
  const scrollContainerRef = useRef<HTMLDivElement>(null);
  const tableContainerRef = useRef<HTMLDivElement>(null);

  const table = useReactTable({
    data,
    columns,
    getCoreRowModel: getCoreRowModel(),
  });

  const { rows } = table.getRowModel();

  // Row virtualizer
  const rowVirtualizer = useVirtualizer({
    count: rows.length,
    getScrollElement: () => scrollContainerRef.current,
    estimateSize: () => 48, // Estimated row height
    overscan: 10,
  });

  const virtualRows = rowVirtualizer.getVirtualItems();
  const totalSize = rowVirtualizer.getTotalSize();

  // Handle infinite scroll
  const lastItem = virtualRows[virtualRows.length - 1];
  useEffect(() => {
    if (!lastItem) return;

    if (
      lastItem.index >= rows.length - 1 &&
      hasNextPage &&
      !isFetchingNextPage
    ) {
      fetchNextPage();
    }
  }, [hasNextPage, fetchNextPage, isFetchingNextPage, lastItem, rows.length]);

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
                  <TableRow
                    key={row.id}
                    className={cn(
                      "cursor-pointer hover:bg-muted/50",
                      "data-[state=selected]:bg-muted"
                    )}
                    onClick={(e) => {
                      // Don't trigger row click if clicking on a clickable element
                      const target = e.target as HTMLElement;
                      if (target.closest('[data-no-row-click]')) {
                        return;
                      }
                      onRowClick?.(row.original);
                    }}
                  >
                    {row.getVisibleCells().map((cell) => (
                      <TableCell 
                        key={cell.id}
                        style={{ width: cell.column.getSize() }}
                      >
                        {flexRender(cell.column.columnDef.cell, cell.getContext())}
                      </TableCell>
                    ))}
                  </TableRow>
                );
              })}
              {paddingBottom > 0 && (
                <tr>
                  <td style={{ height: `${paddingBottom}px` }} />
                </tr>
              )}
            </TableBody>
          </Table>
          
          {/* Loading indicator at the bottom */}
          {isFetchingNextPage && (
            <div className="h-20 flex items-center justify-center">
              <div className="flex items-center gap-2 text-muted-foreground">
                <Loader2 className="h-4 w-4 animate-spin" />
                <span className="text-sm">Loading more logs...</span>
              </div>
            </div>
          )}
          {!hasNextPage && data.length > 0 && (
            <div className="h-20 flex items-center justify-center">
              <span className="text-sm text-muted-foreground">No more logs to load</span>
            </div>
          )}
        </div>
      </div>
    </div>
  );
}