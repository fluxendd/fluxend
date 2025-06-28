import { useState, useCallback, useEffect, useMemo } from "react";
import { zodResolver } from "@hookform/resolvers/zod";
import { useForm } from "react-hook-form";
import * as z from "zod";
import {
  X,
  Search,
  FilterX,
  ChevronDown,
  ChevronUp,
  AlertCircle,
} from "lucide-react";
import {
  Form,
  FormControl,
  FormField,
  FormItem,
  FormLabel,
  FormMessage,
} from "~/components/ui/form";
import { Input } from "~/components/ui/input";
import { Button } from "~/components/ui/button";
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
import {
  Accordion,
  AccordionContent,
  AccordionItem,
  AccordionTrigger,
} from "~/components/ui/accordion";
import { cn } from "~/lib/utils";
import { Badge } from "~/components/ui/badge";
import { Checkbox } from "~/components/ui/checkbox";
import { format } from "date-fns";
import { Calendar } from "~/components/ui/calendar";
import { useQueryClient } from "@tanstack/react-query";
// import { Alert, AlertDescription, AlertTitle } from "~/components/ui/alert";

// Define column types
enum ColumnDataType {
  // Basic types
  Integer = "integer",
  BigInt = "bigint",
  SmallInt = "smallint",
  Text = "text",
  Varchar = "character varying",
  Char = "character",
  Boolean = "boolean",

  // Numeric
  Numeric = "numeric",
  Decimal = "decimal",
  Real = "real",
  DoublePrecision = "double precision",

  // Date/Time
  Date = "date",
  Time = "time",
  Timestamp = "timestamp",
  TimestampTZ = "timestamptz",

  // Other common types
  UUID = "uuid",
  JSON = "json",
  JSONB = "jsonb",
  Array = "array",
}

// Define supported operators
export enum FilterOperator {
  // Equality
  Equals = "eq",
  NotEquals = "neq",

  // Comparison
  GreaterThan = "gt",
  GreaterThanOrEqual = "gte",
  LessThan = "lt",
  LessThanOrEqual = "lte",

  // Text
  Like = "like",
  ILike = "ilike",
  Match = "match",
  IMatch = "imatch",

  // Full Text Search
  FTS = "fts",
  PlainFTS = "plfts",
  PhraseFTS = "phfts",
  WebFTS = "wfts",

  // Lists
  In = "in",

  // Boolean
  Is = "is",

  // Arrays/JSON
  Contains = "cs",
  ContainedIn = "cd",
  Overlap = "ov",

  // Range Operators
  StrictlyLeft = "sl",
  StrictlyRight = "sr",
  NotExtendRight = "nxr",
  NotExtendLeft = "nxl",
  Adjacent = "adj",
}

// Group operators by data type for UI presentation
const OPERATOR_BY_TYPE = {
  text: [
    { value: FilterOperator.Equals, label: "Equals" },
    { value: FilterOperator.NotEquals, label: "Not Equals" },
    { value: FilterOperator.Like, label: "Like" },
    { value: FilterOperator.ILike, label: "Like (case insensitive)" },
    { value: FilterOperator.Match, label: "Match Regex" },
    { value: FilterOperator.IMatch, label: "Match Regex (case insensitive)" },
    { value: FilterOperator.FTS, label: "Full-Text Search" },
    { value: FilterOperator.PlainFTS, label: "Plain Text Search" },
    { value: FilterOperator.PhraseFTS, label: "Phrase Search" },
    { value: FilterOperator.WebFTS, label: "Web Search" },
    { value: FilterOperator.In, label: "In List" },
  ],
  number: [
    { value: FilterOperator.Equals, label: "Equals" },
    { value: FilterOperator.NotEquals, label: "Not Equals" },
    { value: FilterOperator.GreaterThan, label: "Greater Than" },
    {
      value: FilterOperator.GreaterThanOrEqual,
      label: "Greater Than or Equal",
    },
    { value: FilterOperator.LessThan, label: "Less Than" },
    { value: FilterOperator.LessThanOrEqual, label: "Less Than or Equal" },
    { value: FilterOperator.In, label: "In List" },
  ],
  boolean: [{ value: FilterOperator.Is, label: "Is" }],
  date: [
    { value: FilterOperator.Equals, label: "Equals" },
    { value: FilterOperator.NotEquals, label: "Not Equals" },
    { value: FilterOperator.GreaterThan, label: "After" },
    { value: FilterOperator.GreaterThanOrEqual, label: "On or After" },
    { value: FilterOperator.LessThan, label: "Before" },
    { value: FilterOperator.LessThanOrEqual, label: "On or Before" },
  ],
  array: [
    { value: FilterOperator.Contains, label: "Contains" },
    { value: FilterOperator.ContainedIn, label: "Contained In" },
    { value: FilterOperator.Overlap, label: "Overlaps" },
  ],
  json: [{ value: FilterOperator.Contains, label: "Contains" }],
  range: [
    { value: FilterOperator.StrictlyLeft, label: "Strictly Left Of" },
    { value: FilterOperator.StrictlyRight, label: "Strictly Right Of" },
    { value: FilterOperator.NotExtendRight, label: "Does Not Extend Right Of" },
    { value: FilterOperator.NotExtendLeft, label: "Does Not Extend Left Of" },
    { value: FilterOperator.Adjacent, label: "Adjacent To" },
    { value: FilterOperator.Overlap, label: "Overlaps" },
  ],
};

// Map PostgreSQL column types to our simplified categories
const mapColumnTypeToCategory = (
  columnType: string
): keyof typeof OPERATOR_BY_TYPE => {
  const type = columnType.toLowerCase();

  if (
    type.includes(ColumnDataType.Integer) ||
    type.includes(ColumnDataType.BigInt) ||
    type.includes(ColumnDataType.SmallInt) ||
    type.includes(ColumnDataType.Numeric) ||
    type.includes(ColumnDataType.Decimal) ||
    type.includes(ColumnDataType.Real) ||
    type.includes(ColumnDataType.DoublePrecision)
  ) {
    return "number";
  } else if (
    type.includes(ColumnDataType.Text) ||
    type.includes(ColumnDataType.Varchar) ||
    type.includes(ColumnDataType.Char) ||
    type.includes(ColumnDataType.UUID)
  ) {
    return "text";
  } else if (
    type.includes(ColumnDataType.Date) ||
    type.includes(ColumnDataType.Timestamp) ||
    type.includes(ColumnDataType.TimestampTZ)
  ) {
    return "date";
  } else if (type.includes(ColumnDataType.Boolean)) {
    return "boolean";
  } else if (
    type.includes(ColumnDataType.JSON) ||
    type.includes(ColumnDataType.JSONB)
  ) {
    return "json";
  } else if (type.includes(ColumnDataType.Array) || type.endsWith("[]")) {
    return "array";
  } else if (
    type.includes("range") ||
    type.includes("tsrange") ||
    type.includes("daterange")
  ) {
    return "range";
  }

  // Default to text for unknown types
  return "text";
};

// Base filter criteria schema
const baseFilterSchema = z.object({
  column: z.string(),
  operator: z.nativeEnum(FilterOperator),
});

// Type-specific value schemas
const textFilterSchema = baseFilterSchema.extend({
  value: z.string().or(z.array(z.string())),
});

const numberFilterSchema = baseFilterSchema.extend({
  value: z.number().or(z.array(z.number())).or(z.string()),
});

const booleanFilterSchema = baseFilterSchema.extend({
  value: z.boolean().or(z.enum(["true", "false", "null"])),
});

const dateFilterSchema = baseFilterSchema.extend({
  value: z.date().or(z.string()),
});

const arrayFilterSchema = baseFilterSchema.extend({
  value: z.array(z.any()).or(z.string()),
});

const jsonFilterSchema = baseFilterSchema.extend({
  value: z.string(), // JSON as string
});

// Union of all filter criteria
export const filterCriteriaSchema = z.discriminatedUnion("type", [
  textFilterSchema.extend({ type: z.literal("text") }),
  numberFilterSchema.extend({ type: z.literal("number") }),
  booleanFilterSchema.extend({ type: z.literal("boolean") }),
  dateFilterSchema.extend({ type: z.literal("date") }),
  arrayFilterSchema.extend({ type: z.literal("array") }),
  jsonFilterSchema.extend({ type: z.literal("json") }),
]);

// Filter list schema
export const filtersSchema = z.object({
  filters: z.array(filterCriteriaSchema),
  logicalOperator: z.enum(["and", "or"]).default("and"),
});

// Type for our filter criteria
export type FilterCriteria = z.infer<typeof filterCriteriaSchema>;
export type Filters = z.infer<typeof filtersSchema>;

// Helper to format filter values for display in the UI
const formatFilterValueForDisplay = (filter: FilterCriteria): string => {
  if (filter.type === "boolean") {
    if (filter.value === true || filter.value === "true") return "True";
    if (filter.value === false || filter.value === "false") return "False";
    return "Null";
  }

  if (filter.type === "date" && filter.value instanceof Date) {
    return format(filter.value, "PPP");
  }

  if (Array.isArray(filter.value)) {
    return filter.value.join(", ");
  }

  return String(filter.value);
};

// Convert filter criteria to PostgREST query parameters
export const buildPostgrestFilterParams = (
  filters: FilterCriteria[],
  logicalOperator: "and" | "or" = "and"
): Record<string, string> => {
  if (!filters || filters.length === 0) {
    return {};
  }

  // Single filter handling
  if (filters.length === 1) {
    const filter = filters[0];
    const { column, operator, value } = filter;

    // Handle different value formats based on operator and type
    if (operator === FilterOperator.In) {
      const values = Array.isArray(value) ? value : value.toString().split(",");
      const formatted = `(${values
        .map((v) => {
          // Quote string values inside the in() operator
          if (filter.type === "text" && typeof v === "string") {
            // Escape quotes inside the string
            const escaped = v.trim().replace(/"/g, '\\"');
            return `"${escaped}"`;
          }
          return v;
        })
        .join(",")})`;
      return { [column]: `${operator}.${formatted}` };
    }

    if (filter.type === "boolean" && operator === FilterOperator.Is) {
      return { [column]: `${operator}.${value}` };
    }

    if (filter.type === "date") {
      const dateValue = value instanceof Date ? value : new Date(value);
      return { [column]: `${operator}.${format(dateValue, "yyyy-MM-dd")}` };
    }

    if (operator === FilterOperator.Like || operator === FilterOperator.ILike) {
      // Replace * with % for PostgREST pattern matching
      const formattedValue = String(value).replace(/\*/g, "%");
      return { [column]: `${operator}.${formattedValue}` };
    }

    // Handle Full-Text Search operators
    if (
      [
        FilterOperator.FTS,
        FilterOperator.PlainFTS,
        FilterOperator.PhraseFTS,
        FilterOperator.WebFTS,
      ].includes(operator as FilterOperator)
    ) {
      // For full-text search, we can optionally specify the language like "fts(english).query"
      const language = "english"; // Default language, could be made configurable
      return { [column]: `${operator}(${language}).${value}` };
    }

    if (
      ["array", "json"].includes(filter.type) &&
      [
        FilterOperator.Contains,
        FilterOperator.ContainedIn,
        FilterOperator.Overlap,
      ].includes(operator as FilterOperator)
    ) {
      // Format for array operators: column=cs.{val1,val2}
      let arrayValue = value;
      if (typeof value === "string") {
        try {
          // Try to parse as JSON if it's a string
          arrayValue = JSON.parse(value);
        } catch (e) {
          // If parsing fails, split by comma
          arrayValue = value.split(",").map((v) => v.trim());
        }
      }

      if (!Array.isArray(arrayValue)) {
        arrayValue = [arrayValue];
      }

      return { [column]: `${operator}.{${arrayValue.join(",")}}` };
    }

    // Default case
    return { [column]: `${operator}.${value}` };
  }

  // Multiple filters handling with AND/OR
  const filterStr = filters
    .map((filter) => {
      const { column, operator, value } = filter;

      // Format each individual filter condition
      if (filter.type === "boolean" && operator === FilterOperator.Is) {
        return `${column}.${operator}.${value}`;
      }

      if (filter.type === "date") {
        const dateValue = value instanceof Date ? value : new Date(value);
        return `${column}.${operator}.${format(dateValue, "yyyy-MM-dd")}`;
      }

      if (operator === FilterOperator.In) {
        const values = Array.isArray(value)
          ? value
          : value.toString().split(",");
        const formatted = `(${values
          .map((v) => {
            if (filter.type === "text" && typeof v === "string") {
              const escaped = v.trim().replace(/"/g, '\\"');
              return `"${escaped}"`;
            }
            return v;
          })
          .join(",")})`;
        return `${column}.${operator}.${formatted}`;
      }

      if (
        operator === FilterOperator.Like ||
        operator === FilterOperator.ILike
      ) {
        const formattedValue = String(value).replace(/\*/g, "%");
        return `${column}.${operator}.${formattedValue}`;
      }

      // Handle Full-Text Search operators in logical combinations
      if (
        [
          FilterOperator.FTS,
          FilterOperator.PlainFTS,
          FilterOperator.PhraseFTS,
          FilterOperator.WebFTS,
        ].includes(operator as FilterOperator)
      ) {
        const language = "english"; // Default language
        return `${column}.${operator}(${language}).${value}`;
      }

      if (
        ["array", "json"].includes(filter.type) &&
        [
          FilterOperator.Contains,
          FilterOperator.ContainedIn,
          FilterOperator.Overlap,
        ].includes(operator as FilterOperator)
      ) {
        let arrayValue = value;
        if (typeof value === "string") {
          try {
            arrayValue = JSON.parse(value);
          } catch (e) {
            arrayValue = value.split(",").map((v) => v.trim());
          }
        }

        if (!Array.isArray(arrayValue)) {
          arrayValue = [arrayValue];
        }

        return `${column}.${operator}.{${arrayValue.join(",")}}`;
      }

      // Default format
      return `${column}.${operator}.${value}`;
    })
    .join(",");

  return { [logicalOperator]: `(${filterStr})` };
};

interface SearchDataTableProps {
  columns: any[];
  projectId: string;
  tableId: string;
  pagination: {
    pageIndex: number;
    pageSize: number;
  };
  onFilterChange: (filters: Record<string, string>) => void;
  className?: string;
}

export function SearchDataTable({
  columns,
  projectId,
  tableId,
  pagination,
  onFilterChange,
  className,
}: SearchDataTableProps) {
  const [isExpanded, setIsExpanded] = useState(false);
  const [activeFilters, setActiveFilters] = useState<FilterCriteria[]>([]);
  const [logicalOperator, setLogicalOperator] = useState<"and" | "or">("and");
  const [error, setError] = useState<string | null>(null);
  const queryClient = useQueryClient();

  const form = useForm<Filters>({
    resolver: zodResolver(filtersSchema),
    defaultValues: {
      filters: [],
      logicalOperator: "and",
    },
  });

  // Create a new filter form
  const addFilterForm = useForm({
    resolver: zodResolver(
      z.object({
        column: z.string(),
        operator: z.nativeEnum(FilterOperator),
        value: z.any(),
      })
    ),
    defaultValues: {
      column: "",
      operator: FilterOperator.Equals,
      value: "",
    },
  });

  // Get selected column data
  const selectedColumnName = addFilterForm.watch("column");
  const selectedColumn = useMemo(() => {
    if (!selectedColumnName) return null;
    return columns.find((col) => col.name === selectedColumnName);
  }, [selectedColumnName, columns]);

  // Get column type based on selected column
  const columnType = useMemo(() => {
    if (!selectedColumn) return "text";
    return mapColumnTypeToCategory(selectedColumn.type || "text");
  }, [selectedColumn]);

  // Available operators for the current column type
  const availableOperators = useMemo(() => {
    return OPERATOR_BY_TYPE[columnType] || OPERATOR_BY_TYPE.text;
  }, [columnType]);

  // Reset operator and value when column changes
  useEffect(() => {
    if (selectedColumn) {
      addFilterForm.setValue("operator", availableOperators[0].value);

      // Set default value based on column type
      if (columnType === "boolean") {
        addFilterForm.setValue("value", "true");
      } else if (columnType === "date") {
        addFilterForm.setValue("value", new Date());
      } else {
        addFilterForm.setValue("value", "");
      }
    }
  }, [selectedColumn, columnType, addFilterForm, availableOperators]);

  // Add a new filter
  const handleAddFilter = useCallback(() => {
    const values = addFilterForm.getValues();

    // Skip if empty values or no column selected
    if (!values.column) return;

    // Skip empty values (except for boolean)
    if (values.value === "" && columnType !== "boolean") {
      // toast({ title: "Error", description: "Please enter a value for the filter", variant: "destructive" });
      return;
    }

    // Process the value based on column type
    let processedValue = values.value;

    if (columnType === "number") {
      // Convert to number if possible
      const num = parseFloat(values.value);
      if (!isNaN(num)) {
        processedValue = num;
      }
    } else if (columnType === "boolean") {
      // Convert string boolean to actual boolean
      if (values.value === "true") processedValue = true;
      else if (values.value === "false") processedValue = false;
      else processedValue = null;
    } else if (values.operator === FilterOperator.In) {
      // Convert comma-separated values to array
      if (typeof values.value === "string") {
        processedValue = values.value.split(",").map((v) => v.trim());
      }
    }

    // Create the filter criteria
    const newFilter: FilterCriteria = {
      column: values.column,
      operator: values.operator as FilterOperator,
      value: processedValue,
      type: columnType as any,
    };

    // Update the filters
    const updatedFilters = [...activeFilters, newFilter];
    setActiveFilters(updatedFilters);

    // Update form state
    form.setValue("filters", updatedFilters);
    form.setValue("logicalOperator", logicalOperator);

    // Build PostgREST filter parameters
    const filterParams = buildPostgrestFilterParams(
      updatedFilters,
      logicalOperator
    );
    onFilterChange(filterParams);

    // Show success feedback
    // Use toast notification if available in your UI library
    // toast({ title: "Filter applied", description: `Filtering by ${filter.column} ${filter.operator} ${filter.value}` });

    // Reset the add filter form
    addFilterForm.reset({
      column: "",
      operator: FilterOperator.Equals,
      value: "",
    });
  }, [
    addFilterForm,
    form,
    activeFilters,
    columnType,
    logicalOperator,
    onFilterChange,
  ]);

  // Remove a filter
  const handleRemoveFilter = useCallback(
    (index: number) => {
      const updatedFilters = activeFilters.filter((_, i) => i !== index);
      setActiveFilters(updatedFilters);

      // Update form state
      form.setValue("filters", updatedFilters);

      // Update PostgREST filter parameters
      const filterParams = buildPostgrestFilterParams(
        updatedFilters,
        logicalOperator
      );
      onFilterChange(filterParams);
    },
    [activeFilters, form, logicalOperator, onFilterChange]
  );

  // Toggle logical operator (AND/OR)
  const handleToggleLogicalOperator = useCallback(() => {
    const newOperator = logicalOperator === "and" ? "or" : "and";
    setLogicalOperator(newOperator);

    // Update form state
    form.setValue("logicalOperator", newOperator);

    // Update PostgREST filter parameters
    const filterParams = buildPostgrestFilterParams(activeFilters, newOperator);
    onFilterChange(filterParams);

    // Give feedback
    // toast({ title: "Logical operator changed", description: `Using ${newOperator.toUpperCase()} between filters` });
  }, [logicalOperator, activeFilters, form, onFilterChange]);

  // Clear all filters
  const handleClearAllFilters = useCallback(() => {
    setActiveFilters([]);
    form.reset({
      filters: [],
      logicalOperator: "and",
    });
    setLogicalOperator("and");
    onFilterChange({});

    // Refresh data in the table
    if (projectId && tableId) {
      queryClient.invalidateQueries({
        queryKey: [
          "rows",
          projectId,
          tableId,
          pagination.pageSize,
          pagination.pageIndex,
        ],
      });
    }
  }, [form, onFilterChange, queryClient, projectId, tableId, pagination]);

  // Render value input based on column type and operator
  const renderValueInput = () => {
    const operator = addFilterForm.watch("operator");

    if (columnType === "boolean" && operator === FilterOperator.Is) {
      return (
        <FormField
          control={addFilterForm.control}
          name="value"
          render={({ field }) => (
            <FormItem>
              <Select value={field.value} onValueChange={field.onChange}>
                <SelectTrigger>
                  <SelectValue placeholder="Select value" />
                </SelectTrigger>
                <SelectContent>
                  <SelectItem value="true">True</SelectItem>
                  <SelectItem value="false">False</SelectItem>
                  <SelectItem value="null">Null</SelectItem>
                </SelectContent>
              </Select>
            </FormItem>
          )}
        />
      );
    }

    if (columnType === "date") {
      return (
        <FormField
          control={addFilterForm.control}
          name="value"
          render={({ field }) => (
            <FormItem className="flex flex-col">
              <Popover>
                <PopoverTrigger asChild>
                  <FormControl>
                    <Button
                      variant="outline"
                      className={cn(
                        "w-full pl-3 text-left font-normal",
                        !field.value && "text-muted-foreground"
                      )}
                    >
                      {field.value ? (
                        format(
                          field.value instanceof Date
                            ? field.value
                            : new Date(field.value),
                          "PP"
                        )
                      ) : (
                        <span>Pick a date</span>
                      )}
                    </Button>
                  </FormControl>
                </PopoverTrigger>
                <PopoverContent className="w-auto p-0" align="start">
                  <Calendar
                    mode="single"
                    selected={
                      field.value instanceof Date
                        ? field.value
                        : new Date(field.value)
                    }
                    onSelect={field.onChange}
                    disabled={(date) =>
                      date > new Date() || date < new Date("1900-01-01")
                    }
                  />
                </PopoverContent>
              </Popover>
            </FormItem>
          )}
        />
      );
    }

    if (operator === FilterOperator.In) {
      return (
        <FormField
          control={addFilterForm.control}
          name="value"
          render={({ field }) => (
            <FormItem>
              <FormControl>
                <Input {...field} placeholder="Values separated by commas" />
              </FormControl>
              <FormMessage />
            </FormItem>
          )}
        />
      );
    }

    if (
      ["array", "json"].includes(columnType) &&
      [
        FilterOperator.Contains,
        FilterOperator.ContainedIn,
        FilterOperator.Overlap,
      ].includes(operator as FilterOperator)
    ) {
      return (
        <FormField
          control={addFilterForm.control}
          name="value"
          render={({ field }) => (
            <FormItem>
              <FormControl>
                <Input
                  {...field}
                  placeholder={
                    columnType === "json"
                      ? "JSON object"
                      : "Values separated by commas"
                  }
                />
              </FormControl>
              <FormMessage />
            </FormItem>
          )}
        />
      );
    }

    // Default input
    return (
      <FormField
        control={addFilterForm.control}
        name="value"
        render={({ field }) => (
          <FormItem>
            <FormControl>
              <Input
                {...field}
                type={columnType === "number" ? "number" : "text"}
                placeholder="Value"
              />
            </FormControl>
            <FormMessage />
          </FormItem>
        )}
      />
    );
  };

  return (
    <div className={cn("border rounded-lg p-4", className)}>
      <div className="flex flex-col space-y-4">
        <div className="flex items-center justify-between">
          <Button
            variant="ghost"
            onClick={() => setIsExpanded(!isExpanded)}
            className="flex items-center gap-2"
          >
            <Search className="h-4 w-4" />
            <span>Filter Data</span>
            {isExpanded ? (
              <ChevronUp className="h-4 w-4" />
            ) : (
              <ChevronDown className="h-4 w-4" />
            )}
          </Button>

          {activeFilters.length > 0 && (
            <div className="flex items-center gap-2">
              <span className="text-sm text-muted-foreground">
                {activeFilters.length} filter
                {activeFilters.length !== 1 ? "s" : ""} applied
              </span>
              <Button
                variant="ghost"
                size="sm"
                onClick={handleClearAllFilters}
                className="h-8 px-2 text-destructive"
              >
                <FilterX className="h-4 w-4 mr-1" />
                Clear All
              </Button>
            </div>
          )}
        </div>

        {/* Active filters display */}
        {activeFilters.length > 0 && (
          <div className="flex flex-wrap gap-2 mt-2">
            {activeFilters.map((filter, index) => (
              <Badge
                key={`${filter.column}-${index}`}
                variant="secondary"
                className="flex items-center gap-1 px-2 py-1"
              >
                <span className="font-medium">{filter.column}</span>
                <span className="text-muted-foreground mx-0.5">
                  {filter.operator}
                </span>
                <span>{formatFilterValueForDisplay(filter)}</span>
                <Button
                  variant="ghost"
                  size="sm"
                  onClick={() => handleRemoveFilter(index)}
                  className="h-4 w-4 p-0 ml-1 opacity-70 hover:opacity-100"
                >
                  <X className="h-3 w-3" />
                </Button>
              </Badge>
            ))}

            {activeFilters.length > 1 && (
              <Button
                variant="outline"
                size="sm"
                onClick={handleToggleLogicalOperator}
                className="h-7 px-2 capitalize"
              >
                {logicalOperator}
              </Button>
            )}
          </div>
        )}

        {/* Error message for invalid filters */}
        {isExpanded && activeFilters.length === 0 && (
          <div className="bg-destructive/10 text-destructive flex items-start p-4 mb-4 rounded-lg border border-destructive/20">
            <AlertCircle className="h-4 w-4 mr-2 mt-0.5" />
            <div>
              <h5 className="font-medium mb-1">No filters applied</h5>
              <p className="text-sm">
                Select a column, an operator, and a value to create a filter.
              </p>
            </div>
          </div>
        )}

        {/* Filter builder */}
        {isExpanded && (
          <Accordion
            type="single"
            collapsible
            defaultValue="addFilter"
            className="w-full"
          >
            <AccordionItem value="addFilter" className="border-none">
              <AccordionContent>
                <div className="mt-2 grid gap-4">
                  <div className="grid grid-cols-1 md:grid-cols-3 gap-4">
                    {/* Column selection */}
                    <FormField
                      control={addFilterForm.control}
                      name="column"
                      render={({ field }) => (
                        <FormItem>
                          <FormLabel>Column</FormLabel>
                          <Select
                            value={field.value}
                            onValueChange={field.onChange}
                          >
                            <SelectTrigger>
                              <SelectValue placeholder="Select column" />
                            </SelectTrigger>
                            <SelectContent>
                              {columns.map((column) => (
                                <SelectItem
                                  key={column.name}
                                  value={column.name}
                                >
                                  {column.name}
                                </SelectItem>
                              ))}
                            </SelectContent>
                          </Select>
                        </FormItem>
                      )}
                    />

                    {/* Operator selection */}
                    {selectedColumn && (
                      <FormField
                        control={addFilterForm.control}
                        name="operator"
                        render={({ field }) => (
                          <FormItem>
                            <FormLabel>Operator</FormLabel>
                            <Select
                              value={field.value}
                              onValueChange={field.onChange}
                            >
                              <SelectTrigger>
                                <SelectValue placeholder="Select operator" />
                              </SelectTrigger>
                              <SelectContent>
                                {availableOperators.map((op) => (
                                  <SelectItem key={op.value} value={op.value}>
                                    {op.label}
                                  </SelectItem>
                                ))}
                              </SelectContent>
                            </Select>
                          </FormItem>
                        )}
                      />
                    )}

                    {/* Value input */}
                    {selectedColumn && (
                      <div>
                        <FormLabel>Value</FormLabel>
                        {renderValueInput()}
                      </div>
                    )}
                  </div>

                  <div className="flex justify-between">
                    <div className="text-xs text-muted-foreground">
                      {activeFilters.length > 0 && (
                        <span>
                          Using <strong>{logicalOperator.toUpperCase()}</strong>{" "}
                          operator between filters
                        </span>
                      )}
                    </div>
                    <Button
                      type="button"
                      onClick={handleAddFilter}
                      disabled={
                        !selectedColumn ||
                        (addFilterForm.watch("value") === "" &&
                          columnType !== "boolean")
                      }
                    >
                      Add Filter
                    </Button>
                  </div>
                </div>
              </AccordionContent>
            </AccordionItem>
          </Accordion>
        )}
      </div>
    </div>
  );
}
