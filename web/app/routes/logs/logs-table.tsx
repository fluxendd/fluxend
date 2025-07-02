import { useRef, useCallback, memo, useState, useEffect } from "react";
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
}

// Memoized table row component
const VirtualRow = memo(({ 
  row, 
  onRowClick,
  isSelected,
  onKeyDown
}: { 
  row: ReturnType<ReturnType<typeof useReactTable<LogEntry>>['getRowModel']>['rows'][0]; 
  onRowClick?: (row: LogEntry) => void;
  isSelected?: boolean;
  onKeyDown?: (e: React.KeyboardEvent) => void;
}) => {
  const handleClick = useCallback((e: React.MouseEvent) => {
    const target = e.target as HTMLElement;
    if (target.closest('[data-no-row-click]')) {
      return;
    }
    // Focus the row when clicked for keyboard navigation
    (e.currentTarget as HTMLTableRowElement).focus();
    onRowClick?.(row.original);
  }, [onRowClick, row.original]);

  return (
    <TableRow
      id={`row-${row.id}`}
      tabIndex={0}
      className={cn(
        "cursor-pointer hover:bg-muted/50",
        "focus:outline-none focus:bg-muted/60",
        isSelected && "bg-muted/50"
      )}
      onClick={handleClick}
      onKeyDown={onKeyDown}
      aria-selected={isSelected}
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
}: LogsTableProps) {
  const tbodyRef = useRef<HTMLTableSectionElement>(null);
  const [selectedRowIndex, setSelectedRowIndex] = useState<number | null>(null);

  const table = useReactTable({
    data,
    columns,
    getCoreRowModel: getCoreRowModel(),
  });

  const { rows } = table.getRowModel();

  // Since we're using manual pagination, we don't need complex virtualization
  // Just render all rows from the current pages

  // Handle keyboard navigation
  const handleKeyDown = useCallback((e: React.KeyboardEvent, rowIndex: number) => {
    const tbody = tbodyRef.current;
    if (!tbody) return;

    const currentRow = tbody.children[rowIndex] as HTMLTableRowElement;
    let targetRow: HTMLTableRowElement | null = null;

    switch (e.key) {
      case 'ArrowUp':
        e.preventDefault();
        if (rowIndex > 0) {
          targetRow = tbody.children[rowIndex - 1] as HTMLTableRowElement;
          setSelectedRowIndex(rowIndex - 1);
        }
        break;
      case 'ArrowDown':
        e.preventDefault();
        if (rowIndex < rows.length - 1) {
          targetRow = tbody.children[rowIndex + 1] as HTMLTableRowElement;
          setSelectedRowIndex(rowIndex + 1);
        }
        break;
      case 'Enter':
      case ' ':
        e.preventDefault();
        if (onRowClick && rows[rowIndex]) {
          onRowClick(rows[rowIndex].original);
        }
        break;
      case 'Home':
        e.preventDefault();
        if (rows.length > 0) {
          targetRow = tbody.children[0] as HTMLTableRowElement;
          setSelectedRowIndex(0);
        }
        break;
      case 'End':
        e.preventDefault();
        if (rows.length > 0) {
          targetRow = tbody.children[rows.length - 1] as HTMLTableRowElement;
          setSelectedRowIndex(rows.length - 1);
        }
        break;
    }

    if (targetRow) {
      targetRow.focus();
      // Ensure the focused row is visible
      targetRow.scrollIntoView({ block: 'nearest' });
    }
  }, [rows, onRowClick]);

  // Handle load more click
  const handleLoadMore = useCallback(() => {
    if (!isFetchingNextPage && hasNextPage) {
      fetchNextPage();
    }
  }, [isFetchingNextPage, hasNextPage, fetchNextPage]);

  // Reset selected row when data changes
  useEffect(() => {
    setSelectedRowIndex(null);
  }, [data]);


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
    <div className="h-full rounded-lg border">
      <div className="h-full overflow-y-auto rounded-lg" role="region" aria-label="Logs table">
        <Table role="table" aria-label="Logs entries">
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
          <TableBody ref={tbodyRef}>
              {rows.map((row, index) => (
                <VirtualRow
                  key={row.id}
                  row={row}
                  onRowClick={onRowClick}
                  isSelected={selectedRowIndex === index}
                  onKeyDown={(e) => handleKeyDown(e, index)}
                />
              ))}
              {/* Load more row */}
              {hasNextPage && !error && (
                <tr>
                  <td colSpan={columns.length}>
                    <div 
                      className="py-4 flex items-center justify-center cursor-pointer hover:bg-muted/50 transition-colors"
                      onClick={handleLoadMore}
                      onKeyDown={(e) => {
                        if (e.key === 'Enter' || e.key === ' ') {
                          e.preventDefault();
                          handleLoadMore();
                        }
                      }}
                      tabIndex={0}
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