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
    const fetchOptions: APIRequestOptions = {
      headers: {
        "X-Project": projectId,
        Authorization: `Bearer ${authToken}`,
      },
    };

    const response = await get("/containers", fetchOptions);
    const data = await getTypedResponseData<StorageListResponse<StorageContainer>>(response);
    
    return data;
  };

  const createContainer = async (projectId: string, container: CreateContainerRequest) => {
    const fetchOptions: APIRequestOptions = {
      headers: {
        "X-Project": projectId,
        Authorization: `Bearer ${authToken}`,
      },
    };

    const response = await post("/containers", container, fetchOptions);
    const data = await getTypedResponseData<StorageItemResponse<StorageContainer>>(response);
    
    return data;
  };

  const getContainer = async (projectId: string, containerUuid: string) => {
    const fetchOptions: APIRequestOptions = {
      headers: {
        "X-Project": projectId,
        Authorization: `Bearer ${authToken}`,
      },
    };

    const response = await get(`/containers/${containerUuid}`, fetchOptions);
    const data = await getTypedResponseData<StorageItemResponse<StorageContainer>>(response);
    
    return data;
  };

  const updateContainer = async (
    projectId: string,
    containerUuid: string,
    updates: UpdateContainerRequest
  ) => {
    const fetchOptions: APIRequestOptions = {
      headers: {
        "X-Project": projectId,
        Authorization: `Bearer ${authToken}`,
      },
    };

    const response = await put(`/storage/containers/${containerUuid}`, updates, fetchOptions);
    const data = await getTypedResponseData<StorageItemResponse<StorageContainer>>(response);
    
    return data;
  };

  const deleteContainer = async (projectId: string, containerUuid: string) => {
    const fetchOptions: APIRequestOptions = {
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
    const fetchOptions: APIRequestOptions = {
      headers: {
        "X-Project": projectId,
        Authorization: `Bearer ${authToken}`,
      },
    };

    const response = await post(`/containers/${containerUuid}/files`, fileData, fetchOptions);
    const data = await getTypedResponseData<StorageItemResponse<StorageFile>>(response);
    
    return data;
  };

  const uploadFile = async (
    projectId: string,
    containerUuid: string,
    file: File
  ) => {
    const formData = new FormData();
    formData.append('file', file);
    formData.append('projectUUID', projectId);
    formData.append('full_file_name', file.name);

    const fetchOptions: APIRequestOptions = {
      headers: {
        "X-Project": projectId,
        Authorization: `Bearer ${authToken}`,
      },
    };

    const response = await post(`/containers/${containerUuid}/files`, formData, fetchOptions);
    const data = await getTypedResponseData<StorageItemResponse<StorageFile>>(response);
    
    return data;
  };

  const getFile = async (
    projectId: string,
    containerUuid: string,
    fileUuid: string
  ) => {
    const fetchOptions: APIRequestOptions = {
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

  const downloadFile = async (
    projectId: string,
    containerUuid: string,
    fileUuid: string,
    fileName: string
  ) => {
    const fetchOptions: APIRequestOptions = {
      headers: {
        "X-Project": projectId,
        Authorization: `Bearer ${authToken}`,
      },
    };

    const response = await get(
      `/containers/${containerUuid}/files/${fileUuid}/download`,
      fetchOptions
    );

    if (response.ok) {
      const blob = await response.blob();
      const url = window.URL.createObjectURL(blob);
      const a = document.createElement('a');
      a.href = url;
      a.download = fileName;
      document.body.appendChild(a);
      a.click();
      document.body.removeChild(a);
      window.URL.revokeObjectURL(url);
    }

    return response;
  };

  const deleteFile = async (
    projectId: string,
    containerUuid: string,
    fileUuid: string
  ) => {
    const fetchOptions: APIRequestOptions = {
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
    const fetchOptions: APIRequestOptions = {
      headers: {
        "X-Project": projectId,
        Authorization: `Bearer ${authToken}`,
      },
    };

    const response = await put(
      `/containers/${containerUuid}/files/${fileUuid}`,
      renameData,
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
    uploadFile,
    getFile,
    downloadFile,
    deleteFile,
    renameFile,
  };
}