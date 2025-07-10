import { useEffect } from "react";
import { cn } from "~/lib/utils";
import { DataTable } from "~/routes/tables/data-table";
import { QuerySearchBox } from "./query-search-box";
import type {
  ColumnDef,
  OnChangeFn,
  PaginationState,
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
  tableMeta?: any;
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
  onFilterChange,
  tableMeta,
}: SearchDataTableWrapperProps<TData, TValue>) {
  return (
    <div className="flex flex-col h-full gap-4 px-4">
      <QuerySearchBox columns={rawColumns} onQueryChange={onFilterChange} />
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
        />
      </div>
    </div>
  );
}
