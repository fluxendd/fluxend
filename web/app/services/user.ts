import type { APIResponse } from "~/lib/types";
import { getTypedResponseData } from "~/lib/utils";
import { get, post } from "~/tools/fetch";

export type Organization = {
  uuid: string;
  name: string;
};

export type Project = {
  uuid: string;
  name: string;
  description: string;
  status: "active" | "inactive";
  dbName: string;
  organizationUuid: string;
  createdAt: string;
  updatedAt: string;
};

export type User = {
  bio: string;
  createdAt: string;
  email: string;
  organizationUuid: string;
  roleId: number;
  status: string;
  updatedAt: string;
  username: string;
  uuid: string;
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

  const createUserProject = async (name: string, description: string, organizationId: string) => {
    const fetchOptions: RequestInit = {
      headers: {
        "Content-Type": "application/json",
        Authorization: `Bearer ${authToken}`,
      },
    };

    const response = await post(
      `/projects`,
      { name: name, description: description, organization_uuid: organizationId },
      fetchOptions
    );

    const data = await getTypedResponseData<APIResponse<Project>>(response);
    return data;
  };

  const getCurrentUser = async () => {
    const fetchOptions: RequestInit = {
      headers: {
        "Content-Type": "application/json",
        Authorization: `Bearer ${authToken}`,
      },
    };

    const response = await get(`/users/me`, fetchOptions);

    const data = await getTypedResponseData<APIResponse<User>>(response);
    return data;
  };

  return {
    getUserOrganizations,
    getUserProjects,
    createUserProject,
    getCurrentUser,
  };
};
