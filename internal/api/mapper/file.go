package mapper

import (
	fileDto "fluxend/internal/api/dto/storage/file"
	fileDomain "fluxend/internal/domain/storage/file"
)

func ToFileResource(file *fileDomain.File) fileDto.Response {
	return fileDto.Response{
		Uuid:          file.Uuid,
		ContainerUuid: file.ContainerUuid,
		FullFileName:  file.FullFileName,
		Size:          file.Size,
		MimeType:      file.MimeType,
		CreatedBy:     file.CreatedBy,
		UpdatedBy:     file.UpdatedBy,
		CreatedAt:     file.CreatedAt.Format("2006-01-02 15:04:05"),
		UpdatedAt:     file.UpdatedAt.Format("2006-01-02 15:04:05"),
	}
}

func ToDownloadResource(url string, expiresIn int64) fileDto.DownloadResponse {
	return fileDto.DownloadResponse{
		Url:       url,
		ExpiresIn: expiresIn,
	}
}

func ToFileResourceCollection(files []fileDomain.File) []fileDto.Response {
	resourceContainers := make([]fileDto.Response, len(files))
	for i, currentFile := range files {
		resourceContainers[i] = ToFileResource(&currentFile)
	}

	return resourceContainers
}
