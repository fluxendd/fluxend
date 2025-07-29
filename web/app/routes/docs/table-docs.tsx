import { useOutletContext } from "react-router";
import type { Route } from "./+types/table-docs";
import { Badge } from "~/components/ui/badge";
import { Button } from "~/components/ui/button";
import { cn } from "~/lib/utils";
import { ChevronRight, Database, Play } from "lucide-react";
import type { ProjectLayoutOutletContext } from "~/components/shared/project-layout";
import { ScrollArea } from "~/components/ui/scroll-area";

interface DocsOutletContext extends ProjectLayoutOutletContext {
  openApiSpec: any;
}

export function meta({ params }: Route.MetaArgs) {
  return [
    { title: `${params.table} API Documentation - Fluxend` },
    { name: "description", content: `API documentation for ${params.table} table` },
  ];
}

const getMethodColor = (method: string) => {
  const colors = {
    get: "text-blue-400",
    post: "text-green-400",
    put: "text-yellow-400",
    patch: "text-orange-400",
    delete: "text-red-400",
  };
  return colors[method.toLowerCase() as keyof typeof colors] || "text-gray-400";
};

const getMethodBgColor = (method: string) => {
  const colors = {
    get: "bg-blue-500/10 text-blue-400 border-blue-500/20",
    post: "bg-green-500/10 text-green-400 border-green-500/20",
    put: "bg-yellow-500/10 text-yellow-400 border-yellow-500/20",
    patch: "bg-orange-500/10 text-orange-400 border-orange-500/20",
    delete: "bg-red-500/10 text-red-400 border-red-500/20",
  };
  return colors[method.toLowerCase() as keyof typeof colors] || "bg-gray-500/10 text-gray-400 border-gray-500/20";
};

interface Parameter {
  name: string;
  in: string;
  required?: boolean;
  schema?: {
    type: string;
    format?: string;
    items?: any;
  };
  description?: string;
}

interface EndpointDetailsProps {
  path: string;
  method: string;
  details: any;
}

const EndpointDetails = ({ path, method, details }: EndpointDetailsProps) => {
  // Group parameters by location
  const headers = details.parameters?.filter((p: Parameter) => p.in === 'header') || [];
  const pathParams = details.parameters?.filter((p: Parameter) => p.in === 'path') || [];
  const queryParams = details.parameters?.filter((p: Parameter) => p.in === 'query') || [];

  return (
    <div className="bg-background border rounded-lg overflow-hidden">
      {/* Endpoint Header */}
      <div className="p-6 border-b">
        <div className="flex items-center justify-between mb-4">
          <h2 className="text-2xl font-semibold">{details.summary || `${method.toUpperCase()} ${path}`}</h2>
          <Button className="bg-green-500 hover:bg-green-600 text-white">
            Try it
            <Play className="ml-2 h-4 w-4" />
          </Button>
        </div>
        {details.description && (
          <p className="text-muted-foreground mb-4">{details.description}</p>
        )}
        <div className="flex items-center gap-4">
          <Badge 
            variant="outline" 
            className={cn("uppercase font-mono", getMethodBgColor(method))}
          >
            {method}
          </Badge>
          <code className="text-sm bg-muted px-3 py-1.5 rounded font-mono flex-1">
            {path.split('/').map((segment, index) => {
              if (!segment) return null;
              const isParam = segment.startsWith('{') && segment.endsWith('}');
              return (
                <span key={index}>
                  {index > 0 && <span className="text-muted-foreground">/</span>}
                  <span className={isParam ? "text-green-400" : "text-muted-foreground"}>
                    {segment}
                  </span>
                </span>
              );
            })}
          </code>
        </div>
      </div>

      {/* Headers Section */}
      {headers.length > 0 && (
        <div className="p-6 border-b">
          <h3 className="text-lg font-semibold mb-4">Headers</h3>
          <div className="space-y-4">
            {headers.map((param: Parameter) => (
              <div key={param.name} className="space-y-2">
                <div className="flex items-center gap-3">
                  <code className="text-green-400 font-mono">{param.name}</code>
                  <span className="text-sm text-muted-foreground">string</span>
                  {param.required && (
                    <span className="text-xs text-red-400">required</span>
                  )}
                </div>
                {param.description && (
                  <p className="text-sm text-muted-foreground ml-0">{param.description}</p>
                )}
              </div>
            ))}
          </div>
        </div>
      )}

      {/* Path Parameters Section */}
      {pathParams.length > 0 && (
        <div className="p-6 border-b">
          <h3 className="text-lg font-semibold mb-4">Path Parameters</h3>
          <div className="space-y-4">
            {pathParams.map((param: Parameter) => (
              <div key={param.name} className="space-y-2">
                <div className="flex items-center gap-3">
                  <code className="text-green-400 font-mono">{param.name}</code>
                  <span className="text-sm text-muted-foreground">{param.schema?.type || 'string'}</span>
                  {param.required && (
                    <span className="text-xs text-red-400">required</span>
                  )}
                </div>
                {param.description && (
                  <p className="text-sm text-muted-foreground ml-0">{param.description}</p>
                )}
              </div>
            ))}
          </div>
        </div>
      )}

      {/* Query Parameters Section */}
      {queryParams.length > 0 && (
        <div className="p-6 border-b">
          <h3 className="text-lg font-semibold mb-4">Query Parameters</h3>
          <div className="space-y-4">
            {queryParams.map((param: Parameter) => (
              <div key={param.name} className="space-y-2">
                <div className="flex items-center gap-3">
                  <code className="text-green-400 font-mono">{param.name}</code>
                  <span className="text-sm text-muted-foreground">{param.schema?.type || 'string'}</span>
                  {param.required && (
                    <span className="text-xs text-red-400">required</span>
                  )}
                </div>
                {param.description && (
                  <p className="text-sm text-muted-foreground ml-0">{param.description}</p>
                )}
              </div>
            ))}
          </div>
        </div>
      )}

      {/* Request Body Section */}
      {details.requestBody && (
        <div className="p-6 border-b">
          <h3 className="text-lg font-semibold mb-4">Request Body</h3>
          {Object.entries(details.requestBody.content || {}).map(([contentType, content]: [string, any]) => (
            <div key={contentType} className="space-y-2">
              <code className="text-sm text-muted-foreground">{contentType}</code>
              {content.schema && (
                <pre className="mt-2 p-4 bg-muted/50 rounded-lg text-sm overflow-x-auto">
                  {JSON.stringify(content.schema, null, 2)}
                </pre>
              )}
            </div>
          ))}
        </div>
      )}

      {/* Responses Section */}
      {details.responses && (
        <div className="p-6">
          <h3 className="text-lg font-semibold mb-4">Responses</h3>
          <div className="space-y-4">
            {Object.entries(details.responses).map(([statusCode, response]: [string, any]) => (
              <div key={statusCode} className="space-y-2">
                <div className="flex items-center gap-3">
                  <Badge 
                    variant={statusCode.startsWith('2') ? 'default' : 'destructive'} 
                    className="font-mono"
                  >
                    {statusCode}
                  </Badge>
                  <span className="text-sm text-muted-foreground">{response.description}</span>
                </div>
                {response.content && Object.entries(response.content).map(([contentType, content]: [string, any]) => (
                  <div key={contentType} className="ml-0">
                    <code className="text-sm text-muted-foreground">{contentType}</code>
                    {content.example && (
                      <pre className="mt-2 p-4 bg-muted/50 rounded-lg text-sm overflow-x-auto">
                        {JSON.stringify(content.example, null, 2)}
                      </pre>
                    )}
                  </div>
                ))}
              </div>
            ))}
          </div>
        </div>
      )}
    </div>
  );
};

export default function TableDocs({ params }: Route.ComponentProps) {
  const { projectDetails, openApiSpec } = useOutletContext<DocsOutletContext>();
  const { table } = params;

  if (!openApiSpec) {
    return (
      <div className="flex h-full items-center justify-center">
        <div className="text-center">
          <p className="text-muted-foreground">No API documentation available</p>
        </div>
      </div>
    );
  }

  // Filter paths for this specific table
  const tablePaths = Object.entries(openApiSpec.paths || {}).filter(([path]) => 
    path.startsWith(`/${table}`)
  );

  if (tablePaths.length === 0) {
    return (
      <div className="flex h-full items-center justify-center">
        <div className="text-center">
          <p className="text-muted-foreground">No documentation found for table: {table}</p>
        </div>
      </div>
    );
  }

  return (
    <ScrollArea className="h-full">
      <div className="p-8 max-w-5xl mx-auto">
        {/* Header */}
        <div className="mb-8">
          <div className="flex items-center gap-2 text-sm text-muted-foreground mb-2">
            <span>API Documentation</span>
            <ChevronRight className="w-4 h-4" />
            <span className="text-green-400">{table}</span>
          </div>
          <div className="flex items-center gap-3">
            <div className="p-2 bg-muted rounded-lg">
              <Database className="h-6 w-6" />
            </div>
            <div>
              <h1 className="text-3xl font-bold capitalize">{table}</h1>
              <p className="text-muted-foreground">Manage {table} data through the API</p>
            </div>
          </div>
        </div>

        {/* Endpoints */}
        <div className="space-y-8">
          {tablePaths.map(([path, methods]) => (
            Object.entries(methods).map(([method, details]: [string, any]) => (
              <EndpointDetails 
                key={`${method}-${path}`}
                path={path}
                method={method}
                details={{ ...details, table }}
              />
            ))
          ))}
        </div>
      </div>
    </ScrollArea>
  );
}