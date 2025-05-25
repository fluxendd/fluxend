import { Form, NavLink, redirect, useFetcher, data } from "react-router";
import { useState } from "react";
import { Button } from "~/components/ui/button";
import { Input } from "~/components/ui/input";
import { login } from "~/services/auth";
import {
  sessionCookie,
  organizationCookie,
  projectCookie,
  dbCookie,
} from "~/lib/cookies";
import type { Route } from "./+types/login";
import { LoaderCircle, LogIn } from "lucide-react";
import { getAuthToken } from "~/lib/auth";
import { getUserOrganizations, getUserProjects } from "~/services/user";

export function meta({}: Route.MetaArgs) {
  return [
    { title: "Frontend for Fluxend" },
    { name: "description", content: "Login to Fluxend" },
  ];
}

// Check if session_token is present, redirect it to dashboard
export async function loader({ request }: Route.LoaderArgs) {
  const sessionToken = await getAuthToken(request.headers);

  if (sessionToken) {
    // Parse the project ID cookie using the proper API
    const projectIdValue = await projectCookie.parse(
      request.headers.get("Cookie") || ""
    );

    if (projectIdValue) {
      // Redirect to the specific project
      return redirect(`/projects/${projectIdValue}`);
    } else {
      // If project cookie is missing, throw an error
      throw new Response("Project ID is missing", { status: 400 });
    }
  }
  return null;
}

export async function action({ request }: Route.ActionArgs) {
  let formData = await request.formData();
  let email = formData.get("email")?.toString();
  let password = formData.get("password")?.toString();

  if (!email || !password) {
    return data({ error: "Missing email or password" }, { status: 400 });
  }

  try {
    const { success, errors, content } = await login(email, password);
    const error = errors?.[0];

    if (!success && error) {
      return data({ error }, { status: 401 });
    }

    if (!content) {
      return data({ error: "Unexpected Error Occured" }, { status: 401 });
    }

    // Serialize session cookie
    const sessionTokenCookie = await sessionCookie.serialize(content.token);

    // Get user organizations
    const mockRequestHeaders = {
      headers: { Authorization: `Bearer ${content.token}` },
    };

    let organizationsResponse, organizationsData;
    try {
      organizationsResponse = await getUserOrganizations(mockRequestHeaders);
      organizationsData = await organizationsResponse.json();
    } catch (error) {
      console.error("Failed to fetch organizations:", error);
      return data({ error: "Failed to fetch organizations" }, { status: 500 });
    }

    let organizationId = "default_org";
    let projectId = "default_project";
    let dbId = "default_db";

    if (
      organizationsData.success &&
      organizationsData.content &&
      organizationsData.content.length > 0
    ) {
      organizationId = organizationsData.content[0].uuid;

      // Get user projects for the first organization
      let projectsResponse, projectsData;
      try {
        projectsResponse = await getUserProjects(
          mockRequestHeaders,
          organizationId
        );
        projectsData = await projectsResponse.json();
      } catch (error) {
        console.error("Failed to fetch projects:", error);
        return data({ error: "Failed to fetch projects" }, { status: 500 });
      }

      // Get first project UUID
      if (
        projectsData.success &&
        projectsData.content &&
        projectsData.content.length > 0
      ) {
        projectId = projectsData.content[0].uuid;
        dbId = projectsData.content[0].dbName;
      } else {
        console.warn("No projects found or invalid project data format");
      }
    } else {
      console.warn(
        "No organizations found or invalid organization data format"
      );
    }

    // Serialize organization and project cookies
    const organizationIdCookie = await organizationCookie.serialize(
      organizationId
    );
    const projectIdCookie = await projectCookie.serialize(projectId);
    const dbIdCookie = await dbCookie.serialize(dbId);

    return redirect(`/projects/${projectId}`, {
      headers: {
        "Set-Cookie": [
          sessionTokenCookie,
          organizationIdCookie,
          projectIdCookie,
          dbIdCookie,
        ].join(", "),
      },
    });
  } catch (error) {
    console.error("Unexpected error during login process:", error);
    return data(
      { error: "An unexpected error occurred during login" },
      { status: 500 }
    );
  }
}

export default function Login({}: Route.ComponentProps) {
  const fetcher = useFetcher();
  const [email, setEmail] = useState("");
  const [password, setPassword] = useState("");
  const data = fetcher.data;
  const isLoading = fetcher.state != "idle";

  return (
    <div className="min-h-screen flex items-center justify-center py-12 px-4 sm:px-6 lg:px-8">
      <div className="max-w-md w-full space-y-8">
        <h2 className="text-3xl font-extrabold text-foreground mb-0">
          Sign in to your account
        </h2>
        <span className="text-sm text-muted-foreground">
          Or{" "}
          <Button
            variant="link"
            className="p-0 font-medium text-primary hover:text-primary/80 hover:underline"
          >
            <NavLink to="/signup">create a new account</NavLink>
          </Button>
        </span>

        <div className="mt-8 space-y-6">
          <fetcher.Form method="post" className="space-y-6">
            <div>
              <label
                htmlFor="email"
                className="block text-sm font-medium text-muted-foreground"
              >
                Email address
              </label>
              <div className="mt-1">
                <Input
                  id="email"
                  name="email"
                  type="email"
                  autoComplete="email"
                  required
                  value={email}
                  onChange={(e) => setEmail(e.target.value)}
                  placeholder="your email"
                  className="w-full px-3 py-2  text-sm"
                />
              </div>
            </div>

            <div>
              <label
                htmlFor="password"
                className="block text-sm font-medium text-muted-foreground"
              >
                Password
              </label>
              <div className="mt-1">
                <Input
                  id="password"
                  name="password"
                  type="password"
                  autoComplete="current-password"
                  required
                  value={password}
                  placeholder="your password"
                  onChange={(e) => setPassword(e.target.value)}
                  className="w-full px-3 py-2 text-sm"
                />
              </div>
            </div>

            <div className="flex items-center justify-between">
              <div className="text-sm">
                <Button
                  variant="link"
                  className="font-medium text-muted-foreground p-0"
                >
                  Forgot your password?
                </Button>
              </div>
            </div>
            {data && <div className="text-red-500 text-sm">{data.error}</div>}
            <div>
              <Button
                disabled={isLoading}
                type="submit"
                className="w-full"
                size="lg"
              >
                {isLoading && <LoaderCircle className="loading-icon" />}
                {!isLoading && <LogIn />}
                Sign in
              </Button>
            </div>
          </fetcher.Form>
        </div>
      </div>
    </div>
  );
}
