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
  order?: 'asc' | 'desc';
}

export function createLogsService(authToken: string) {
  const getLogs = async (projectId: string, filters?: LogsFilters) => {
    const fetchOptions: APIRequestOptions = {
      headers: {
        "X-Project": projectId,
        "Content-Type": "application/json",
        Authorization: `Bearer ${authToken}`,
      },
      params: filters,
    };

    const response = await get("/admin/logs", fetchOptions);
    
    if (!response.ok && response.status > 500) {
      throw new Error("Unexpected Error");
    }

    const data = await getTypedResponseData<APIResponse<LogEntry[]>>(response);

    if (!response.ok) {
      throw new Error(data.errors?.[0] || "Unexpected Error");
    }

    // Extract total count from response headers if available
    let totalCount = 0;
    const contentRange = response.headers.get("Content-Range");
    if (contentRange) {
      const match = contentRange.match(/\/(\d+|\*)/);
      if (match && match[1] && match[1] !== "*") {
        totalCount = parseInt(match[1], 10);
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