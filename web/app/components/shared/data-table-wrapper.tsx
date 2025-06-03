import { useEffect } from "react";
import { cn } from "~/lib/utils";
import { DataTable } from "~/routes/collections/data-table";
import type {
  ColumnDef,
  OnChangeFn,
  PaginationState,
} from "@tanstack/react-table";

interface DataTableWrapperProps<TData, TValue> {
  columns: ColumnDef<TData, TValue>[];
  data: TData[];
  isLoading?: boolean;
  className?: string;
  emptyMessage?: string;
  pagination: PaginationState;
  onPaginationChange: OnChangeFn<PaginationState>;
  totalRows: number;
}

export function DataTableWrapper<TData, TValue>({
  columns,
  data,
  isLoading = false,
  className,
  emptyMessage = "No results.",
  pagination,
  onPaginationChange,
  totalRows,
}: DataTableWrapperProps<TData, TValue>) {
  return (
    <div
      className={cn(
        "rounded-md overflow-hidden flex flex-col h-full min-h-0 max-h-full",
        isLoading ? "animate-pulse-border" : "border",
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
      />
    </div>
  );
}
