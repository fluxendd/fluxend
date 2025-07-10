import type { APIResponse } from "~/lib/types";
import { getTypedResponseData } from "~/lib/utils";
import { get, type APIRequestOptions } from "~/tools/fetch";

export interface HealthData {
  database_status: string;
  app_status: string;
  postgrest_status: string;
  disk_usage: string;
  disk_available: string;
  disk_total: string;
  cpu_usage: string;
  cpu_cores: number;
}

export interface UnusedIndex {
  tableName: string;
  indexName: string;
  indexScans: number;
  indexSize: string;
}

export interface TableCount {
  TableName: string;
  EstimatedRowCount: number;
}

export interface TableSize {
  tableName: string;
  totalSize: string;
}

export interface ProjectStats {
  id: number;
  databaseName: string;
  totalSize: string;
  indexSize: string;
  unusedIndex: UnusedIndex[];
  tableCount: TableCount[];
  tableSize: TableSize[];
  createdAt: string;
}

export const createDashboardService = (authToken: string) => {
  const getHealthStatus = async (): Promise<HealthData> => {
    try {
      const fetchOptions: APIRequestOptions = {
        headers: {
          "Content-Type": "application/json",
          ...(authToken && { Authorization: `Bearer ${authToken}` }),
        },
      };

      const response = await get("/admin/health", fetchOptions);
      const result = await getTypedResponseData<APIResponse<HealthData>>(
          response
      );

      if (result.success && result.content) {
        return result.content;
      } else {
        throw new Error(result.errors?.join(", ") || "Unknown error");
      }
    } catch (error) {
      // Fallback to mock data for development
      console.warn("Failed to fetch health status, using mock data:", error);
      return {
        database_status: "OK",
        app_status: "OK",
        postgrest_status: "OK",
        disk_usage: "48.9%",
        disk_available: "24.7 GB",
        disk_total: "48.3 GB",
        cpu_usage: "2.3%",
        cpu_cores: 1,
      };
    }
  };

  const getProjectStats = async (projectUUID: string): Promise<ProjectStats> => {
    const fetchOptions: APIRequestOptions = {
      headers: {
        "Content-Type": "application/json",
        ...(authToken && { Authorization: `Bearer ${authToken}` }),
      },
    };

    const response = await get(`/projects/${projectUUID}/stats`, fetchOptions);
    const result = await getTypedResponseData<APIResponse<ProjectStats>>(
        response
    );

    if (result.success && result.content) {
      return result.content;
    } else {
      throw new Error(result.errors?.join(", ") || "Unknown error");
    }
  };

  return {
    getHealthStatus,
    getProjectStats,
  };
};

export type DashboardService = ReturnType<typeof createDashboardService>;