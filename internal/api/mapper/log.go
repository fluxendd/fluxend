package mapper

import (
	logDto "fluxend/internal/api/dto/logging"
	logDomain "fluxend/internal/domain/logging"
)

func ToLoggingResource(log *logDomain.RequestLog) logDto.Response {
	return logDto.Response{
		Uuid:      log.Uuid,
		UserUuid:  log.UserUuid,
		Method:    log.Method,
		Status:    log.Status,
		Endpoint:  log.Endpoint,
		IPAddress: log.IPAddress,
		UserAgent: log.UserAgent,
		Params:    log.Params,
		Body:      log.Body,
		CreatedAt: log.CreatedAt.Format("2006-01-02 15:04:05"),
	}
}

func ToLoggingResourceCollection(files []logDomain.RequestLog) []logDto.Response {
	resourceContainers := make([]logDto.Response, len(files))
	for i, currentFile := range files {
		resourceContainers[i] = ToLoggingResource(&currentFile)
	}

	return resourceContainers
}
