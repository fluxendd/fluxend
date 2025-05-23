import { index, type RouteConfigEntry } from "@react-router/dev/routes";
import { existsSync } from "node:fs";
import { join } from "node:path";
import { fileURLToPath } from "node:url";

export const routeFolder = (
  path: string,
  routeFolder: string
): RouteConfigEntry => {
  const projectRoot = new URL("../../", import.meta.url);
  const absoluteRoutePath = join(
    fileURLToPath(projectRoot),
    "app",
    routeFolder
  );

  // check if sidebar.tsx exists (fixed typo from tex to tsx in comment)
  const sidebarFile = join(absoluteRoutePath, "sidebar.tsx");

  const sidebarExists = existsSync(sidebarFile);

  if (sidebarExists) {
    return {
      path,
      file: routeFolder + "sidebar.tsx", // Keep original format for the file path
      children: [index(routeFolder + "page.tsx")], // Keep original format for index
    };
  }

  return {
    path,
    file: routeFolder + "page.tsx",
    children: [],
  };
};
