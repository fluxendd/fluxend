import { useEffect } from "react";
import { cn } from "~/lib/utils";
import { DataTable } from "~/routes/tables/data-table";
import { QuerySearchBox } from "./query-search-box";
import type {
  ColumnDef,
  OnChangeFn,
  PaginationState,
  RowSelectionState,
} from "@tanstack/react-table";

interface SearchDataTableWrapperProps<TData, TValue> {
  columns: ColumnDef<TData, TValue>[];
  rawColumns: any[]; // Original column metadata from API
  data: TData[];
  isFetching?: boolean;
  className?: string;
  emptyMessage?: string;
  pagination: PaginationState;
  onPaginationChange: OnChangeFn<PaginationState>;
  totalRows: number;
  projectId: string;
  tableId: string;
  onFilterChange: (filters: Record<string, string>) => void;
  searchQuery?: string;
  onSearchQueryChange?: (query: string) => void;
  tableMeta?: any;
  onRowSelectionChange?: (selectedRows: TData[]) => void;
  rowSelection?: RowSelectionState;
  onRowSelectionStateChange?: OnChangeFn<RowSelectionState>;
}

export function SearchDataTableWrapper<TData, TValue>({
  columns,
  rawColumns,
  data,
  isFetching = false,
  className,
  emptyMessage = "No results.",
  pagination,
  onPaginationChange,
  totalRows,
  projectId,
  tableId,
  onFilterChange,
  searchQuery,
  onSearchQueryChange,
  tableMeta,
  onRowSelectionChange,
  rowSelection,
  onRowSelectionStateChange,
}: SearchDataTableWrapperProps<TData, TValue>) {
  return (
    <div className="flex flex-col h-full gap-4 px-4">
      <QuerySearchBox 
        key={`query-search-${tableId}`}
        columns={rawColumns} 
        onQueryChange={onFilterChange}
        value={searchQuery}
        onChange={onSearchQueryChange}
      />
      <div
        className={cn(
          "rounded-lg overflow-hidden flex flex-col h-full min-h-0 max-h-full flex-1",
          isFetching ? "animate-pulse-border" : "border",
          className
        )}
      >
        <DataTable
          columns={columns}
          data={data}
          emptyMessage={emptyMessage}
          pagination={pagination}
          onPaginationChange={onPaginationChange}
          totalRows={totalRows}
          tableMeta={tableMeta}
          onRowSelectionChange={onRowSelectionChange}
          rowSelection={rowSelection}
          onRowSelectionStateChange={onRowSelectionStateChange}
        />
      </div>
    </div>
  );
}
