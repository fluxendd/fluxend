import React from "react";
import { AppHeader } from "~/components/shared/header";
import type { Route } from "./+types/page";

export function meta({}: Route.MetaArgs) {
  return [
    { title: "Backups - Fluxend" },
    {
      name: "description",
      content: "Automatic backups synced to your storage",
    },
  ];
}

export default function Backups() {
  return (
    <>
      <AppHeader title="Backups" />
      <div className="flex flex-col items-center justify-center min-h-[70vh] px-4">
        <div className="text-center max-w-md">
          <div className="mb-6">
            <div className="mx-auto w-16 h-16 bg-gradient-to-br from-amber-400 to-amber-600 rounded-full flex items-center justify-center">
              <svg
                className="w-8 h-8 text-white"
                fill="none"
                stroke="currentColor"
                viewBox="0 0 24 24"
              >
                <path
                  strokeLinecap="round"
                  strokeLinejoin="round"
                  strokeWidth={2}
                  d="M4 16v1a3 3 0 003 3h10a3 3 0 003-3v-1m-4-8l-4-4m0 0L8 8m4-4v12"
                />
              </svg>
            </div>
          </div>

          <h2 className="text-2xl font-bold text-foreground mb-3">
            Automatic Backups Coming Soon
          </h2>

          <p className="text-muted-foreground mb-6 leading-relaxed">
            Set up automatic backups that seamlessly sync with your selected
            storage driver. Your data will be safe and always accessible!
          </p>

          <div className="inline-flex items-center px-4 py-2 bg-amber-50 text-amber-700 rounded-full text-sm font-medium">
            <div className="w-2 h-2 bg-amber-500 rounded-full mr-2 animate-pulse"></div>
            In Development
          </div>
        </div>
      </div>
    </>
  );
}

