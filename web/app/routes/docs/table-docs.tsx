import { useOutletContext, useNavigate, useSearchParams } from "react-router";
import type { Route } from "./+types/table-docs";
import { Badge } from "~/components/ui/badge";
import { Button } from "~/components/ui/button";
import { cn } from "~/lib/utils";
import { ChevronRight, ChevronDown, Database, Play } from "lucide-react";
import type { ProjectLayoutOutletContext } from "~/components/shared/project-layout";
import { ScrollArea } from "~/components/ui/scroll-area";
import { useState, useEffect, useMemo } from "react";
import {
  Collapsible,
  CollapsibleContent,
  CollapsibleTrigger,
} from "~/components/ui/collapsible";

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
  isSelected: boolean;
  onSelect: () => void;
}

const EndpointDetails = ({ path, method, details, isSelected, onSelect }: EndpointDetailsProps) => {
  // Group parameters by location
  const headers = details.parameters?.filter((p: Parameter) => p.in === 'header') || [];
  const pathParams = details.parameters?.filter((p: Parameter) => p.in === 'path') || [];
  const queryParams = details.parameters?.filter((p: Parameter) => p.in === 'query') || [];

  if (!isSelected) {
    return null;
  }

  return (
    <div className="mt-6 space-y-6">
      {/* Endpoint Header */}
      <div className="flex items-center justify-between">
        <div className="flex-1">
          <h3 className="text-xl font-semibold mb-2">{details.summary || `${method.toUpperCase()} ${path}`}</h3>
          {details.description && (
            <p className="text-muted-foreground mb-4">{details.description}</p>
          )}
          <code className="text-sm bg-muted px-3 py-1.5 rounded font-mono inline-block">
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
        <Button className="bg-green-500 hover:bg-green-600 text-white">
          Try it
          <Play className="ml-2 h-4 w-4" />
        </Button>
      </div>

      {/* Headers Section */}
      {headers.length > 0 && (
        <div className="space-y-4">
          <h4 className="text-lg font-semibold">Headers</h4>
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
        <div className="space-y-4">
          <h4 className="text-lg font-semibold">Path Parameters</h4>
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
        <div className="space-y-4">
          <h4 className="text-lg font-semibold">Query Parameters</h4>
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
        <div className="space-y-4">
          <h4 className="text-lg font-semibold">Request Body</h4>
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
        <div className="space-y-4">
          <h4 className="text-lg font-semibold">Responses</h4>
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

const MethodBadge = ({ method, summary }: { method: string; summary?: string }) => {
  const methodColors = {
    get: "bg-blue-500 text-white",
    post: "bg-green-500 text-white",
    put: "bg-yellow-500 text-white",
    patch: "bg-orange-500 text-white",
    delete: "bg-red-500 text-white",
  };

  return (
    <div className="flex items-center gap-3">
      <Badge 
        className={cn("uppercase font-mono text-xs px-2 py-0.5", methodColors[method.toLowerCase() as keyof typeof methodColors] || "bg-gray-500 text-white")}
      >
        {method}
      </Badge>
      <span className="text-sm">{summary || `${method.charAt(0).toUpperCase() + method.slice(1)} ${method === 'get' ? 'records' : method === 'post' ? 'record' : method === 'put' ? 'record' : method === 'delete' ? 'record' : 'operation'}`}</span>
    </div>
  );
};

export default function TableDocs({ params }: Route.ComponentProps) {
  const { projectDetails, openApiSpec } = useOutletContext<DocsOutletContext>();
  const { table } = params;
  const [searchParams, setSearchParams] = useSearchParams();
  const navigate = useNavigate();
  
  const [expandedTable, setExpandedTable] = useState(true);
  const [selectedMethod, setSelectedMethod] = useState<string | null>(null);

  if (!openApiSpec) {
    return (
      <div className="flex h-full items-center justify-center">
        <div className="text-center">
          <p className="text-muted-foreground">No API documentation available</p>
        </div>
      </div>
    );
  }

  // Group paths by method for this table
  const tableEndpoints = useMemo(() => {
    const endpoints: Record<string, any> = {};
    
    Object.entries(openApiSpec.paths || {}).forEach(([path, methods]: [string, any]) => {
      if (path.startsWith(`/${table}`)) {
        Object.entries(methods).forEach(([method, details]) => {
          if (!endpoints[method]) {
            endpoints[method] = [];
          }
          endpoints[method].push({ path, details });
        });
      }
    });
    
    return endpoints;
  }, [openApiSpec, table]);

  // Get the currently selected method from URL or set the first one
  useEffect(() => {
    const methodFromUrl = searchParams.get('method');
    const methods = Object.keys(tableEndpoints);
    
    if (methodFromUrl && methods.includes(methodFromUrl)) {
      setSelectedMethod(methodFromUrl);
    } else if (methods.length > 0 && !selectedMethod) {
      // Auto-select the first method
      const firstMethod = methods[0];
      setSelectedMethod(firstMethod);
      setSearchParams({ method: firstMethod });
    }
  }, [tableEndpoints, searchParams, selectedMethod, setSearchParams]);

  const handleMethodSelect = (method: string) => {
    setSelectedMethod(method);
    setSearchParams({ method });
  };

  if (Object.keys(tableEndpoints).length === 0) {
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

        {/* Table Accordion */}
        <div className="bg-background border rounded-lg overflow-hidden">
          <Collapsible open={expandedTable} onOpenChange={setExpandedTable}>
            <CollapsibleTrigger className="w-full p-4 hover:bg-muted/50 transition-colors">
              <div className="flex items-center justify-between">
                <div className="flex items-center gap-2">
                  {expandedTable ? <ChevronDown className="h-4 w-4" /> : <ChevronRight className="h-4 w-4" />}
                  <span className="font-semibold text-lg capitalize">{table}</span>
                </div>
                <span className="text-sm text-muted-foreground">{Object.keys(tableEndpoints).length} endpoints</span>
              </div>
            </CollapsibleTrigger>
            
            <CollapsibleContent>
              <div className="border-t">
                {/* Method List */}
                <div className="p-4 space-y-2">
                  {Object.entries(tableEndpoints).map(([method, endpoints]) => {
                    const firstEndpoint = endpoints[0];
                    const isSelected = selectedMethod === method;
                    
                    return (
                      <button
                        key={method}
                        onClick={() => handleMethodSelect(method)}
                        className={cn(
                          "w-full text-left p-3 rounded-lg transition-colors",
                          isSelected 
                            ? "bg-muted border border-border" 
                            : "hover:bg-muted/50"
                        )}
                      >
                        <MethodBadge 
                          method={method} 
                          summary={firstEndpoint.details.summary}
                        />
                      </button>
                    );
                  })}
                </div>

                {/* Selected Method Details */}
                {selectedMethod && tableEndpoints[selectedMethod] && (
                  <div className="border-t p-6">
                    {tableEndpoints[selectedMethod].map((endpoint: any, index: number) => (
                      <EndpointDetails
                        key={`${selectedMethod}-${endpoint.path}-${index}`}
                        path={endpoint.path}
                        method={selectedMethod}
                        details={{ ...endpoint.details, table }}
                        isSelected={true}
                        onSelect={() => {}}
                      />
                    ))}
                  </div>
                )}
              </div>
            </CollapsibleContent>
          </Collapsible>
        </div>
      </div>
    </ScrollArea>
  );
}