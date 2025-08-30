import { AppHeader } from "~/components/shared/header";
import type { Route } from "./+types/page";
import { getServerAuthToken } from "~/lib/auth";
import { initializeServices } from "~/services";
import { Form, redirect, useFetcher, useRevalidator } from "react-router";
import { organizationCookie } from "~/lib/cookies";
import ProjectsList from "./projects-list";
import { Button } from "~/components/ui/button";
import { File, PackagePlus, RefreshCw, Trash2 } from "lucide-react";
import { RefreshButton } from "~/components/shared/refresh-button";
import {
  Dialog,
  DialogClose,
  DialogContent,
  DialogDescription,
  DialogFooter,
  DialogHeader,
  DialogTitle,
  DialogTrigger,
} from "~/components/ui/dialog";

export function meta({}: Route.MetaArgs) {
  return [
    { title: "Projects - Fluxend" },
    { name: "description", content: "Manage your Fluxend projects" },
  ];
}
import { Label } from "~/components/ui/label";
import { Input } from "~/components/ui/input";
import { useEffect, useState } from "react";
import { Textarea } from "~/components/ui/textarea";

export async function loader({ request }: Route.LoaderArgs) {
  const authToken = await getServerAuthToken(request.headers);

  if (!authToken) {
    return redirect("/logout");
  }

  const organizationId = await organizationCookie.parse(
    request.headers.get("Cookie")
  );

  if (!organizationId) {
    return redirect("/logout");
  }

  const services = initializeServices(authToken);

  const { content: projects, ok } = await services.user.getUserProjects(
    organizationId
  );

  if (!ok) {
    return redirect("/logout");
  }

  return { title: "Projects", projects: projects || [] };
}

const CreateProjectDialog = ({ children }: { children: React.ReactNode }) => {
  const fetcher = useFetcher();
  const [isOpen, setIsOpen] = useState(false);
  let busy = fetcher.state !== "idle";

  useEffect(() => {
    if (fetcher.data?.ok) {
      setIsOpen(false);
    }
  }, [fetcher.data]);

  return (
    <Dialog open={isOpen} onOpenChange={setIsOpen}>
      {children}
      <DialogContent className="sm:max-w-[425px]">
        <fetcher.Form
          method="post"
          action="/projects/create"
          className="grid gap-4"
        >
          <DialogHeader>
            <DialogTitle>Create Projects</DialogTitle>
            <DialogDescription>
              Describe your project and its purpose
            </DialogDescription>
          </DialogHeader>
          <div className="grid gap-4">
            <div className="grid gap-3">
              <Label htmlFor="name-1">Name</Label>
              <Input id="name-1" name="name" placeholder="Dragonstone" />
            </div>
          </div>
          <div className="grid gap-4">
            <div className="grid gap-3">
              <Label htmlFor="description">Description</Label>
              <Textarea
                id="description"
                name="description"
                placeholder="What is this project about?"
              />
            </div>
          </div>
          <DialogFooter>
            <DialogClose asChild>
              <Button variant="outline">Cancel</Button>
            </DialogClose>
            <Button type="submit" className="cursor-pointer" disabled={busy}>
              Create
            </Button>
          </DialogFooter>
        </fetcher.Form>
      </DialogContent>
    </Dialog>
  );
};

const NoProjectMessage = () => (
  <div className="relative flex flex-col items-center justify-center py-12 px-4">
    <File className="size-10 mb-3" />
    <h3 className="text-lg font-medium text-gray-900 dark:text-gray-100 mb-2">
      No projects yet
    </h3>
    <p className="text-sm text-gray-500 dark:text-gray-400 text-center max-w-sm">
      Get started by creating your first project. Projects help you organize
      your work and collaborate with your team.
    </p>
    <CreateProjectDialog>
      <DialogTrigger asChild>
        <Button className="cursor-pointer mt-4">
          <PackagePlus className="size-4" /> Create Project
        </Button>
      </DialogTrigger>
    </CreateProjectDialog>
  </div>
);

const ProjectsPage = ({ loaderData }: Route.ComponentProps) => {
  const { title, projects } = loaderData;
  const revalidator = useRevalidator();

  return (
    <div>
      <AppHeader title={title}>
        <RefreshButton
          onRefresh={() => {
            revalidator.revalidate();
          }}
        />
        <CreateProjectDialog>
          <DialogTrigger asChild>
            <Button size="icon" className="cursor-pointer">
              <PackagePlus className="size-4" />
            </Button>
          </DialogTrigger>
        </CreateProjectDialog>
      </AppHeader>
      {!projects.length && <NoProjectMessage />}
      <ProjectsList projects={projects} />
    </div>
  );
};

export default ProjectsPage;
