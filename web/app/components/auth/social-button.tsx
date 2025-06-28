import { cn } from "~/lib/utils";
import React from "react";
import { Button } from "../ui/button";

interface SocialButtonProps {
  icon: React.ReactNode;
  onClick: () => void;
  children?: React.ReactNode;
  className?: string;
}

const SocialButton = ({
  icon,
  onClick,
  children,
  className,
}: SocialButtonProps) => {
  return (
    <Button
      type="button"
      onClick={onClick}
      size="lg"
      className={cn(
        "w-full flex justify-center py-2 px-4 border border-gray-300",
        "rounded-lg shadow-sm bg-white text-sm font-medium text-gray-700",
        "hover:bg-gray-50 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-indigo-500"
      )}
    >
      {icon}
      {children}
    </Button>
  );
};

export default SocialButton;
