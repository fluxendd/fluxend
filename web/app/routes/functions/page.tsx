import { AppHeader } from "~/components/shared/header";
import type { Route } from "./+types/page";

export function meta({}: Route.MetaArgs) {
  return [
    { title: "Functions - Fluxend" },
    { name: "description", content: "Manage your serverless functions" },
  ];
}

export default function Functions() {
  return (
    <>
      <AppHeader title="Functions" />
    </>
  );
}
