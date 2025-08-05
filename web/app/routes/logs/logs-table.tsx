import { useRef, useCallback, memo } from "react";
import { flexRender, getCoreRowModel, useReactTable } from "@tanstack/react-table";
import type { ColumnDef } from "@tanstack/react-table";
import { Table, TableBody, TableCell, TableHead, TableHeader, TableRow } from "~/components/ui/table";
import { cn } from "~/lib/utils";
import { Loader2 } from "lucide-react";
import type { LogEntry } from "~/services/logs";

interface LogsTableProps {
  columns: ColumnDef<LogEntry>[];
  data: LogEntry[];
  onRowClick?: (row: LogEntry) => void;
  fetchNextPage: () => void;
  hasNextPage: boolean;
  isFetchingNextPage: boolean;
  isLoading: boolean;
  error?: Error | null;
  hasRemovedPages?: boolean;
  onLoadPreviousPage?: () => void;
  isFetchingPreviousPage?: boolean;
}

// Memoized table row component
const VirtualRow = memo(({ 
  row, 
  onRowClick 
}: { 
  row: ReturnType<ReturnType<typeof useReactTable<LogEntry>>['getRowModel']>['rows'][0]; 
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
});

VirtualRow.displayName = "VirtualRow";

export function LogsTable({
  columns,
  data,
  onRowClick,
  fetchNextPage,
  hasNextPage,
  isFetchingNextPage,
  isLoading,
  error,
  hasRemovedPages,
  onLoadPreviousPage,
  isFetchingPreviousPage,
}: LogsTableProps) {

  const table = useReactTable({
    data,
    columns,
    getCoreRowModel: getCoreRowModel(),
  });

  const { rows } = table.getRowModel();

  // Since we're using manual pagination, we don't need complex virtualization
  // Just render all rows from the current pages

  // Handle load more click
  const handleLoadMore = useCallback(() => {
    if (!isFetchingNextPage && hasNextPage) {
      fetchNextPage();
    }
  }, [isFetchingNextPage, hasNextPage, fetchNextPage]);

  // Handle load previous click
  const handleLoadPrevious = useCallback(() => {
    if (!isFetchingPreviousPage && onLoadPreviousPage) {
      onLoadPreviousPage();
    }
  }, [isFetchingPreviousPage, onLoadPreviousPage]);


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
    <div className="h-full rounded-lg border overflow-hidden">
      <div className="h-full overflow-auto">
        <Table>
          <TableHeader className="sticky top-0 z-10 bg-background border-b">
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
              {/* Load previous logs row when pages have been removed */}
              {hasRemovedPages && onLoadPreviousPage && (
                <tr>
                  <td colSpan={columns.length}>
                    <div 
                      className="py-4 flex items-center justify-center cursor-pointer hover:bg-muted/50 transition-colors bg-muted/20"
                      onClick={handleLoadPrevious}
                    >
                      {isFetchingPreviousPage ? (
                        <div className="flex items-center gap-2">
                          <Loader2 className="h-4 w-4 animate-spin text-primary" />
                          <span className="text-sm font-medium">Loading previous logs...</span>
                        </div>
                      ) : (
                        <span className="text-sm text-primary font-medium hover:underline">
                          Load previous logs
                        </span>
                      )}
                    </div>
                  </td>
                </tr>
              )}
              {rows.map((row) => (
                <VirtualRow
                  key={row.id}
                  row={row}
                  onRowClick={onRowClick}
                />
              ))}
              {/* Load more row */}
              {hasNextPage && !error && (
                <tr>
                  <td colSpan={columns.length}>
                    <div 
                      className="py-4 flex items-center justify-center cursor-pointer hover:bg-muted/50 transition-colors"
                      onClick={handleLoadMore}
                    >
                      {isFetchingNextPage ? (
                        <div className="flex items-center gap-2">
                          <Loader2 className="h-4 w-4 animate-spin text-primary" />
                          <span className="text-sm font-medium">Loading more logs...</span>
                        </div>
                      ) : (
                        <span className="text-sm text-primary font-medium hover:underline">
                          Load more logs
                        </span>
                      )}
                    </div>
                  </td>
                </tr>
              )}
              {/* No more logs message */}
              {!hasNextPage && data.length > 0 && (
                <tr>
                  <td colSpan={columns.length}>
                    <div className="py-3 flex items-center justify-center">
                      <span className="text-sm text-muted-foreground">No more logs to load</span>
                    </div>
                  </td>
                </tr>
              )}
              {/* Error row */}
              {error && (
                <tr>
                  <td colSpan={columns.length}>
                    <div className="py-4 flex flex-col items-center justify-center gap-2">
                      <span className="text-sm text-destructive">Failed to load more logs</span>
                      <button
                        onClick={handleLoadMore}
                        className="text-sm text-primary hover:underline"
                      >
                        Try again
                      </button>
                    </div>
                  </td>
                </tr>
              )}
          </TableBody>
        </Table>
      </div>
    </div>
  );
}