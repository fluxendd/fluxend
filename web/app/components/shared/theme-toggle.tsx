import { Monitor, Moon, Sun } from "lucide-react";
import { useTheme } from "~/hooks/use-theme";
import {
  DropdownMenuItem,
  DropdownMenuPortal,
  DropdownMenuSub,
  DropdownMenuSubContent,
  DropdownMenuSubTrigger,
} from "~/components/ui/dropdown-menu";

export function ThemeToggle() {
  const { theme, setTheme } = useTheme();

  return (
    <DropdownMenuSub>
      <DropdownMenuSubTrigger className="gap-2" aria-label="Theme selection">
        <div className="relative h-4 w-4">
          <Sun className="absolute h-4 w-4 rotate-0 scale-100 transition-all dark:-rotate-90 dark:scale-0" aria-hidden="true" />
          <Moon className="absolute h-4 w-4 rotate-90 scale-0 transition-all dark:rotate-0 dark:scale-100" aria-hidden="true" />
        </div>
        Theme
      </DropdownMenuSubTrigger>
      <DropdownMenuPortal>
        <DropdownMenuSubContent>
          <DropdownMenuItem 
            onClick={() => setTheme("light")}
            role="menuitemradio"
            aria-checked={theme === "light"}
          >
            <Sun className="h-4 w-4" aria-hidden="true" />
            Light
            {theme === "light" && (
              <span className="ml-auto text-xs" aria-label="Selected">✓</span>
            )}
          </DropdownMenuItem>
          <DropdownMenuItem 
            onClick={() => setTheme("dark")}
            role="menuitemradio"
            aria-checked={theme === "dark"}
          >
            <Moon className="h-4 w-4" aria-hidden="true" />
            Dark
            {theme === "dark" && (
              <span className="ml-auto text-xs" aria-label="Selected">✓</span>
            )}
          </DropdownMenuItem>
          <DropdownMenuItem 
            onClick={() => setTheme("system")}
            role="menuitemradio"
            aria-checked={theme === "system"}
          >
            <Monitor className="h-4 w-4" aria-hidden="true" />
            System
            {theme === "system" && (
              <span className="ml-auto text-xs" aria-label="Selected">✓</span>
            )}
          </DropdownMenuItem>
        </DropdownMenuSubContent>
      </DropdownMenuPortal>
    </DropdownMenuSub>
  );
}