import { useEffect, useLayoutEffect, useState } from "react";

export function useTheme() {
  const [theme, setTheme] = useState<"light" | "dark">("light");

  useEffect(() => {
    // Check initial preference
    const mediaQuery = window.matchMedia("(prefers-color-scheme: dark)");
    setTheme(mediaQuery.matches ? "dark" : "light");

    // Listen for changes
    const handleChange = (e: MediaQueryListEvent) => {
      setTheme(e.matches ? "dark" : "light");
    };

    mediaQuery.addEventListener("change", handleChange);

    return () => {
      mediaQuery.removeEventListener("change", handleChange);
    };
  }, []);

  // useLayoutEffect(() => {
  //   // Apply theme to document
  //   const root = document.documentElement;
  //   if (theme === "dark") {
  //     root.classList.add("dark");
  //   } else {
  //     root.classList.remove("dark");
  //   }
  // }, [theme]);

  return theme;
}

