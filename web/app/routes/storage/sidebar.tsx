import { useState, useMemo } from "react";
import { NavLink } from "react-router";
import {
  Sidebar,
  SidebarContent,
  SidebarGroup,
  SidebarGroupContent,
  SidebarHeader,
  SidebarInput,
  SidebarMenu,
  SidebarMenuButton,
  SidebarMenuItem,
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
              <SidebarMenu>
                {filteredContainers.map((container) => (
                  <SidebarMenuItem key={container.uuid}>
                    <NavLink
                      to={`/projects/${projectId}/storage/${container.uuid}`}
                      className={({ isActive }) =>
                        cn(
                          "block w-full",
                          isActive && "bg-accent"
                        )
                      }
                    >
                      <SidebarMenuButton
                        isActive={activeContainerId === container.uuid}
                        className="w-full justify-start"
                      >
                        <Package2 className="h-4 w-4 mr-2" />
                        <div className="flex-1 text-left">
                          <div className="font-medium">{container.name}</div>
                          {container.description && (
                            <div className="text-xs text-muted-foreground truncate">
                              {container.description}
                            </div>
                          )}
                          <div className="text-xs text-muted-foreground">
                            {container.totalFiles} file{container.totalFiles !== 1 ? 's' : ''}
                            {container.isPublic && ' • Public'}
                          </div>
                        </div>
                      </SidebarMenuButton>
                    </NavLink>
                  </SidebarMenuItem>
                ))}
              </SidebarMenu>
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