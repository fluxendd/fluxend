import type { Route } from "./+types/sidebar";
import { data, Outlet } from "react-router";
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
import {
  LoaderCircle,
  PlusCircle,
  Type,
  Hash,
  Calendar,
  CheckSquare,
  Trash2,
  Plus,
  ChevronDown,
  ChevronUp,
  Edit,
  KeyRound,
  Lock,
} from "lucide-react";
import { useQueryClient } from "@tanstack/react-query";
import { RefreshButton } from "~/components/shared/refresh-button";
import { CollectionListSkeleton } from "~/components/shared/collection-list-skeleton";
import { Button } from "~/components/ui/button";
import { cn } from "~/lib/utils";
import {
  Sheet,
  SheetContent,
  SheetDescription,
  SheetFooter,
  SheetHeader,
  SheetTitle,
  SheetTrigger,
} from "~/components/ui/sheet";
import { Input } from "~/components/ui/input";
import { Label } from "~/components/ui/label";

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
            <RefreshButton onRefresh={() => {}} disabled />
          </div>
        </SidebarHeader>
        <SidebarContent>
          <SidebarGroup className="px-0">
            <SidebarGroupContent>
              <CollectionListFallback />
            </SidebarGroupContent>
          </SidebarGroup>
          <div className="mt-auto p-4 border-t">
            <Button disabled className="w-full" size="sm">
              <PlusCircle className="mr-2 size-4" />
              Create Collection
            </Button>
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

function CollectionListFallback() {
  return <CollectionListSkeleton count={8} />;
}

function AddCollectionButton({
  projectId,
  disabled = false,
}: {
  projectId?: string;
  disabled?: boolean;
}) {
  const [isEditing, setIsEditing] = useState(false);
  const [collectionName, setCollectionName] = useState("");
  const [columns, setColumns] = useState<
    Array<{
      name: string;
      type: string;
      isEditing?: boolean;
      isDefault?: boolean;
    }>
  >([
    { name: "id", type: "string", isDefault: true },
    { name: "", type: "string" },
  ]);
  const [showTypeSelector, setShowTypeSelector] = useState<number | null>(null);

  const handleAddColumn = () => {
    setColumns([...columns, { name: "", type: "string" }]);
  };

  const handleColumnNameChange = (index: number, value: string) => {
    const newColumns = [...columns];
    newColumns[index].name = value;
    setColumns(newColumns);
  };

  const handleColumnTypeChange = (index: number, type: string) => {
    const newColumns = [...columns];
    newColumns[index].type = type;
    setColumns(newColumns);
    setShowTypeSelector(null);
  };

  const handleRemoveColumn = (index: number) => {
    // Don't allow removing default fields
    if (columns[index].isDefault) return;

    const newColumns = [...columns];
    newColumns.splice(index, 1);
    setColumns(newColumns);
  };

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();

    // Here you would implement the actual API call to create a collection
    console.log("Creating collection:", {
      name: collectionName,
      projectId,
      columns,
    });

    // Reset form
    setCollectionName("");
    setColumns([
      { name: "id", type: "string", isDefault: true },
      { name: "", type: "string" },
    ]);
    setIsEditing(false);
  };

  const getTypeIcon = (type: string, isDefault: boolean = false) => {
    if (isDefault && type === "string" && columns[0]?.name === "id") {
      return <KeyRound className="size-4" />;
    }

    switch (type) {
      case "string":
        return <Type className="size-4" />;
      case "number":
        return <Hash className="size-4" />;
      case "boolean":
        return <CheckSquare className="size-4" />;
      case "date":
        return <Calendar className="size-4" />;
      default:
        return <Type className="size-4" />;
    }
  };

  const getTypeName = (type: string) => {
    switch (type) {
      case "string":
        return "Text";
      case "number":
        return "Number";
      case "boolean":
        return "Boolean";
      case "date":
        return "Date";
      default:
        return type;
    }
  };

  return (
    <Sheet open={isEditing} onOpenChange={setIsEditing}>
      <SheetTrigger asChild>
        <Button className="w-full" size="sm" disabled={disabled}>
          <PlusCircle className="mr-2 size-4" />
          Create Collection
        </Button>
      </SheetTrigger>
      <SheetContent className="sm:max-w-md md:max-w-lg">
        <SheetHeader>
          <SheetTitle>Add New Collection</SheetTitle>
          <SheetDescription>
            Create a new collection and define its columns.
          </SheetDescription>
        </SheetHeader>

        <form onSubmit={handleSubmit} className="space-y-6 py-6 px-4">
          <div className="space-y-3">
            <Label htmlFor="collection-name" className="text-sm font-medium">
              Collection Name
            </Label>
            <Input
              id="collection-name"
              value={collectionName}
              onChange={(e) => setCollectionName(e.target.value)}
              placeholder="Enter collection name"
              className="w-full"
              required
            />
          </div>

          <div className="space-y-4">
            <div className="flex items-center justify-between">
              <Label className="text-sm font-medium">Collection Fields</Label>
            </div>

            <div className="border rounded-md overflow-hidden max-h-[350px] overflow-y-auto">
              <div className="bg-secondary/30 border-b px-3 py-2 text-xs font-medium grid grid-cols-[1fr,1fr,auto] gap-2">
                <div>Field Name</div>
                <div>Type</div>
                <div>Actions</div>
              </div>

              {columns.map((column, index) => (
                <div
                  key={index}
                  className={cn(
                    "border-b last:border-0",
                    column.isEditing ? "bg-secondary/10" : "",
                    column.isDefault ? "bg-secondary/5" : ""
                  )}
                >
                  {/* Collapsed row view (table-like) */}
                  <div
                    className={cn(
                      "grid grid-cols-[1fr,1fr,auto] gap-2 px-3 py-3 items-center",
                      column.isEditing ? "border-b" : ""
                    )}
                  >
                    <div className="flex items-center gap-2">
                      {column.isDefault && (
                        <Lock className="size-3.5 text-muted-foreground" />
                      )}
                      <span
                        className={cn(
                          "truncate",
                          column.isDefault ? "font-medium" : "",
                          !column.name ? "text-muted-foreground italic" : ""
                        )}
                      >
                        {column.name || "Unnamed field"}
                      </span>
                    </div>

                    <div className="flex items-center gap-1.5">
                      {getTypeIcon(column.type, column.isDefault)}
                      <span className="text-sm">
                        {getTypeName(column.type)}
                      </span>
                    </div>

                    <div className="flex items-center gap-1">
                      <Button
                        type="button"
                        variant="ghost"
                        size="icon"
                        className="h-7 w-7"
                        onClick={() => {
                          const newColumns = [...columns];
                          newColumns[index].isEditing =
                            !newColumns[index].isEditing;
                          setColumns(newColumns);
                        }}
                        title={column.isEditing ? "Collapse" : "Edit field"}
                      >
                        {column.isEditing ? (
                          <ChevronUp className="size-4" />
                        ) : (
                          <Edit className="size-4" />
                        )}
                      </Button>

                      {!column.isDefault && (
                        <Button
                          type="button"
                          variant="ghost"
                          size="icon"
                          className="h-7 w-7 text-muted-foreground hover:text-destructive"
                          onClick={() => handleRemoveColumn(index)}
                          disabled={columns.length <= 2}
                          title="Delete field"
                        >
                          <Trash2 className="size-4" />
                        </Button>
                      )}
                    </div>
                  </div>

                  {/* Expanded edit view */}
                  {column.isEditing && (
                    <div className="px-3 pb-3 space-y-3">
                      <div className="space-y-1.5">
                        <Label
                          htmlFor={`column-name-${index}`}
                          className="text-xs"
                        >
                          Field Name
                        </Label>
                        <Input
                          id={`column-name-${index}`}
                          value={column.name}
                          onChange={(e) =>
                            handleColumnNameChange(index, e.target.value)
                          }
                          placeholder="Field name"
                          disabled={column.isDefault}
                          required
                        />
                      </div>

                      <div className="space-y-1.5">
                        <Label className="text-xs">Field Type</Label>
                        <div className="relative">
                          <Button
                            type="button"
                            variant="outline"
                            className="w-full justify-start text-left font-normal"
                            onClick={() =>
                              setShowTypeSelector(
                                showTypeSelector === index ? null : index
                              )
                            }
                            disabled={column.isDefault}
                          >
                            <span className="flex items-center gap-2">
                              {getTypeIcon(column.type, column.isDefault)}
                              <span>{getTypeName(column.type)}</span>
                            </span>
                          </Button>

                          {showTypeSelector === index && (
                            <div className="absolute top-full left-0 z-10 w-full mt-1 bg-background border rounded-md shadow-md">
                              <Button
                                type="button"
                                variant="ghost"
                                className="w-full justify-start font-normal rounded-none hover:bg-secondary/50"
                                onClick={() =>
                                  handleColumnTypeChange(index, "string")
                                }
                              >
                                <span className="flex items-center gap-2">
                                  <Type className="size-4" />
                                  <span>Text</span>
                                </span>
                              </Button>
                              <Button
                                type="button"
                                variant="ghost"
                                className="w-full justify-start font-normal rounded-none hover:bg-secondary/50"
                                onClick={() =>
                                  handleColumnTypeChange(index, "number")
                                }
                              >
                                <span className="flex items-center gap-2">
                                  <Hash className="size-4" />
                                  <span>Number</span>
                                </span>
                              </Button>
                              <Button
                                type="button"
                                variant="ghost"
                                className="w-full justify-start font-normal rounded-none hover:bg-secondary/50"
                                onClick={() =>
                                  handleColumnTypeChange(index, "boolean")
                                }
                              >
                                <span className="flex items-center gap-2">
                                  <CheckSquare className="size-4" />
                                  <span>Boolean</span>
                                </span>
                              </Button>
                              <Button
                                type="button"
                                variant="ghost"
                                className="w-full justify-start font-normal rounded-none hover:bg-secondary/50"
                                onClick={() =>
                                  handleColumnTypeChange(index, "date")
                                }
                              >
                                <span className="flex items-center gap-2">
                                  <Calendar className="size-4" />
                                  <span>Date</span>
                                </span>
                              </Button>
                            </div>
                          )}
                        </div>
                      </div>
                    </div>
                  )}
                </div>
              ))}
            </div>

            <Button
              type="button"
              variant="outline"
              className="w-full border-dashed"
              onClick={handleAddColumn}
            >
              <Plus className="mr-2 size-4" />
              Add New Field
            </Button>
          </div>

          <SheetFooter className="p-0">
            <Button type="submit" className="w-full sm:w-auto">
              Create Collection
            </Button>
          </SheetFooter>
        </form>
      </SheetContent>
    </Sheet>
  );
}

export default function CollectionSidebar({ params }: Route.ComponentProps) {
  const { projectId } = params;
  const [searchTerm, setSearchTerm] = useState("");
  const queryClient = useQueryClient();

  const handleSearch = (e: React.ChangeEvent<HTMLInputElement>) => {
    setSearchTerm(e.target.value);
  };

  const handleRefresh = async () => {
    return queryClient.invalidateQueries({
      queryKey: ["collections", projectId],
    });
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
              <RefreshButton
                onRefresh={handleRefresh}
                title="Refresh collections"
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
              <AddCollectionButton projectId={projectId} />
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
