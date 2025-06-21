import { Skeleton } from "~/components/ui/skeleton";
import { FolderIcon, HashIcon } from "lucide-react";

interface TableListSkeletonProps {
  count?: number;
}

export function TableListSkeleton({ count = 5 }: TableListSkeletonProps) {
  return (
    <>
      {Array.from({ length: count }).map((_, index) => (
        <div
          key={`collection-skeleton-${index}`}
          className="flex flex-col items-start gap-8 whitespace-nowrap p-2 mx-2  text-sm leading-tight"
        >
          <div className="flex w-full items-center gap-1">
            <HashIcon size={12} className="text-muted-foreground" />
            <Skeleton className="h-4 w-24" />
            <Skeleton className="ml-auto h-4 w-8" />
          </div>
        </div>
      ))}
    </>
  );
}
