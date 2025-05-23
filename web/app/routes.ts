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
  layout("./components/shared/layout.tsx", [
    ...prefix("projects/:projectId", [
      index("./routes/dashboard/page.tsx"),
      routeFolder("collections/:collectionId?", "./routes/collections/"),
      routeFolder("functions", "./routes/functions/"),
      routeFolder("storage", "./routes/storage/"),
      routeFolder("logs", "./routes/logs/"),
      routeFolder("settings", "./routes/settings/"),
    ]),
  ]),
] satisfies RouteConfig;
