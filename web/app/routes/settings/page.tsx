import React from "react";
import { AppHeader } from "~/components/shared/header";
import type { Route } from "./+types/page";

export function meta({}: Route.MetaArgs) {
  return [
    { title: "Settings - Fluxend" },
    { name: "description", content: "Manage your application settings" },
  ];
}

export default function Settings() {
  return (
      <>
        <AppHeader title="Settings" />
        <div className="flex flex-col items-center justify-center min-h-[70vh] px-4">
          <div className="text-center max-w-md">
            <div className="mb-6">
              <div className="mx-auto w-16 h-16 bg-gradient-to-br from-blue-500 to-blue-700 rounded-full flex items-center justify-center">
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
                      d="M10.325 4.317c.426-1.756 2.924-1.756 3.35 0a1.724 1.724 0 002.573 1.066c1.543-.94 3.31.826 2.37 2.37a1.724 1.724 0 001.065 2.572c1.756.426 1.756 2.924 0 3.35a1.724 1.724 0 00-1.066 2.573c.94 1.543-.826 3.31-2.37 2.37a1.724 1.724 0 00-2.572 1.065c-.426 1.756-2.924 1.756-3.35 0a1.724 1.724 0 00-2.573-1.066c-1.543.94-3.31-.826-2.37-2.37a1.724 1.724 0 00-1.065-2.572c-1.756-.426-1.756-2.924 0-3.35a1.724 1.724 0 001.066-2.573c-.94-1.543.826-3.31 2.37-2.37.996.608 2.296.07 2.572-1.065z"
                  />
                  <path
                      strokeLinecap="round"
                      strokeLinejoin="round"
                      strokeWidth={2}
                      d="M15 12a3 3 0 11-6 0 3 3 0 016 0z"
                  />
                </svg>
              </div>
            </div>

            <h2 className="text-2xl font-bold text-foreground mb-3">
              Settings Coming Soon
            </h2>

            <p className="text-muted-foreground mb-6 leading-relaxed">
              We're working hard to bring you comprehensive application settings
              with support for <strong>Mail Configuration</strong>,{" "}
              <strong>Storage Drivers</strong>, <strong>API Management</strong>,
              and <strong>Security Settings</strong>
            </p>

            <div className="inline-flex items-center px-4 py-2 bg-blue-50 text-blue-700 rounded-full text-sm font-medium">
              <div className="w-2 h-2 bg-blue-500 rounded-full mr-2 animate-pulse"></div>
              In Development
            </div>
          </div>
        </div>
      </>
  );
}