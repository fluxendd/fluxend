import { useOutletContext } from "react-router";
import type { Route } from "./+types/overview";
import { FileText, Server, Shield, Package } from "lucide-react";
import { Badge } from "~/components/ui/badge";

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
        <div className="bg-background border rounded-lg mb-6">
          <div className="p-6 border-b">
            <h3 className="text-lg font-semibold flex items-center gap-2">
              <Server className="h-5 w-5" />
              Base URL
            </h3>
          </div>
          <div className="p-6">
            {openApiSpec.servers.map((server: any, idx: number) => (
              <div key={idx} className="flex items-center gap-2 mb-2">
                <code className="text-sm bg-muted px-3 py-1.5 rounded font-mono">{server.url}</code>
                {server.description && (
                  <span className="text-sm text-muted-foreground">- {server.description}</span>
                )}
              </div>
            ))}
          </div>
        </div>
      )}

      {/* Statistics */}
      <div className="grid gap-4 md:grid-cols-2 lg:grid-cols-4 mb-8">
        <div className="bg-background border rounded-lg p-6">
          <div className="flex flex-row items-center justify-between space-y-0 pb-2">
            <h3 className="text-sm font-medium text-muted-foreground">Total Tables</h3>
            <Package className="h-4 w-4 text-muted-foreground" />
          </div>
          <div className="mt-2">
            <div className="text-2xl font-bold">{tables.size}</div>
            <p className="text-xs text-muted-foreground">Database tables</p>
          </div>
        </div>
        <div className="bg-background border rounded-lg p-6">
          <div className="flex flex-row items-center justify-between space-y-0 pb-2">
            <h3 className="text-sm font-medium text-muted-foreground">Total Endpoints</h3>
            <FileText className="h-4 w-4 text-muted-foreground" />
          </div>
          <div className="mt-2">
            <div className="text-2xl font-bold">{pathCount}</div>
            <p className="text-xs text-muted-foreground">API endpoints</p>
          </div>
        </div>
        <div className="bg-background border rounded-lg p-6">
          <div className="flex flex-row items-center justify-between space-y-0 pb-2">
            <h3 className="text-sm font-medium text-muted-foreground">GET Endpoints</h3>
            <Shield className="h-4 w-4 text-muted-foreground" />
          </div>
          <div className="mt-2">
            <div className="text-2xl font-bold">{methodCounts.get || 0}</div>
            <p className="text-xs text-muted-foreground">Read operations</p>
          </div>
        </div>
        <div className="bg-background border rounded-lg p-6">
          <div className="flex flex-row items-center justify-between space-y-0 pb-2">
            <h3 className="text-sm font-medium text-muted-foreground">POST Endpoints</h3>
            <Shield className="h-4 w-4 text-muted-foreground" />
          </div>
          <div className="mt-2">
            <div className="text-2xl font-bold">{methodCounts.post || 0}</div>
            <p className="text-xs text-muted-foreground">Create operations</p>
          </div>
        </div>
      </div>

      {/* Authentication Info */}
      <div className="bg-background border rounded-lg">
        <div className="p-6 border-b">
          <h3 className="text-lg font-semibold">Authentication</h3>
          <p className="text-sm text-muted-foreground mt-1">All API requests require authentication</p>
        </div>
        <div className="p-6">
          <div className="space-y-6">
            <div className="space-y-2">
              <div className="flex items-center gap-3">
                <code className="text-green-400 font-mono">Authorization</code>
                <span className="text-sm text-muted-foreground">string</span>
                <span className="text-xs text-red-400">required</span>
              </div>
              <p className="text-sm text-muted-foreground">Bearer Token</p>
              <code className="text-sm bg-muted px-3 py-1.5 rounded inline-block font-mono">
                Authorization: Bearer YOUR_API_TOKEN
              </code>
            </div>
            <div className="space-y-2">
              <div className="flex items-center gap-3">
                <code className="text-green-400 font-mono">X-Project</code>
                <span className="text-sm text-muted-foreground">string</span>
                <span className="text-xs text-red-400">required</span>
              </div>
              <p className="text-sm text-muted-foreground">Project UUID</p>
              <code className="text-sm bg-muted px-3 py-1.5 rounded inline-block font-mono">
                X-Project: YOUR_PROJECT_UUID
              </code>
            </div>
          </div>
        </div>
      </div>
    </div>
  );
}