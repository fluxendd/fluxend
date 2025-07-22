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
  FormField,
  FormItem,
  FormLabel,
  FormMessage,
} from "~/components/ui/form";
import { Input } from "~/components/ui/input";
import { Button } from "~/components/ui/button";
import type { StorageFile } from "~/types/storage";

const renameFileSchema = z.object({
  name: z.string().min(1, "File name is required"),
});

type RenameFileFormData = z.infer<typeof renameFileSchema>;

interface RenameFileDialogProps {
  open: boolean;
  onOpenChange: (open: boolean) => void;
  file: StorageFile;
  onSubmit: (newName: string) => Promise<void>;
}

export function RenameFileDialog({
  open,
  onOpenChange,
  file,
  onSubmit,
}: RenameFileDialogProps) {
  const [isSubmitting, setIsSubmitting] = useState(false);

  const form = useForm<RenameFileFormData>({
    resolver: zodResolver(renameFileSchema),
    defaultValues: {
      name: file.fullFileName,
    },
  });

  // Reset form when file changes
  useEffect(() => {
    form.reset({
      name: file.fullFileName,
    });
  }, [file, form]);

  const handleSubmit = async (data: RenameFileFormData) => {
    if (data.name === file.fullFileName) {
      onOpenChange(false);
      return;
    }

    setIsSubmitting(true);
    try {
      await onSubmit(data.name);
    } finally {
      setIsSubmitting(false);
    }
  };

  return (
    <Dialog open={open} onOpenChange={onOpenChange}>
      <DialogContent className="sm:max-w-[425px]">
        <DialogHeader>
          <DialogTitle>Rename File</DialogTitle>
          <DialogDescription>
            Enter a new name for the file "{file.fullFileName}".
          </DialogDescription>
        </DialogHeader>
        <Form {...form}>
          <form onSubmit={form.handleSubmit(handleSubmit)} className="space-y-4">
            <FormField
              control={form.control}
              name="name"
              render={({ field }) => (
                <FormItem>
                  <FormLabel>File Name</FormLabel>
                  <FormControl>
                    <Input
                      placeholder="filename.txt"
                      {...field}
                      disabled={isSubmitting}
                      autoFocus
                    />
                  </FormControl>
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
                {isSubmitting ? "Renaming..." : "Rename"}
              </Button>
            </DialogFooter>
          </form>
        </Form>
      </DialogContent>
    </Dialog>
  );
}