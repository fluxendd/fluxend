import { get, type APIRequestOptions } from "~/tools/fetch";

export type Project = {
  uuid: string;
  name: string;
  status: "active" | "inactive";
  dbName: string;
  organizationUuid: string;
  createdAt: string;
  updatedAt: string;
};

export const createProjectsService = (authToken: string) => {
  const getProjectDetails = async (projectId: string): Promise<Project> => {
    try {
      const fetchOptions: APIRequestOptions = {
        headers: {
          "Content-Type": "application/json",
          ...(authToken && { Authorization: `Bearer ${authToken}` }),
        },
      };

      const response = await get(`/projects/${projectId}`, fetchOptions);
      const result = await response.json();

      if (result.success && result.content) {
        return result.content;
      } else {
        throw new Error(result.errors?.join(", ") || "Unknown error");
      }
    } catch (error) {
      // Fallback to mock data for development
      console.warn("Failed to fetch health status, using mock data:", error);
      throw new Error("Failed to fetch project details");
    }
  };

  return {
    getProjectDetails,
  };
};

export type ProjectService = ReturnType<typeof createProjectsService>;
