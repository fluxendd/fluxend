import { useOutletContext } from "react-router";
import type { Route } from "./+types/overview";
import { FileText, Server, Shield, Package } from "lucide-react";
import { Badge } from "~/components/ui/badge";
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "~/components/ui/card";

interface DocsOutletContext {
  openApiSpec: any;
}

export function meta({}: Route.MetaArgs) {
  return [
    { title: "API Documentation - Fluxend" },
    { name: "description", content: "API documentation overview for your project" },
  ];
}

export default function DocsOverview({ params }: Route.ComponentProps) {
  const { openApiSpec } = useOutletContext<DocsOutletContext>();

  if (!openApiSpec) {
    return (
      <div className="flex h-full items-center justify-center">
        <div className="text-center">
          <p className="text-muted-foreground">No API documentation available</p>
        </div>
      </div>
    );
  }

  // Extract statistics
  const pathCount = Object.keys(openApiSpec.paths || {}).length;
  const methodCounts = Object.entries(openApiSpec.paths || {}).reduce((acc, [path, methods]: [string, any]) => {
    Object.keys(methods).forEach(method => {
      acc[method] = (acc[method] || 0) + 1;
    });
    return acc;
  }, {} as Record<string, number>);

  const tables = new Set<string>();
  Object.keys(openApiSpec.paths || {}).forEach(path => {
    const match = path.match(/^\/([^\/]+)(?:\/|$)/);
    if (match && match[1] !== 'rpc') {
      tables.add(match[1]);
    }
  });

  return (
    <div className="p-8 max-w-6xl mx-auto">
      {/* Header */}
      <div className="mb-8">
        <h1 className="text-3xl font-bold mb-2">{openApiSpec.info?.title || 'API Documentation'}</h1>
        {openApiSpec.info?.description && (
          <p className="text-muted-foreground">{openApiSpec.info.description}</p>
        )}
        <div className="flex items-center gap-2 mt-4">
          <Badge variant="secondary">Version {openApiSpec.info?.version || '1.0.0'}</Badge>
          <Badge variant="secondary">OpenAPI {openApiSpec.openapi || '3.0.0'}</Badge>
        </div>
      </div>

      {/* Server Info */}
      {openApiSpec.servers && openApiSpec.servers.length > 0 && (
        <Card className="mb-6">
          <CardHeader>
            <CardTitle className="flex items-center gap-2">
              <Server className="h-5 w-5" />
              Base URL
            </CardTitle>
          </CardHeader>
          <CardContent>
            {openApiSpec.servers.map((server: any, idx: number) => (
              <div key={idx} className="flex items-center gap-2 mb-2">
                <code className="text-sm bg-muted px-2 py-1 rounded">{server.url}</code>
                {server.description && (
                  <span className="text-sm text-muted-foreground">- {server.description}</span>
                )}
              </div>
            ))}
          </CardContent>
        </Card>
      )}

      {/* Statistics */}
      <div className="grid gap-4 md:grid-cols-2 lg:grid-cols-4 mb-8">
        <Card>
          <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
            <CardTitle className="text-sm font-medium">Total Tables</CardTitle>
            <Package className="h-4 w-4 text-muted-foreground" />
          </CardHeader>
          <CardContent>
            <div className="text-2xl font-bold">{tables.size}</div>
            <p className="text-xs text-muted-foreground">Database tables</p>
          </CardContent>
        </Card>
        <Card>
          <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
            <CardTitle className="text-sm font-medium">Total Endpoints</CardTitle>
            <FileText className="h-4 w-4 text-muted-foreground" />
          </CardHeader>
          <CardContent>
            <div className="text-2xl font-bold">{pathCount}</div>
            <p className="text-xs text-muted-foreground">API endpoints</p>
          </CardContent>
        </Card>
        <Card>
          <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
            <CardTitle className="text-sm font-medium">GET Endpoints</CardTitle>
            <Shield className="h-4 w-4 text-muted-foreground" />
          </CardHeader>
          <CardContent>
            <div className="text-2xl font-bold">{methodCounts.get || 0}</div>
            <p className="text-xs text-muted-foreground">Read operations</p>
          </CardContent>
        </Card>
        <Card>
          <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
            <CardTitle className="text-sm font-medium">POST Endpoints</CardTitle>
            <Shield className="h-4 w-4 text-muted-foreground" />
          </CardHeader>
          <CardContent>
            <div className="text-2xl font-bold">{methodCounts.post || 0}</div>
            <p className="text-xs text-muted-foreground">Create operations</p>
          </CardContent>
        </Card>
      </div>

      {/* Authentication Info */}
      <Card>
        <CardHeader>
          <CardTitle>Authentication</CardTitle>
          <CardDescription>All API requests require authentication</CardDescription>
        </CardHeader>
        <CardContent>
          <div className="space-y-2">
            <div className="flex items-start gap-2">
              <Badge variant="outline">Bearer Token</Badge>
              <div className="text-sm">
                <p>Include your API token in the Authorization header:</p>
                <code className="text-xs bg-muted px-2 py-1 rounded mt-1 inline-block">
                  Authorization: Bearer YOUR_API_TOKEN
                </code>
              </div>
            </div>
            <div className="flex items-start gap-2 mt-4">
              <Badge variant="outline">X-Project Header</Badge>
              <div className="text-sm">
                <p>Include your project ID in the X-Project header:</p>
                <code className="text-xs bg-muted px-2 py-1 rounded mt-1 inline-block">
                  X-Project: YOUR_PROJECT_ID
                </code>
              </div>
            </div>
          </div>
        </CardContent>
      </Card>
    </div>
  );
}