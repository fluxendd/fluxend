export type ColumnType =
  | "integer"
  | "serial"
  | "varchar"
  | "text"
  | "boolean"
  | "date"
  | "timestamp"
  | "float"
  | "uuid"
  | "json";

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
  { value: "integer", label: "Integer", description: "32-bit integer" },
  { value: "serial", label: "Serial", description: "Auto-incrementing integer" },
  { value: "varchar", label: "Varchar", description: "Variable-length character string" },
  { value: "text", label: "Text", description: "Variable-length text" },
  { value: "boolean", label: "Boolean", description: "True/false value" },
  { value: "date", label: "Date", description: "Date only" },
  { value: "timestamp", label: "Timestamp", description: "Date and time" },
  { value: "float", label: "Float", description: "Floating-point number" },
  { value: "uuid", label: "UUID", description: "Universally unique identifier" },
  { value: "json", label: "JSON", description: "JSON data" },
];