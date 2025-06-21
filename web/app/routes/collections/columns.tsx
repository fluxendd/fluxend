import {
  Hash,
  AlignLeft,
  Clock,
  Text,
  MoreVertical,
  RotateCcw,
  Circle,
  ToggleLeft,
  Calendar,
  FileJson,
  Shuffle,
} from "lucide-react";
import type { ColumnDef, PaginationState, Row } from "@tanstack/react-table";
import { Button } from "~/components/ui/button";
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuTrigger,
  DropdownMenuSeparator,
  DropdownMenuLabel,
} from "~/components/ui/dropdown-menu";

import React, { useState } from "react";
import { Switch } from "~/components/ui/switch";
import {
  Tooltip,
  TooltipContent,
  TooltipTrigger,
} from "~/components/ui/tooltip";
import { formatTimestamp, getTypedResponseData } from "~/lib/utils";
import type { APIResponse } from "~/lib/types";
import { getClientAuthToken } from "~/lib/auth";
import { initializeServices } from "~/services";

export const columnsQuery = (projectId: string, collectionId: string) => ({
  queryKey: ["columns", projectId, collectionId],
  queryFn: async () => {
    const authToken = await getClientAuthToken();
    if (!authToken) {
      throw new Error("Unauthorized");
    }

    const services = initializeServices(authToken);

    const res = await services.collections.getTableColumns(
      projectId,
      collectionId
    );

    if (!res.ok && res.status > 500) {
      throw new Error("Unexpected Error");
    }

    const data = await getTypedResponseData<APIResponse<any>>(res);

    if (!res.ok) {
      throw new Error(data.errors?.[0] || "Unexpected Error");
    }

    return data.content;
  },
});

// Mock function to handle row deletion
const mockDeleteRow = (rowId: string) => {
  console.log(`Delete row with ID: ${rowId}`);
  // In a real implementation, you would call an API endpoint here
  // and potentially use queryClient.invalidateQueries to refresh data
};

enum ColumnType {
  Integer = "integer",
  Serial = "serial",
  Varchar = "varchar",
  Text = "text",
  CharacterVarying = "character varying(255)",
  Boolean = "boolean",
  Date = "date",
  Timestamp = "timestamp",
  TimestampWithoutTimeZone = "timestamp without time zone",
  Float = "float",
  UUID = "uuid",
  JSON = "json",
}

interface ColumnIconProps {
  type: ColumnType | string;
}

const ColumnIcon: React.FC<ColumnIconProps> = ({ type }) => {
  let Icon = Circle;

  console.log("Column type:", type);

  switch (type) {
    case ColumnType.Integer:
      Icon = Hash;
      break;
    case ColumnType.Serial:
      Icon = Hash;
      break;
    case ColumnType.Varchar:
      Icon = Text;
      break;
    case ColumnType.CharacterVarying:
      Icon = Text;
      break;
    case ColumnType.Text:
      Icon = AlignLeft;
      break;
    case ColumnType.Boolean:
      Icon = ToggleLeft;
      break;
    case ColumnType.Date:
      Icon = Calendar;
      break;
    case ColumnType.Timestamp:
      Icon = Clock;
      break;
    case ColumnType.TimestampWithoutTimeZone:
      Icon = Clock;
      break;
    case ColumnType.Float:
      Icon = Hash;
      break;
    case ColumnType.UUID:
      Icon = Shuffle;
      break;
    case ColumnType.JSON:
      Icon = FileJson;
      break;
    default:
      Icon = Circle;
      break;
  }

  return <Icon className="mr-2 h-3 w-3 flex-shrink-0 translate-y-[0.5px]" />;
};

export const prepareColumns = (
  columns: any[] | undefined | null,
  collectionName?: string
): ColumnDef<any>[] => {
  if (!columns || !Array.isArray(columns) || columns.length === 0) {
    return [];
  }

  // // Create the actions column with column visibility controls
  const ActionsColumnHeader = React.memo(({ table }: { table: any }) => {
    // State to control dropdown open/close
    const [dropdownOpen, setDropdownOpen] = useState(false);

    return (
      <div className="flex items-center justify-end pr-6">
        <DropdownMenu
          open={dropdownOpen}
          onOpenChange={setDropdownOpen}
          modal={false}
        >
          <DropdownMenuTrigger asChild>
            <Button
              variant="ghost"
              size="sm"
              className="h-8 w-8 p-0 rounded-full hover:bg-muted data-[state=open]:bg-muted hover:shadow-sm"
            >
              <MoreVertical className="h-4 w-4" />
              <span className="sr-only">Column options</span>
            </Button>
          </DropdownMenuTrigger>
          <DropdownMenuContent
            align="end"
            className="w-[220px]"
            onCloseAutoFocus={(e) => {
              e.preventDefault();
              // Prevent closing when interacting with content
              e.stopPropagation();
            }}
            onEscapeKeyDown={() => setDropdownOpen(false)}
            forceMount
          >
            <div className="flex items-center justify-between px-2 py-1.5">
              <DropdownMenuLabel className="px-0 py-0">
                Column Visibility
              </DropdownMenuLabel>
              <Button
                variant="ghost"
                size="sm"
                className="h-8 w-8 p-0 rounded-full hover:bg-muted"
                onClick={(e) => {
                  table.resetColumnVisibility();
                  e.stopPropagation();
                }}
                title="Reset all columns to visible"
              >
                <RotateCcw className="h-3.5 w-3.5" />
              </Button>
            </div>
            <DropdownMenuSeparator />
            {Array.isArray(table.getAllColumns())
              ? table.getAllColumns().map((column) => {
                  if (!column) return null;

                  return (
                    <DropdownMenuItem
                      key={column.id}
                      className="flex items-center justify-between py-2 px-2"
                      onSelect={(e) => {
                        e.preventDefault();
                        e.stopPropagation();
                      }}
                    >
                      <div className="flex items-center gap-2">
                        <span>{column.id}</span>
                      </div>
                      <Switch
                        defaultChecked={column.getIsVisible()}
                        disabled={!column.getCanHide()}
                        onCheckedChange={column.toggleVisibility}
                      />
                    </DropdownMenuItem>
                  );
                })
              : null}
          </DropdownMenuContent>
        </DropdownMenu>
      </div>
    );
  });

  const actionsColumn: ColumnDef<any, unknown> = {
    id: "actions",
    accessorKey: "actions",
    header: ({ table }) => (
      <ActionsColumnHeader key="actions-column-header" table={table} />
    ),
    enableHiding: false,
    cell: ({ row }) => {
      return (
        <div className="flex justify-end pr-6">
          <DropdownMenu>
            <DropdownMenuTrigger asChild>
              <Button
                variant="ghost"
                size="sm"
                className="h-8 w-8 p-0 rounded-full hover:bg-muted data-[state=open]:bg-muted hover:shadow-sm"
              >
                <MoreVertical className="h-4 w-4" />
                <span className="sr-only">Open menu</span>
              </Button>
            </DropdownMenuTrigger>
            <DropdownMenuContent align="end" className="w-[160px]">
              <DropdownMenuItem
                onClick={() =>
                  mockDeleteRow(String(row.original?.id || row.id))
                }
                className="text-destructive focus:text-destructive"
              >
                Delete
              </DropdownMenuItem>
            </DropdownMenuContent>
          </DropdownMenu>
        </div>
      );
    },
    meta: {
      isSticky: true,
      isEven: false,
    },
  };

  // Create the data columns
  const dataColumns = Array.isArray(columns)
    ? columns.filter(Boolean).map((column, index) => {
        const isEvenColumn = index % 2 === 0;
        const dataType = column.type || "";
        const isTimestamp = dataType.includes("timestamp");
        const columnId = column.name;

        return {
          id: columnId,
          accessorKey: columnId,
          accessorFn: (row: any) =>
            row && column && column.name ? row[column.name] : null,
          header: () => (
            <div className="inline-flex items-center whitespace-nowrap">
              <ColumnIcon type={column.type} />
              <span className="truncate font-medium">
                {column && column.name ? column.name : ""}
              </span>
            </div>
          ),
          cell: ({ row }: { row: Row<any> }) => {
            // Get value safely
            let value = null;
            try {
              if (row && column && column.name) {
                value = row.getValue(column.name);
              }
            } catch (e) {
              // Handle error silently
            }

            // Format timestamp values
            if (isTimestamp && value !== null && value !== undefined) {
              const { date, time, fullDate, relativeTime } = formatTimestamp(
                String(value)
              );
              return (
                <Tooltip>
                  <TooltipTrigger asChild>
                    <div className="flex flex-col cursor-default hover:bg-muted/50 p-1 rounded-sm transition-colors">
                      <span className="text-xs font-medium text-foreground whitespace-nowrap">
                        {date}
                      </span>
                      <span className="text-xs text-muted-foreground/80 flex gap-1 items-center whitespace-nowrap">
                        <Clock className="h-2.5 w-2.5 inline-block opacity-70 flex-shrink-0" />
                        {time}
                      </span>
                    </div>
                  </TooltipTrigger>
                  <TooltipContent
                    sideOffset={5}
                    className="bg-popover text-popover-foreground border border-border shadow-md p-3 text-xs max-w-[240px]"
                  >
                    <div className="font-semibold mb-1 flex items-center gap-1">
                      {date} {time}
                      <Clock className="h-3 w-3 opacity-50" />
                    </div>
                    <div className="text-muted-foreground text-[11px]">
                      {relativeTime}
                    </div>
                    <div className="text-[10px] text-muted-foreground/70 mt-1 break-all">
                      {fullDate}
                    </div>
                  </TooltipContent>
                </Tooltip>
              );
            }

            return (
              <span>
                {value !== null && value !== undefined ? String(value) : ""}
              </span>
            );
          },
          meta: {
            isEven: isEvenColumn,
            collectionName,
          },
        };
      })
    : [];

  return [...dataColumns, actionsColumn];
};

export const rowsQuery = (
  projectId: string,
  dbId: string,
  collectionName: string,
  pagination: PaginationState = { pageIndex: 0, pageSize: 50 },
  filters: Record<string, string> = {}
) => ({
  queryKey: [
    "rows",
    projectId,
    collectionName,
    pagination.pageSize,
    pagination.pageIndex,
    filters,
  ],
  queryFn: async () => {
    if (!collectionName) {
      return { rows: [], totalCount: 0 };
    }

    const baseDomain = import.meta.env.VITE_FLX_BASE_DOMAIN;
    const httpScheme = import.meta.env.VITE_FLX_HTTP_SCHEME;

    const offset = pagination.pageIndex * pagination.pageSize;
    const limit = pagination.pageSize;

    const authToken = await getClientAuthToken();

    if (!authToken) {
      throw new Error("No auth token");
    }

    const services = initializeServices(authToken);

    try {
      const res = await services.collections.getTableRows(
        projectId,
        collectionName,
        {
          headers: {
            Prefer: "count=exact",
          },
          params: {
            limit,
            offset,
            ...filters,
          },
          baseUrl: `${httpScheme}://${dbId}.${baseDomain}/`,
        }
      );

      if (!res.ok) {
        const responseData = await res.json();
        const errorMessage = responseData?.errors[0] || "Unknown error";

        if (res.status === 401) {
          // throw new UnauthorizedError(errorMessage);
        } else {
          throw new Error(errorMessage);
        }
      }

      // Extract total count from Content-Range header
      let totalCount = 0;
      const contentRange = res.headers.get("Content-Range");
      if (contentRange) {
        // Format is typically like "0-49/100" or "0-49/*"
        const match = contentRange.match(/\/(\d+|\*)/);
        if (match && match[1] && match[1] !== "*") {
          totalCount = parseInt(match[1], 10);
        }
      }

      const data = await res.json();
      return {
        rows: Array.isArray(data) ? data : [],
        totalCount,
      };
    } catch (error) {
      console.error("Error fetching rows:", error);
      return { rows: [], totalCount: 0 };
    }
  },
});
