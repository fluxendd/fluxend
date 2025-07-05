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
                      d="M7 16a4 4 0 01-.88-7.903A5 5 0 1115.9 6L16 6a5 5 0 011 9.9M9 19l3 3m0 0l3-3m-3 3V10"
                  />
                </svg>
              </div>
            </div>

            <h2 className="text-2xl font-bold text-gray-900 mb-3">
              Storage Coming Soon
            </h2>

            <p className="text-gray-600 mb-6 leading-relaxed">
              We're working hard to bring you powerful file storage and management capabilities with support for <strong>S3</strong>, <strong>Dropbox</strong> and <strong>Backblaze</strong>
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