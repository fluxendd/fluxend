import React from "react";
import { AppHeader } from "~/components/shared/header";
import type { Route } from "./+types/page";

export function meta({}: Route.MetaArgs) {
  return [
    { title: "Storage - Fluxend" },
    { name: "description", content: "Manage your file storage" },
  ];
}

export default function Storage() {
  return (
    <>
      <AppHeader title="Storage" />
    </>
  );
}
