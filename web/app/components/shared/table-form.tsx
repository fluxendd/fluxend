import { useForm, useFieldArray } from "react-hook-form";
import { zodResolver } from "@hookform/resolvers/zod";
import { z } from "zod";
import { Button } from "~/components/ui/button";
import { Input } from "~/components/ui/input";
import { Label } from "~/components/ui/label";
import {
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
import { Plus, Trash2 } from "lucide-react";
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

  const isEdit = mode === "edit";
  const submitButtonText = isSubmitting
    ? isEdit
      ? "Updating..."
      : "Creating..."
    : isEdit
    ? "Update Table"
    : "Create Table";

  return (
    <Form {...form}>
      <form onSubmit={form.handleSubmit(onSubmit)} className="space-y-8">
        <CardHeader>
          <CardTitle>{isEdit ? "Edit Table" : "Create New Table"}</CardTitle>
          <CardDescription>
            {isEdit
              ? "Modify your table structure. Be careful when changing column types as it may affect existing data."
              : "Define your table name and column structure"}
          </CardDescription>
        </CardHeader>
        <CardContent className="space-y-6">
          <FormField
            control={form.control}
            name="tableName"
            render={({ field }) => (
              <FormItem>
                <FormLabel>Table Name</FormLabel>
                <FormControl>
                  <Input
                    placeholder="e.g., users, products, orders"
                    {...field}
                    disabled={isEdit}
                  />
                </FormControl>
                <FormMessage />
              </FormItem>
            )}
          />

          <div className="space-y-8">
            <div className="flex items-center justify-between">
              <Label className="text-base font-medium">Columns</Label>
              <Button
                type="button"
                variant="outline"
                size="sm"
                className="cursor-pointer"
                onClick={addColumn}
              >
                <Plus className="w-4 h-4 mr-2" />
                Add Column
              </Button>
            </div>

            <div className="space-y-4">
              <div className="grid grid-cols-12 items-end gap-4 pb-4">
                <div className="col-span-4 ">
                  <FormLabel>Column Name</FormLabel>
                </div>
                <div className="col-span-4 ">
                  <FormLabel>Data Type</FormLabel>
                </div>
                <div className="col-span-2 ">
                  <FormLabel>Primary</FormLabel>
                </div>
                <div className="col-span-2 ">
                  <FormLabel>Actions</FormLabel>
                </div>
                {fields.map((field, index) => (
                  <div className="col-span-12" key={field.id}>
                    <div className="grid grid-cols-12 border-b-1 gap-4 pb-4 items-center">
                      <div className="col-span-4">
                        <FormField
                          control={form.control}
                          name={`columns.${index}.name`}
                          render={({ field }) => (
                            <FormItem>
                              <FormControl>
                                <Input
                                  placeholder="Column name"
                                  {...field}
                                  disabled={
                                    isEdit && index === 0 && field.value === "id"
                                  }
                                  className={
                                    isEdit && index === 0 && field.value === "id"
                                      ? "bg-muted"
                                      : ""
                                  }
                                />
                              </FormControl>
                              <FormMessage />
                            </FormItem>
                          )}
                        />
                      </div>
                      <div className="col-span-4">
                        <FormField
                          control={form.control}
                          name={`columns.${index}.type`}
                          render={({ field }) => (
                            <FormItem className="w-full">
                              <Select
                                onValueChange={field.onChange}
                                defaultValue={field.value}
                                disabled={
                                  isEdit &&
                                  index === 0 &&
                                  form.getValues(`columns.${index}.name`) === "id"
                                }
                              >
                                <FormControl className="w-[90%]">
                                  <SelectTrigger
                                    className={
                                      isEdit &&
                                      index === 0 &&
                                      form.getValues(`columns.${index}.name`) === "id"
                                        ? "bg-muted"
                                        : ""
                                    }
                                  >
                                    <SelectValue placeholder="Select type" />
                                  </SelectTrigger>
                                </FormControl>
                                <SelectContent>
                                  {COLUMN_TYPE_OPTIONS.map((option) => (
                                    <SelectItem
                                      key={option.value}
                                      value={option.value}
                                    >
                                      <div>
                                        <div className="font-medium">
                                          {option.label}
                                        </div>
                                        <div className="text-xs text-muted-foreground">
                                          {option.description}
                                        </div>
                                      </div>
                                    </SelectItem>
                                  ))}
                                </SelectContent>
                              </Select>
                              <FormMessage />
                            </FormItem>
                          )}
                        />
                      </div>
                      <div className="col-span-2">
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
                                    disabled={
                                      isEdit &&
                                      index === 0 &&
                                      form.getValues(`columns.${index}.name`) === "id"
                                    }
                                    className={
                                      isEdit &&
                                      index === 0 &&
                                      form.getValues(`columns.${index}.name`) === "id"
                                        ? ""
                                        : "cursor-pointer"
                                    }
                                  />
                                </div>
                              </FormControl>
                              <FormMessage />
                            </FormItem>
                          )}
                        />
                      </div>
                      <div className="col-span-2">
                        <Button
                          type="button"
                          variant="ghost"
                          size="icon"
                          onClick={() => removeColumn(index)}
                          disabled={
                            (isEdit &&
                              index === 0 &&
                              form.getValues(`columns.${index}.name`) === "id") ||
                            fields.length === 1
                          }
                          className={
                            (isEdit &&
                              index === 0 &&
                              form.getValues(`columns.${index}.name`) === "id") ||
                            fields.length === 1
                              ? ""
                              : "cursor-pointer"
                          }
                        >
                          <Trash2 className="w-4 h-4" />
                        </Button>
                      </div>
                    </div>
                  </div>
                ))}
              </div>
            </div>
          </div>

          <div className="flex justify-end">
            <Button
              type="submit"
              disabled={isSubmitting}
              className="cursor-pointer"
            >
              {submitButtonText}
            </Button>
          </div>
        </CardContent>
      </form>
    </Form>
  );
}