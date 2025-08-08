import { useState, useEffect } from "react";
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
import type { StorageContainer } from "~/types/storage";
import { bytesToMB, mbToBytes } from "~/lib/utils";

const updateContainerSchema = z.object({
  name: z.string().min(1, "Container name is required"),
  description: z.string().optional(),
  is_public: z.boolean(),
  max_file_size: z.number().min(0.001, "Max file size must be at least 0.001 MB"),
});

type UpdateContainerFormData = z.infer<typeof updateContainerSchema>;

interface UpdateContainerDialogProps {
  open: boolean;
  onOpenChange: (open: boolean) => void;
  container: StorageContainer;
  onSubmit: (data: UpdateContainerFormData) => Promise<void>;
}

export function UpdateContainerDialog({
  open,
  onOpenChange,
  container,
  onSubmit,
}: UpdateContainerDialogProps) {
  const [isSubmitting, setIsSubmitting] = useState(false);

  const form = useForm<UpdateContainerFormData>({
    resolver: zodResolver(updateContainerSchema),
    defaultValues: {
      name: container.name,
      description: container.description || "",
      is_public: container.isPublic,
      max_file_size: bytesToMB(container.maxFileSize),
    },
  });

  // Reset form when container changes
  useEffect(() => {
    form.reset({
      name: container.name,
      description: container.description || "",
      is_public: container.isPublic,
      max_file_size: bytesToMB(container.maxFileSize),
    });
  }, [container, form]);

  const handleSubmit = async (data: UpdateContainerFormData) => {
    setIsSubmitting(true);
    try {
      // Convert MB to bytes before submitting
      await onSubmit({
        ...data,
        max_file_size: mbToBytes(data.max_file_size),
      });
    } finally {
      setIsSubmitting(false);
    }
  };

  return (
    <Dialog open={open} onOpenChange={onOpenChange}>
      <DialogContent className="sm:max-w-[425px]">
        <DialogHeader>
          <DialogTitle>Update Container</DialogTitle>
          <DialogDescription>
            Update the settings for your storage container.
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
                {isSubmitting ? "Updating..." : "Update Container"}
              </Button>
            </DialogFooter>
          </form>
        </Form>
      </DialogContent>
    </Dialog>
  );
}