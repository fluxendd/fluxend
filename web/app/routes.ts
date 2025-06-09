import {
  type RouteConfig,
  index,
  layout,
  prefix,
  route,
} from "@react-router/dev/routes";
import { routeFolder } from "./lib/router";

export default [
  index("./routes/auth/login.tsx"),
  route("signup", "./routes/auth/signup.tsx"),
  route("logout", "./routes/auth/logout.tsx"),
  layout("./components/shared/app-layout.tsx", [
    route("projects", "./routes/projects/page.tsx"),
  ]),
  ...prefix("projects/:projectId", [
    layout("./components/shared/project-layout.tsx", [
      route("dashboard", "./routes/dashboard/page.tsx"),
      route("collections", "./routes/collections/sidebar.tsx", [
        route("create", "./routes/collections/create.tsx"),
        route(":collectionId", "./routes/collections/page.tsx"),
        // route(":collectionId/edit", "./routes/collections/edit.tsx"),
      ]),
      // routeFolder("collections/:collectionId", "./routes/collections/"),
      routeFolder("functions", "./routes/functions/"),
      routeFolder("storage", "./routes/storage/"),
      routeFolder("logs", "./routes/logs/"),
      routeFolder("settings", "./routes/settings/"),
    ]),
  ]),
] satisfies RouteConfig;
