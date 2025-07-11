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
  Tooltip,
  TooltipContent,
  TooltipTrigger,
} from "~/components/ui/tooltip";
import {
  X,
  Search,
  CalendarIcon,
  Clock2Icon,
  ChevronDownIcon,
  Info,
} from "lucide-react";
import { format, subDays } from "date-fns";
import { fromZonedTime } from "date-fns-tz";
import type { DateRange } from "react-day-picker";
import { cn } from "~/lib/utils";
import type { LogsFilters } from "~/services/logs";

interface LogFiltersProps {
  onFiltersChange: (filters: LogsFilters) => void;
  initialFilters?: LogsFilters;
}

const HTTP_METHODS = ["GET", "POST", "PUT", "DELETE", "PATCH"] as const;
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

export const LogFilters = memo(({ onFiltersChange, initialFilters }: LogFiltersProps) => {
  // Get user's timezone
  const userTimezone = Intl.DateTimeFormat().resolvedOptions().timeZone;
  
  // Initialize with today's date range in user's timezone
  const today = new Date();
  today.setHours(0, 0, 0, 0);
  
  const [filters, setFilters] = useState<LogsFilters>(() => {
    if (initialFilters && Object.keys(initialFilters).length > 0) {
      return initialFilters;
    }
    // Calculate initial filters with today's date
    // Calculate UTC timestamps for start and end of day in user's timezone
    const startOfDayLocal = new Date(today);
    startOfDayLocal.setHours(0, 0, 0, 0);
    const endOfDayLocal = new Date(today);
    endOfDayLocal.setHours(23, 59, 59, 999);
    
    const startOfDayUTC = fromZonedTime(startOfDayLocal, userTimezone);
    const endOfDayUTC = fromZonedTime(endOfDayLocal, userTimezone);
    
    return {
      startTime: Math.floor(startOfDayUTC.getTime() / 1000),
      endTime: Math.floor(endOfDayUTC.getTime() / 1000)
    };
  });

  // Initialize states from filters
  const getInitialDateRange = (): DateRange | undefined => {
    if (filters.startTime && filters.endTime) {
      const fromDate = new Date(filters.startTime * 1000);
      const toDate = new Date(filters.endTime * 1000);
      // Set to start of day for date comparison
      fromDate.setHours(0, 0, 0, 0);
      toDate.setHours(0, 0, 0, 0);
      return { from: fromDate, to: toDate };
    }
    return { from: today, to: today };
  };
  
  const getTimeFromTimestamp = (timestamp: number | undefined, defaultTime: string): string => {
    if (!timestamp) return defaultTime;
    // Convert Unix timestamp to local date
    // The timestamp is already in UTC, so creating a Date object will automatically convert to local timezone
    const date = new Date(timestamp * 1000);
    return `${date.getHours().toString().padStart(2, '0')}:${date.getMinutes().toString().padStart(2, '0')}:${date.getSeconds().toString().padStart(2, '0')}`;
  };
  
  const [endpointSearch, setEndpointSearch] = useState(filters.endpoint || "");
  const [dateRange, setDateRange] = useState<DateRange | undefined>(getInitialDateRange());
  const [startTime, setStartTime] = useState(getTimeFromTimestamp(filters.startTime, "00:00:00"));
  const [endTime, setEndTime] = useState(getTimeFromTimestamp(filters.endTime, "23:59:59"));
  const [popoverOpen, setPopoverOpen] = useState(false);
  const [pendingDateRange, setPendingDateRange] = useState<DateRange | undefined>(getInitialDateRange());
  const [pendingStartTime, setPendingStartTime] = useState(getTimeFromTimestamp(filters.startTime, "00:00:00"));
  const [pendingEndTime, setPendingEndTime] = useState(getTimeFromTimestamp(filters.endTime, "23:59:59"));
  const searchTimeoutRef = useRef<NodeJS.Timeout | null>(null);

  const handleFilterChange = useCallback(
    (key: keyof LogsFilters, value: string | number | undefined) => {
      setFilters((prevFilters) => {
        const newFilters = { ...prevFilters };

        if (value !== undefined) {
          (newFilters as any)[key] = value;
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
        handleFilterChange("endpoint", value || undefined);
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
            // Calculate Unix timestamp for startTime
            // Create date in user's timezone then convert to UTC
            const [hours, minutes, seconds] = pendingStartTime.split(":").map(Number);
            const startDateTime = new Date(pendingDateRange.from);
            startDateTime.setHours(hours, minutes, seconds, 0);
            
            // Convert to UTC timestamp
            const utcStartTime = fromZonedTime(startDateTime, userTimezone);
            newFilters.startTime = Math.floor(utcStartTime.getTime() / 1000);
          } else {
            delete newFilters.startTime;
          }

          if (pendingDateRange?.to) {
            // Calculate Unix timestamp for endTime
            // Create date in user's timezone then convert to UTC
            const [hours, minutes, seconds] = pendingEndTime.split(":").map(Number);
            const endDateTime = new Date(pendingDateRange.to);
            endDateTime.setHours(hours, minutes, seconds, 0);
            
            // Convert to UTC timestamp
            const utcEndTime = fromZonedTime(endDateTime, userTimezone);
            newFilters.endTime = Math.floor(utcEndTime.getTime() / 1000);
          } else {
            delete newFilters.endTime;
          }

          setFilters(newFilters);
          onFiltersChange(newFilters);
        }
      }
    },
    [pendingDateRange, pendingStartTime, pendingEndTime, dateRange, startTime, endTime, filters, onFiltersChange, userTimezone]
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
    
    // Calculate UTC timestamps for start and end of day in user's timezone
    const startOfDayLocal = new Date(today);
    startOfDayLocal.setHours(0, 0, 0, 0);
    const endOfDayLocal = new Date(today);
    endOfDayLocal.setHours(23, 59, 59, 999);
    
    const startOfDayUTC = fromZonedTime(startOfDayLocal, userTimezone);
    const endOfDayUTC = fromZonedTime(endOfDayLocal, userTimezone);
    
    const defaultFilters = {
      startTime: Math.floor(startOfDayUTC.getTime() / 1000),
      endTime: Math.floor(endOfDayUTC.getTime() / 1000)
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
  }, [onFiltersChange, userTimezone]);

  // Trigger initial filter change on mount
  useEffect(() => {
    onFiltersChange(filters);
  }, []); // eslint-disable-line react-hooks/exhaustive-deps

  // Check if we have non-default filters
  const hasActiveFilters = useMemo(() => {
    // Calculate today's default time range for comparison
    const today = new Date();
    today.setHours(0, 0, 0, 0);
    
    const startOfDayLocal = new Date(today);
    startOfDayLocal.setHours(0, 0, 0, 0);
    const endOfDayLocal = new Date(today);
    endOfDayLocal.setHours(23, 59, 59, 999);
    
    const startOfDayUTC = fromZonedTime(startOfDayLocal, userTimezone);
    const endOfDayUTC = fromZonedTime(endOfDayLocal, userTimezone);
    
    const defaultStartTime = Math.floor(startOfDayUTC.getTime() / 1000);
    const defaultEndTime = Math.floor(endOfDayUTC.getTime() / 1000);
    
    // Check if current filters differ from defaults
    const hasNonDefaultTimeFilter = filters.startTime !== defaultStartTime || filters.endTime !== defaultEndTime;
    
    // Has other filters beyond just times
    const hasOtherFilters = !!(filters.method || filters.status || filters.ipAddress || filters.endpoint || filters.userUuid);
    
    return hasNonDefaultTimeFilter || hasOtherFilters;
  }, [filters, userTimezone]);

  // Date preset handlers
  const handleDatePreset = useCallback((preset: 'today' | 'yesterday' | 'last3days') => {
    const now = new Date();
    let fromDate: Date;
    let toDate: Date = new Date();
    
    // Set toDate to end of today
    toDate.setHours(23, 59, 59, 999);
    
    switch (preset) {
      case 'today':
        fromDate = new Date();
        fromDate.setHours(0, 0, 0, 0);
        break;
      case 'yesterday':
        fromDate = subDays(now, 1);
        fromDate.setHours(0, 0, 0, 0);
        toDate = subDays(now, 1);
        toDate.setHours(23, 59, 59, 999);
        break;
      case 'last3days':
        fromDate = subDays(now, 2); // 3 days including today
        fromDate.setHours(0, 0, 0, 0);
        break;
    }
    
    // Update date range and times
    const newDateRange: DateRange = { from: fromDate, to: toDate };
    setDateRange(newDateRange);
    setPendingDateRange(newDateRange);
    setStartTime("00:00:00");
    setPendingStartTime("00:00:00");
    setEndTime("23:59:59");
    setPendingEndTime("23:59:59");
    
    // Calculate timestamps
    const startDateTime = new Date(fromDate);
    startDateTime.setHours(0, 0, 0, 0);
    const endDateTime = new Date(toDate);
    endDateTime.setHours(23, 59, 59, 999);
    
    const utcStartTime = fromZonedTime(startDateTime, userTimezone);
    const utcEndTime = fromZonedTime(endDateTime, userTimezone);
    
    const newFilters = {
      ...filters,
      startTime: Math.floor(utcStartTime.getTime() / 1000),
      endTime: Math.floor(utcEndTime.getTime() / 1000)
    };
    
    setFilters(newFilters);
    onFiltersChange(newFilters);
  }, [filters, onFiltersChange, userTimezone]);

  // Check which preset is active
  const activePreset = useMemo(() => {
    if (!dateRange?.from || !dateRange?.to) return null;
    
    const now = new Date();
    const today = new Date();
    today.setHours(0, 0, 0, 0);
    const yesterday = subDays(today, 1);
    const threeDaysAgo = subDays(today, 2);
    
    // Check if dates match today
    if (
      dateRange.from.toDateString() === today.toDateString() &&
      dateRange.to.toDateString() === today.toDateString()
    ) {
      return 'today';
    }
    
    // Check if dates match yesterday
    if (
      dateRange.from.toDateString() === yesterday.toDateString() &&
      dateRange.to.toDateString() === yesterday.toDateString()
    ) {
      return 'yesterday';
    }
    
    // Check if dates match last 3 days
    if (
      dateRange.from.toDateString() === threeDaysAgo.toDateString() &&
      dateRange.to.toDateString() === today.toDateString()
    ) {
      return 'last3days';
    }
    
    return null;
  }, [dateRange]);

  return (
    <div className="flex flex-wrap items-center gap-2 px-4 py-2 border-b isolate">
      {/* Date presets */}
      <div className="flex items-center gap-1">
        <Button
          size="sm"
          variant={activePreset === 'today' ? 'default' : 'outline'}
          onClick={() => handleDatePreset('today')}
          className="h-9"
        >
          Today
        </Button>
        <Button
          size="sm"
          variant={activePreset === 'yesterday' ? 'default' : 'outline'}
          onClick={() => handleDatePreset('yesterday')}
          className="h-9"
        >
          Yesterday
        </Button>
        <Button
          size="sm"
          variant={activePreset === 'last3days' ? 'default' : 'outline'}
          onClick={() => handleDatePreset('last3days')}
          className="h-9"
        >
          Last 3 Days
        </Button>
      </div>
      
      <div className="w-px h-6 bg-border" />
      
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
              <div className="flex items-center gap-1 text-xs text-muted-foreground">
                <Tooltip>
                  <TooltipTrigger asChild>
                    <div className="flex items-center gap-1 cursor-help">
                      <Info className="h-3 w-3" />
                      <span>Times are in your local timezone ({userTimezone})</span>
                    </div>
                  </TooltipTrigger>
                  <TooltipContent className="max-w-xs">
                    <p>All times are displayed in your local timezone but sent to the server as UTC timestamps for consistent filtering across all users.</p>
                  </TooltipContent>
                </Tooltip>
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