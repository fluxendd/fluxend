import { Button } from "~/components/ui/button";
import { Trash2 } from "lucide-react";

type AppHeaderProps = {
  title: string;
  onDelete?: () => void;
  deleteLabel?: string;
  showDelete?: boolean;
};

export function AppHeader({ 
  title, 
  onDelete, 
  deleteLabel = "Delete", 
  showDelete = false 
}: AppHeaderProps) {
  if (!title) {
    return null;
  }

  return (
    <header className="group-has-data-[collapsible=icon]/sidebar-wrapper:h-12 flex h-12 shrink-0 items-center gap-2 border-b transition-[width,height] ease-linear px-4">
      <h1 className="text-base font-medium flex-1">{title}</h1>
      {showDelete && onDelete && (
        <Button
          variant="destructive"
          size="sm"
          onClick={onDelete}
          className="ml-auto"
        >
          <Trash2 className="w-4 h-4 mr-2" />
          {deleteLabel}
        </Button>
      )}
    </header>
  );
}
