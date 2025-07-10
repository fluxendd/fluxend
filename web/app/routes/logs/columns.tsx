import {
  FileText,
  Globe,
  Hash,
  Clock,
  AlertCircle,
  CheckCircle,
  XCircle,
  Check,
  Copy,
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
import type { LogEntry, HttpMethod, HttpStatusCode } from "~/services/logs";
import { useState, useEffect, useRef } from "react";
import { Button } from "~/components/ui/button";

// Copy indicator component
const CopyIndicator = ({ text, label }: { text: string; label: string }) => {
  const [copied, setCopied] = useState(false);
  const timeoutRef = useRef<NodeJS.Timeout | null>(null);
  
  const handleCopy = async (e: React.MouseEvent) => {
    e.stopPropagation();
    try {
      await navigator.clipboard.writeText(text);
      setCopied(true);
      
      if (timeoutRef.current) {
        clearTimeout(timeoutRef.current);
      }
      
      timeoutRef.current = setTimeout(() => setCopied(false), 2000);
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
    <div className="group flex items-center gap-1 -m-2 p-2">
      {label === "IP Address" ? (
        <Badge variant="outline" className="font-mono">
          {text}
        </Badge>
      ) : (
        <span className="font-mono text-sm truncate">
          {text}
        </span>
      )}
      <Button
        variant="ghost"
        size="icon"
        className="h-6 w-6 p-0 transition-all duration-200 opacity-0 scale-75 group-hover:opacity-100 group-hover:scale-100"
        onClick={handleCopy}
        data-no-row-click
      >
        {copied ? (
          <Check className="h-3 w-3 text-green-600" />
        ) : (
          <Copy className="h-3 w-3" />
        )}
      </Button>
    </div>
  );
};

// Type guard to check if value is an object
const isObject = (value: unknown): value is Record<string, any> => {
  return typeof value === 'object' && value !== null && !Array.isArray(value);
};

// Helper function to safely parse JSON strings
const parseJsonSafely = (value: string | Record<string, any>): Record<string, any> | null => {
  if (isObject(value)) return value;
  if (typeof value !== 'string') return null;
  
  try {
    const parsed = JSON.parse(value);
    return isObject(parsed) ? parsed : null;
  } catch {
    return null;
  }
};

// Helper function to get status color and icon
const getStatusInfo = (status: HttpStatusCode) => {
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
      const timestamp = row.getValue<string>("createdAt");
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
      const method = row.getValue<HttpMethod>("method");
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
      const endpoint = row.getValue<string>("endpoint");
      
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
      const status = row.getValue<HttpStatusCode>("status");
      const statusInfo = getStatusInfo(status);
      
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
      const ip = row.getValue<string>("ipAddress");
      return (
        <div className="-m-2 p-2">
          <CopyIndicator text={ip} label="IP Address" />
        </div>
      );
    },
  },
  {
    accessorKey: "userAgent",
    size: 200,
    header: "User Agent",
    cell: ({ row }) => {
      const userAgent = row.getValue<string>("userAgent");
      const shortAgent = userAgent?.split(" ")[0] || userAgent || "Unknown";
      
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
            <div className="text-xs break-all">{userAgent || "No user agent"}</div>
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
      const bodyData = parseJsonSafely(row.original.body);
      const paramsData = parseJsonSafely(row.original.params);
      
      const hasBody = bodyData && Object.keys(bodyData).length > 0;
      const hasParams = paramsData && Object.keys(paramsData).length > 0;
      
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
                  {JSON.stringify(bodyData, null, 2)}
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
                  {JSON.stringify(paramsData, null, 2)}
                </pre>
              </TooltipContent>
            </Tooltip>
          )}
        </div>
      );
    },
  },
];