import {
  BadgeCheck,
  Bell,
  ChevronsUpDown,
  ChartSpline,
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
  Cloudy,
  CloudUpload,
  Send,
  Captions,
  SendHorizontal,
  MessageCircleCode,
  MessageCircleCodeIcon,
  DatabaseBackupIcon, LucideDatabaseBackup, HelpCircle, GithubIcon,
} from "lucide-react";
import { href, NavLink, useOutletContext, type Params } from "react-router";

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
import type { User } from "~/services/user";
import { ThemeToggle } from "./theme-toggle";

type AppSidebarItem = {
  title: string;
  url: string;
  Icon: React.ComponentType;
  isActive?: boolean;
};

const items = [
  { title: "Dashboard", url: "dashboard", Icon: LayoutDashboard },
  { title: "Tables", url: "tables", Icon: Database, isActive: true },
  { title: "Logs", url: "logs", Icon: ChartSpline },
  { title: "Storage", url: "storage", Icon: CloudUpload },
  { title: "Forms", url: "forms", Icon: MessageCircleCodeIcon },
  { title: "Backups", url: "backups", Icon: LucideDatabaseBackup },
] as const satisfies readonly AppSidebarItem[];

type ProjectSidebarProps = {
  projectId: string;
  userDetails: User;
};

export const ProjectSidebar = memo(
  ({ projectId, userDetails }: ProjectSidebarProps) => (
    <Sidebar collapsible="icon">
      <SidebarHeader>
        <SidebarMenu>
          <SidebarMenuItem>
            <SidebarMenuButton size="lg" asChild className="md:h-8 md:p-0">
              <NavLink to="/projects">
                <div className="flex aspect-square size-8 items-center justify-center rounded-lg bg-muted text-black">
                  <Logo className="size-4" />
                </div>
              </NavLink>
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
                            className="absolute inset-0 bg-primary/30 dark:bg-primary/10 rounded-lg"
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
                          tooltip={item.title}
                          className="isolate hover:bg-transparent active:bg-transparent data-[active=true]:bg-transparent dark:active:text-primary dark:data-[active=true]:text-primary hover:text-foreground/70 transition-colors duration-300"
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
                <DropdownMenuLabel className="p-0 font-normal">
                  <div className="flex items-center gap-2 px-1 py-1.5 text-left text-sm">
                    <div className="grid flex-1 text-left text-sm leading-tight">
                      <span className="truncate font-medium">
                        {userDetails.username}
                      </span>
                      <span className="text-muted-foreground truncate text-xs">
                        {userDetails.email}
                      </span>
                    </div>
                  </div>
                </DropdownMenuLabel>
                <DropdownMenuSeparator />
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
                <ThemeToggle />
                <DropdownMenuSeparator />
                <DropdownMenuItem asChild>
                  <NavLink to="https://docs.fluxend.app" target="_blank">
                    <HelpCircle />
                    Documentation
                  </NavLink>
                </DropdownMenuItem>
                <DropdownMenuItem asChild>
                  <NavLink to="https://github.com/fluxendd/fluxend" target="_blank">
                    <GithubIcon />
                    Github
                  </NavLink>
                </DropdownMenuItem>
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
  )
);
