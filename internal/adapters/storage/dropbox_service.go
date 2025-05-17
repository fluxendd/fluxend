package storage

import (
	"encoding/json"
	"fluxton/pkg"
	"fluxton/pkg/errors"
	"fmt"
	"os"
	"resty.dev/v3"
	"strings"
	"time"
)

const (
	dropboxActionCreateFolder = "CREATE_FOLDER"
	dropboxActionListFolders  = "LIST_FOLDERS"
	dropboxActionShowFolder   = "SHOW_FOLDER"
	dropboxActionDeleteFolder = "DELETE_FOLDER"
	dropboxActionUploadFile   = "UPLOAD_FILE"
	dropboxActionRenameFile   = "RENAME_FILE"
	dropboxActionDownloadFile = "DOWNLOAD_FILE"
	dropboxActionDeleteFile   = "DELETE_FILE"
)

type DropboxServiceImpl struct {
	accessToken string
	client      *resty.Client
	apiBase     string
	contentBase string
}

type FolderMetadata struct {
	Path        string `json:"path_display"`
	Name        string `json:"name"`
	ID          string `json:"id"`
	IsShared    bool   `json:"sharing_info,omitempty"`
	DateCreated string `json:"server_modified,omitempty"`
}

type FileMetadata struct {
	Path        string `json:"path_display"`
	Name        string `json:"name"`
	ID          string `json:"id"`
	Size        int64  `json:"size"`
	ContentHash string `json:"content_hash,omitempty"`
	DateCreated string `json:"server_modified"`
}

type ListFolderResult struct {
	Entries []struct {
		Tag         string `json:".tag"`
		Name        string `json:"name"`
		PathDisplay string `json:"path_display"`
		ID          string `json:"id"`
	} `json:"entries"`
	Cursor  string `json:"cursor"`
	HasMore bool   `json:"has_more"`
}

func NewDropboxService() (StorageInterface, error) {
	// TODO: go with refreshable tokens instead
	accessToken := os.Getenv("DROPBOX_ACCESS_TOKEN")
	if accessToken == "" {
		return nil, fmt.Errorf("DROPBOX_ACCESS_TOKEN is not set")
	}

	client := resty.New().
		SetAuthToken(accessToken).
		SetRetryCount(3).
		SetRetryWaitTime(2*time.Second).
		SetRetryMaxWaitTime(20*time.Second).
		SetTimeout(30*time.Second).
		SetHeader("Content-Type", "application/json")

	return &DropboxServiceImpl{
		accessToken: accessToken,
		client:      client,
		apiBase:     "https://api.dropboxapi.com/2",
		contentBase: "https://content.dropboxapi.com/2",
	}, nil
}

func (d *DropboxServiceImpl) CreateContainer(path string) (string, error) {
	path = normalizePath(path)

	payload := map[string]interface{}{
		"path":       path,
		"autorename": false,
	}

	resp, err := d.executeAPIRequest("POST", "/files/create_folder_v2", payload)
	if err != nil {
		return "", err
	}

	pkg.DumpJSON(resp.StatusCode())
	// TODO: check if created folder's path or something else can be included
	return "", d.handleAPIError(resp, dropboxActionCreateFolder)
}

func (d *DropboxServiceImpl) ContainerExists(path string) bool {
	metadata, err := d.ShowContainer(path)
	return err == nil && metadata != nil
}

func (d *DropboxServiceImpl) ListContainers(input ListContainersInput) ([]string, string, error) {
	var endpoint string
	var payload map[string]interface{}

	if input.Token == "" {
		endpoint = "/files/list_folder"
		payload = map[string]interface{}{
			"path":                                input.Path,
			"recursive":                           false,
			"include_media_info":                  false,
			"include_deleted":                     false,
			"include_has_explicit_shared_members": false,
			"include_mounted_folders":             true,
		}

		if input.Limit > 0 {
			payload["limit"] = input.Limit
		}
	} else {
		endpoint = "/files/list_folder/continue"
		payload = map[string]interface{}{
			"cursor": input.Token,
		}
	}

	resp, err := d.executeAPIRequest("POST", endpoint, payload)
	if err != nil {
		return nil, "", err
	}

	if err := d.handleAPIError(resp, dropboxActionListFolders); err != nil {
		return nil, "", err
	}

	var result ListFolderResult
	if err := json.Unmarshal(resp.Bytes(), &result); err != nil {
		return nil, "", fmt.Errorf("unable to decode response: %v", err)
	}

	var folderPaths []string
	for _, entry := range result.Entries {
		if entry.Tag == "folder" {
			folderPaths = append(folderPaths, entry.PathDisplay)
		}
	}

	var nextCursor string
	if result.HasMore {
		nextCursor = result.Cursor
	}

	return folderPaths, nextCursor, nil
}

func (d *DropboxServiceImpl) ShowContainer(path string) (*ContainerMetadata, error) {
	path = normalizePath(path)

	payload := map[string]interface{}{
		"path": path,
	}

	resp, err := d.executeAPIRequest("POST", "/files/get_metadata", payload)
	if err != nil {
		return nil, err
	}

	if err := d.handleAPIError(resp, dropboxActionShowFolder); err != nil {
		return nil, err
	}

	var metadata map[string]interface{}
	if err := json.Unmarshal(resp.Bytes(), &metadata); err != nil {
		return nil, fmt.Errorf("unable to decode response: %v", err)
	}

	tag, ok := metadata[".tag"].(string)
	if !ok || tag != "folder" {
		return nil, fmt.Errorf("path is not a folder")
	}

	return &ContainerMetadata{
		Identifier: metadata["id"].(string),
		Name:       metadata["name"].(string),
		Path:       metadata["path_display"].(string),
	}, nil
}

func (d *DropboxServiceImpl) DeleteContainer(path string) error {
	path = normalizePath(path)

	payload := map[string]interface{}{
		"path": path,
	}

	resp, err := d.executeAPIRequest("POST", "/files/delete_v2", payload)
	if err != nil {
		return err
	}

	return d.handleAPIError(resp, dropboxActionDeleteFolder)
}

func (d *DropboxServiceImpl) UploadFile(input UploadFileInput) error {
	path := normalizePath(fmt.Sprintf("%s/%s", input.ContainerName, input.FileName))

	apiArg := map[string]interface{}{
		"path":       path,
		"mode":       "overwrite",
		"autorename": false,
		"mute":       false,
	}

	resp, err := d.executeContentRequest("POST", "/files/upload", apiArg, input.FileBytes)
	if err != nil {
		return err
	}

	return d.handleAPIError(resp, dropboxActionUploadFile)
}

func (d *DropboxServiceImpl) RenameFile(input RenameFileInput) error {
	payload := map[string]interface{}{
		"from_path":                normalizePath(fmt.Sprintf("%s/%s", input.ContainerName, input.FileName)),
		"to_path":                  normalizePath(input.NewFileName),
		"allow_shared_folder":      false,
		"autorename":               false,
		"allow_ownership_transfer": false,
	}

	resp, err := d.executeAPIRequest("POST", "/files/move_v2", payload)
	if err != nil {
		return err
	}

	return d.handleAPIError(resp, dropboxActionRenameFile)
}

func (d *DropboxServiceImpl) DownloadFile(input FileInput) ([]byte, error) {
	apiArg := map[string]string{
		"path": normalizePath(input.FileName),
	}

	resp, err := d.executeContentRequest("POST", "/files/download", apiArg, nil)
	if err != nil {
		return nil, err
	}

	if err := d.handleAPIError(resp, dropboxActionDownloadFile); err != nil {
		return nil, err
	}

	return resp.Bytes(), nil
}

func (d *DropboxServiceImpl) DeleteFile(input FileInput) error {
	payload := map[string]interface{}{
		"path": normalizePath(input.FileName),
	}

	resp, err := d.executeAPIRequest("POST", "/files/delete_v2", payload)
	if err != nil {
		return err
	}

	return d.handleAPIError(resp, dropboxActionDeleteFile)
}

func (d *DropboxServiceImpl) handleAPIError(resp *resty.Response, operation string) error {
	if resp.IsSuccess() {
		return nil
	}

	return d.transformError(fmt.Errorf("%s failed: %s, %s", operation, resp.Status(), resp.String()))
}

func (d *DropboxServiceImpl) executeAPIRequest(method, endpoint string, payload interface{}) (*resty.Response, error) {
	req := d.client.R()

	if payload != nil {
		req.SetBody(payload)
	}

	var resp *resty.Response
	var err error

	switch method {
	case "POST":
		resp, err = req.Post(d.apiBase + endpoint)
	case "GET":
		resp, err = req.Get(d.apiBase + endpoint)
	default:
		return nil, fmt.Errorf("unsupported HTTP method: %s", method)
	}

	if err != nil {
		return nil, fmt.Errorf("request failed: %v", err)
	}

	return resp, nil
}

func (d *DropboxServiceImpl) executeContentRequest(method, endpoint string, apiArg interface{}, body []byte) (*resty.Response, error) {
	apiArgJson, err := json.Marshal(apiArg)
	if err != nil {
		return nil, fmt.Errorf("unable to marshal JSON: %v", err)
	}

	req := d.client.R().
		SetHeader("Content-Type", "application/octet-stream").
		SetHeader("Dropbox-API-Arg", string(apiArgJson))

	if body != nil {
		req.SetBody(body)
	}

	var resp *resty.Response

	switch method {
	case "POST":
		resp, err = req.Post(d.contentBase + endpoint)
	case "GET":
		resp, err = req.Get(d.contentBase + endpoint)
	default:
		return nil, fmt.Errorf("unsupported HTTP method: %s", method)
	}

	if err != nil {
		return nil, fmt.Errorf("request failed: %v", err)
	}

	return resp, nil
}

func normalizePath(path string) string {
	if !strings.HasPrefix(path, "/") {
		return "/" + path
	}
	return path
}

func (d *DropboxServiceImpl) transformError(err error) error {
	if err == nil {
		return nil
	}

	errorString := err.Error()

	if strings.Contains(errorString, "path/not_found") {
		return errors.NewNotFoundError("dropbox.error.pathNotFound")
	}

	if strings.Contains(errorString, "path/conflict") {
		return errors.NewBadRequestError("dropbox.error.pathConflict")
	}

	if strings.Contains(errorString, "insufficient_space") {
		return errors.NewBadRequestError("dropbox.error.insufficientSpace")
	}

	if strings.Contains(errorString, "too_many_write_operations") {
		return errors.NewBadRequestError("dropbox.error.tooManyWriteOperations")
	}

	if strings.Contains(errorString, "too_many_files") {
		return errors.NewBadRequestError("dropbox.error.tooManyFiles")
	}

	return err
}
