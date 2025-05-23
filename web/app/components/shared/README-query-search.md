# QuerySearchBox Component

The QuerySearchBox is a SQL-like query interface for filtering data in collections using PostgREST's powerful filtering capabilities. This component allows users to write simple query expressions that get translated to PostgREST filter parameters.

## Usage

```tsx
import { QuerySearchBox } from "~/components/shared/query-search-box";

// In your component
<QuerySearchBox
  columns={columnMetadata}
  onQueryChange={handleFilterChange}
  className="mb-2"
/>
```

## Features

- SQL-like query syntax (`column = 'value'`)
- Column name autocomplete suggestions
- Support for all PostgREST operators
- Real-time feedback and validation
- Interactive examples and help

## Query Syntax

Queries follow a simple format: `column operator value`

### Examples

- `name = 'John'` - Equal to exact value
- `age > 30` - Greater than
- `price >= 10.5` - Greater than or equal
- `status != 'completed'` - Not equal
- `description like '*important*'` - Pattern matching (use * instead of %)
- `category in (1,2,3)` - In a list of values
- `is_active is true` - Boolean comparison
- `manager is null` - NULL check
- `tags cs {sports,outdoor}` - Array contains
- `content fts 'search terms'` - Full-text search

### Supported Operators

| Query Syntax | PostgREST | Description |
|--------------|-----------|-------------|
| `=` | `eq` | Equal |
| `!=` | `neq` | Not equal |
| `>` | `gt` | Greater than |
| `>=` | `gte` | Greater than or equal |
| `<` | `lt` | Less than |
| `<=` | `lte` | Less than or equal |
| `like` | `like` | LIKE pattern matching |
| `ilike` | `ilike` | Case-insensitive LIKE |
| `in` | `in` | In a list of values |
| `is` | `is` | For boolean and NULL values |
| `cs` | `cs` | Array contains |
| `cd` | `cd` | Array contained in |
| `ov` | `ov` | Array overlap |
| `fts` | `fts` | Full-text search |
| `match` | `match` | Regular expression match |
| `imatch` | `imatch` | Case-insensitive regex match |

## Implementation Details

The component transforms user queries into PostgREST filter parameters, which are then applied to API requests. For example, a query like `name = 'John'` becomes a parameter like `?name=eq.John` in the final API request.

## Accessibility

- Keyboard navigation supported (Enter to execute, Escape to clear)
- Clear visual feedback for active filters
- Error messages for invalid queries

## Column Suggestions

As users type, the component provides autocomplete suggestions for column names based on the available columns in the current collection.