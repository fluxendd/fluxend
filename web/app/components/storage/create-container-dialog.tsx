import { useState } from "react";
import { useForm } from "react-hook-form";
import { zodResolver } from "@hookform/resolvers/zod";
import { z } from "zod";
import {
  Dialog,
  DialogContent,
  DialogDescription,
  DialogFooter,
  DialogHeader,
  DialogTitle,
} from "~/components/ui/dialog";
import {
  Form,
  FormControl,
  FormDescription,
  FormField,
  FormItem,
  FormLabel,
  FormMessage,
} from "~/components/ui/form";
import { Input } from "~/components/ui/input";
import { Textarea } from "~/components/ui/textarea";
import { Button } from "~/components/ui/button";
import { Switch } from "~/components/ui/switch";
import { bytesToMB, mbToBytes } from "~/lib/utils";

const createContainerSchema = z.object({
  name: z.string().min(1, "Container name is required"),
  description: z.string().optional(),
  is_public: z.boolean(),
  max_file_size: z.number().min(0.001, "Max file size must be at least 0.001 MB"),
});

type CreateContainerFormData = z.infer<typeof createContainerSchema>;

interface CreateContainerDialogProps {
  open: boolean;
  onOpenChange: (open: boolean) => void;
  onSubmit: (data: CreateContainerFormData) => Promise<void>;
}

export function CreateContainerDialog({
  open,
  onOpenChange,
  onSubmit,
}: CreateContainerDialogProps) {
  const [isSubmitting, setIsSubmitting] = useState(false);

  const form = useForm<CreateContainerFormData>({
    resolver: zodResolver(createContainerSchema),
    defaultValues: {
      name: "",
      description: "",
      is_public: false,
      max_file_size: 10, // 10MB default
    },
  });

  const handleSubmit = async (data: CreateContainerFormData) => {
    setIsSubmitting(true);
    try {
      // Convert MB to bytes before submitting
      await onSubmit({
        ...data,
        max_file_size: mbToBytes(data.max_file_size),
      });
      // Only reset form on successful submission
      form.reset();
    } catch (error) {
      // Keep form data on error to allow user to modify
      console.error("Failed to create container:", error);
    } finally {
      setIsSubmitting(false);
    }
  };

  return (
    <Dialog open={open} onOpenChange={onOpenChange}>
      <DialogContent className="sm:max-w-[425px]">
        <DialogHeader>
          <DialogTitle>Create Storage Container</DialogTitle>
          <DialogDescription>
            Create a new storage container to organize your files.
          </DialogDescription>
        </DialogHeader>
        <Form {...form}>
          <form onSubmit={form.handleSubmit(handleSubmit)} className="space-y-4">
            <FormField
              control={form.control}
              name="name"
              render={({ field }) => (
                <FormItem>
                  <FormLabel>Container Name</FormLabel>
                  <FormControl>
                    <Input
                      placeholder="my-container"
                      {...field}
                      disabled={isSubmitting}
                    />
                  </FormControl>
                  <FormMessage />
                </FormItem>
              )}
            />

            <FormField
              control={form.control}
              name="description"
              render={({ field }) => (
                <FormItem>
                  <FormLabel>Description</FormLabel>
                  <FormControl>
                    <Textarea
                      placeholder="Container for storing project assets..."
                      className="resize-none"
                      {...field}
                      disabled={isSubmitting}
                    />
                  </FormControl>
                  <FormMessage />
                </FormItem>
              )}
            />

            <FormField
              control={form.control}
              name="is_public"
              render={({ field }) => (
                <FormItem className="flex flex-row items-center justify-between rounded-lg border p-3">
                  <div className="space-y-0.5">
                    <FormLabel>Public Container</FormLabel>
                    <FormDescription>
                      Allow public access to files in this container
                    </FormDescription>
                  </div>
                  <FormControl>
                    <Switch
                      checked={field.value}
                      onCheckedChange={field.onChange}
                      disabled={isSubmitting}
                    />
                  </FormControl>
                </FormItem>
              )}
            />

            <FormField
              control={form.control}
              name="max_file_size"
              render={({ field }) => (
                <FormItem>
                  <FormLabel>Max File Size (MB)</FormLabel>
                  <FormControl>
                    <Input
                      type="number"
                      step="0.1"
                      {...field}
                      onChange={(e) => field.onChange(parseFloat(e.target.value))}
                      disabled={isSubmitting}
                    />
                  </FormControl>
                  <FormDescription>
                    Maximum size allowed for individual files in megabytes
                  </FormDescription>
                  <FormMessage />
                </FormItem>
              )}
            />

            <DialogFooter>
              <Button
                type="button"
                variant="outline"
                onClick={() => onOpenChange(false)}
                disabled={isSubmitting}
              >
                Cancel
              </Button>
              <Button type="submit" disabled={isSubmitting}>
                {isSubmitting ? "Creating..." : "Create Container"}
              </Button>
            </DialogFooter>
          </form>
        </Form>
      </DialogContent>
    </Dialog>
  );
}