import { useForm, useFieldArray } from "react-hook-form";
import { zodResolver } from "@hookform/resolvers/zod";
import { z } from "zod";
import { Button } from "~/components/ui/button";
import { Input } from "~/components/ui/input";
import { Label } from "~/components/ui/label";
import {
  Card,
  CardContent,
  CardDescription,
  CardHeader,
  CardTitle,
} from "~/components/ui/card";
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from "~/components/ui/select";
import {
  Form,
  FormControl,
  FormField,
  FormItem,
  FormLabel,
  FormMessage,
} from "~/components/ui/form";
import { Checkbox } from "~/components/ui/checkbox";
import { Plus, Trash2, Database, Settings, Table } from "lucide-react";
import { Badge } from "~/components/ui/badge";
import { Separator } from "~/components/ui/separator";
import {
  COLUMN_TYPE_OPTIONS,
  type CreateTableRequest,
  type ColumnType,
} from "~/types/table";

const createTableSchema = z.object({
  tableName: z
      .string()
      .min(1, "Table name is required")
      .regex(
          /^[a-zA-Z_][a-zA-Z0-9_]*$/,
          "Table name must start with letter or underscore and contain only letters, numbers, and underscores"
      ),
  columns: z
      .array(
          z.object({
            name: z
                .string()
                .min(1, "Column name is required")
                .regex(
                    /^[a-zA-Z_][a-zA-Z0-9_]*$/,
                    "Column name must start with letter or underscore and contain only letters, numbers, and underscores"
                ),
            type: z.string() as z.ZodType<ColumnType>,
            primary: z.boolean(),
          })
      )
      .min(1, "At least one column is required"),
});

type CreateTableFormData = z.infer<typeof createTableSchema>;

interface TableFormProps {
  mode: "create" | "edit";
  initialData?: {
    tableName: string;
    columns: Array<{
      name: string;
      type: ColumnType;
      primary: boolean;
    }>;
  };
  onSubmit: (data: CreateTableFormData) => Promise<void>;
  isSubmitting: boolean;
}

export function TableForm({
                            mode,
                            initialData,
                            onSubmit,
                            isSubmitting,
                          }: TableFormProps) {
  const form = useForm<CreateTableFormData>({
    resolver: zodResolver(createTableSchema),
    defaultValues: initialData || {
      tableName: "",
      columns: [{ name: "id", type: "serial", primary: true }],
    },
  });

  const { fields, append, remove } = useFieldArray({
    control: form.control,
    name: "columns",
  });

  const addColumn = () => {
    append({ name: "", type: "text", primary: false });
  };

  const removeColumn = (index: number) => {
    if (fields.length > 1) {
      remove(index);
    }
  };

  const getTypeInfo = (type: string) => {
    const option = COLUMN_TYPE_OPTIONS.find(opt => opt.value === type);
    return option || { label: type, description: "", value: type };
  };

  const isEdit = mode === "edit";
  const submitButtonText = isSubmitting
      ? isEdit
          ? "Updating Table..."
          : "Creating Table..."
      : isEdit
          ? "Update Table"
          : "Create Table";

  return (
      <div className="max-w-4xl mx-auto space-y-6">
        <Form {...form}>
          <form onSubmit={form.handleSubmit(onSubmit)} className="space-y-6">
            {/* Header Section */}
            <Card>
              <CardHeader className="pb-4">
                <div className="flex items-center space-x-2">
                  <Database className="h-5 w-5 text-primary" />
                  <CardTitle className="text-xl">
                    {isEdit ? "Edit Table Schema" : "Create New Table"}
                  </CardTitle>
                </div>
                <CardDescription>
                  {isEdit
                      ? "Modify your table structure. Changes to column types may affect existing data."
                      : "Define your table name and column structure below."}
                </CardDescription>
              </CardHeader>

              <CardContent className="space-y-6">
                <FormField
                    control={form.control}
                    name="tableName"
                    render={({ field }) => (
                        <FormItem>
                          <FormLabel className="text-base font-medium">Table Name</FormLabel>
                          <FormControl>
                            <Input
                                placeholder="Enter table name (e.g., users, products, orders)"
                                {...field}
                                disabled={isEdit}
                                className="text-base"
                            />
                          </FormControl>
                          {isEdit && (
                              <p className="text-sm text-muted-foreground">
                                Table name cannot be changed after creation
                              </p>
                          )}
                          <FormMessage />
                        </FormItem>
                    )}
                />
              </CardContent>
            </Card>

            {/* Columns Section */}
            <Card>
              <CardHeader className="pb-4">
                <div className="flex items-center justify-between">
                  <div className="flex items-center space-x-2">
                    <Table className="h-5 w-5 text-primary" />
                    <CardTitle className="text-lg">Column Configuration</CardTitle>
                    <Badge variant="secondary" className="ml-2">
                      {fields.length} {fields.length === 1 ? 'column' : 'columns'}
                    </Badge>
                  </div>
                  <Button
                      type="button"
                      variant="outline"
                      size="sm"
                      onClick={addColumn}
                      className="flex items-center space-x-1"
                  >
                    <Plus className="h-4 w-4" />
                    <span>Add Column</span>
                  </Button>
                </div>
              </CardHeader>

              <CardContent className="space-y-4">
                {fields.map((field, index) => {
                  const currentType = form.watch(`columns.${index}.type`);
                  const currentName = form.watch(`columns.${index}.name`);
                  const isPrimary = form.watch(`columns.${index}.primary`);
                  const isIdColumn = isEdit && index === 0 && currentName === "id";
                  const typeInfo = getTypeInfo(currentType);

                  return (
                      <Card key={field.id} className="border-border/50 bg-muted/20">
                        <CardContent className="p-6">
                          <div className="flex items-start justify-between mb-4">
                            <div className="flex items-center space-x-2">
                              <Settings className="h-4 w-4 text-muted-foreground" />
                              <span className="font-medium text-sm text-muted-foreground">
                            Column {index + 1}
                          </span>
                              {isPrimary && (
                                  <Badge variant="default" className="text-xs px-2 py-0">
                                    Primary Key
                                  </Badge>
                              )}
                              {isIdColumn && (
                                  <Badge variant="secondary" className="text-xs px-2 py-0">
                                    System Column
                                  </Badge>
                              )}
                            </div>
                            {fields.length > 1 && !isIdColumn && (
                                <Button
                                    type="button"
                                    variant="ghost"
                                    size="sm"
                                    onClick={() => removeColumn(index)}
                                    className="text-muted-foreground hover:text-destructive"
                                >
                                  <Trash2 className="h-4 w-4" />
                                </Button>
                            )}
                          </div>

                          <div className="grid grid-cols-1 md:grid-cols-3 gap-4">
                            {/* Column Name */}
                            <div className="space-y-2">
                              <FormField
                                  control={form.control}
                                  name={`columns.${index}.name`}
                                  render={({ field }) => (
                                      <FormItem>
                                        <FormLabel className="text-sm font-medium">
                                          Column Name
                                        </FormLabel>
                                        <FormControl>
                                          <Input
                                              placeholder="Enter column name"
                                              {...field}
                                              disabled={isIdColumn}
                                              className={isIdColumn ? "bg-muted" : ""}
                                          />
                                        </FormControl>
                                        <FormMessage />
                                      </FormItem>
                                  )}
                              />
                            </div>

                            {/* Data Type */}
                            <div className="space-y-2">
                              <FormField
                                  control={form.control}
                                  name={`columns.${index}.type`}
                                  render={({ field }) => (
                                      <FormItem>
                                        <FormLabel className="text-sm font-medium">
                                          Data Type
                                        </FormLabel>
                                        <Select
                                            onValueChange={field.onChange}
                                            defaultValue={field.value}
                                            disabled={isIdColumn}
                                        >
                                          <FormControl>
                                            <SelectTrigger className={isIdColumn ? "bg-muted" : ""}>
                                              <SelectValue placeholder="Select data type" />
                                            </SelectTrigger>
                                          </FormControl>
                                          <SelectContent>
                                            {COLUMN_TYPE_OPTIONS.map((option) => (
                                                <SelectItem
                                                    key={option.value}
                                                    value={option.value}
                                                    className="py-2"
                                                >
                                                  {option.label}
                                                </SelectItem>
                                            ))}
                                          </SelectContent>
                                        </Select>
                                        <FormMessage />
                                      </FormItem>
                                  )}
                              />
                            </div>

                            {/* Options */}
                            <div className="space-y-2">
                              <FormLabel className="text-sm font-medium">Options</FormLabel>
                              <FormField
                                  control={form.control}
                                  name={`columns.${index}.primary`}
                                  render={({ field }) => (
                                      <FormItem>
                                        <FormControl>
                                          <div className="flex items-center space-x-2">
                                            <Checkbox
                                                checked={field.value}
                                                onCheckedChange={field.onChange}
                                                disabled={isIdColumn}
                                                id={`primary-${index}`}
                                            />
                                            <Label
                                                htmlFor={`primary-${index}`}
                                                className="text-sm font-normal cursor-pointer"
                                            >
                                              Primary Key
                                            </Label>
                                          </div>
                                        </FormControl>
                                        <FormMessage />
                                      </FormItem>
                                  )}
                              />
                            </div>
                          </div>

                          {/* Type Description */}
                          {currentType && typeInfo.description && (
                              <>
                                <Separator className="my-4" />
                                <div className="bg-muted/50 rounded-md p-3">
                                  <p className="text-sm text-muted-foreground">
                                    <strong>{typeInfo.label}:</strong> {typeInfo.description}
                                  </p>
                                </div>
                              </>
                          )}
                        </CardContent>
                      </Card>
                  );
                })}
              </CardContent>
            </Card>

            {/* Submit Section */}
            <Card>
              <CardContent className="pt-6">
                <div className="flex items-center justify-between">
                  <div className="text-sm text-muted-foreground">
                    {isEdit
                        ? "Review your changes before updating the table structure."
                        : "Ready to create your table? Click the button below to proceed."
                    }
                  </div>
                  <Button
                      type="submit"
                      disabled={isSubmitting}
                      className="min-w-[140px]"
                  >
                    {submitButtonText}
                  </Button>
                </div>
              </CardContent>
            </Card>
          </form>
        </Form>
      </div>
  );
}