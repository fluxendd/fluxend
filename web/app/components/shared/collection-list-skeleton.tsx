import { Skeleton } from "~/components/ui/skeleton";
import { FolderIcon } from "lucide-react";

interface CollectionListSkeletonProps {
  count?: number;
}

export function CollectionListSkeleton({ count = 5 }: CollectionListSkeletonProps) {
  return (
    <>
      {Array.from({ length: count }).map((_, index) => (
        <div
          key={`collection-skeleton-${index}`}
          className="flex flex-col items-start gap-2 whitespace-nowrap border-b p-4 text-sm leading-tight last:border-b-0"
        >
          <div className="flex w-full items-center gap-2">
            <FolderIcon size={20} className="text-muted-foreground" />
            <Skeleton className="h-4 w-24" />
            <Skeleton className="ml-auto h-3 w-8" />
          </div>
        </div>
      ))}
    </>
  );
}