import {
  FileText,
  Globe,
  Hash,
  Clock,
  AlertCircle,
  CheckCircle,
  XCircle,
  Check,
} from "lucide-react";
import type { ColumnDef } from "@tanstack/react-table";
import { Badge } from "~/components/ui/badge";
import {
  Tooltip,
  TooltipContent,
  TooltipTrigger,
} from "~/components/ui/tooltip";
import { formatTimestamp } from "~/lib/utils";
import { cn } from "~/lib/utils";
import type { LogEntry } from "~/services/logs";
import { useState, useEffect, useRef } from "react";

// Copy indicator component
const CopyIndicator = ({ text, label }: { text: string; label: string }) => {
  const [copied, setCopied] = useState(false);
  const timeoutRef = useRef<NodeJS.Timeout | null>(null);
  
  const handleClick = async (e: React.MouseEvent) => {
    e.stopPropagation();
    try {
      await navigator.clipboard.writeText(text);
      setCopied(true);
      
      if (timeoutRef.current) {
        clearTimeout(timeoutRef.current);
      }
      
      timeoutRef.current = setTimeout(() => setCopied(false), 5000);
    } catch (err) {
      console.error('Failed to copy:', err);
    }
  };
  
  useEffect(() => {
    return () => {
      if (timeoutRef.current) {
        clearTimeout(timeoutRef.current);
      }
    };
  }, []);
  
  return (
    <div className="flex items-center gap-1 -m-2 p-2" data-no-row-click>
      {label === "IP Address" ? (
        <Badge variant="outline" className="font-mono cursor-pointer" onClick={handleClick}>
          {text}
        </Badge>
      ) : (
        <span 
          className="font-mono text-sm truncate cursor-pointer"
          onClick={handleClick}
        >
          {text}
        </span>
      )}
      <div className="w-3 h-3 flex-shrink-0 relative">
        <Check 
          className={cn(
            "h-3 w-3 text-green-600 absolute inset-0 transition-all duration-300",
            copied 
              ? "opacity-100 scale-100" 
              : "opacity-0 scale-50"
          )}
        />
      </div>
    </div>
  );
};

// Helper function to get status color and icon
const getStatusInfo = (status: number) => {
  if (status >= 200 && status < 300) {
    return { color: "text-green-600", bg: "bg-green-50", Icon: CheckCircle };
  } else if (status >= 400 && status < 500) {
    return { color: "text-yellow-600", bg: "bg-yellow-50", Icon: AlertCircle };
  } else if (status >= 500) {
    return { color: "text-red-600", bg: "bg-red-50", Icon: XCircle };
  }
  return { color: "text-gray-600", bg: "bg-gray-50", Icon: Hash };
};


export const createLogsColumns = (): ColumnDef<LogEntry>[] => [
  {
    accessorKey: "createdAt",
    size: 150,
    header: () => (
      <div className="flex items-center gap-2">
        <Clock className="h-3 w-3" />
        <span>Timestamp</span>
      </div>
    ),
    cell: ({ row }) => {
      const timestamp = row.getValue("createdAt") as string;
      const { date, time, relativeTime } = formatTimestamp(timestamp);
      
      return (
        <Tooltip>
          <TooltipTrigger asChild>
            <div className="-m-2 p-2">
              <div className="text-sm font-medium">{date}</div>
              <div className="text-xs text-muted-foreground">{time}</div>
            </div>
          </TooltipTrigger>
          <TooltipContent>
            <div className="text-xs">{relativeTime}</div>
          </TooltipContent>
        </Tooltip>
      );
    },
  },
  {
    accessorKey: "method",
    size: 80,
    header: "Method",
    cell: ({ row }) => {
      const method = row.getValue("method") as string;
      return (
        <div className="-m-2 p-2">
          <Badge variant="outline" className="font-mono">
            {method}
          </Badge>
        </div>
      );
    },
  },
  {
    accessorKey: "endpoint",
    size: 300,
    header: () => (
      <div className="flex items-center gap-2">
        <FileText className="h-3 w-3" />
        <span>Endpoint</span>
      </div>
    ),
    cell: ({ row }) => {
      const endpoint = row.getValue("endpoint") as string;
      
      return (
        <div className="-m-2 p-2">
          <CopyIndicator text={endpoint} label="Endpoint" />
        </div>
      );
    },
  },
  {
    accessorKey: "status",
    size: 80,
    header: "Status",
    cell: ({ row }) => {
      const status = row.getValue("status") as number;
      
      return (
        <div className="-m-2 p-2">
          <Badge 
            variant="outline" 
            className={cn(
              "font-mono",
              status >= 200 && status < 300 && "border-green-600 text-green-600",
              status >= 400 && status < 500 && "border-yellow-600 text-yellow-600",
              status >= 500 && "border-red-600 text-red-600"
            )}
          >
            {status}
          </Badge>
        </div>
      );
    },
  },
  {
    accessorKey: "ipAddress",
    size: 120,
    header: () => (
      <div className="flex items-center gap-2">
        <Globe className="h-3 w-3" />
        <span>IP Address</span>
      </div>
    ),
    cell: ({ row }) => {
      const ip = row.getValue("ipAddress") as string;
      return <CopyIndicator text={ip} label="IP Address" />;
    },
  },
  {
    accessorKey: "userAgent",
    size: 200,
    header: "User Agent",
    cell: ({ row }) => {
      const userAgent = row.getValue("userAgent") as string;
      const shortAgent = userAgent?.split(" ")[0] || userAgent;
      
      return (
        <Tooltip>
          <TooltipTrigger asChild>
            <div 
              className="text-sm truncate max-w-[200px] cursor-pointer -m-2 p-2"
            >
              {shortAgent}
            </div>
          </TooltipTrigger>
          <TooltipContent className="max-w-[500px]">
            <div className="text-xs break-all">{userAgent}</div>
          </TooltipContent>
        </Tooltip>
      );
    },
  },
  {
    id: "details",
    size: 100,
    header: "Details",
    cell: ({ row }) => {
      const hasBody = row.original.body && Object.keys(row.original.body).length > 0;
      const hasParams = row.original.params && Object.keys(row.original.params).length > 0;
      
      if (!hasBody && !hasParams) {
        return (
          <div 
            className="cursor-pointer -m-2 p-2"
          >
            <span className="text-xs text-muted-foreground">No data</span>
          </div>
        );
      }
      
      return (
        <div 
          className="flex gap-2 cursor-pointer -m-2 p-2"
        >
          {hasBody && (
            <Tooltip>
              <TooltipTrigger asChild>
                <Badge variant="outline" className="text-xs">
                  Body
                </Badge>
              </TooltipTrigger>
              <TooltipContent className="max-w-[400px]">
                <pre className="text-xs">
                  {JSON.stringify(row.original.body, null, 2)}
                </pre>
              </TooltipContent>
            </Tooltip>
          )}
          {hasParams && (
            <Tooltip>
              <TooltipTrigger asChild>
                <Badge variant="outline" className="text-xs">
                  Params
                </Badge>
              </TooltipTrigger>
              <TooltipContent className="max-w-[400px]">
                <pre className="text-xs">
                  {JSON.stringify(row.original.params, null, 2)}
                </pre>
              </TooltipContent>
            </Tooltip>
          )}
        </div>
      );
    },
  },
];