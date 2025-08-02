import { FileText, Database } from "lucide-react";
import { useMemo } from "react";
import { href, NavLink } from "react-router";
import { motion } from "motion/react";
import { InfoMessage } from "~/components/shared/info-message";

type DocsListProps = {
  tables: string[];
  projectId: string;
  searchTerm?: string;
};

export const DocsList = ({
  tables,
  projectId,
  searchTerm = "",
}: DocsListProps) => {
  const filteredTables = useMemo(() => {
    if (!searchTerm) {
      return tables;
    }

    return tables.filter((table) =>
      table.toLowerCase().includes(searchTerm.toLowerCase())
    );
  }, [searchTerm, tables]);

  if (tables.length === 0) {
    return <InfoMessage message="No API documentation found" />;
  }

  if (searchTerm && filteredTables.length === 0) {
    return <InfoMessage message="No matching documentation found" />;
  }

  return (
    <div className="h-full overflow-y-auto flex flex-col">
      {/* All Tables Link */}
      <NavLink
        to={href("/projects/:projectId/docs", {
          projectId: projectId,
        })}
        end
        className={({ isActive }) =>
          [
            `relative flex flex-col items-start p-2 rounded-lg mx-2 text-sm leading-tight hover:text-foreground/70 cursor-pointer`,
            isActive ? "dark:text-primary" : "",
          ].join(" ")
        }
      >
        {({ isActive }) => (
          <>
            <div className="flex w-full items-center gap-2">
              {isActive && (
                <motion.div
                  layoutId="docsId"
                  className="absolute inset-0 bg-primary/30 dark:bg-primary/10 rounded-lg"
                  transition={{
                    type: "spring",
                    bounce: 0.2,
                    duration: 0.3,
                    delay: 0.1,
                  }}
                />
              )}
              <FileText size={14} className="flex-shrink-0" />
              <span className="truncate flex-1 min-w-0 font-medium">All Tables</span>
              <span className="ml-auto text-xs flex-shrink-0 text-muted-foreground">
                {tables.length}
              </span>
            </div>
          </>
        )}
      </NavLink>

      <div className="mx-2 my-2 border-t"></div>

      {/* Individual Table Links */}
      {filteredTables.map((table) => (
        <NavLink
          to={href("/projects/:projectId/docs/:table", {
            projectId: projectId,
            table: table,
          })}
          key={table}
          className={({ isActive }) =>
            [
              `relative flex flex-col items-start p-2 rounded-lg mx-2 text-sm leading-tight hover:text-foreground/70 cursor-pointer`,
              isActive ? "dark:text-primary" : "",
            ].join(" ")
          }
        >
          {({ isActive }) => (
            <>
              <div className="flex w-full items-center gap-2">
                {isActive && (
                  <motion.div
                    layoutId="docsId"
                    className="absolute inset-0 bg-primary/30 dark:bg-primary/10 rounded-lg"
                    transition={{
                      type: "spring",
                      bounce: 0.2,
                      duration: 0.3,
                      delay: 0.1,
                    }}
                  />
                )}
                <Database size={14} className="flex-shrink-0" />
                <span className="truncate flex-1 min-w-0">{table}</span>
              </div>
            </>
          )}
        </NavLink>
      ))}
    </div>
  );
};