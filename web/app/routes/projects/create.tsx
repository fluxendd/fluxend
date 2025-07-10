import { organizationCookie, sessionCookie } from "~/lib/cookies";
import type { Route } from "./+types/create";
import { initializeServices } from "~/services";

export async function clientAction({ request }: Route.ClientActionArgs) {
  let formData = await request.formData();

  const name = formData.get("name") as string;
  if (!name) throw new Error("Name is required");

  const description = formData.get("description") as string;

  const organizationId = await organizationCookie.parse(document.cookie);
  const authToken = await sessionCookie.parse(document.cookie);

  const services = initializeServices(authToken);

  const data = await services.user.createUserProject(name, description, organizationId);
  return data;
}
