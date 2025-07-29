import { useQuery } from "@tanstack/react-query";
import { useSearchParams, useNavigate, useOutletContext } from "react-router";
import type { Route } from "./+types/page";
import type { ProjectLayoutOutletContext } from "~/components/shared/project-layout";
import { ChevronRight, FileText, ExternalLink } from "lucide-react";
import { useState, useEffect, useMemo } from "react";
import type { OpenAPIResponse } from "~/services/openapi";
import { Button } from "~/components/ui/button";
import { ScrollArea } from "~/components/ui/scroll-area";
import { Badge } from "~/components/ui/badge";
import { cn } from "~/lib/utils";

export function meta({}: Route.MetaArgs) {
  return [
    { title: "API Documentation - Fluxend" },
    { name: "description", content: "API documentation for your project" },
  ];
}

interface OpenAPIPath {
  [method: string]: {
    tags?: string[];
    summary?: string;
    description?: string;
    parameters?: Array<{
      name: string;
      in: string;
      required?: boolean;
      schema: {
        type: string;
        format?: string;
        items?: any;
      };
      description?: string;
    }>;
    requestBody?: {
      content: {
        [contentType: string]: {
          schema: any;
        };
      };
    };
    responses?: {
      [statusCode: string]: {
        description: string;
        content?: {
          [contentType: string]: {
            schema: any;
            example?: any;
          };
        };
      };
    };
  };
}

interface OpenAPISpec {
  openapi: string;
  info: {
    title: string;
    version: string;
    description?: string;
  };
  servers?: Array<{
    url: string;
    description?: string;
  }>;
  paths: {
    [path: string]: OpenAPIPath;
  };
  components?: {
    schemas?: {
      [name: string]: any;
    };
  };
}

const getMethodColor = (method: string) => {
  const colors = {
    get: "bg-blue-500",
    post: "bg-green-500",
    put: "bg-yellow-500",
    patch: "bg-orange-500",
    delete: "bg-red-500",
  };
  return colors[method.toLowerCase() as keyof typeof colors] || "bg-gray-500";
};

const ParametersList = ({ parameters }: { parameters: any[] }) => {
  if (!parameters || parameters.length === 0) return null;

  return (
    <div className="space-y-2">
      <h4 className="text-sm font-semibold text-gray-700 dark:text-gray-300">Parameters</h4>
      <div className="space-y-2">
        {parameters.map((param, idx) => (
          <div key={idx} className="border-l-2 border-gray-200 dark:border-gray-700 pl-3 py-1">
            <div className="flex items-center gap-2">
              <code className="text-sm font-mono text-blue-600 dark:text-blue-400">{param.name}</code>
              {param.required && <Badge variant="outline" className="text-xs">required</Badge>}
              <Badge variant="secondary" className="text-xs">{param.in}</Badge>
              <Badge variant="secondary" className="text-xs">{param.schema?.type}</Badge>
            </div>
            {param.description && (
              <p className="text-sm text-gray-600 dark:text-gray-400 mt-1">{param.description}</p>
            )}
          </div>
        ))}
      </div>
    </div>
  );
};

const ResponsesList = ({ responses }: { responses: any }) => {
  if (!responses) return null;

  return (
    <div className="space-y-2">
      <h4 className="text-sm font-semibold text-gray-700 dark:text-gray-300">Responses</h4>
      <div className="space-y-2">
        {Object.entries(responses).map(([statusCode, response]: [string, any]) => (
          <div key={statusCode} className="border-l-2 border-gray-200 dark:border-gray-700 pl-3 py-1">
            <div className="flex items-center gap-2">
              <Badge variant={statusCode.startsWith('2') ? 'default' : 'destructive'} className="text-xs">
                {statusCode}
              </Badge>
              <span className="text-sm text-gray-600 dark:text-gray-400">{response.description}</span>
            </div>
            {response.content && (
              <div className="mt-2">
                {Object.entries(response.content).map(([contentType, content]: [string, any]) => (
                  <div key={contentType} className="text-xs">
                    <span className="text-gray-500 dark:text-gray-400">{contentType}</span>
                    {content.example && (
                      <pre className="mt-1 p-2 bg-gray-100 dark:bg-gray-800 rounded text-xs overflow-x-auto">
                        {JSON.stringify(content.example, null, 2)}
                      </pre>
                    )}
                  </div>
                ))}
              </div>
            )}
          </div>
        ))}
      </div>
    </div>
  );
};

export default function DocsPage({ params }: Route.ComponentProps) {
  const { projectId } = params;
  const [searchParams] = useSearchParams();
  const navigate = useNavigate();
  const selectedTable = searchParams.get("table");
  const { projectDetails, services } = useOutletContext<ProjectLayoutOutletContext>();
  
  const [availableTables, setAvailableTables] = useState<string[]>([]);

  const { data, isLoading, error } = useQuery({
    queryKey: ["openapi", projectId, selectedTable],
    queryFn: async () => {
      const response = await services.openapi.getProjectOpenAPI(projectId, selectedTable || undefined);
      return response;
    },
  });

  const openApiSpec = useMemo(() => {
    if (!data?.content) return null;
    
    try {
      if (typeof data.content === 'object') {
        return data.content as OpenAPISpec;
      }
      return JSON.parse(data.content) as OpenAPISpec;
    } catch {
      return null;
    }
  }, [data?.content]);

  // Extract available tables from paths
  useEffect(() => {
    if (openApiSpec && !selectedTable) {
      const tables = new Set<string>();
      Object.keys(openApiSpec.paths).forEach(path => {
        const match = path.match(/^\/([^\/]+)(?:\/|$)/);
        if (match && match[1] !== 'rpc') {
          tables.add(match[1]);
        }
      });
      setAvailableTables(Array.from(tables).sort());
    }
  }, [openApiSpec, selectedTable]);

  if (isLoading) {
    return (
      <div className="flex h-full items-center justify-center">
        <div className="animate-pulse">Loading documentation...</div>
      </div>
    );
  }

  if (error || !openApiSpec) {
    return (
      <div className="flex h-full items-center justify-center">
        <div className="text-center">
          <p className="text-red-500 mb-2">Failed to load API documentation</p>
          <p className="text-sm text-gray-500">Please try again later</p>
        </div>
      </div>
    );
  }

  return (
    <div className="flex h-full">
      {/* Sidebar */}
      <div className="w-64 border-r bg-gray-50 dark:bg-gray-900">
        <div className="p-4">
          <h2 className="text-lg font-semibold mb-4">API Documentation</h2>
          <nav className="space-y-1">
            <button
              onClick={() => navigate(`/projects/${projectId}/docs`)}
              className={cn(
                "w-full text-left px-3 py-2 rounded-md text-sm transition-colors",
                !selectedTable
                  ? "bg-primary text-primary-foreground"
                  : "hover:bg-gray-100 dark:hover:bg-gray-800"
              )}
            >
              <FileText className="inline-block w-4 h-4 mr-2" />
              All Tables
            </button>
            {availableTables.map(table => (
              <button
                key={table}
                onClick={() => navigate(`/projects/${projectId}/docs?table=${table}`)}
                className={cn(
                  "w-full text-left px-3 py-2 rounded-md text-sm transition-colors",
                  selectedTable === table
                    ? "bg-primary text-primary-foreground"
                    : "hover:bg-gray-100 dark:hover:bg-gray-800"
                )}
              >
                {table}
              </button>
            ))}
          </nav>
        </div>
      </div>

      {/* Main content */}
      <ScrollArea className="flex-1">
        <div className="p-8 max-w-5xl mx-auto">
          {/* Header */}
          <div className="mb-8">
            <div className="flex items-center gap-2 text-sm text-gray-500 mb-2">
              <span>Projects</span>
              <ChevronRight className="w-4 h-4" />
              <span>{projectDetails?.name}</span>
              <ChevronRight className="w-4 h-4" />
              <span>API Documentation</span>
              {selectedTable && (
                <>
                  <ChevronRight className="w-4 h-4" />
                  <span>{selectedTable}</span>
                </>
              )}
            </div>
            <h1 className="text-3xl font-bold">{openApiSpec.info.title}</h1>
            {openApiSpec.info.description && (
              <p className="text-gray-600 dark:text-gray-400 mt-2">{openApiSpec.info.description}</p>
            )}
            <div className="flex items-center gap-4 mt-4">
              <Badge variant="secondary">Version {openApiSpec.info.version}</Badge>
              <Badge variant="secondary">OpenAPI {openApiSpec.openapi}</Badge>
            </div>
          </div>

          {/* Server info */}
          {openApiSpec.servers && openApiSpec.servers.length > 0 && (
            <div className="mb-8 p-4 bg-gray-100 dark:bg-gray-800 rounded-lg">
              <h3 className="text-sm font-semibold mb-2">Base URL</h3>
              {openApiSpec.servers.map((server, idx) => (
                <div key={idx} className="flex items-center gap-2">
                  <code className="text-sm">{server.url}</code>
                  {server.description && (
                    <span className="text-sm text-gray-600 dark:text-gray-400">
                      - {server.description}
                    </span>
                  )}
                </div>
              ))}
            </div>
          )}

          {/* Endpoints */}
          <div className="space-y-6">
            {Object.entries(openApiSpec.paths).map(([path, methods]) => {
              // Filter paths based on selected table
              if (selectedTable && !path.startsWith(`/${selectedTable}`)) {
                return null;
              }

              return (
                <div key={path} className="border rounded-lg p-6">
                  <div className="mb-4">
                    <code className="text-lg font-mono">{path}</code>
                  </div>
                  <div className="space-y-4">
                    {Object.entries(methods).map(([method, details]) => (
                      <div key={method} className="border-l-4 border-gray-200 dark:border-gray-700 pl-4">
                        <div className="flex items-center gap-3 mb-2">
                          <Badge className={cn("uppercase", getMethodColor(method))}>
                            {method}
                          </Badge>
                          {details.summary && (
                            <span className="text-sm font-medium">{details.summary}</span>
                          )}
                        </div>
                        {details.description && (
                          <p className="text-sm text-gray-600 dark:text-gray-400 mb-3">
                            {details.description}
                          </p>
                        )}
                        {details.parameters && (
                          <div className="mb-3">
                            <ParametersList parameters={details.parameters} />
                          </div>
                        )}
                        {details.responses && (
                          <ResponsesList responses={details.responses} />
                        )}
                      </div>
                    ))}
                  </div>
                </div>
              );
            })}
          </div>
        </div>
      </ScrollArea>
    </div>
  );
}