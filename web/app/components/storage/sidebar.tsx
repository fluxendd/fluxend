import { useState, useMemo } from "react";
import { NavLink } from "react-router";
import {
  Sidebar,
  SidebarContent,
  SidebarGroup,
  SidebarGroupContent,
  SidebarHeader,
  SidebarInput,
} from "~/components/ui/sidebar";
import { Package2, Plus } from "lucide-react";
import { TableListSkeleton } from "~/components/shared/collection-list-skeleton";
import { Button } from "~/components/ui/button";
import type { StorageContainer } from "~/types/storage";
import { cn } from "~/lib/utils";

interface StorageSidebarProps {
  containers: StorageContainer[];
  activeContainerId?: string;
  isLoading: boolean;
  projectId: string;
  onCreateContainer: () => void;
}

export function StorageSidebar({
  containers,
  activeContainerId,
  isLoading,
  projectId,
  onCreateContainer,
}: StorageSidebarProps) {
  const [searchValue, setSearchValue] = useState("");

  const filteredContainers = useMemo(() => {
    if (!searchValue) return containers;
    
    const searchLower = searchValue.toLowerCase();
    return containers.filter(
      (container) =>
        container.name.toLowerCase().includes(searchLower) ||
        container.description?.toLowerCase().includes(searchLower)
    );
  }, [containers, searchValue]);

  return (
    <Sidebar
      collapsible="none"
      className="hidden md:flex h-full flex-shrink-0 border-r"
      variant="inset"
    >
      <SidebarHeader className="gap-3 border-b p-2 mb-2 flex-shrink-0">
        <div className="flex items-center gap-2">
          <SidebarInput
            placeholder="Search containers..."
            value={searchValue}
            onChange={(e) => setSearchValue(e.target.value)}
            className="flex-1 rounded-lg"
          />
        </div>
      </SidebarHeader>
      <SidebarContent className="flex-1 min-h-0 flex flex-col">
        <SidebarGroup className="p-0 flex-1 overflow-hidden">
          <SidebarGroupContent className="h-full overflow-y-auto">
            {isLoading ? (
              <TableListSkeleton count={5} />
            ) : filteredContainers.length === 0 ? (
              <div className="px-4 py-8 text-center text-muted-foreground text-sm">
                {searchValue
                  ? "No containers found matching your search"
                  : "No containers yet"}
              </div>
            ) : (
              <div className="h-full overflow-y-auto flex flex-col">
                {filteredContainers.map((container) => (
                  <NavLink
                    to={`/projects/${projectId}/storage/${container.uuid}`}
                    key={container.uuid}
                    className={({ isActive }) =>
                      cn(
                        "relative flex flex-col items-start p-2 rounded-lg mx-2 text-sm leading-tight hover:text-foreground/70 cursor-pointer",
                        isActive && "dark:text-primary"
                      )
                    }
                  >
                    {({ isActive }) => (
                      <>
                        <div className="flex w-full items-center gap-1">
                          {isActive && (
                            <div className="absolute inset-0 bg-primary/30 dark:bg-primary/10 rounded-lg" />
                          )}
                          <Package2 size={12} className="flex-shrink-0" />
                          <span className="truncate flex-1 min-w-0">{container.name}</span>
                          <span className="ml-auto text-xs flex-shrink-0">
                            {container.totalFiles} {container.totalFiles === 1 ? 'file' : 'files'}
                          </span>
                        </div>
                      </>
                    )}
                  </NavLink>
                ))}
              </div>
            )}
          </SidebarGroupContent>
        </SidebarGroup>
        <div className="p-4 border-t flex-shrink-0">
          <Button
            className="w-full"
            size="sm"
            onClick={onCreateContainer}
          >
            <Plus className="mr-1 size-4" />
            Create Container
          </Button>
        </div>
      </SidebarContent>
    </Sidebar>
  );
}