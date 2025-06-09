import { AppHeader } from "~/components/shared/header";
import type { Route } from "./+types/page";
import { getServerAuthToken } from "~/lib/auth";
import { initializeServices } from "~/services";
import { redirect } from "react-router";
import { organizationCookie } from "~/lib/cookies";
import ProjectsList from "./projects-list";

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

  const { content: projects } = await services.user.getUserProjects(
    organizationId
  );

  return { title: "Projects", projects: projects || [] };
}

const ProjectsPage = ({ loaderData }: Route.ComponentProps) => {
  const { title, projects } = loaderData;

  return (
    <div>
      <AppHeader title={title} />
      <ProjectsList projects={projects} />
    </div>
  );
};

export default ProjectsPage;
