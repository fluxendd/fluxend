import React, { 
  memo, 
  useState, 
  useRef, 
  useEffect, 
  useCallback
} from "react";
import { Input } from "~/components/ui/input";
import { Textarea } from "~/components/ui/textarea";
import { Checkbox } from "~/components/ui/checkbox";
import { Calendar } from "~/components/ui/calendar";
import {
  Popover,
  PopoverContent,
  PopoverTrigger,
} from "~/components/ui/popover";
import { Button } from "~/components/ui/button";
import { CalendarIcon, Check } from "lucide-react";
import { format } from "date-fns";
import { cn } from "~/lib/utils";
import { formatTimestamp } from "~/lib/utils";
import {
  Tooltip,
  TooltipContent,
  TooltipTrigger,
} from "~/components/ui/tooltip";
import { Clock } from "lucide-react";
import { toast } from "sonner";
import type { TableColumnMetadata, TableRow, CellValue } from "~/types/table";

interface OptimisticTableCellProps {
  value: CellValue;
  column: TableColumnMetadata;
  rowId: string;
  rowData: TableRow;
  isTimestamp: boolean;
  onUpdate: (rowId: string, columnName: string, value: CellValue, rowData: TableRow) => Promise<boolean>;
}

export const OptimisticTableCell = memo(({ 
  value: initialValue, 
  column, 
  rowId, 
  rowData,
  isTimestamp,
  onUpdate 
}: OptimisticTableCellProps) => {
  const [isEditing, setIsEditing] = useState(false);
  const [localValue, setLocalValue] = useState(initialValue);
  const [optimisticValue, setOptimisticValue] = useState(initialValue);
  const [jsonError, setJsonError] = useState("");
  const [isPending, setIsPending] = useState(false);
  const [showSuccess, setShowSuccess] = useState(false);
  const successTimeoutRef = useRef<NodeJS.Timeout | null>(null);
  
  const inputRef = useRef<HTMLInputElement | HTMLTextAreaElement | null>(null);
  const containerRef = useRef<HTMLDivElement>(null);
  
  const dataType = (column.type || "").toLowerCase();
  const isNullable = column.is_nullable === "YES";
  const isReadOnly = column.name === "id" || 
                    column.name === "created_at" || 
                    column.name === "updated_at";

  // Cleanup timeout on unmount
  useEffect(() => {
    return () => {
      if (successTimeoutRef.current) {
        clearTimeout(successTimeoutRef.current);
      }
    };
  }, []);

  useEffect(() => {
    // Sync with server value when it changes
    setOptimisticValue(initialValue);
    if (!isEditing) {
      setLocalValue(initialValue);
    }
  }, [initialValue]);

  useEffect(() => {
    if (isEditing && inputRef.current) {
      inputRef.current.focus();
      if ('select' in inputRef.current) {
        inputRef.current.select();
      }
    }
  }, [isEditing]);

  useEffect(() => {
    if (isEditing) {
      const handleClickOutside = (event: MouseEvent) => {
        if (containerRef.current && !containerRef.current.contains(event.target as Node)) {
          handleSave();
        }
      };

      document.addEventListener('mousedown', handleClickOutside);
      return () => {
        document.removeEventListener('mousedown', handleClickOutside);
      };
    }
  }, [isEditing]); // eslint-disable-line react-hooks/exhaustive-deps
  // Note: handleSave is intentionally not in deps to avoid re-registering

  const handleDoubleClick = () => {
    if (!isReadOnly && !isEditing && !isPending) {
      setIsEditing(true);
      setLocalValue(optimisticValue);
    }
  };

  const saveValue = useCallback(async (valueToSave: CellValue) => {
    // Validate required fields
    if (!isNullable && (valueToSave === null || valueToSave === undefined || valueToSave === '')) {
      toast.error(`${column.name} is required and cannot be empty`);
      setIsEditing(false);
      setLocalValue(optimisticValue);
      return;
    }

    // Validate JSON fields
    if (dataType.includes("json") && valueToSave) {
      try {
        JSON.parse(typeof valueToSave === 'string' ? valueToSave : JSON.stringify(valueToSave));
      } catch (e) {
        const error = e as Error;
        toast.error(`Invalid JSON: ${error.message}`);
        return;
      }
    }

    // Skip if value hasn't changed (deep comparison for objects)
    const hasChanged = JSON.stringify(valueToSave) !== JSON.stringify(initialValue);
    if (!hasChanged) {
      return;
    }

    if (!onUpdate) {
      toast.error("Unable to update: Table is not properly configured");
      return;
    }

    // Optimistically update the value immediately
    setOptimisticValue(valueToSave);
    setLocalValue(valueToSave);
    setIsPending(true);

    // Make the API call
    onUpdate(rowId, column.name, valueToSave, rowData)
      .then(success => {
        setIsPending(false);
        if (!success) {
          // Revert on failure
          setOptimisticValue(initialValue);
          setLocalValue(initialValue);
          // Error toast is now handled in the parent component
        } else {
          // Show success indicator
          setShowSuccess(true);
          if (successTimeoutRef.current) {
            clearTimeout(successTimeoutRef.current);
          }
          successTimeoutRef.current = setTimeout(() => {
            setShowSuccess(false);
            successTimeoutRef.current = null;
          }, 2000);
        }
      })
      .catch(error => {
        setIsPending(false);
        // Revert on error
        setOptimisticValue(initialValue);
        setLocalValue(initialValue);
        // Error handling is done in the parent component
      });
  }, [dataType, initialValue, onUpdate, rowId, column.name, rowData, isNullable]);

  const handleSave = () => {
    if (!isEditing) return;
    setIsEditing(false);
    // Make sure we save the current local value
    const valueToSave = localValue;
    saveValue(valueToSave);
  };

  const handleCancel = () => {
    setIsEditing(false);
    setLocalValue(optimisticValue);
    setJsonError("");
  };

  const handleKeyDown = (e: React.KeyboardEvent) => {
    if (e.key === 'Enter' && !e.shiftKey) {
      e.preventDefault();
      handleSave();
    } else if (e.key === 'Escape') {
      handleCancel();
    }
  };

  // For timestamp fields in non-edit mode
  if (!isEditing && isTimestamp && optimisticValue !== null && optimisticValue !== undefined) {
    const { date, time, fullDate, relativeTime } = formatTimestamp(String(optimisticValue));
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

  // Display mode
  if (!isEditing) {
    if (dataType === "boolean" || dataType.includes("bool")) {
      return (
        <div 
          onDoubleClick={handleDoubleClick}
          className={cn(
            "cursor-pointer select-none -m-2 p-2 min-h-[2.5rem] flex items-center",
            !isReadOnly && "hover:bg-muted/50 cursor-text",
            isReadOnly && "cursor-not-allowed",
            isPending && "opacity-50"
          )}
        >
          <Checkbox
            checked={!!optimisticValue}
            disabled
            className="pointer-events-none"
          />
        </div>
      );
    }

    if (dataType.includes("json") && optimisticValue) {
      return (
        <div 
          onDoubleClick={handleDoubleClick}
          className={cn(
            "text-sm font-mono cursor-pointer select-none -m-2 p-2 min-h-[2.5rem]",
            !isReadOnly && "hover:bg-muted/50 cursor-text",
            isReadOnly && "cursor-not-allowed",
            isPending && "opacity-50"
          )}
        >
          {typeof optimisticValue === "string" ? optimisticValue : JSON.stringify(optimisticValue)}
        </div>
      );
    }

    return (
      <div 
        onDoubleClick={handleDoubleClick}
        className={cn(
          "text-sm cursor-pointer select-none -m-2 p-2 min-h-[2.5rem] flex items-center justify-between",
          !isReadOnly && "hover:bg-muted/50 cursor-text",
          isReadOnly && "cursor-not-allowed",
          isPending && "opacity-50"
        )}
      >
        <span>{String(optimisticValue || "-")}</span>
        {showSuccess && <Check className="h-3 w-3 text-green-600 ml-2" />}
      </div>
    );
  }

  // Edit mode
  const handleChange = (newValue: any) => {
    setLocalValue(newValue);
  };

  return (
    <div ref={containerRef} className="w-full">
      {dataType === "boolean" || dataType.includes("bool") ? (
        <Checkbox
          checked={!!localValue}
          onCheckedChange={(checked) => {
            setLocalValue(checked);
            setIsEditing(false);
            saveValue(checked);
          }}
          className="ml-1"
        />
      ) : dataType.includes("date") || dataType.includes("timestamp") ? (
        <Popover open={isEditing} onOpenChange={(open) => {
          if (!open) handleSave();
        }}>
          <PopoverTrigger asChild>
            <Button
              variant="outline"
              size="sm"
              className={cn(
                "h-8 w-full justify-start text-left font-normal",
                !localValue && "text-muted-foreground"
              )}
            >
              <CalendarIcon className="mr-2 h-3 w-3" />
              {localValue ? format(new Date(String(localValue)), "PPP") : "Pick a date"}
            </Button>
          </PopoverTrigger>
          <PopoverContent className="w-auto p-0" align="start">
            <Calendar
              mode="single"
              selected={localValue ? new Date(String(localValue)) : undefined}
              onSelect={(date) => {
                const newValue = date?.toISOString();
                setLocalValue(newValue);
                setIsEditing(false);
                saveValue(newValue);
              }}
            />
          </PopoverContent>
        </Popover>
      ) : dataType.includes("text") ? (
        <Textarea
          ref={inputRef as React.RefObject<HTMLTextAreaElement>}
          value={String(localValue || "")}
          onChange={(e) => handleChange(e.target.value)}
          onKeyDown={handleKeyDown}
          placeholder={isNullable ? "Optional" : "Required"}
          rows={2}
          className="min-h-[32px] resize-none"
        />
      ) : dataType.includes("json") ? (
        <>
          <Textarea
            ref={inputRef as React.RefObject<HTMLTextAreaElement>}
            value={
              typeof localValue === "string"
                ? localValue
                : JSON.stringify(localValue, null, 2)
            }
            onChange={(e) => {
              const value = e.target.value;
              setLocalValue(value);
              if (value) {
                try {
                  JSON.parse(value);
                  setJsonError("");
                } catch (e) {
                  const error = e as Error;
                  setJsonError(`Invalid JSON: ${error.message}`);
                }
              } else {
                setJsonError("");
              }
            }}
            onKeyDown={handleKeyDown}
            placeholder={isNullable ? "Optional JSON" : "Required JSON"}
            rows={2}
            className={cn("min-h-[32px] resize-none font-mono text-xs", !isNullable && !localValue && "border-destructive", jsonError && "border-destructive")}
          />
          {jsonError && (
            <span className="text-xs text-destructive mt-1">{jsonError}</span>
          )}
        </>
      ) : (
        <Input
          ref={inputRef as React.RefObject<HTMLInputElement>}
          value={String(localValue || "")}
          onChange={(e) => handleChange(e.target.value)}
          onKeyDown={handleKeyDown}
          type={
            dataType.includes("int") ||
            dataType.includes("serial") ||
            dataType.includes("float") ||
            dataType.includes("numeric")
              ? "number"
              : "text"
          }
          placeholder={isNullable ? "Optional" : "Required"}
          className={cn("h-8", !isNullable && !localValue && "border-destructive")}
        />
      )}
    </div>
  );
}, (prevProps, nextProps) => {
  // Custom comparison function for memoization
  return (
    prevProps.value === nextProps.value &&
    prevProps.column === nextProps.column &&
    prevProps.rowId === nextProps.rowId &&
    prevProps.isTimestamp === nextProps.isTimestamp &&
    prevProps.onUpdate === nextProps.onUpdate
  );
});

OptimisticTableCell.displayName = 'OptimisticTableCell';