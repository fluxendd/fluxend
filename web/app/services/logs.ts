import type { APIResponse } from "~/lib/types";
import { getTypedResponseData } from "~/lib/utils";
import { get, type APIRequestOptions } from "~/tools/fetch";

export interface LogEntry {
  uuid: string;
  userUuid: string;
  method: string;
  endpoint: string;
  status: number;
  ipAddress: string;
  userAgent: string;
  body: any;
  params: any;
  createdAt: string;
}

export interface LogsResponse {
  content: LogEntry[];
  totalCount: number;
}

export interface LogsFilters {
  userUuid?: string;
  status?: string;
  method?: string;
  endpoint?: string;
  ipAddress?: string;
  page?: number;
  limit?: number;
  offset?: number;
  sort?: string;
  order?: 'asc' | 'desc' | string; // Can be 'asc'/'desc' or PostgREST format "column.desc"
}

export function createLogsService(authToken: string) {
  const getLogs = async (projectId: string, filters?: LogsFilters) => {
    // Build headers with PostgREST Prefer header for count
    const headers: Record<string, string> = {
      "X-Project": projectId,
      "Content-Type": "application/json",
      Authorization: `Bearer ${authToken}`,
      Prefer: "count=exact", // Request exact count from PostgREST
    };

    // Build query params
    const params: any = { ...filters };
    
    // Convert sort/order to PostgREST format
    if (params.sort && params.order) {
      // PostgREST expects format like "createdAt.desc" or "createdAt.asc"
      params.order = `${params.sort}.${params.order}`;
      delete params.sort;
    }
    
    // If using offset/limit, try Range header approach
    // But keep the params as fallback in case server doesn't support Range headers
    if (filters?.offset !== undefined && filters?.limit !== undefined) {
      const start = filters.offset;
      const end = filters.offset + filters.limit - 1;
      headers["Range-Unit"] = "items";
      headers["Range"] = `${start}-${end}`;
    }

    const fetchOptions: APIRequestOptions = {
      headers,
      params,
    };

    const response = await get("/admin/logs", fetchOptions);
    
    if (!response.ok && response.status > 500) {
      throw new Error("Unexpected Error");
    }

    const data = await getTypedResponseData<APIResponse<LogEntry[]>>(response);

    if (!response.ok) {
      throw new Error(data.errors?.[0] || "Unexpected Error");
    }

    // Extract total count from Content-Range header
    // PostgREST format: "start-end/total" or "start-end/*"
    let totalCount = 0;
    const contentRange = response.headers.get("Content-Range");
    if (contentRange) {
      // Try to match the full format first
      const fullMatch = contentRange.match(/(\d+)-(\d+)\/(\d+|\*)/);
      if (fullMatch) {
        const [, start, end, total] = fullMatch;
        if (total !== "*") {
          totalCount = parseInt(total, 10);
        }
      } else {
        // Fallback to simpler format
        const simpleMatch = contentRange.match(/\/(\d+|\*)/);
        if (simpleMatch && simpleMatch[1] !== "*") {
          totalCount = parseInt(simpleMatch[1], 10);
        }
      }
    }

    return {
      content: data.content || [],
      totalCount,
    };
  };

  return {
    getLogs,
  };
}

export type LogsService = ReturnType<typeof createLogsService>;