import { useState } from "react";
import { useForm, useFieldArray } from "react-hook-form";
import { zodResolver } from "@hookform/resolvers/zod";
import { z } from "zod";
import { Button } from "~/components/ui/button";
import { Input } from "~/components/ui/input";
import { Label } from "~/components/ui/label";
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "~/components/ui/card";
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from "~/components/ui/select";
import { Form, FormControl, FormField, FormItem, FormLabel, FormMessage } from "~/components/ui/form";
import { Checkbox } from "~/components/ui/checkbox";
import { AppHeader } from "~/components/shared/header";
import { Plus, Trash2 } from "lucide-react";
import { useParams, useNavigate } from "react-router";
import { useQueryClient } from "@tanstack/react-query";
import { createCollection } from "~/services/collections";
import { COLUMN_TYPE_OPTIONS, type CreateTableRequest, type ColumnType } from "~/types/table";

const createTableSchema = z.object({
  tableName: z.string().min(1, "Table name is required").regex(/^[a-zA-Z_][a-zA-Z0-9_]*$/, "Table name must start with letter or underscore and contain only letters, numbers, and underscores"),
  columns: z.array(z.object({
    name: z.string().min(1, "Column name is required").regex(/^[a-zA-Z_][a-zA-Z0-9_]*$/, "Column name must start with letter or underscore and contain only letters, numbers, and underscores"),
    type: z.string() as z.ZodType<ColumnType>,
    primary: z.boolean()
  })).min(1, "At least one column is required")
});

type CreateTableFormData = z.infer<typeof createTableSchema>;

export default function CreateCollection() {
  const { projectId } = useParams();
  const navigate = useNavigate();
  const queryClient = useQueryClient();
  const [isSubmitting, setIsSubmitting] = useState(false);

  const form = useForm<CreateTableFormData>({
    resolver: zodResolver(createTableSchema),
    defaultValues: {
      tableName: "",
      columns: [
        { name: "id", type: "serial", primary: true }
      ]
    }
  });

  const { fields, append, remove } = useFieldArray({
    control: form.control,
    name: "columns"
  });

  const addColumn = () => {
    append({ name: "", type: "text", primary: false });
  };

  const removeColumn = (index: number) => {
    if (fields.length > 1) {
      remove(index);
    }
  };

  const onSubmit = async (data: CreateTableFormData) => {
    if (!projectId) {
      alert("Project ID is required");
      return;
    }

    setIsSubmitting(true);
    try {
      const requestBody: CreateTableRequest = {
        name: data.tableName,
        columns: data.columns
      };

      const response = await createCollection({ headers: {} }, projectId, requestBody);
      
      if (response.ok) {
        const responseData = await response.json();
        const newTableName = responseData.content?.name || data.tableName;
        
        // Invalidate collections query to refresh the sidebar
        await queryClient.invalidateQueries({
          queryKey: ["collections", projectId],
        });
        
        // Redirect to the new table
        navigate(`/projects/${projectId}/collections/${newTableName}`);
      } else {
        const errorData = await response.json().catch(() => ({}));
        throw new Error(errorData.message || "Failed to create table");
      }
    } catch (error) {
      console.error("Error creating table:", error);
      alert("Failed to create table. Please try again.");
    } finally {
      setIsSubmitting(false);
    }
  };



  return (
    <div className="flex flex-col h-full">
      <AppHeader title="Create Collection" />
      
      <div className="flex-1 p-6">
        <div className="max-w-4xl mx-auto">
          <Form {...form}>
              <form onSubmit={form.handleSubmit(onSubmit)} className="space-y-8">
                <Card>
                  <CardHeader>
                    <CardTitle>Create New Collection</CardTitle>
                    <CardDescription>
                      Define your table name and column structure
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
                            />
                          </FormControl>
                          <FormMessage />
                        </FormItem>
                      )}
                    />

                    <div className="space-y-4">
                      <div className="flex items-center justify-between">
                        <Label className="text-base font-medium">Columns</Label>
                        <Button
                          type="button"
                          variant="outline"
                          size="sm"
                          onClick={addColumn}
                        >
                          <Plus className="w-4 h-4 mr-2" />
                          Add Column
                        </Button>
                      </div>

                      <div className="space-y-4">
                        {fields.map((field, index) => (
                          <div key={field.id} className="grid grid-cols-12 gap-4 items-end p-4 border rounded-lg">
                            <div className="col-span-4">
                              <FormField
                                control={form.control}
                                name={`columns.${index}.name`}
                                render={({ field }) => (
                                  <FormItem>
                                    <FormLabel>Column Name</FormLabel>
                                    <FormControl>
                                      <Input placeholder="column_name" {...field} />
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
                                  <FormItem>
                                    <FormLabel>Data Type</FormLabel>
                                    <Select onValueChange={field.onChange} defaultValue={field.value}>
                                      <FormControl>
                                        <SelectTrigger>
                                          <SelectValue placeholder="Select type" />
                                        </SelectTrigger>
                                      </FormControl>
                                      <SelectContent>
                                        {COLUMN_TYPE_OPTIONS.map((option) => (
                                          <SelectItem key={option.value} value={option.value}>
                                            <div>
                                              <div className="font-medium">{option.label}</div>
                                              <div className="text-xs text-muted-foreground">{option.description}</div>
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
                                    <FormLabel>Primary Key</FormLabel>
                                    <FormControl>
                                      <div className="flex items-center space-x-2">
                                        <Checkbox
                                          checked={field.value}
                                          onCheckedChange={field.onChange}
                                        />
                                        <Label className="text-sm">Primary</Label>
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
                                variant="outline"
                                size="sm"
                                onClick={() => removeColumn(index)}
                                disabled={fields.length === 1}
                                className="w-full"
                              >
                                <Trash2 className="w-4 h-4" />
                              </Button>
                            </div>
                          </div>
                        ))}
                      </div>
                    </div>

                    <div className="flex justify-end">
                      <Button type="submit" disabled={isSubmitting}>
                        {isSubmitting ? "Creating..." : "Create Table"}
                      </Button>
                    </div>
                  </CardContent>
                </Card>
              </form>
            </Form>
        </div>
      </div>
    </div>
  );
}