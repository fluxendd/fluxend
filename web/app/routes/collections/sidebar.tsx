import type { Route } from "./+types/sidebar";
import { data, Outlet, useNavigate } from "react-router";
import { useState } from "react";
import {
  Sidebar,
  SidebarContent,
  SidebarGroup,
  SidebarGroupContent,
  SidebarHeader,
  SidebarInput,
  SidebarInset,
  SidebarProvider,
} from "~/components/ui/sidebar";
import { CollectionList } from "./collection-list";
import { PlusCircle } from "lucide-react";
import { CollectionListSkeleton } from "~/components/shared/collection-list-skeleton";
import { Button } from "~/components/ui/button";

export function HydrateFallback() {
  return (
    <SidebarProvider>
      <Sidebar
        collapsible="none"
        className="hidden md:flex h-screen"
        variant="inset"
      >
        <SidebarHeader className="gap-3 border-b p-2">
          <div className="flex items-center gap-2">
            <SidebarInput
              placeholder="Type to search..."
              disabled
              className="flex-1"
            />
          </div>
        </SidebarHeader>
        <SidebarContent>
          <SidebarGroup className="px-0">
            <SidebarGroupContent>
              <CollectionListSkeleton count={8} />{" "}
            </SidebarGroupContent>
          </SidebarGroup>
          <div className="mt-auto p-4 border-t">
            <CreateCollectionButton disabled={true} />
          </div>
        </SidebarContent>
      </Sidebar>
      <SidebarInset className="overflow-hidden flex flex-col">
        <div className="p-4 flex-1 overflow-auto">
          <Outlet />
        </div>
      </SidebarInset>
    </SidebarProvider>
  );
}

const CreateCollectionButton = ({
  disabled = false,
}: {
  disabled?: boolean;
}) => {
  const navigate = useNavigate();

  const handleClick = () => {
    navigate("create");
  };

  return (
    <Button
      className="w-full"
      size="sm"
      disabled={disabled}
      onClick={handleClick}
    >
      <PlusCircle className="mr-2 size-4" />
      Create Collection
    </Button>
  );
};

export default function CollectionSidebar({ params }: Route.ComponentProps) {
  const { projectId } = params;
  const [searchTerm, setSearchTerm] = useState("");

  const handleSearch = (e: React.ChangeEvent<HTMLInputElement>) => {
    setSearchTerm(e.target.value);
  };

  return (
    <SidebarProvider>
      <div className="flex h-screen w-full overflow-hidden">
        <Sidebar
          collapsible="none"
          className="hidden md:flex h-full flex-shrink-0"
          variant="inset"
        >
          <SidebarHeader className="gap-3 border-b p-2 mb-2 flex-shrink-0">
            <div className="flex items-center gap-2">
              <SidebarInput
                placeholder="Type to search..."
                value={searchTerm}
                onChange={handleSearch}
                className="flex-1"
              />
            </div>
          </SidebarHeader>
          <SidebarContent className="flex-1 min-h-0 flex flex-col">
            <SidebarGroup className="p-0 flex-1 overflow-hidden">
              <SidebarGroupContent className="h-full overflow-y-auto">
                <CollectionList projectId={projectId} searchTerm={searchTerm} />
              </SidebarGroupContent>
            </SidebarGroup>
            <div className="p-4 border-t flex-shrink-0">
              <CreateCollectionButton />
            </div>
          </SidebarContent>
        </Sidebar>
        <SidebarInset className="flex-1 overflow-hidden">
          <div className="h-full overflow-auto">
            <Outlet />
          </div>
        </SidebarInset>
      </div>
    </SidebarProvider>
  );
}
