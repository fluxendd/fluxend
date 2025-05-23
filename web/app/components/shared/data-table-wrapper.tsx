import { useEffect } from "react";
import { cn } from "~/lib/utils";
import { DataTable } from "~/routes/collections/data-table";
import type {
  ColumnDef,
  OnChangeFn,
  PaginationState,
} from "@tanstack/react-table";

// Add global styles for the border animation
const addBorderAnimationStyles = () => {
  if (
    typeof document !== "undefined" &&
    !document.getElementById("data-table-animation-style")
  ) {
    const style = document.createElement("style");
    style.id = "data-table-animation-style";
    style.textContent = `
      @keyframes pulse-border {
        0% { border-color: rgba(99, 102, 241, 0.3); }
        50% { border-color: rgba(99, 102, 241, 0.9); }
        100% { border-color: rgba(99, 102, 241, 0.3); }
      }
      .loading-border {
        animation: pulse-border 1.5s infinite ease-in-out;
        border: 1px solid rgba(99, 102, 241, 0.7);
      }
    `;
    document.head.appendChild(style);
  }
};

// Add styles when the component is first imported
if (typeof document !== "undefined") {
  addBorderAnimationStyles();
}

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
  // Ensure styles are added when component is mounted
  useEffect(() => {
    addBorderAnimationStyles();
  }, []);

  return (
    <div
      className={cn(
        "rounded-md overflow-hidden flex flex-col h-full min-h-0 max-h-full",
        isLoading ? "loading-border" : "border",
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
        isLoading={isLoading}
      />
    </div>
  );
}
