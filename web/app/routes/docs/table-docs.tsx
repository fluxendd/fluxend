import { useOutletContext } from "react-router";
import type { Route } from "./+types/table-docs";
import { Badge } from "~/components/ui/badge";
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "~/components/ui/card";
import { cn } from "~/lib/utils";
import { ChevronRight, Database } from "lucide-react";
import type { ProjectLayoutOutletContext } from "~/components/shared/project-layout";

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
      <h4 className="text-sm font-semibold mb-2">Parameters</h4>
      <div className="space-y-2">
        {parameters.map((param, idx) => (
          <div key={idx} className="flex items-start gap-2 text-sm">
            <div className="flex-1">
              <div className="flex items-center gap-2 mb-1">
                <code className="font-mono text-blue-600 dark:text-blue-400">{param.name}</code>
                {param.required && <Badge variant="outline" className="text-xs h-5">required</Badge>}
                <Badge variant="secondary" className="text-xs h-5">{param.in}</Badge>
                <Badge variant="secondary" className="text-xs h-5">{param.schema?.type}</Badge>
              </div>
              {param.description && (
                <p className="text-muted-foreground ml-0">{param.description}</p>
              )}
            </div>
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
      <h4 className="text-sm font-semibold mb-2">Responses</h4>
      <div className="space-y-2">
        {Object.entries(responses).map(([statusCode, response]: [string, any]) => (
          <div key={statusCode} className="text-sm">
            <div className="flex items-center gap-2 mb-1">
              <Badge variant={statusCode.startsWith('2') ? 'default' : 'destructive'} className="text-xs">
                {statusCode}
              </Badge>
              <span className="text-muted-foreground">{response.description}</span>
            </div>
            {response.content && Object.entries(response.content).map(([contentType, content]: [string, any]) => (
              <div key={contentType} className="ml-6 mt-1">
                <code className="text-xs text-muted-foreground">{contentType}</code>
                {content.example && (
                  <pre className="mt-1 p-2 bg-muted rounded text-xs overflow-x-auto">
                    {JSON.stringify(content.example, null, 2)}
                  </pre>
                )}
              </div>
            ))}
          </div>
        ))}
      </div>
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
    <div className="p-8 max-w-5xl mx-auto">
      {/* Header */}
      <div className="mb-8">
        <div className="flex items-center gap-2 text-sm text-muted-foreground mb-2">
          <span>API Documentation</span>
          <ChevronRight className="w-4 h-4" />
          <span>{table}</span>
        </div>
        <div className="flex items-center gap-3">
          <Database className="h-8 w-8 text-muted-foreground" />
          <div>
            <h1 className="text-3xl font-bold">{table}</h1>
            <p className="text-muted-foreground">Database table API endpoints</p>
          </div>
        </div>
      </div>

      {/* Endpoints */}
      <div className="space-y-6">
        {tablePaths.map(([path, methods]) => (
          <Card key={path}>
            <CardHeader>
              <CardTitle className="font-mono text-lg">{path}</CardTitle>
            </CardHeader>
            <CardContent className="space-y-6">
              {Object.entries(methods).map(([method, details]: [string, any]) => (
                <div key={method} className="space-y-4">
                  <div className="flex items-center gap-3">
                    <Badge className={cn("uppercase", getMethodColor(method))}>
                      {method}
                    </Badge>
                    {details.summary && (
                      <span className="text-sm font-medium">{details.summary}</span>
                    )}
                  </div>
                  
                  {details.description && (
                    <p className="text-sm text-muted-foreground">{details.description}</p>
                  )}
                  
                  {details.parameters && (
                    <ParametersList parameters={details.parameters} />
                  )}
                  
                  {details.requestBody && (
                    <div className="space-y-2">
                      <h4 className="text-sm font-semibold">Request Body</h4>
                      {Object.entries(details.requestBody.content || {}).map(([contentType, content]: [string, any]) => (
                        <div key={contentType} className="text-sm">
                          <code className="text-xs text-muted-foreground">{contentType}</code>
                          {content.schema && (
                            <pre className="mt-1 p-2 bg-muted rounded text-xs overflow-x-auto">
                              {JSON.stringify(content.schema, null, 2)}
                            </pre>
                          )}
                        </div>
                      ))}
                    </div>
                  )}
                  
                  {details.responses && (
                    <ResponsesList responses={details.responses} />
                  )}
                  
                  {method !== Object.keys(methods)[Object.keys(methods).length - 1] && (
                    <div className="border-t pt-4"></div>
                  )}
                </div>
              ))}
            </CardContent>
          </Card>
        ))}
      </div>
    </div>
  );
}