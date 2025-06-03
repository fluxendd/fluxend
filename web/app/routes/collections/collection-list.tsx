import { useQuery, useQueryClient } from "@tanstack/react-query";
import {
  BadgeAlert,
  FolderIcon,
  HashIcon,
  LoaderCircle,
  MessageCircleWarning,
} from "lucide-react";
import { useCallback, useEffect, useLayoutEffect, useMemo } from "react";
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
import { motion } from "motion/react";

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
  <div className="flex items-center p-4 text-muted-foreground">
    <BadgeAlert className="mr-2" />
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
  return <CollectionListSkeleton count={5} />;
}

type CollectionListProps = {
  projectId: string;
  searchTerm?: string;
};

export const CollectionList = ({
  projectId,
  searchTerm = "",
}: CollectionListProps) => {
  const navigate = useNavigate();

  const { isLoading, isFetching, isError, data, error } = useQuery<
    Collection[]
  >(collectionsQuery(projectId));

  useEffect(() => {
    if (error?.name === "UnauthorizedError") {
      navigate(href("/logout"));
    }
  }, [error]);

  const filteredData = useMemo(() => {
    if (!data) {
      return [];
    }

    if (!searchTerm) {
      return data;
    }

    return data.filter((table) =>
      table.name.toLowerCase().includes(searchTerm.toLowerCase())
    );
  }, [searchTerm, data]);

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
          className={({ isActive }) =>
            [
              `relative flex flex-col items-start whitespace-nowrap p-2 rounded-md mx-2  text-sm leading-tight hover:text-primary cursor-pointer ${
                isFetching ? "opacity-60" : ""
              }`,
              isActive ? "text-primary" : "",
            ].join(" ")
          }
        >
          {({ isActive }) => (
            <>
              <div className="flex w-full items-center gap-1">
                {isActive && (
                  <motion.div
                    layoutId="collectionId"
                    className="absolute inset-0 bg-primary/10 rounded-md"
                    transition={{
                      type: "spring",
                      bounce: 0.2,
                      duration: 0.3,
                      delay: 0.1,
                    }}
                  />
                )}
                <HashIcon size={12} />
                <span className="">{table.name}</span>
                <span className="ml-auto text-xs">{table.totalSize}</span>
              </div>
            </>
          )}
        </NavLink>
      ))}
    </div>
  );
};
