import type { APIResponse } from "~/lib/types";
import { getTypedResponseData } from "~/lib/utils";
import { get } from "~/tools/fetch";

export type Organization = {
  uuid: string;
  name: string;
};

export type Project = {
  uuid: string;
  name: string;
  status: "active" | "inactive";
  dbName: string;
  organizationUuid: string;
  createdAt: string;
  updatedAt: string;
};

export const createUserService = (authToken: string) => {
  const getUserOrganizations = async () => {
    const fetchOptions: RequestInit = {
      headers: {
        "Content-Type": "application/json",
        Authorization: `Bearer ${authToken}`,
      },
    };

    const response = await get("/organizations", fetchOptions);
    const data = await getTypedResponseData<APIResponse<Organization[]>>(
      response
    );

    return data;
  };

  const getUserProjects = async (organizationId: string) => {
    const fetchOptions: RequestInit = {
      headers: {
        "Content-Type": "application/json",
        Authorization: `Bearer ${authToken}`,
      },
    };

    const response = await get(
      `/projects?${new URLSearchParams({ organization_uuid: organizationId })}`,
      fetchOptions
    );

    const data = await getTypedResponseData<APIResponse<Project[]>>(response);

    return data;
  };

  return {
    getUserOrganizations,
    getUserProjects,
  };
};
