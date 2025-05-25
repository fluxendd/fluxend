import {
  Command,
  Database,
  LayoutDashboard,
  LogOutIcon,
  PackageOpen,
  Parentheses,
  Scroll,
  Settings2,
} from "lucide-react";
import { href, Link, useLocation, useParams } from "react-router";

import {
  Sidebar,
  SidebarContent,
  SidebarFooter,
  SidebarGroup,
  SidebarGroupContent,
  SidebarHeader,
  SidebarMenu,
  SidebarMenuButton,
  SidebarMenuItem,
} from "~/components/ui/sidebar";

type AppSidebarItem = {
  title: string;
  url: string;
  icon: React.ComponentType;
  isActive?: boolean;
};

// Menu items.
const items: AppSidebarItem[] = [
  {
    title: "Dashboard",
    url: "",
    icon: LayoutDashboard,
  },
  {
    title: "Collections",
    url: "collections",
    icon: Database,
    isActive: true,
  },
  {
    title: "Functions",
    url: "functions",
    icon: Parentheses,
  },
  { title: "Storage", url: "storage", icon: PackageOpen },
  {
    title: "Logs",
    url: "logs",
    icon: Scroll,
  },
  {
    title: "Settings",
    url: "settings",
    icon: Settings2,
  },
];

export function AppSidebar() {
  const { projectId = "not-found" } = useParams();
  const location = useLocation();

  const checkIsActive = (url: string) => {
    if (url === "") {
      return location.pathname.endsWith("/");
    }

    return location.pathname.includes(url);
  };

  return (
    <Sidebar collapsible="icon">
      <SidebarHeader>
        <SidebarMenu>
          <SidebarMenuItem>
            <SidebarMenuButton size="lg" asChild className="md:h-8 md:p-0">
              <a href="#">
                <div className="flex aspect-square size-8 items-center justify-center rounded-lg bg-sidebar-primary text-sidebar-primary-foreground">
                  <Command className="size-4" />
                </div>
                <div className="grid flex-1 text-left text-sm leading-tight">
                  <span className="truncate font-semibold">Fluxend</span>
                  <span className="truncate text-xs">Enterprise</span>
                </div>
              </a>
            </SidebarMenuButton>
          </SidebarMenuItem>
        </SidebarMenu>
      </SidebarHeader>
      <SidebarContent>
        <SidebarGroup>
          <SidebarGroupContent>
            <SidebarMenu>
              {items.map((item) => (
                <SidebarMenuItem key={item.title}>
                  <SidebarMenuButton
                    asChild
                    tooltip={item.title}
                    isActive={checkIsActive(item.url)}
                  >
                    <Link
                      to={href(`projects/:projectId/${item.url}`, {
                        projectId,
                      })}
                    >
                      <item.icon />
                      <span>{item.title}</span>
                    </Link>
                  </SidebarMenuButton>
                </SidebarMenuItem>
              ))}
            </SidebarMenu>
          </SidebarGroupContent>
        </SidebarGroup>
      </SidebarContent>
      <SidebarFooter>
        <SidebarMenuButton asChild tooltip={"Logout"}>
          <Link to={href("/logout")} relative="route">
            <LogOutIcon />
          </Link>
        </SidebarMenuButton>
      </SidebarFooter>
    </Sidebar>
  );
}
