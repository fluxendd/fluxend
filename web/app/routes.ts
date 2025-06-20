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
    route("/settings", "./routes/settings/page.tsx"),
    route("projects", "./routes/projects/page.tsx"),
    route("projects/create", "./routes/projects/create.tsx"),
  ]),
  ...prefix("projects/:projectId", [
    layout("./components/shared/project-layout.tsx", [
      route("dashboard", "./routes/dashboard/page.tsx"),
      route("tables/create", "./routes/tables/create.tsx"),
      route("tables", "./routes/tables/sidebar.tsx", [
        route(":tableId", "./routes/tables/page.tsx", [
          route("edit", "./routes/tables/edit.tsx"),
        ]),
      ]),
      routeFolder("functions", "./routes/functions/"),
      routeFolder("storage", "./routes/storage/"),
      routeFolder("logs", "./routes/logs/"),
    ]),
  ]),
] satisfies RouteConfig;
