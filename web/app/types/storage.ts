export interface StorageContainer {
  createdAt: string;
  createdBy: string;
  description: string;
  isPublic: boolean;
  maxFileSize: number;
  name: string;
  projectUuid: string;
  totalFiles: number;
  updatedAt: string;
  updatedBy: string;
  url: string;
  uuid: string;
}

export interface StorageFile {
  containerUuid: string;
  createdAt: string;
  createdBy: string;
  fullFileName: string;
  mimeType: string;
  size: number;
  updatedAt: string;
  updatedBy: string;
  uuid: string;
}

export interface CreateContainerRequest {
  description: string;
  is_public: boolean;
  max_file_size: number;
  name: string;
  projectUUID: string;
}

export interface UpdateContainerRequest {
  description?: string;
  is_public?: boolean;
  max_file_size?: number;
  name?: string;
  projectUUID?: string;
}

export interface CreateFileRequest {
  projectUUID: string;
}

export interface RenameFileRequest {
  full_file_name: string;
  projectUUID: string;
}

export interface StorageListResponse<T> {
  content: T[];
  errors: string[];
  metadata: any;
  success: boolean;
}

export interface StorageItemResponse<T> {
  content: T;
  errors: string[];
  metadata: any;
  success: boolean;
}

export interface FileListParams {
  page?: number;
  limit?: number;
  sort?: string;
  order?: 'asc' | 'desc';
}