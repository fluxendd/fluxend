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
import { href, NavLink, useHref, useLocation, useParams } from "react-router";

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
import { cn } from "~/lib/utils";
import { memo, useMemo } from "react";
import { useQueryClient } from "@tanstack/react-query";

type AppSidebarItem = {
  title: string;
  url: string;
  Icon: React.ComponentType;
  isActive?: boolean;
};

const items: AppSidebarItem[] = [
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
  {
    title: "Functions",
    url: "functions",
    Icon: Parentheses,
  },
  { title: "Storage", url: "storage", Icon: PackageOpen },
  {
    title: "Logs",
    url: "logs",
    Icon: Scroll,
  },
  {
    title: "Settings",
    url: "settings",
    Icon: Settings2,
  },
];

export const AppSidebar = memo(
  ({ projectId }: { projectId: string | undefined }) => {
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
                      to={href(`projects/:projectId/${item.url}`, {
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
          <SidebarMenuButton asChild tooltip={"Logout"}>
            <NavLink to={href("/logout")} relative="route">
              <LogOutIcon />
            </NavLink>
          </SidebarMenuButton>
        </SidebarFooter>
      </Sidebar>
    );
  }
);
