import { createContext, useContext, useEffect, useState } from "react";
import { THEME_STORAGE_KEY, THEME_VALUES, type Theme } from "~/lib/theme-constants";

interface ThemeContextType {
  theme: Theme;
  setTheme: (theme: Theme) => void;
  resolvedTheme: "light" | "dark";
}

const ThemeContext = createContext<ThemeContextType | undefined>(undefined);

export function ThemeProvider({ children }: { children: React.ReactNode }) {
  const [theme, setThemeState] = useState<Theme>(() => {
    // Check localStorage first
    if (typeof window !== "undefined") {
      try {
        const stored = localStorage.getItem(THEME_STORAGE_KEY);
        if (stored === "light" || stored === "dark" || stored === "system") {
          return stored;
        }
      } catch (error) {
        console.error("Failed to read theme from localStorage:", error);
      }
    }
    return "system";
  });

  // Initialize with a safe default, will be updated in useEffect
  const [resolvedTheme, setResolvedTheme] = useState<"light" | "dark">("light");

  // Handle theme changes
  useEffect(() => {
    const root = document.documentElement;
    let mediaQuery: MediaQueryList | null = null;
    let handleChange: ((e: MediaQueryListEvent) => void) | null = null;

    if (theme === "system") {
      // Check if matchMedia is supported
      if (typeof window !== "undefined" && window.matchMedia) {
        mediaQuery = window.matchMedia("(prefers-color-scheme: dark)");
        handleChange = (e: MediaQueryListEvent) => {
          const newTheme = e.matches ? "dark" : "light";
          setResolvedTheme(newTheme);
          root.classList.remove("light", "dark");
          root.classList.add(newTheme);
        };

        // Set initial theme
        const systemTheme = mediaQuery.matches ? "dark" : "light";
        setResolvedTheme(systemTheme);
        root.classList.remove("light", "dark");
        root.classList.add(systemTheme);

        // Listen for changes
        mediaQuery.addEventListener("change", handleChange);
      } else {
        // Fallback to light theme if matchMedia not supported
        setResolvedTheme("light");
        root.classList.remove("light", "dark");
        root.classList.add("light");
      }
    } else {
      // Manual theme
      setResolvedTheme(theme);
      root.classList.remove("light", "dark");
      root.classList.add(theme);
    }

    // Cleanup function
    return () => {
      if (mediaQuery && handleChange) {
        mediaQuery.removeEventListener("change", handleChange);
      }
    };
  }, [theme]);

  const setTheme = (newTheme: Theme) => {
    setThemeState(newTheme);
    try {
      localStorage.setItem(THEME_STORAGE_KEY, newTheme);
    } catch (error) {
      console.error("Failed to save theme to localStorage:", error);
    }
  };

  return (
    <ThemeContext.Provider value={{ theme, setTheme, resolvedTheme }}>
      {children}
    </ThemeContext.Provider>
  );
}

export function useTheme() {
  const context = useContext(ThemeContext);
  if (context === undefined) {
    throw new Error("useTheme must be used within a ThemeProvider");
  }
  return context;
}