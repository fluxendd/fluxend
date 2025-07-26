import type { APIResponse } from "~/lib/types";
import { getTypedResponseData } from "~/lib/utils";
import { get, post, del, put, type APIRequestOptions } from "~/tools/fetch";
import type {
  StorageContainer,
  StorageFile,
  CreateContainerRequest,
  UpdateContainerRequest,
  CreateFileRequest,
  RenameFileRequest,
  StorageListResponse,
  StorageItemResponse,
  FileListParams,
} from "~/types/storage";

export function createStorageService(authToken: string) {
  // Container operations
  const listContainers = async (projectId: string) => {
    const fetchOptions: RequestInit = {
      headers: {
        "X-Project": projectId,
        Authorization: `Bearer ${authToken}`,
      },
    };

    const response = await get("/storage", fetchOptions);
    const data = await getTypedResponseData<StorageListResponse<StorageContainer>>(response);
    
    return data;
  };

  const createContainer = async (projectId: string, container: CreateContainerRequest) => {
    const fetchOptions: RequestInit = {
      headers: {
        "X-Project": projectId,
        "Content-Type": "application/json",
        Authorization: `Bearer ${authToken}`,
      },
      body: JSON.stringify(container),
    };

    const response = await post("/storage", fetchOptions);
    const data = await getTypedResponseData<StorageItemResponse<StorageContainer>>(response);
    
    return data;
  };

  const getContainer = async (projectId: string, containerUuid: string) => {
    const fetchOptions: RequestInit = {
      headers: {
        "X-Project": projectId,
        Authorization: `Bearer ${authToken}`,
      },
    };

    const response = await get(`/storage/containers/${containerUuid}`, fetchOptions);
    const data = await getTypedResponseData<StorageItemResponse<StorageContainer>>(response);
    
    return data;
  };

  const updateContainer = async (
    projectId: string,
    containerUuid: string,
    updates: UpdateContainerRequest
  ) => {
    const fetchOptions: RequestInit = {
      headers: {
        "X-Project": projectId,
        "Content-Type": "application/json",
        Authorization: `Bearer ${authToken}`,
      },
      body: JSON.stringify(updates),
    };

    const response = await put(`/storage/containers/${containerUuid}`, fetchOptions);
    const data = await getTypedResponseData<StorageItemResponse<StorageContainer>>(response);
    
    return data;
  };

  const deleteContainer = async (projectId: string, containerUuid: string) => {
    const fetchOptions: RequestInit = {
      headers: {
        "X-Project": projectId,
        Authorization: `Bearer ${authToken}`,
      },
    };

    const response = await del(`/storage/containers/${containerUuid}`, fetchOptions);
    return response;
  };

  // File operations
  const listFiles = async (
    projectId: string,
    containerUuid: string,
    params?: FileListParams
  ) => {
    const fetchOptions: APIRequestOptions = {
      headers: {
        "X-Project": projectId,
        Authorization: `Bearer ${authToken}`,
      },
      params,
    };

    const response = await get(`/containers/${containerUuid}/files`, fetchOptions);
    const data = await getTypedResponseData<StorageListResponse<StorageFile>>(response);
    
    return data;
  };

  const createFile = async (
    projectId: string,
    containerUuid: string,
    fileData: CreateFileRequest
  ) => {
    const fetchOptions: RequestInit = {
      headers: {
        "X-Project": projectId,
        "Content-Type": "application/json",
        Authorization: `Bearer ${authToken}`,
      },
      body: JSON.stringify(fileData),
    };

    const response = await post(`/containers/${containerUuid}/files`, fetchOptions);
    const data = await getTypedResponseData<StorageItemResponse<StorageFile>>(response);
    
    return data;
  };

  const getFile = async (
    projectId: string,
    containerUuid: string,
    fileUuid: string
  ) => {
    const fetchOptions: RequestInit = {
      headers: {
        "X-Project": projectId,
        Authorization: `Bearer ${authToken}`,
      },
    };

    const response = await get(
      `/containers/${containerUuid}/files/${fileUuid}`,
      fetchOptions
    );
    const data = await getTypedResponseData<StorageItemResponse<StorageFile>>(response);
    
    return data;
  };

  const deleteFile = async (
    projectId: string,
    containerUuid: string,
    fileUuid: string
  ) => {
    const fetchOptions: RequestInit = {
      headers: {
        "X-Project": projectId,
        Authorization: `Bearer ${authToken}`,
      },
    };

    const response = await del(
      `/containers/${containerUuid}/files/${fileUuid}`,
      fetchOptions
    );
    return response;
  };

  const renameFile = async (
    projectId: string,
    containerUuid: string,
    fileUuid: string,
    renameData: RenameFileRequest
  ) => {
    const fetchOptions: RequestInit = {
      headers: {
        "X-Project": projectId,
        "Content-Type": "application/json",
        Authorization: `Bearer ${authToken}`,
      },
      body: JSON.stringify(renameData),
    };

    const response = await put(
      `/containers/${containerUuid}/files/${fileUuid}/rename`,
      fetchOptions
    );
    const data = await getTypedResponseData<StorageItemResponse<StorageFile>>(response);
    
    return data;
  };

  return {
    // Container operations
    listContainers,
    createContainer,
    getContainer,
    updateContainer,
    deleteContainer,
    // File operations
    listFiles,
    createFile,
    getFile,
    deleteFile,
    renameFile,
  };
}