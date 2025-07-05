import React from "react";
import { AppHeader } from "~/components/shared/header";
import type { Route } from "./+types/page";

export function meta({}: Route.MetaArgs) {
  return [
    { title: "Forms - Fluxend" },
    { name: "description", content: "Create dynamic forms with validation" },
  ];
}

export default function Forms() {
  return (
      <>
        <AppHeader title="Forms" />
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
                      d="M9 12h6m-6 4h6m2 5H7a2 2 0 01-2-2V5a2 2 0 012-2h5.586a1 1 0 01.707.293l5.414 5.414a1 1 0 01.293.707V19a2 2 0 01-2 2z"
                  />
                </svg>
              </div>
            </div>

            <h2 className="text-2xl font-bold text-gray-900 mb-3">
              Dynamic Forms Coming Soon
            </h2>

            <p className="text-gray-600 mb-6 leading-relaxed">
              Create custom forms with validation rules and get instant endpoints for your HTML forms.
              Powerful form management is on the way!
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