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
  sort?: string;
  order?: "asc" | "desc";
}

export function createLogsService(authToken: string) {
  const getLogs = async (projectId: string, filters?: LogsFilters) => {
    // Build headers
    const headers: Record<string, string> = {
      "X-Project": projectId,
      "Content-Type": "application/json",
      Authorization: `Bearer ${authToken}`,
    };

    // Build query params - API expects page, limit, sort, order
    const params: any = { ...filters };
    
    // Ensure limit is set
    if (!params.limit) {
      params.limit = 100; // Default limit
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

    // Since we're using page-based pagination, we don't have exact count
    // You might need to get this from a separate endpoint or response metadata
    return {
      content: data.content || [],
      totalCount: 0, // TODO: Get actual count from API response or separate endpoint
    };
  };

  return {
    getLogs,
  };
}

export type LogsService = ReturnType<typeof createLogsService>;

