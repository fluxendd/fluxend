import { getAuthToken } from "~/lib/auth";
import fetch, { get, type APIRequestOptions } from "~/tools/fetch";

export const getAllCollections = async (request: any, projectId: string) => {
  const authToken = await getAuthToken(request.headers);
  console.log("Token", authToken);

  const fetchOptions: RequestInit = {
    headers: {
      "X-Project": projectId,
      "Content-Type": "application/json",
      Authorization: `Bearer ${authToken}`,
    },
  };

  const response = get<any>("/tables", fetchOptions);

  return response;
};

export const getCollectionColumn = async (
  request: any,
  projectId: string,
  collectionName: string
) => {
  const authToken = await getAuthToken(request.headers);

  const fetchOptions: RequestInit = {
    headers: {
      "X-Project": projectId,
      "Content-Type": "application/json",
      Authorization: `Bearer ${authToken}`,
    },
  };

  const response = get<any>(
    `/tables/public.${collectionName}/columns`,
    fetchOptions
  );

  return response;
};

export const getCollectionRows = async (
  request: any,
  projectId: string,
  collectionName: string
) => {
  const authToken = await getAuthToken(request.headers);

  const fetchOptions: APIRequestOptions = {
    headers: {
      "X-Project": projectId,
      "Content-Type": "application/json",
      Authorization: `Bearer ${authToken}`,
      ...request?.headers,
    },
    params: {
      ...request?.params,
    },
    baseUrl: request.baseUrl,
  };

  const response = get<any>(collectionName, fetchOptions);

  return response;
};
