import { data, NavLink, redirect, useFetcher, useNavigate } from "react-router";
import type { Route } from "./+types/signup";
import { useState } from "react";
import { Button } from "~/components/ui/button";
import { Input } from "~/components/ui/input";
import { sessionCookie } from "~/lib/cookies";
import { signup } from "~/services/auth";
import { LoaderCircle, LogIn } from "lucide-react";

export function meta({}: Route.MetaArgs) {
  return [
    { title: "Frontend for Fluxton" },
    { name: "description", content: "Login to Fluxton" },
  ];
}

export async function loader({ request }: Route.LoaderArgs) {
  const cookieHeader = request.headers.get("Cookie");
  const sessionToken = await sessionCookie.parse(cookieHeader);

  if (sessionToken) {
    return redirect("/projects/123");
  }

  return null;
}

export async function action({ request }: Route.ActionArgs) {
  let formData = await request.formData();
  let email = formData.get("email")?.toString();
  let username = formData.get("username")?.toString();
  let password = formData.get("password")?.toString();

  if (!email || !password || !username) {
    return data(
      { error: "Missing email, username or  password" },
      { status: 400 }
    );
  }

  const { success, errors, content } = await signup(email, username, password);
  const error = errors?.[0];

  if (!success && error) {
    return data({ error }, { status: 401 });
  }

  if (!content) {
    return data({ error: "Unexpected Error Occured" }, { status: 401 });
  }

  return data({
    success: true,
    error: undefined,
    content,
  });
}

export default function SignUp() {
  const fetcher = useFetcher();
  const navigate = useNavigate();
  const [email, setEmail] = useState("");
  const [username, setUsername] = useState("");
  const [password, setPassword] = useState("");

  const isLoading = fetcher.state != "idle";
  const data = fetcher.data;

  if (data?.success) {
    console.log(data.success, "success");
    // Redirect to login after 3 seconds
    setTimeout(() => {
      navigate("/");
    }, 3000);
  }

  return (
    <div className="min-h-screen flex items-center justify-center py-12 px-4 sm:px-6 lg:px-8">
      <div className="max-w-md w-full space-y-8">
        <h2 className="text-3xl font-extrabold text-foreground mb-0">
          Register new account
        </h2>
        <span className="text-sm text-muted-foreground">
          Or{" "}
          <Button
            variant="link"
            className="p-0 font-medium text-primary hover:text-primary/80 hover:underline"
          >
            <NavLink to="/">sign in to your account</NavLink>
          </Button>
        </span>

        <div className="mt-8 space-y-6">
          <fetcher.Form className="space-y-6" method="post">
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
                htmlFor="username"
                className="block text-sm font-medium text-muted-foreground"
              >
                Username
              </label>
              <div className="mt-1">
                <Input
                  id="username"
                  name="username"
                  type="text"
                  autoComplete="email"
                  required
                  value={username}
                  onChange={(e) => setUsername(e.target.value)}
                  placeholder="your username"
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
            {data && data.error && (
              <div className="text-red-500 text-sm">{data.error}</div>
            )}
            {data && data.success && (
              <div className="text-green-600 text-sm">
                Signed Up Successfully, Redirecting to login...
              </div>
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
                Sign up
              </Button>
            </div>
          </fetcher.Form>
        </div>
      </div>
    </div>
  );
}
