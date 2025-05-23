import { getAuthToken } from "~/lib/auth";
import { get } from "~/tools/fetch";

export const getUserOrganizations = async (request: any) => {
  const fetchOptions: RequestInit = {
    headers: {
      "Content-Type": "application/json",
      ...request?.headers,
    },
  };

  const response = get<any>("/organizations", fetchOptions);

  return response;
};

export const getUserProjects = async (request: any, organizationId: string) => {
  const fetchOptions: RequestInit = {
    headers: {
      "Content-Type": "application/json",
      ...request?.headers,
    },
  };

  const response = get<any>(
    `/projects?${new URLSearchParams({ organization_uuid: organizationId })}`,
    fetchOptions
  );

  return response;
};
