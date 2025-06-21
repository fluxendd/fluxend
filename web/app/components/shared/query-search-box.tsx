import { useState, useRef, useEffect, useCallback } from "react";
import { cn } from "~/lib/utils";
import { Button } from "~/components/ui/button";
import { Input } from "~/components/ui/input";
import {
  Search,
  X,
  ArrowRight,
  Filter,
  ChevronsUpDown,
  Info,
} from "lucide-react";
import { ScrollArea } from "~/components/ui/scroll-area";
import { Badge } from "~/components/ui/badge";
import {
  Command,
  CommandGroup,
  CommandItem,
  CommandList,
} from "~/components/ui/command";
import {
  Popover,
  PopoverContent,
  PopoverTrigger,
} from "~/components/ui/popover";
import { motion } from "motion/react";

interface QuerySearchBoxProps {
  columns: any[];
  onQueryChange: (params: Record<string, string>) => void;
  className?: string;
}

export function QuerySearchBox({
  columns,
  onQueryChange,
  className,
}: QuerySearchBoxProps) {
  const [query, setQuery] = useState<string>("");
  const [activeQuery, setActiveQuery] = useState<string | null>(null);
  const [isShowingHelp, setIsShowingHelp] = useState(false);
  const [errorMessage, setErrorMessage] = useState<string | null>(null);
  const [isColumnPickerOpen, setIsColumnPickerOpen] = useState(false);
  const [columnSuggestions, setColumnSuggestions] = useState<string[]>([]);
  const [suggestionPosition, setSuggestionPosition] = useState<{
    top: number;
    left: number;
    width: number;
    maxHeight: number;
  }>({ top: 0, left: 0, width: 0, maxHeight: 200 });
  const inputRef = useRef<HTMLInputElement>(null);
  const suggestionsRef = useRef<HTMLDivElement>(null);

  // Update suggestion dropdown position
  const updateSuggestionPosition = useCallback(() => {
    if (inputRef.current) {
      const inputRect = inputRef.current.getBoundingClientRect();
      const windowHeight = window.innerHeight;
      const spaceBelow = windowHeight - inputRect.bottom;
      const maxHeight = Math.min(200, spaceBelow - 10); // Leave 10px padding

      setSuggestionPosition({
        top: inputRect.bottom + window.scrollY + 2, // Add small gap
        left: inputRect.left + window.scrollX,
        width: inputRect.width,
        maxHeight: maxHeight,
      });
    }
  }, []);

  // Focus the input when the component mounts
  useEffect(() => {
    if (inputRef.current) {
      inputRef.current.focus();
    }

    // Add click outside listener to close suggestions
    const handleClickOutside = (event: MouseEvent) => {
      if (
        suggestionsRef.current &&
        !suggestionsRef.current.contains(event.target as Node) &&
        inputRef.current &&
        !inputRef.current.contains(event.target as Node)
      ) {
        setIsColumnPickerOpen(false);
      }
    };

    // Handle window resize to update dropdown position
    const handleResize = () => {
      if (isColumnPickerOpen) {
        updateSuggestionPosition();
      }
    };

    document.addEventListener("mousedown", handleClickOutside);
    window.addEventListener("resize", handleResize);
    window.addEventListener("scroll", handleResize);

    return () => {
      document.removeEventListener("mousedown", handleClickOutside);
      window.removeEventListener("resize", handleResize);
      window.removeEventListener("scroll", handleResize);
    };
  }, [isColumnPickerOpen, updateSuggestionPosition]);

  const handleQueryChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    const newQuery = e.target.value;
    setQuery(newQuery);
    setErrorMessage(null);

    // Check if we're at the beginning of a query or after a space
    // to suggest column names
    const lastSpaceIndex = newQuery.lastIndexOf(" ");
    const currentToken =
      lastSpaceIndex === -1 ? newQuery : newQuery.substring(lastSpaceIndex + 1);

    if (
      !currentToken.includes("=") &&
      !currentToken.includes(">") &&
      !currentToken.includes("<") &&
      !currentToken.includes(" ")
    ) {
      const matchingColumns = columns
        .map((col) => col.name)
        .filter((name) =>
          name.toLowerCase().includes(currentToken.toLowerCase())
        );

      setColumnSuggestions(matchingColumns);

      // Only open the picker if we have suggestions and the user has typed something
      if (matchingColumns.length > 0 && currentToken.length > 0) {
        // Get input position for dropdown positioning
        setIsColumnPickerOpen(true);
        updateSuggestionPosition();
      }
    }
  };

  const clearQuery = () => {
    setQuery("");
    setActiveQuery(null);
    setErrorMessage(null);
    onQueryChange({});
    if (inputRef.current) {
      inputRef.current.focus();
    }
  };

  const parseQuery = useCallback(() => {
    if (!query.trim()) {
      // Clear any existing filters
      onQueryChange({});
      return;
    }

    try {
      // Check for AND/OR operators in the query
      const trimmedQuery = query.trim();
      const andParts = trimmedQuery.split(/\s+AND\s+/i);
      const orParts = trimmedQuery.split(/\s+OR\s+/i);

      // If we have multiple conditions with AND or OR
      if (andParts.length > 1 || orParts.length > 1) {
        // We'll use whichever operator appears in the query
        const isAnd = andParts.length > orParts.length;
        const parts = isAnd ? andParts : orParts;
        const logicalOp = isAnd ? "and" : "or";

        // Parse each condition separately
        const conditions: string[] = [];

        for (const part of parts) {
          // Basic pattern for individual conditions
          const conditionPattern =
            /^([a-zA-Z0-9_]+)\s*(=|!=|>|>=|<|<=|like|ilike|in|is|cs|cd|ov|fts|plfts|phfts|wfts|match|imatch)\s*(['"]?)(.*?)\3$/i;

          const match = part.match(conditionPattern);
          if (!match) {
            setErrorMessage(`Invalid condition format: ${part}`);
            return;
          }

          const [_, columnName, operator, quote, value] = match;

          // Verify column exists
          const columnExists = columns.some((col) => col.name === columnName);
          if (!columnExists) {
            setErrorMessage(`Column "${columnName}" not found`);
            return;
          }

          // Map the operator to PostgREST format
          const operatorMap: Record<string, string> = {
            "=": "eq",
            "!=": "neq",
            ">": "gt",
            ">=": "gte",
            "<": "lt",
            "<=": "lte",
            like: "like",
            ilike: "ilike",
            in: "in",
            is: "is",
            cs: "cs",
            cd: "cd",
            ov: "ov",
            fts: "fts",
            plfts: "plfts",
            phfts: "phfts",
            wfts: "wfts",
            match: "match",
            imatch: "imatch",
          };

          const pgOperator = operatorMap[operator.toLowerCase()];
          if (!pgOperator) {
            setErrorMessage(`Unsupported operator: ${operator}`);
            return;
          }

          // Process the value based on operator
          let processedValue = value.trim();

          // Handle NULL value for 'is' operator
          if (pgOperator === "is" && processedValue.toLowerCase() === "null") {
            processedValue = "null";
          }

          // Handle boolean values for 'is' operator
          if (
            pgOperator === "is" &&
            ["true", "false"].includes(processedValue.toLowerCase())
          ) {
            processedValue = processedValue.toLowerCase();
          }

          // Handle 'in' operator
          if (pgOperator === "in") {
            if (
              processedValue.startsWith("(") &&
              processedValue.endsWith(")")
            ) {
              // Already in correct format
            } else {
              processedValue = `(${processedValue})`;
            }
          }

          // For like/ilike, replace * with %
          if (["like", "ilike"].includes(pgOperator)) {
            processedValue = processedValue.replace(/\*/g, "%");
          }

          // Handle array operators (cs, cd, ov)
          if (["cs", "cd", "ov"].includes(pgOperator)) {
            if (
              !(processedValue.startsWith("{") && processedValue.endsWith("}"))
            ) {
              if (
                processedValue.startsWith("(") &&
                processedValue.endsWith(")")
              ) {
                processedValue = processedValue.substring(
                  1,
                  processedValue.length - 1
                );
              }
              processedValue = `{${processedValue}}`;
            }
          }

          // Handle full-text search operators
          if (["fts", "plfts", "phfts", "wfts"].includes(pgOperator)) {
            if (processedValue.includes(":")) {
              const [lang, terms] = processedValue.split(":", 2);
              processedValue = `(${lang}).${terms}`;
            } else {
              processedValue = `(english).${processedValue}`;
            }
          }

          // Add the condition to our list
          conditions.push(`${columnName}.${pgOperator}.${processedValue}`);
        }

        // Create the filter param with logical operator
        const filterParam = {
          [logicalOp]: `(${conditions.join(",")})`,
        };

        // Apply the filter
        onQueryChange(filterParam);
        setActiveQuery(query.trim());
        setErrorMessage(null);
      } else {
        // Single condition case (original code path)
        // Basic pattern for supported queries like:
        // column = 'value'
        // column > 10
        // column like '*pattern*'
        const queryPattern =
          /^([a-zA-Z0-9_]+)\s*(=|!=|>|>=|<|<=|like|ilike|in|is|cs|cd|ov|fts|plfts|phfts|wfts|match|imatch)\s*(['"]?)(.*?)\3$/i;

        const match = trimmedQuery.match(queryPattern);

        if (!match) {
          setErrorMessage(
            "Invalid query format. Try: column = 'value' or column > 10"
          );
          return;
        }

        const [_, columnName, operator, quote, value] = match;

        // Verify column exists
        const columnExists = columns.some((col) => col.name === columnName);
        if (!columnExists) {
          setErrorMessage(`Column "${columnName}" not found`);
          return;
        }

        // Map the operator to PostgREST format
        const operatorMap: Record<string, string> = {
          "=": "eq",
          "!=": "neq",
          ">": "gt",
          ">=": "gte",
          "<": "lt",
          "<=": "lte",
          like: "like",
          ilike: "ilike",
          in: "in",
          is: "is",
          cs: "cs",
          cd: "cd",
          ov: "ov",
          fts: "fts",
          plfts: "plfts",
          phfts: "phfts",
          wfts: "wfts",
          match: "match",
          imatch: "imatch",
        };

        const pgOperator = operatorMap[operator.toLowerCase()];

        if (!pgOperator) {
          setErrorMessage(`Unsupported operator: ${operator}`);
          return;
        }

        // Handle special cases for values
        let processedValue = value.trim();

        // Handle NULL value for 'is' operator
        if (pgOperator === "is" && processedValue.toLowerCase() === "null") {
          processedValue = "null";
        }

        // Handle boolean values for 'is' operator
        if (
          pgOperator === "is" &&
          ["true", "false"].includes(processedValue.toLowerCase())
        ) {
          processedValue = processedValue.toLowerCase();
        }

        // Handle 'in' operator
        if (pgOperator === "in") {
          // If the value is in parentheses, use them directly
          if (processedValue.startsWith("(") && processedValue.endsWith(")")) {
            // We already have the correct format
          } else {
            // Wrap it in parentheses
            processedValue = `(${processedValue})`;
          }
        }

        // For like/ilike, replace * with %
        if (["like", "ilike"].includes(pgOperator)) {
          processedValue = processedValue.replace(/\*/g, "%");
        }

        // Handle array operators (cs, cd, ov)
        if (["cs", "cd", "ov"].includes(pgOperator)) {
          // If the value is already in curly braces format
          if (
            !(processedValue.startsWith("{") && processedValue.endsWith("}"))
          ) {
            // Convert comma-separated values to curly brace format
            // Check if it's already in parentheses
            if (
              processedValue.startsWith("(") &&
              processedValue.endsWith(")")
            ) {
              // Extract the content within parentheses
              processedValue = processedValue.substring(
                1,
                processedValue.length - 1
              );
            }
            // Wrap with curly braces
            processedValue = `{${processedValue}}`;
          }
        }

        // Handle full-text search operators
        if (["fts", "plfts", "phfts", "wfts"].includes(pgOperator)) {
          // Check if language is specified
          if (processedValue.includes(":")) {
            const [lang, terms] = processedValue.split(":", 2);
            processedValue = `(${lang}).${terms}`;
          } else {
            // Use English as default language
            processedValue = `(english).${processedValue}`;
          }
        }

        // Create the filter param
        const filterParam = { [columnName]: `${pgOperator}.${processedValue}` };

        // Apply the filter
        onQueryChange(filterParam);
        setActiveQuery(query.trim());
        setErrorMessage(null);
      }
    } catch (error) {
      console.error("Error parsing query:", error);
      setErrorMessage("Invalid query format. See help for examples.");
    }
  }, [query, columns, onQueryChange]);

  const handleKeyDown = (e: React.KeyboardEvent<HTMLInputElement>) => {
    if (e.key === "Enter") {
      e.preventDefault();
      parseQuery();
    } else if (e.key === "Escape") {
      clearQuery();
    }
  };

  // Insert a column name suggestion into the query
  const insertColumnName = (columnName: string) => {
    const lastSpaceIndex = query.lastIndexOf(" ");

    // If we're at the beginning of the query or after a space,
    // replace the current token with the column name
    if (lastSpaceIndex === -1) {
      setQuery(columnName + " ");
    } else {
      setQuery(query.substring(0, lastSpaceIndex + 1) + columnName + " ");
    }

    // Add a suggested operator based on column type
    const selectedColumn = columns.find((col) => col.name === columnName);

    if (selectedColumn && selectedColumn.type) {
      const type = selectedColumn.type.toLowerCase();
      let suggestedOperator = "= ";

      if (
        type.includes("text") ||
        type.includes("char") ||
        type.includes("varchar")
      ) {
        suggestedOperator = "= '";
      } else if (
        type.includes("int") ||
        type.includes("numeric") ||
        type.includes("decimal")
      ) {
        suggestedOperator = "= ";
      } else if (type.includes("bool")) {
        suggestedOperator = "is ";
      } else if (type.includes("array")) {
        suggestedOperator = "cs ";
      }

      setQuery((prevQuery) => prevQuery + suggestedOperator);
    }

    setIsColumnPickerOpen(false);
    if (inputRef.current) {
      inputRef.current.focus();
    }
  };

  return (
    <div className={cn("flex flex-col space-y-2", className)}>
      <div className="relative">
        <div className="flex items-center">
          <div className="relative w-full">
            <Search className="absolute left-2 top-1/2 h-4 w-4 -translate-y-1/2 text-muted-foreground" />
            <Input
              ref={inputRef}
              value={query}
              onChange={handleQueryChange}
              onKeyDown={handleKeyDown}
              placeholder="Query..."
              className="pl-8 pr-[76px] w-full border focus-visible:ring-0 focus-visible:ring-offset-0 rounded-lg"
            />
            <div className="absolute right-1 top-1/2 -translate-y-1/2 flex items-center space-x-1">
              {query && (
                <Button
                  variant="ghost"
                  size="icon"
                  className="h-7 w-7 hover:bg-accent cursor-pointer"
                  onClick={clearQuery}
                  title="Clear"
                >
                  <X className="h-4 w-4" />
                </Button>
              )}
              <Button
                className={cn("h-7 w-7 hover:bg-accent cursor-pointer", {
                  "bg-accent": isShowingHelp,
                })}
                size="icon"
                variant="ghost"
                onClick={() => setIsShowingHelp(!isShowingHelp)}
                title="Query Examples"
              >
                <Info />
              </Button>
              <Button
                className="h-7 w-7 hover:bg-accent"
                size="icon"
                variant="ghost"
                onClick={parseQuery}
                title="Search"
              >
                <ArrowRight className="h-4 w-4" />
              </Button>
            </div>
          </div>
        </div>
        {isColumnPickerOpen && columnSuggestions.length > 0 && (
          <div
            ref={suggestionsRef}
            style={{
              position: "fixed",
              zIndex: 9999,
              top: `${suggestionPosition.top}px`,
              left: `${suggestionPosition.left}px`,
              width: `${suggestionPosition.width}px`,
              maxHeight: `${suggestionPosition.maxHeight}px`,
            }}
            className="mt-1 bg-background border rounded-md shadow-lg overflow-auto"
          >
            <div className="p-1">
              <div className="px-2 py-1 text-xs text-muted-foreground font-medium bg-muted/50 sticky top-0">
                Column suggestions
              </div>
              {columnSuggestions.slice(0, 8).map((column) => (
                <div
                  key={column}
                  onClick={() => {
                    insertColumnName(column);
                    setIsColumnPickerOpen(false);
                  }}
                  className="px-2 py-1.5 text-sm rounded-sm hover:bg-accent hover:text-accent-foreground cursor-pointer"
                >
                  {column}
                </div>
              ))}
            </div>
          </div>
        )}
      </div>

      {errorMessage && (
        <div className="text-sm text-destructive">{errorMessage}</div>
      )}

      {activeQuery && (
        <div className="flex items-center gap-2">
          <Badge
            variant="secondary"
            className="flex items-center gap-1.5 px-2 py-1 bg-accent/20 border border-accent/30 text-accent-foreground"
          >
            <Filter className="h-3 w-3" />
            <span className="font-normal">{activeQuery}</span>
            <Button
              variant="ghost"
              size="sm"
              onClick={clearQuery}
              className="h-4 w-4 p-0 ml-1 rounded-full opacity-70 hover:opacity-100 hover:bg-accent/30"
            >
              <X className="h-3 w-3" />
            </Button>
          </Badge>
        </div>
      )}

      {isShowingHelp && (
        <motion.div
          initial={{ height: 0, opacity: 0 }}
          animate={{
            height: "auto",
            opacity: 1,
            transition: {
              type: "ease",
              ease: "easeInOut",
              duration: 0.3,
            },
          }}
        >
          <div className="bg-muted/30 rounded-md border p-3 border-border/40">
            <h4 className="text-sm font-medium mb-2">Query Examples:</h4>
            <ScrollArea className="h-[150px]">
              <ul className="text-xs space-y-2 text-muted-foreground">
                <li>
                  <code className="bg-muted p-0.5 rounded">name = 'John'</code>{" "}
                  - Exact match
                </li>
                <li>
                  <code className="bg-muted p-0.5 rounded">age &gt; 30</code> -
                  Greater than
                </li>
                <li>
                  <code className="bg-muted p-0.5 rounded">age &gt;= 30</code> -
                  Greater than or equal
                </li>
                <li>
                  <code className="bg-muted p-0.5 rounded">
                    name like '*Smith*'
                  </code>{" "}
                  - Pattern match (use * instead of %)
                </li>
                <li>
                  <code className="bg-muted p-0.5 rounded">status is true</code>{" "}
                  - Boolean value
                </li>
                <li>
                  <code className="bg-muted p-0.5 rounded">
                    manager is null
                  </code>{" "}
                  - NULL value
                </li>
                <li>
                  <code className="bg-muted p-0.5 rounded">id in (1,2,3)</code>{" "}
                  - List of values
                </li>
                <li>
                  <code className="bg-muted p-0.5 rounded">
                    tags cs {"{"}"sport,outdoor{"}"}
                  </code>{" "}
                  - Array contains
                </li>
                <li>
                  <code className="bg-muted p-0.5 rounded">
                    data cd {"{"}"1,2{"}"}
                  </code>{" "}
                  - Array contained in
                </li>
                <li>
                  <code className="bg-muted p-0.5 rounded">
                    ranges ov {"{"}"20,30{"}"}
                  </code>{" "}
                  - Array overlap
                </li>
                <li>
                  <code className="bg-muted p-0.5 rounded">
                    description fts 'search terms'
                  </code>{" "}
                  - Full-text search
                </li>
                <li>
                  <code className="bg-muted p-0.5 rounded">
                    content fts english:'query'
                  </code>{" "}
                  - Full-text search with language
                </li>
                <li>
                  <code className="bg-muted p-0.5 rounded">
                    name = 'John' AND age &gt; 30
                  </code>{" "}
                  - Multiple conditions with AND
                </li>
                <li>
                  <code className="bg-muted p-0.5 rounded">
                    status = 'active' OR status = 'pending'
                  </code>{" "}
                  - Multiple conditions with OR
                </li>
              </ul>
            </ScrollArea>
          </div>
        </motion.div>
      )}
    </div>
  );
}
