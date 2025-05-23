import { useQuery, useQueryClient } from "@tanstack/react-query";
import { FolderIcon, LoaderCircle, MessageCircleWarning } from "lucide-react";
import { useEffect, useLayoutEffect } from "react";
import { CollectionListSkeleton } from "~/components/shared/collection-list-skeleton";
import {
  href,
  NavLink,
  redirect,
  useNavigate,
  useNavigation,
  useParams,
} from "react-router";
import { getAllCollections } from "~/services/collections";

// Define collection type
interface Collection {
  name: string;
  totalSize: string;
  [key: string]: any;
}

class UnauthorizedError extends Error {
  constructor(message: string) {
    super(message);
    this.name = "UnauthorizedError";
  }
}

const InfoMessage = ({ message }: { message: string }) => (
  <div className="flex items-center justify-center p-4 text-muted-foreground h-full">
    <MessageCircleWarning className="mr-2" />
    <div className="text-md">{message}</div>
  </div>
);

const collectionsQuery = (projectId: string) => ({
  queryKey: ["collections", projectId],
  queryFn: async () => {
    const res = await getAllCollections({ headers: {} }, projectId);

    if (!res.ok) {
      const responseData = await res.json();
      const errorMessage = responseData?.errors[0] || "Unknown error";
      if (res.status === 401) {
        throw new UnauthorizedError(errorMessage);
      } else {
        throw new Error(errorMessage);
      }
    }

    const data = await res.json();
    return data.content;
  },
});

function CollectionListFallback() {
  return <CollectionListSkeleton count={8} />;
}

type CollectionListProps = {
  projectId: string;
  searchTerm?: string;
};

export const CollectionList = ({
  projectId,
  searchTerm = "",
}: CollectionListProps) => {
  const queryClient = useQueryClient();
  const navigate = useNavigate();
  const { isLoading, isFetching, isError, data, error } = useQuery<Collection[]>(
    collectionsQuery(projectId)
  );

  const { collectionId } = useParams();

  useLayoutEffect(() => {
    if (error?.name === "UnauthorizedError") {
      navigate(href("/logout"));
    }
  }, [error]);

  // Navigate to the first collection when no collection is selected
  useEffect(() => {
    if (!collectionId && !isLoading && data && data.length > 0) {
      navigate(
        `/projects/${projectId}/collections/${data[0].name}`,
        { replace: true }
      );
    }
  }, [collectionId, isLoading, data, projectId, navigate]);

  if (isLoading) {
    return <CollectionListFallback />;
  }

  if (isError) {
    return (
      <InfoMessage message={error.message || "Failed to load collections"} />
    );
  }

  if (!data || data.length === 0) {
    return <InfoMessage message="No collections found" />;
  }

  const filteredData = searchTerm
    ? data.filter((table) =>
        table.name.toLowerCase().includes(searchTerm.toLowerCase())
      )
    : data;

  if (filteredData.length === 0) {
    return <InfoMessage message="No matching collections found" />;
  }

  return (
    <div className="h-full overflow-y-auto flex flex-col">
      {filteredData.map((table) => (
    <NavLink
      to={href(`/projects/:projectId/collections/:collectionId?`, {
        projectId: projectId,
        collectionId: table.name,
      })}
      key={table.name}
      className={({ isActive, isPending, isTransitioning }) =>
        [
          `flex flex-col items-start gap-2 whitespace-nowrap border-b p-4 text-sm leading-tight last:border-b-0 hover:bg-sidebar-accent hover:text-sidebar-accent-foreground ${
            isFetching ? "opacity-60" : ""
          }`,
          isActive ? "bg-sidebar-accent text-sidebar-accent-foreground" : "",
        ].join(" ")
      }
    >
      <div className="flex w-full items-center gap-2">
        <FolderIcon size={20} />
        <span className="font-medium">{table.name}</span>
        <span className="ml-auto text-xs">{table.totalSize}</span>
      </div>
    </NavLink>
      ))}
    </div>
  );
};
