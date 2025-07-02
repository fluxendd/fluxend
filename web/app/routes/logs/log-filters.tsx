import { useState, useCallback, memo, useEffect, useMemo, useRef } from "react";
import { Button } from "~/components/ui/button";
import { Input } from "~/components/ui/input";
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from "~/components/ui/select";
import {
  Popover,
  PopoverContent,
  PopoverTrigger,
} from "~/components/ui/popover";
import { Calendar } from "~/components/ui/calendar";
import { Card, CardContent, CardFooter } from "~/components/ui/card";
import { Label } from "~/components/ui/label";
import {
  X,
  Search,
  CalendarIcon,
  Clock2Icon,
  ChevronDownIcon,
} from "lucide-react";
import { format } from "date-fns";
import type { DateRange } from "react-day-picker";
import { cn } from "~/lib/utils";
import type { LogsFilters } from "~/services/logs";

interface LogFiltersProps {
  onFiltersChange: (filters: LogsFilters) => void;
}

const HTTP_METHODS = ["GET", "POST", "PUT", "DELETE", "PATCH"];
const STATUS_CODES = [
  { value: "200", label: "200 OK" },
  { value: "201", label: "201 Created" },
  { value: "204", label: "204 No Content" },
  { value: "400", label: "400 Bad Request" },
  { value: "401", label: "401 Unauthorized" },
  { value: "403", label: "403 Forbidden" },
  { value: "404", label: "404 Not Found" },
  { value: "500", label: "500 Server Error" },
  { value: "502", label: "502 Bad Gateway" },
  { value: "503", label: "503 Service Unavailable" },
];

export const LogFilters = memo(({ onFiltersChange }: LogFiltersProps) => {
  // Initialize with today's date range
  const today = new Date();
  today.setHours(0, 0, 0, 0);
  const initialDateRange: DateRange = { from: today, to: today };
  
  const [filters, setFilters] = useState<LogsFilters>(() => {
    // Calculate initial filters with today's date
    const todayFormatted = format(today, "MM-dd-yyyy");
    const startOfDay = Math.floor(today.getTime() / 1000);
    const endOfDay = Math.floor(new Date(today.getTime() + 24 * 60 * 60 * 1000 - 1).getTime() / 1000);
    
    return {
      dateStart: todayFormatted,
      dateEnd: todayFormatted,
      startTime: startOfDay,
      endTime: endOfDay
    };
  });
  const [endpointSearch, setEndpointSearch] = useState("");
  const [dateRange, setDateRange] = useState<DateRange | undefined>(initialDateRange);
  const [startTime, setStartTime] = useState("00:00:00");
  const [endTime, setEndTime] = useState("23:59:59");
  const [popoverOpen, setPopoverOpen] = useState(false);
  const [pendingDateRange, setPendingDateRange] = useState<DateRange | undefined>(initialDateRange);
  const [pendingStartTime, setPendingStartTime] = useState("00:00:00");
  const [pendingEndTime, setPendingEndTime] = useState("23:59:59");
  const searchTimeoutRef = useRef<NodeJS.Timeout | null>(null);

  const handleFilterChange = useCallback(
    (key: keyof LogsFilters, value: string | undefined) => {
      setFilters((prevFilters) => {
        const newFilters = { ...prevFilters };

        if (value) {
          newFilters[key] = value;
        } else {
          delete newFilters[key];
        }

        onFiltersChange(newFilters);
        return newFilters;
      });
    },
    [onFiltersChange]
  );

  const handleEndpointSearch = useCallback(() => {
    handleFilterChange("endpoint", endpointSearch || undefined);
  }, [handleFilterChange, endpointSearch]);

  const handleEndpointInputChange = useCallback((value: string) => {
    setEndpointSearch(value);
    
    // Clear existing timeout
    if (searchTimeoutRef.current) {
      clearTimeout(searchTimeoutRef.current);
    }
    
    // Debounce search
    if (value) {
      searchTimeoutRef.current = setTimeout(() => {
        handleFilterChange("endpoint", value);
      }, 500);
    } else {
      handleFilterChange("endpoint", undefined);
    }
  }, [handleFilterChange]);

  // Cleanup timeout on unmount
  useEffect(() => {
    return () => {
      if (searchTimeoutRef.current) {
        clearTimeout(searchTimeoutRef.current);
      }
    };
  }, []);

  const handleDateRangeChange = useCallback(
    (range: DateRange | undefined) => {
      setPendingDateRange(range);
    },
    []
  );

  const handlePopoverOpenChange = useCallback(
    (open: boolean) => {
      setPopoverOpen(open);
      
      if (open) {
        // When opening, sync pending values with current values
        setPendingDateRange(dateRange);
        setPendingStartTime(startTime);
        setPendingEndTime(endTime);
      } else {
        // When closing the popover, apply the pending changes
        if (pendingDateRange !== dateRange || pendingStartTime !== startTime || pendingEndTime !== endTime) {
          setDateRange(pendingDateRange);
          setStartTime(pendingStartTime);
          setEndTime(pendingEndTime);
          
          const newFilters = { ...filters };

          if (pendingDateRange?.from) {
            // Format date as MM-DD-YYYY for dateStart
            newFilters.dateStart = format(pendingDateRange.from, "MM-dd-yyyy");

            // Calculate Unix timestamp for startTime
            const [hours, minutes, seconds] = pendingStartTime.split(":").map(Number);
            const startDateTime = new Date(pendingDateRange.from);
            startDateTime.setHours(hours, minutes, seconds, 0);
            newFilters.startTime = Math.floor(startDateTime.getTime() / 1000);
          } else {
            delete newFilters.dateStart;
            delete newFilters.startTime;
          }

          if (pendingDateRange?.to) {
            // Format date as MM-DD-YYYY for dateEnd
            newFilters.dateEnd = format(pendingDateRange.to, "MM-dd-yyyy");

            // Calculate Unix timestamp for endTime
            const [hours, minutes, seconds] = pendingEndTime.split(":").map(Number);
            const endDateTime = new Date(pendingDateRange.to);
            endDateTime.setHours(hours, minutes, seconds, 0);
            newFilters.endTime = Math.floor(endDateTime.getTime() / 1000);
          } else {
            delete newFilters.dateEnd;
            delete newFilters.endTime;
          }

          setFilters(newFilters);
          onFiltersChange(newFilters);
        }
      }
    },
    [pendingDateRange, pendingStartTime, pendingEndTime, dateRange, startTime, endTime, filters, onFiltersChange]
  );

  const handleTimeChange = useCallback(
    (type: "start" | "end", time: string) => {
      if (type === "start") {
        setPendingStartTime(time);
      } else {
        setPendingEndTime(time);
      }
    },
    []
  );

  const clearFilters = useCallback(() => {
    // Reset to today's date
    const today = new Date();
    today.setHours(0, 0, 0, 0);
    const todayRange: DateRange = { from: today, to: today };
    const todayFormatted = format(today, "MM-dd-yyyy");
    const startOfDay = Math.floor(today.getTime() / 1000);
    const endOfDay = Math.floor(new Date(today.getTime() + 24 * 60 * 60 * 1000 - 1).getTime() / 1000);
    
    const defaultFilters = {
      dateStart: todayFormatted,
      dateEnd: todayFormatted,
      startTime: startOfDay,
      endTime: endOfDay
    };
    
    setFilters(defaultFilters);
    setEndpointSearch("");
    setDateRange(todayRange);
    setStartTime("00:00:00");
    setEndTime("23:59:59");
    setPendingDateRange(todayRange);
    setPendingStartTime("00:00:00");
    setPendingEndTime("23:59:59");
    onFiltersChange(defaultFilters);
  }, [onFiltersChange]);

  // Trigger initial filter change on mount
  useEffect(() => {
    onFiltersChange(filters);
  }, []); // eslint-disable-line react-hooks/exhaustive-deps

  // Check if we have non-default filters
  const hasActiveFilters = useMemo(() => {
    const filterKeys = Object.keys(filters);
    // Only date filters with default values
    if (filterKeys.length === 4 && 
        filters.dateStart && 
        filters.dateEnd && 
        filters.startTime !== undefined && 
        filters.endTime !== undefined &&
        !filters.method &&
        !filters.status &&
        !filters.ipAddress &&
        !filters.endpoint &&
        !filters.userUuid) {
      return false;
    }
    // Has other filters beyond just dates
    return filterKeys.some(key => 
      key !== 'dateStart' && 
      key !== 'dateEnd' && 
      key !== 'startTime' && 
      key !== 'endTime'
    );
  }, [filters]);

  return (
    <div className="flex flex-wrap items-center gap-2 px-4 py-2 border-b isolate">
      <Select
        value={filters.method || "all"}
        onValueChange={(value) =>
          handleFilterChange("method", value === "all" ? undefined : value)
        }
      >
        <SelectTrigger className="w-[140px]">
          <SelectValue placeholder="Method" />
        </SelectTrigger>
        <SelectContent>
          <SelectItem value="all">All Methods</SelectItem>
          {HTTP_METHODS.map((method) => (
            <SelectItem key={method} value={method}>
              {method}
            </SelectItem>
          ))}
        </SelectContent>
      </Select>

      <Select
        value={filters.status || "all"}
        onValueChange={(value) =>
          handleFilterChange("status", value === "all" ? undefined : value)
        }
      >
        <SelectTrigger className="w-[180px]">
          <SelectValue placeholder="Status Code" />
        </SelectTrigger>
        <SelectContent>
          <SelectItem value="all">All Status Codes</SelectItem>
          {STATUS_CODES.map((status) => (
            <SelectItem key={status.value} value={status.value}>
              {status.label}
            </SelectItem>
          ))}
        </SelectContent>
      </Select>

      <Input
        placeholder="IP Address"
        value={filters.ipAddress || ""}
        onChange={(e) =>
          handleFilterChange("ipAddress", e.target.value || undefined)
        }
        className="w-[150px]"
      />

      <div className="flex items-center gap-1">
        <Input
          placeholder="Search endpoint..."
          value={endpointSearch}
          onChange={(e) => handleEndpointInputChange(e.target.value)}
          onKeyDown={(e) => {
            if (e.key === "Enter") {
              e.preventDefault();
              if (searchTimeoutRef.current) {
                clearTimeout(searchTimeoutRef.current);
              }
              handleEndpointSearch();
            }
          }}
          className="w-[250px]"
        />
      </div>

      <Popover open={popoverOpen} onOpenChange={handlePopoverOpenChange}>
        <PopoverTrigger asChild>
          <Button
            variant="outline"
            className={cn(
              "w-[280px] justify-between font-normal",
              !dateRange && "text-muted-foreground"
            )}
          >
            <div className="flex items-center">
              <CalendarIcon className="mr-2 h-4 w-4" />
              {dateRange?.from && dateRange?.to
                ? `${format(dateRange.from, "MM-dd-yyyy")} - ${format(
                    dateRange.to,
                    "MM-dd-yyyy"
                  )}`
                : dateRange?.from
                ? format(dateRange.from, "MM-dd-yyyy")
                : "Select date range"}
            </div>
            <ChevronDownIcon className="ml-auto h-4 w-4 opacity-50" />
          </Button>
        </PopoverTrigger>
        <PopoverContent className="w-fit p-0" align="start">
          <Card className="w-fit border-0 p-0 gap-0">
            <CardContent className="p-3">
              <Calendar
                mode="range"
                selected={pendingDateRange}
                onSelect={handleDateRangeChange}
                captionLayout="dropdown"
                defaultMonth={dateRange?.from || new Date()}
                className="bg-transparent p-0"
              />
            </CardContent>
            <CardFooter className="flex flex-col gap-3 border-t px-3 py-3">
              <div className="grid grid-cols-2 gap-4 w-full">
                <div className="flex flex-col gap-2">
                  <Label htmlFor="start-time">Start Time</Label>
                  <div className="relative flex items-center">
                    <Clock2Icon className="pointer-events-none absolute left-2.5 size-4 select-none text-muted-foreground" />
                    <Input
                      id="start-time"
                      type="time"
                      step="1"
                      value={pendingStartTime}
                      onChange={(e) =>
                        handleTimeChange("start", e.target.value)
                      }
                      className="appearance-none pl-8 [&::-webkit-calendar-picker-indicator]:hidden [&::-webkit-calendar-picker-indicator]:appearance-none"
                    />
                  </div>
                </div>
                <div className="flex flex-col gap-2">
                  <Label htmlFor="end-time">End Time</Label>
                  <div className="relative flex items-center">
                    <Clock2Icon className="pointer-events-none absolute left-2.5 size-4 select-none text-muted-foreground" />
                    <Input
                      id="end-time"
                      type="time"
                      step="1"
                      value={pendingEndTime}
                      onChange={(e) => handleTimeChange("end", e.target.value)}
                      className="appearance-none pl-8 [&::-webkit-calendar-picker-indicator]:hidden [&::-webkit-calendar-picker-indicator]:appearance-none"
                    />
                  </div>
                </div>
              </div>
            </CardFooter>
          </Card>
        </PopoverContent>
      </Popover>

      {hasActiveFilters && (
        <Button
          size="sm"
          variant="ghost"
          onClick={clearFilters}
          className="ml-auto"
        >
          <X className="h-4 w-4 mr-1" />
          Clear Filters
        </Button>
      )}
    </div>
  );
});

LogFilters.displayName = "LogFilters";