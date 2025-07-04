import { AppHeader } from "~/components/shared/header";
import type { Route } from "./+types/page";

export function meta({}: Route.MetaArgs) {
  return [
    { title: "Settings - Fluxend" },
    { name: "description", content: "Manage your account settings" },
  ];
}

export default function Settings() {
  return (
    <>
      <AppHeader title="Settings"></AppHeader>
    </>
  );
}
