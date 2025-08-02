import type { APIResponse } from "~/lib/types";
import { getTypedResponseData } from "~/lib/utils";
import { get } from "~/tools/fetch";

export interface OpenAPIResponse {
  content: any; // Can be either string or already parsed object
  errors: string[];
  metadata: Record<string, any>;
  success: boolean;
}

export const createOpenAPIService = (authToken: string) => {
  const getProjectOpenAPI = async (projectUUID: string, table?: string) => {
    const fetchOptions: RequestInit = {
      headers: {
        "Content-Type": "application/json",
        Authorization: `Bearer ${authToken}`,
      },
    };

    const params = new URLSearchParams();
    if (table) {
      params.append("tables", table);
    }

    const url = `/projects/${projectUUID}/openapi${params.toString() ? `?${params.toString()}` : ""}`;
    const response = await get(url, fetchOptions);
    const data = await getTypedResponseData<OpenAPIResponse>(response);

    return data;
  };

  return {
    getProjectOpenAPI,
  };
};