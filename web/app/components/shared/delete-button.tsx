import { Trash2 } from "lucide-react";
import { Button } from "~/components/ui/button";
import { useEffect, useRef, useState } from "react";
import { cn } from "~/lib/utils";

// Add global styles for the rotation animation
// This will be done once when the component is first imported
const addRotationStyles = () => {
  if (
    typeof document !== "undefined" &&
    !document.getElementById("delete-button-animation-style")
  ) {
    const style = document.createElement("style");
    style.id = "delete-button-animation-style";
    style.textContent = `
      @keyframes rotate {
        from { transform: rotate(0deg); }
        to { transform: rotate(360deg); }
      }
      .rotate-animation {
        animation: rotate 0.2s ease-in-out;
      }
    `;
    document.head.appendChild(style);
  }
};

// Add styles when the component is first imported
if (typeof document !== "undefined") {
  addRotationStyles();
}

export interface DeleteButtonProps
  extends React.ButtonHTMLAttributes<HTMLButtonElement> {
  onDelete: () => Promise<any> | void;
  title?: string;
  icon?: React.ReactNode;
  size?: "default" | "sm" | "lg" | "icon";
  variant?:
    | "default"
    | "destructive"
    | "outline"
    | "secondary"
    | "ghost"
    | "link";
}

export function DeleteButton({
  onDelete,
  title = "Delete",
  icon,
  size = "icon",
  variant = "ghost",
  className,
  ...props
}: DeleteButtonProps) {
  const deleteIconRef = useRef<HTMLDivElement>(null);
  const [isDeleting, setIsDeleting] = useState(false);

  // Ensure styles are added when component is mounted
  useEffect(() => {
    addRotationStyles();
  }, []);

  const handleDelete = async () => {
    if (isDeleting) return;

    // Start animation
    setIsDeleting(true);
    if (deleteIconRef.current) {
      deleteIconRef.current.classList.add("rotate-animation");
    }

    try {
      // Call the delete function
      await onDelete();
    } catch (error) {
      console.error("Delete failed:", error);
    } finally {
      // End animation after 1 second or when delete is done, whichever is longer
      setTimeout(() => {
        if (deleteIconRef.current) {
          deleteIconRef.current.classList.remove("rotate-animation");
        }
        setIsDeleting(false);
      }, 200);
    }
  };

  return (
    <Button
      type="button"
      variant={variant}
      size={size}
      onClick={handleDelete}
      title={title}
      className={cn("h-8 w-8 cursor-pointer", className)}
      disabled={isDeleting}
      {...props}
    >
      <div ref={deleteIconRef}>
        {icon || <Trash2 className="h-4 w-4" />}
      </div>
    </Button>
  );
}