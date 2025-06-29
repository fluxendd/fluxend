import { format } from "date-fns";
import {
  Sheet,
  SheetContent,
  SheetDescription,
  SheetHeader,
  SheetTitle,
} from "~/components/ui/sheet";
import { Badge } from "~/components/ui/badge";
import { Separator } from "~/components/ui/separator";
import { ScrollArea } from "~/components/ui/scroll-area";
import { Button } from "~/components/ui/button";
import { Copy, Check } from "lucide-react";
import type { LogEntry } from "~/services/logs";
import { cn } from "~/lib/utils";
import { useState, useEffect, useRef } from "react";

interface LogDetailSheetProps {
  log: LogEntry | null;
  open: boolean;
  onOpenChange: (open: boolean) => void;
}

const getStatusInfo = (status: number) => {
  if (status >= 200 && status < 300) {
    return { color: "text-green-600", bg: "bg-green-100" };
  } else if (status >= 400 && status < 500) {
    return { color: "text-yellow-600", bg: "bg-yellow-100" };
  } else if (status >= 500) {
    return { color: "text-red-600", bg: "bg-red-100" };
  }
  return { color: "text-gray-600", bg: "bg-gray-100" };
};

// Copy button component for sheet
const SheetCopyButton = ({ text, label }: { text: string; label: string }) => {
  const [copied, setCopied] = useState(false);
  const timeoutRef = useRef<NodeJS.Timeout | null>(null);
  
  const handleCopy = async () => {
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
    <button
      className="inline-flex items-center justify-center h-5 w-5 ml-2 rounded-sm transition-opacity cursor-pointer"
      onClick={handleCopy}
      aria-label={`Copy ${label}`}
    >
      <div className="relative w-3 h-3">
        <Copy 
          className={cn(
            "h-3 w-3 absolute inset-0 transition-all duration-200 opacity-70",
            copied ? "opacity-0 scale-50" : "opacity-100 scale-100"
          )}
        />
        <Check 
          className={cn(
            "h-3 w-3 text-green-600 absolute inset-0 transition-all duration-200",
            copied ? "opacity-100 scale-100" : "opacity-0 scale-50"
          )}
        />
      </div>
    </button>
  );
};

// Helper function to parse JSON body if it's a string
const parseJsonBody = (body: any) => {
  if (!body) return null;
  
  // If it's already an object, return as is
  if (typeof body === 'object') {
    return body;
  }
  
  // If it's a string, try to parse it
  if (typeof body === 'string') {
    try {
      return JSON.parse(body);
    } catch (e) {
      // If parsing fails, return the original string
      return body;
    }
  }
  
  return body;
};

export function LogDetailSheet({ log, open, onOpenChange }: LogDetailSheetProps) {
  if (!log) return null;

  const statusInfo = getStatusInfo(log.status);
  const formattedDate = format(new Date(log.createdAt), "PPpp");
  const parsedBody = parseJsonBody(log.body);
  const parsedParams = parseJsonBody(log.params);

  return (
    <Sheet open={open} onOpenChange={onOpenChange}>
      <SheetContent className="sm:max-w-[600px]">
        <SheetHeader className="px-6">
          <SheetTitle>Log Details</SheetTitle>
          <SheetDescription>
            {formattedDate}
          </SheetDescription>
        </SheetHeader>

        <ScrollArea className="h-[calc(100vh-8rem)] mt-6">
          <div className="space-y-6 px-6">
            {/* Request Info */}
            <div>
              <h3 className="text-sm font-semibold mb-3">Request Information</h3>
              <div className="rounded-lg border">
                <table className="w-full">
                  <tbody>
                    <tr className="border-b">
                      <td className="px-4 py-3 text-sm font-medium text-muted-foreground w-32">
                        Method
                      </td>
                      <td className="px-4 py-3">
                        <Badge variant="outline" className="font-mono">
                          {log.method}
                        </Badge>
                      </td>
                    </tr>
                    <tr className="border-b">
                      <td className="px-4 py-3 text-sm font-medium text-muted-foreground w-32">
                        Endpoint
                      </td>
                      <td className="px-4 py-3">
                        <div className="flex items-center">
                          <Badge variant="outline" className="font-mono max-w-full overflow-hidden">
                            <span className="truncate block">{log.endpoint}</span>
                          </Badge>
                          <SheetCopyButton text={log.endpoint} label="endpoint" />
                        </div>
                      </td>
                    </tr>
                    <tr className="border-b">
                      <td className="px-4 py-3 text-sm font-medium text-muted-foreground w-32">
                        Status
                      </td>
                      <td className="px-4 py-3">
                        <Badge 
                          variant="outline"
                          className={cn(
                            "font-mono",
                            log.status >= 200 && log.status < 300 && "border-green-600 text-green-600",
                            log.status >= 400 && log.status < 500 && "border-yellow-600 text-yellow-600",
                            log.status >= 500 && "border-red-600 text-red-600"
                          )}
                        >
                          {log.status}
                        </Badge>
                      </td>
                    </tr>
                    <tr>
                      <td className="px-4 py-3 text-sm font-medium text-muted-foreground w-32">
                        Identifier
                      </td>
                      <td className="px-4 py-3">
                        <div className="flex items-center">
                          <Badge variant="outline" className="font-mono text-xs">
                            {log.uuid}
                          </Badge>
                          <SheetCopyButton text={log.uuid} label="identifier" />
                        </div>
                      </td>
                    </tr>
                  </tbody>
                </table>
              </div>
            </div>

            <Separator />

            {/* Client Info */}
            <div>
              <h3 className="text-sm font-semibold mb-3">Client Information</h3>
              <div className="rounded-lg border">
                <table className="w-full">
                  <tbody>
                    <tr className="border-b">
                      <td className="px-4 py-3 text-sm font-medium text-muted-foreground w-32">
                        IP Address
                      </td>
                      <td className="px-4 py-3">
                        <div className="flex items-center">
                          <Badge variant="outline" className="font-mono">
                            {log.ipAddress}
                          </Badge>
                          <SheetCopyButton text={log.ipAddress} label="IP address" />
                        </div>
                      </td>
                    </tr>
                    {log.userUuid && (
                      <tr className="border-b">
                        <td className="px-4 py-3 text-sm font-medium text-muted-foreground w-32">
                          User ID
                        </td>
                        <td className="px-4 py-3">
                          <div className="flex items-center">
                            <Badge variant="outline" className="font-mono text-xs max-w-full overflow-hidden">
                              <span className="truncate block">{log.userUuid}</span>
                            </Badge>
                            <SheetCopyButton text={log.userUuid} label="user ID" />
                          </div>
                        </td>
                      </tr>
                    )}
                    <tr>
                      <td className="px-4 py-3 text-sm font-medium text-muted-foreground w-32">
                        User Agent
                      </td>
                      <td className="px-4 py-3">
                        <div className="text-sm break-words overflow-wrap-anywhere">{log.userAgent}</div>
                      </td>
                    </tr>
                  </tbody>
                </table>
              </div>
            </div>

            {/* Request Body */}
            {parsedBody && (typeof parsedBody === 'object' ? Object.keys(parsedBody).length > 0 : true) && (
              <>
                <Separator />
                <div>
                  <h3 className="text-sm font-semibold mb-3">Request Body</h3>
                  <pre className="bg-muted p-3 rounded-lg overflow-x-auto text-xs whitespace-pre-wrap break-all">
                    {typeof parsedBody === 'string' ? parsedBody : JSON.stringify(parsedBody, null, 2)}
                  </pre>
                </div>
              </>
            )}

            {/* Query Parameters */}
            {parsedParams && (typeof parsedParams === 'object' ? Object.keys(parsedParams).length > 0 : true) && (
              <>
                <Separator />
                <div>
                  <h3 className="text-sm font-semibold mb-3">Query Parameters</h3>
                  <pre className="bg-muted p-3 rounded-lg overflow-x-auto text-xs whitespace-pre-wrap break-all">
                    {typeof parsedParams === 'string' ? parsedParams : JSON.stringify(parsedParams, null, 2)}
                  </pre>
                </div>
              </>
            )}

          </div>
        </ScrollArea>
      </SheetContent>
    </Sheet>
  );
}