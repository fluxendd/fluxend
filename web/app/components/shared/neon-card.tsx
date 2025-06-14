import { useState, type ReactNode } from "react";
import { cn } from "~/lib/utils";
import { Card } from "../ui/card";
// NeonCard component with flowing neon effect
interface NeonCardProps {
  children: ReactNode;
  className?: string;
  [key: string]: any;
}

// NeonCard component with flowing muted border trail
export function NeonCard({ children, className, ...props }: NeonCardProps) {
  const [hoverState, setHoverState] = useState(false);

  return (
    <div
      className={cn("relative group", className)}
      onMouseEnter={() => setHoverState(true)}
      onMouseLeave={() => setHoverState(false)}
      {...props}
    >
      {/* Card container with border */}
      <div className="card-with-trail rounded-xl">
        {/* Animated muted trail that follows the border */}
        <div className="trail" />

        {/* Inner card with content */}
        <Card className="relative bg-card rounded-xl overflow-clip border-transparent shadow-sm">
          {children}
        </Card>
      </div>
    </div>
  );
}
