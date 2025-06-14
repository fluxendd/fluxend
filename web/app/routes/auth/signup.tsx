import {
  data,
  Link,
  NavLink,
  redirect,
  useFetcher,
  useNavigate,
} from "react-router";
import type { Route } from "./+types/signup";
import { useState, useEffect } from "react";
import { Button } from "~/components/ui/button";
import { Input } from "~/components/ui/input";
import { signup } from "~/services/auth";
import { LoaderCircle, LogIn } from "lucide-react";
import { getAuthToken } from "~/lib/auth";
import { cn } from "~/lib/utils";
import { Logo } from "~/components/shared/logo";
import {
  CardContent,
  CardDescription,
  CardHeader,
  CardTitle,
} from "~/components/ui/card";
import { NeonCard } from "~/components/shared/neon-card";

export function meta({}: Route.MetaArgs) {
  return [
    { title: "Frontend for Fluxend" },
    { name: "description", content: "Login to Fluxend" },
  ];
}

// Check if session_token is present, redirect it to dashboard
// Its also duplicate in login.tsx, make sure to make changes there as well.
export async function loader({ request }: Route.LoaderArgs) {
  const sessionToken = await getAuthToken(request.headers);

  if (sessionToken) {
    return redirect(`/projects`);
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
    success,
    error,
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
    <div className="relative">
      {/* <div className="border-2 absolute inset-0 mx-20" />
      <div className="border-2 absolute inset-0 my-20" /> */}
      <div className="bg-background flex min-h-svh  items-center justify-center gap-6 p-6 md:p-10">
        <div className="flex w-full max-w-sm flex-col gap-6">
          <div className="flex items-center gap-2 self-center font-medium">
            <div className="bg-muted flex py-1 px-2 items-center justify-center rounded-md">
              <Logo className="size-4" />
              <p className="ml-1">Fluxend</p>
            </div>
          </div>
          <NeonCard className="shadow-md">
            <CardHeader className="text-center">
              <CardTitle className="text-xl">Sign in to your account</CardTitle>
              <CardDescription>
                Or{" "}
                <NavLink to="/">
                  <Button
                    variant="link"
                    className="p-0 font-medium text-primary hover:text-primary/80 hover:underline cusror-pointer"
                  >
                    sign in to your account
                  </Button>
                </NavLink>
              </CardDescription>
            </CardHeader>
            <CardContent>
              <div className={cn("flex flex-col gap-6")}>
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
