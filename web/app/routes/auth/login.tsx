import { NavLink, useFetcher, data, Link, redirect } from "react-router";
import { useState } from "react";
import { Button } from "~/components/ui/button";
import { Input } from "~/components/ui/input";
import { login } from "~/services/auth";
import { sessionCookie, organizationCookie } from "~/lib/cookies";
import type { Route } from "./+types/login";
import { GalleryVerticalEnd, LoaderCircle, LogIn } from "lucide-react";
import { getServerAuthToken } from "~/lib/auth";
import { initializeServices } from "~/services";
import { cn } from "~/lib/utils";
import {
  Card,
  CardContent,
  CardDescription,
  CardHeader,
  CardTitle,
} from "~/components/ui/card";
import { Logo } from "~/components/shared/logo";
import { NeonCard } from "~/components/shared/neon-card";

export function meta({}: Route.MetaArgs) {
  return [
    { title: "Fluxend - The only backend you'll ever need" },
    { name: "description", content: "Login to Fluxend" },
  ];
}

// Check if session_token is present, redirect it to dashboard
// Its also duplicate in signup.tsx, make sure to make changes there as well.
export async function loader({ request }: Route.LoaderArgs) {
  const sessionToken = await getServerAuthToken(request.headers);

  if (sessionToken) {
    return redirect(`/projects`);
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

    if (!success && errors?.[0]) {
      return data({ error: errors?.[0] }, { status: 401 });
    }

    if (!content || !content.token) {
      return data({ error: "Unexpected Error Occured" }, { status: 401 });
    }

    const services = initializeServices(content.token);

    const {
      success: orgSuccess,
      errors: orgErrors,
      content: orgContent,
    } = await services.user.getUserOrganizations();

    if (!orgSuccess && orgErrors?.[0]) {
      return data({ error: orgErrors[0] }, { status: 401 });
    }

    if (!orgContent || !orgContent.length) {
      return data({ error: "No organization found" }, { status: 401 });
    }

    const organization = orgContent?.[0];

    const sessionTokenCookie = await sessionCookie.serialize(content.token);
    const organizationIdCookie = await organizationCookie.serialize(
      organization.uuid
    );

    // Create a Response object that allows setting multiple cookies
    const response = new Response(null, {
      status: 302,
      headers: {
        Location: `/projects`,
      },
    });

    // Append each cookie separately
    response.headers.append("Set-Cookie", sessionTokenCookie);
    response.headers.append("Set-Cookie", organizationIdCookie);

    return response;
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
    <div className="bg-background flex min-h-svh flex-col items-center justify-center gap-6 p-6 md:p-10">
      <div className="flex w-full max-w-sm flex-col gap-6">
        <div className="flex items-center gap-2 self-center font-medium">
          <div className="bg-muted flex py-1 px-2 items-center justify-center rounded-lg">
            <Logo className="size-4" />
            <p className="ml-1">Fluxend</p>
          </div>
        </div>
        <div className={cn("flex flex-col gap-6")}>
          <NeonCard className="shadow-md">
            <CardHeader className="text-center">
              <CardTitle className="text-xl">Sign in to your account</CardTitle>
              <CardDescription>
                Or{" "}
                <NavLink to="/signup">
                  <Button
                    variant="link"
                    className="p-0 font-medium text-primary hover:text-primary/80 hover:underline cusror-pointer"
                  >
                    create a new account
                  </Button>
                </NavLink>
              </CardDescription>
            </CardHeader>
            <CardContent>
              <div className="grid gap-6">
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
                        disabled={isLoading}
                      >
                        Forgot your password?
                      </Button>
                    </div>
                  </div>
                  {data && (
                    <div className="text-red-500 text-sm">{data.error}</div>
                  )}
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
            </CardContent>
          </NeonCard>
          <div className="text-muted-foreground *:[a]:hover:text-primary text-center text-xs text-balance *:[a]:underline *:[a]:underline-offset-4">
            By clicking continue, you agree to our{" "}
            <Link to="/">Terms of Service</Link> and{" "}
            <Link to="/">Privacy Policy</Link>.
          </div>
        </div>
      </div>
    </div>
  );
}
