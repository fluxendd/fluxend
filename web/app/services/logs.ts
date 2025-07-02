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
  body: string;
  params: string;
  createdAt: string;
}

export interface LogsMetadata {
  limit: number;
  page: number;
  total: number;
}

export interface LogsApiResponse {
  success: boolean;
  content: LogEntry[];
  metadata: LogsMetadata;
}

export interface LogsResponse {
  content: LogEntry[];
  metadata: LogsMetadata;
  hasMore: boolean;
}

export interface LogsFilters {
  userUuid?: string;
  status?: string;
  method?: string;
  endpoint?: string;
  ipAddress?: string;
  dateStart?: string;
  dateEnd?: string;
  startTime?: number; // Unix timestamp
  endTime?: number; // Unix timestamp
  page?: number;
  limit?: number;
  sort?: string;
  order?: "asc" | "desc";
}

export function createLogsService(authToken: string) {
  const getLogs = async (projectId: string, filters?: LogsFilters) => {
    // Build headers
    const headers: Record<string, string> = {
      "Content-Type": "application/json",
      Authorization: `Bearer ${authToken}`,
    };

    // Build query params with proper typing
    const params: Record<string, string | number | undefined> = filters ? {
      ...filters,
      limit: filters.limit || 100, // Default limit
      // Convert numbers to strings for query params
      ...(filters.startTime && { startTime: filters.startTime.toString() }),
      ...(filters.endTime && { endTime: filters.endTime.toString() }),
    } : {
      limit: 100,
    };

    const fetchOptions: APIRequestOptions = {
      headers,
      params,
    };

    const response = await get(`/projects/${projectId}/logs`, fetchOptions);

    if (!response.ok) {
      if (response.status >= 500) {
        throw new Error("Server error occurred while fetching logs");
      }
      throw new Error(`Failed to fetch logs: ${response.status} ${response.statusText}`);
    }

    const data = await getTypedResponseData<LogsApiResponse>(response);

    if (!data.success) {
      throw new Error("Invalid response format from logs API");
    }

    // Calculate if there are more pages
    const currentPage = data.metadata.page;
    const totalPages = Math.ceil(data.metadata.total / data.metadata.limit);
    const hasMore = currentPage < totalPages;

    return {
      content: data.content || [],
      metadata: data.metadata,
      hasMore,
    };
  };

  return {
    getLogs,
  };
}

export type LogsService = ReturnType<typeof createLogsService>;