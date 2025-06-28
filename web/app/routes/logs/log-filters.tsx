import { useState } from "react";
import { Button } from "~/components/ui/button";
import { Input } from "~/components/ui/input";
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from "~/components/ui/select";
import { X, Search } from "lucide-react";
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

export function LogFilters({ onFiltersChange }: LogFiltersProps) {
  const [filters, setFilters] = useState<LogsFilters>({});
  const [endpointSearch, setEndpointSearch] = useState("");

  const handleFilterChange = (key: keyof LogsFilters, value: string | undefined) => {
    const newFilters = { ...filters };
    
    if (value) {
      (newFilters as any)[key] = value;
    } else {
      delete newFilters[key];
    }
    
    setFilters(newFilters);
    onFiltersChange(newFilters);
  };

  const handleEndpointSearch = () => {
    handleFilterChange("endpoint", endpointSearch || undefined);
  };

  const clearFilters = () => {
    setFilters({});
    setEndpointSearch("");
    onFiltersChange({});
  };

  const hasActiveFilters = Object.keys(filters).length > 0;

  return (
    <div className="flex flex-wrap items-center gap-2 px-4 py-2 border-b">
      <Select
        value={filters.method || "all"}
        onValueChange={(value) => handleFilterChange("method", value === "all" ? undefined : value)}
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
        onValueChange={(value) => handleFilterChange("status", value === "all" ? undefined : value)}
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
        onChange={(e) => handleFilterChange("ipAddress", e.target.value || undefined)}
        className="w-[150px]"
      />

      <div className="flex items-center gap-1">
        <Input
          placeholder="Search endpoint..."
          value={endpointSearch}
          onChange={(e) => setEndpointSearch(e.target.value)}
          onKeyDown={(e) => {
            if (e.key === "Enter") {
              handleEndpointSearch();
            }
          }}
          className="w-[250px]"
        />
        <Button
          size="sm"
          variant="secondary"
          onClick={handleEndpointSearch}
        >
          <Search className="h-4 w-4" />
        </Button>
      </div>

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
}