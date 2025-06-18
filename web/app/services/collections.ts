import type { APIResponse } from "~/lib/types";
import { getTypedResponseData } from "~/lib/utils";
import type { Collection } from "~/routes/collections/collection-list";
import { get, post, del, type APIRequestOptions } from "~/tools/fetch";

export function createCollectionsService(authToken: string) {
  const getAllCollections = async (projectId: string) => {
    const fetchOptions: RequestInit = {
      headers: {
        "X-Project": projectId,
        "Content-Type": "application/json",
        Authorization: `Bearer ${authToken}`,
      },
    };

    const response = await get("/tables", fetchOptions);
    const data = await getTypedResponseData<APIResponse<any>>(response);

    return data;
  };

  const getCollectionColumns = async (
    projectId: string,
    collectionName: string
  ) => {
    const fetchOptions: RequestInit = {
      headers: {
        "X-Project": projectId,
        "Content-Type": "application/json",
        Authorization: `Bearer ${authToken}`,
      },
    };

    return get(`/tables/public.${collectionName}/columns`, fetchOptions);
  };

  const getCollectionRows = async (
    projectId: string,
    collectionName: string,
    options?: {
      headers?: HeadersInit;
      params?: Record<string, any>;
      baseUrl?: string;
    }
  ) => {
    const fetchOptions: APIRequestOptions = {
      headers: {
        "X-Project": projectId,
        "Content-Type": "application/json",
        Authorization: `Bearer ${authToken}`,
        ...options?.headers,
      },
      params: options?.params,
      baseUrl: options?.baseUrl,
    };

    return get(collectionName, fetchOptions);
  };

  const createCollection = async (projectId: string, data: any) => {
    const fetchOptions: RequestInit = {
      headers: {
        "X-Project": projectId,
        "Content-Type": "application/json",
        Authorization: `Bearer ${authToken}`,
      },
    };

    return post("/tables", data, fetchOptions);
  };

  const deleteCollection = async (projectId: string, tableName: string) => {
    const fetchOptions: RequestInit = {
      headers: {
        "X-Project": projectId,
        "Content-Type": "application/json",
        Authorization: `Bearer ${authToken}`,
      },
    };

    const response = await del(`/tables/public.${tableName}`, fetchOptions);
    const data = await getTypedResponseData<APIResponse<null>>(response);
    return data;
  };

  return {
    getAllCollections,
    getCollectionColumns,
    getCollectionRows,
    createCollection,
    deleteCollection,
  };
}

export type CollectionsService = ReturnType<typeof createCollectionsService>;
