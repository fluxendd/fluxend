export type ColumnType =
  | "serial"
  | "text"
  | "varchar"
  | "integer"
  | "bigint"
  | "boolean"
  | "timestamp"
  | "timestamptz"
  | "date"
  | "decimal"
  | "numeric"
  | "real"
  | "double"
  | "json"
  | "jsonb"
  | "uuid"
  | "bytea";

export interface TableColumn {
  name: string;
  type: ColumnType;
  primary: boolean;
}

export interface CreateTableRequest {
  name: string;
  columns: TableColumn[];
}

export interface CreateTableFormData {
  tableName: string;
  columns: TableColumn[];
}

export const COLUMN_TYPE_OPTIONS: { value: ColumnType; label: string; description: string }[] = [
  { value: "serial", label: "Serial", description: "Auto-incrementing integer" },
  { value: "text", label: "Text", description: "Variable-length text" },
  { value: "varchar", label: "Varchar", description: "Variable-length character string" },
  { value: "integer", label: "Integer", description: "32-bit integer" },
  { value: "bigint", label: "Big Integer", description: "64-bit integer" },
  { value: "boolean", label: "Boolean", description: "True/false value" },
  { value: "timestamp", label: "Timestamp", description: "Date and time" },
  { value: "timestamptz", label: "Timestamp with Timezone", description: "Date and time with timezone" },
  { value: "date", label: "Date", description: "Date only" },
  { value: "decimal", label: "Decimal", description: "Exact numeric value" },
  { value: "numeric", label: "Numeric", description: "Exact numeric value" },
  { value: "real", label: "Real", description: "Single precision floating-point" },
  { value: "double", label: "Double", description: "Double precision floating-point" },
  { value: "json", label: "JSON", description: "JSON data" },
  { value: "jsonb", label: "JSONB", description: "Binary JSON data" },
  { value: "uuid", label: "UUID", description: "Universally unique identifier" },
  { value: "bytea", label: "Bytea", description: "Binary data" },
];