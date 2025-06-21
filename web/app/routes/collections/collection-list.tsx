import { useQuery } from "@tanstack/react-query";
import { HashIcon } from "lucide-react";
import { useEffect, useMemo } from "react";
import { TableListSkeleton } from "~/components/shared/collection-list-skeleton";
import { href, NavLink, useNavigate, useOutletContext } from "react-router";
import { motion } from "motion/react";
import { InfoMessage } from "~/components/shared/info-message";
import type { ProjectLayoutOutletContext } from "~/components/shared/project-layout";
import type { Services } from "~/services";

export type Table = {
  name: string;
  totalSize: string;
  [key: string]: any;
};

class UnauthorizedError extends Error {
  constructor(message: string) {
    super(message);
    this.name = "UnauthorizedError";
  }
}

const collectionsQuery = (services: Services, projectId: string) => ({
  queryKey: ["collections", projectId],
  queryFn: async () => {
    const { success, errors, content, ok, status } =
      await services.collections.getAllTables(projectId);

    if (!ok) {
      const errorMessage = errors?.[0] || "Unknown error";
      if (status === 401) {
        throw new UnauthorizedError(errorMessage);
      } else {
        throw new Error(errorMessage);
      }
    }

    return content;
  },
});

function TableListFallback() {
  return <TableListSkeleton count={5} />;
}

type TableListProps = {
  initialData: Table[];
  projectId: string;
  searchTerm?: string;
};

export const TableList = ({
  initialData,
  projectId,
  searchTerm = "",
}: TableListProps) => {
  const navigate = useNavigate();
  const { services } = useOutletContext<ProjectLayoutOutletContext>();

  const { isLoading, isFetching, isError, data, error } = useQuery<Table[]>({
    initialData: initialData,
    ...collectionsQuery(services, projectId),
  });

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
    return <TableListFallback />;
  }

  if (isError) {
    return (
      <InfoMessage message={error.message || "Failed to load collections"} />
    );
  }

  if (!data || data.length === 0) {
    return <InfoMessage message="No collections found" />;
  }

  if (data.length === 0) {
    return <InfoMessage message="No matching collections found" />;
  }

  return (
    <div className="h-full overflow-y-auto flex flex-col">
      {filteredData.map((table) => (
        <NavLink
          to={href("/projects/:projectId/collections/:collectionId", {
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
                    className="absolute inset-0 bg-primary/10 rounded-lg"
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
