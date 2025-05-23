import { RefreshCw } from "lucide-react";
import { Button } from "~/components/ui/button";
import { useEffect, useRef, useState } from "react";
import { cn } from "~/lib/utils";

// Add global styles for the rotation animation
// This will be done once when the component is first imported
const addRotationStyles = () => {
  if (
    typeof document !== "undefined" &&
    !document.getElementById("refresh-button-animation-style")
  ) {
    const style = document.createElement("style");
    style.id = "refresh-button-animation-style";
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

export interface RefreshButtonProps
  extends React.ButtonHTMLAttributes<HTMLButtonElement> {
  onRefresh: () => Promise<any> | void;
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

export function RefreshButton({
  onRefresh,
  title = "Refresh",
  icon,
  size = "icon",
  variant = "ghost",
  className,
  ...props
}: RefreshButtonProps) {
  const refreshIconRef = useRef<HTMLDivElement>(null);
  const [isRefreshing, setIsRefreshing] = useState(false);

  // Ensure styles are added when component is mounted
  useEffect(() => {
    addRotationStyles();
  }, []);

  const handleRefresh = async () => {
    if (isRefreshing) return;

    // Start animation
    setIsRefreshing(true);
    if (refreshIconRef.current) {
      refreshIconRef.current.classList.add("rotate-animation");
    }

    try {
      // Call the refresh function
      await onRefresh();
    } catch (error) {
      console.error("Refresh failed:", error);
    } finally {
      // End animation after 1 second or when refresh is done, whichever is longer
      setTimeout(() => {
        if (refreshIconRef.current) {
          refreshIconRef.current.classList.remove("rotate-animation");
        }
        setIsRefreshing(false);
      }, 200);
    }
  };

  return (
    <Button
      type="button"
      variant={variant}
      size={size}
      onClick={handleRefresh}
      title={title}
      className={cn("h-8 w-8 cursor-pointer", className)}
      disabled={isRefreshing}
      {...props}
    >
      <div ref={refreshIconRef}>
        {icon || <RefreshCw className="h-4 w-4" />}
      </div>
    </Button>
  );
}
