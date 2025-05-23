# SearchDataTable Component

The SearchDataTable component is a powerful tool for filtering data in tables using PostgREST's vertical filtering capabilities. It allows users to build complex queries without needing to know the underlying API structure.

## Overview

This component integrates with the existing DataTableWrapper to add powerful filtering capabilities. It supports most PostgreSQL data types and a wide range of operators.

## Features

- Filter by a wide range of operators: =, >, <, >=, <=, LIKE, ILIKE, etc.
- Support for text, numbers, dates, booleans, arrays, and JSON data
- Combine multiple filters with AND/OR logic
- Full-text search capabilities
- Range operators for date ranges and number ranges
- Array operators: contains, contained in, overlaps

## Usage

The component is designed to be used with the existing table components:

```tsx
<SearchDataTableWrapper
  columns={columns}
  rawColumns={rawColumns}
  data={rows}
  isLoading={isLoading}
  emptyMessage="No results found."
  className="w-full h-full"
  pagination={pagination}
  totalRows={totalCount}
  projectId={projectId}
  collectionId={collectionId}
  onFilterChange={handleFilterChange}
  onPaginationChange={handlePaginationChange}
/>
```

## Filter Operators

The component supports the following operators based on PostgREST's capabilities:

### Equality & Comparison
- `eq` - equals
- `neq` - not equals
- `gt` - greater than
- `gte` - greater than or equal
- `lt` - less than
- `lte` - less than or equal

### Text Operators
- `like` - LIKE pattern matching (% can be replaced with *)
- `ilike` - Case-insensitive LIKE
- `match` - Regular expression match (~)
- `imatch` - Case-insensitive regular expression match (~*)

### Boolean Operators
- `is` - IS comparison for true, false, null

### Array & JSON Operators
- `cs` - Contains (@>)
- `cd` - Contained in (<@)
- `ov` - Overlap (&&)

### Full-Text Search
- `fts` - Full-text search using to_tsquery
- `plfts` - Full-text search using plainto_tsquery
- `phfts` - Full-text search using phraseto_tsquery
- `wfts` - Full-text search using websearch_to_tsquery

### Range Operators
- `sl` - Strictly left of (<<)
- `sr` - Strictly right of (>>)
- `nxr` - Does not extend to the right of (&<)
- `nxl` - Does not extend to the left of (&>)
- `adj` - Is adjacent to (-|-)

## How Filters Are Applied

Filters are converted to PostgREST query parameters. For example:

```
?name=eq.John
?age=gt.30
?or=(age.lt.18,age.gt.65)
?tags=cs.{project,important}
```

## Implementation Details

The component uses React Hook Form for form validation and zod for type safety. It converts UI filter selections into PostgREST compatible query parameters.