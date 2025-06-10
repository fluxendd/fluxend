import {
  BadgeCheck,
  Bell,
  ChevronsUpDown,
  CreditCard,
  Database,
  EllipsisIcon,
  EllipsisVertical,
  Files,
  LayoutDashboard,
  LogOut,
  LogOutIcon,
  PackageOpen,
  Parentheses,
  Scroll,
  Settings,
  Settings2,
  Sparkles,
} from "lucide-react";
import { href, NavLink, type Params } from "react-router";

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
import { Logo } from "./logo";
import { motion } from "motion/react";
import { memo } from "react";
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuGroup,
  DropdownMenuItem,
  DropdownMenuLabel,
  DropdownMenuSeparator,
  DropdownMenuTrigger,
} from "../ui/dropdown-menu";

type AppSidebarItem = {
  title: string;
  url: string;
  Icon: React.ComponentType;
  isActive?: boolean;
};

const items = [
  {
    title: "Dashboard",
    url: "dashboard",
    Icon: LayoutDashboard,
  },
  {
    title: "Collections",
    url: "collections",
    Icon: Database,
    isActive: true,
  },
  // {
  //   title: "Functions",
  //   url: "functions",
  //   Icon: Parentheses,
  // },
  // { title: "Storage", url: "storage", Icon: PackageOpen },
  // {
  //   title: "Logs",
  //   url: "logs",
  //   Icon: Scroll,
  // },
  // {
  //   title: "Settings",
  //   url: "settings",
  //   Icon: Settings2,
  // },
] as const satisfies readonly AppSidebarItem[];

export const ProjectSidebar = memo(
  ({ projectId }: { projectId: string | undefined }) => {
    if (!projectId) {
      return null;
    }

    return (
      <Sidebar collapsible="icon">
        <SidebarHeader>
          <SidebarMenu>
            <SidebarMenuItem>
              <SidebarMenuButton size="lg" asChild className="md:h-8 md:p-0">
                <>
                  <div className="flex aspect-square size-8 items-center justify-center rounded-lg bg-muted text-black">
                    <Logo className="size-4" />
                  </div>
                </>
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
                    <NavLink
                      to={href(`/projects/:projectId/${item.url}`, {
                        projectId,
                      })}
                    >
                      {({ isActive }) => (
                        <>
                          {isActive && (
                            <motion.div
                              layoutId="sidebarItemId"
                              className="absolute inset-0 bg-sidebar-accent text-sidebar-accent-foreground rounded-md"
                              transition={{
                                type: "spring",
                                bounce: 0.2,
                                duration: 0.3,
                              }}
                            />
                          )}
                          <SidebarMenuButton
                            asChild
                            isActive={isActive}
                            // tooltip={item.title}
                            className="isolate hover:bg-transparent active:bg-transparent data-[active=true]:bg-transparent active:text-primary data-[active=true]:text-primary hover:text-primary transition-colors duration-300"
                          >
                            <item.Icon />
                          </SidebarMenuButton>
                        </>
                      )}
                    </NavLink>
                  </SidebarMenuItem>
                ))}
              </SidebarMenu>
            </SidebarGroupContent>
          </SidebarGroup>
        </SidebarContent>
        <SidebarFooter>
          <SidebarMenu>
            <SidebarMenuItem>
              <DropdownMenu>
                <DropdownMenuTrigger asChild>
                  <SidebarMenuButton className="data-[state=open]:bg-sidebar-accent data-[state=open]:text-sidebar-accent-foreground">
                    <EllipsisVertical />
                  </SidebarMenuButton>
                </DropdownMenuTrigger>
                <DropdownMenuContent
                  className="w-(--radix-dropdown-menu-trigger-width) min-w-56 rounded-lg"
                  // side={isMobile ? "bottom" : "right"}
                  side="right"
                  align="end"
                  sideOffset={4}
                >
                  <DropdownMenuGroup>
                    <DropdownMenuItem asChild>
                      <NavLink to="/projects">
                        <Files />
                        View All Projects
                      </NavLink>
                    </DropdownMenuItem>
                    <DropdownMenuItem asChild>
                      <NavLink to="/settings">
                        <Settings />
                        Settings
                      </NavLink>
                    </DropdownMenuItem>
                  </DropdownMenuGroup>
                  <DropdownMenuSeparator />
                  <DropdownMenuItem asChild>
                    <NavLink to="/logout">
                      <LogOut />
                      Log out
                    </NavLink>
                  </DropdownMenuItem>
                </DropdownMenuContent>
              </DropdownMenu>
            </SidebarMenuItem>
          </SidebarMenu>
        </SidebarFooter>
      </Sidebar>
    );
  }
);
