package storage

import (
	"encoding/base64"
	"encoding/json"
	"fluxton/pkg/errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"

	"github.com/guregu/null/v6"
)

// B2 API response structures
type B2AuthorizeAccountResponse struct {
	AccountID               string `json:"accountId"`
	AuthorizationToken      string `json:"authorizationToken"`
	APIURL                  string `json:"apiUrl"`
	DownloadURL             string `json:"downloadUrl"`
	RecommendedPartSize     int64  `json:"recommendedPartSize"`
	AbsoluteMinimumPartSize int64  `json:"absoluteMinimumPartSize"`
}

type B2CreateBucketRequest struct {
	AccountID      string            `json:"accountId"`
	BucketName     string            `json:"bucketName"`
	BucketType     string            `json:"bucketType"`
	BucketInfo     map[string]string `json:"bucketInfo,omitempty"`
	CorsRules      []interface{}     `json:"corsRules,omitempty"`
	LifecycleRules []interface{}     `json:"lifecycleRules,omitempty"`
}

type B2CreateBucketResponse struct {
	BucketID       string            `json:"bucketId"`
	AccountID      string            `json:"accountId"`
	BucketName     string            `json:"bucketName"`
	BucketType     string            `json:"bucketType"`
	BucketInfo     map[string]string `json:"bucketInfo"`
	CorsRules      []interface{}     `json:"corsRules"`
	LifecycleRules []interface{}     `json:"lifecycleRules"`
	RevisionNumber int               `json:"revision"`
	Options        []string          `json:"options,omitempty"`
}

type B2ListBucketsRequest struct {
	AccountID   string   `json:"accountId"`
	BucketID    string   `json:"bucketId,omitempty"`
	BucketName  string   `json:"bucketName,omitempty"`
	BucketTypes []string `json:"bucketTypes,omitempty"`
}

type B2ListBucketsResponse struct {
	Buckets []B2BucketItem `json:"buckets"`
}

type B2BucketItem struct {
	BucketID       string            `json:"bucketId"`
	AccountID      string            `json:"accountId"`
	BucketName     string            `json:"bucketName"`
	BucketType     string            `json:"bucketType"`
	BucketInfo     map[string]string `json:"bucketInfo"`
	CorsRules      []interface{}     `json:"corsRules"`
	LifecycleRules []interface{}     `json:"lifecycleRules"`
	RevisionNumber int               `json:"revision"`
	Options        []string          `json:"options,omitempty"`
}

type B2GetUploadURLRequest struct {
	BucketID string `json:"bucketId"`
}

type B2GetUploadURLResponse struct {
	BucketID           string `json:"bucketId"`
	UploadURL          string `json:"uploadUrl"`
	AuthorizationToken string `json:"authorizationToken"`
}

type B2FileMetadata struct {
	FileID          string            `json:"fileId"`
	FileName        string            `json:"fileName"`
	ContentLength   int64             `json:"contentLength"`
	ContentType     string            `json:"contentType"`
	ContentSha1     string            `json:"contentSha1"`
	FileInfo        map[string]string `json:"fileInfo,omitempty"`
	BucketID        string            `json:"bucketId"`
	Action          string            `json:"action"`
	UploadTimestamp int64             `json:"uploadTimestamp"`
}

type B2ListFilesRequest struct {
	BucketID      string `json:"bucketId"`
	StartFileName string `json:"startFileName,omitempty"`
	MaxFileCount  int    `json:"maxFileCount,omitempty"`
	Prefix        string `json:"prefix,omitempty"`
	Delimiter     string `json:"delimiter,omitempty"`
}

type B2ListFilesResponse struct {
	Files        []B2FileMetadata `json:"files"`
	NextFileName string           `json:"nextFileName"`
}

type B2DeleteFileVersionRequest struct {
	FileID   string `json:"fileId"`
	FileName string `json:"fileName"`
}

type BackblazeServiceImpl struct {
	apiBase            string
	accountID          string
	authorizationToken string
	apiURL             string
	downloadURL        string
	httpClient         *http.Client
}

func NewBackblazeService() (StorageInterface, error) {
	applicationKeyID := os.Getenv("BACKBLAZE_KEY_ID")
	applicationKey := os.Getenv("BACKBLAZE_APPLICATION_KEY")

	if applicationKeyID == "" || applicationKey == "" {
		return nil, fmt.Errorf("Backblaze credentials not found in environment variables")
	}

	service := &BackblazeServiceImpl{
		apiBase: "https://api.backblazeb2.com/b2api/v2",
		httpClient: &http.Client{
			Timeout: time.Second * 30,
		},
	}

	// Authorize the account first
	err := service.authorizeAccount(applicationKeyID, applicationKey)
	if err != nil {
		return nil, fmt.Errorf("unable to authorize Backblaze account: %v", err)
	}

	return service, nil
}

func (b *BackblazeServiceImpl) authorizeAccount(applicationKeyID, applicationKey string) error {
	req, err := http.NewRequest("GET", b.apiBase+"/b2_authorize_account", nil)
	if err != nil {
		return fmt.Errorf("error creating authorization request: %w", err)
	}

	// Add authorization header with base64 encoded credentials
	authString := base64.StdEncoding.EncodeToString([]byte(applicationKeyID + ":" + applicationKey))
	req.Header.Add("Authorization", "Basic "+authString)

	// Make the request
	resp, err := b.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("error making authorization request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("authorization failed with status %d: %s", resp.StatusCode, string(body))
	}

	// Parse the response
	var authResponse B2AuthorizeAccountResponse
	if err := json.NewDecoder(resp.Body).Decode(&authResponse); err != nil {
		return fmt.Errorf("error parsing authorization response: %w", err)
	}

	// Store auth details
	b.accountID = authResponse.AccountID
	b.authorizationToken = authResponse.AuthorizationToken
	b.apiURL = authResponse.APIURL
	b.downloadURL = authResponse.DownloadURL

	return nil
}

func (b *BackblazeServiceImpl) makeAuthorizedRequest(method, endpoint string, body io.Reader) (*http.Response, error) {
	requestUrl := b.apiURL + endpoint
	req, err := http.NewRequest(method, requestUrl, body)
	if err != nil {
		return nil, fmt.Errorf("error creating request: %w", err)
	}

	req.Header.Add("Authorization", b.authorizationToken)
	if body != nil {
		req.Header.Add("Content-Type", "application/json")
	}

	resp, err := b.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error making request: %w", err)
	}

	if resp.StatusCode >= 400 {
		defer resp.Body.Close()
		bodyBytes, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("request failed with status %d: %s", resp.StatusCode, string(bodyBytes))
	}

	return resp, nil
}

func (b *BackblazeServiceImpl) CreateContainer(bucketName string) (string, error) {
	requestBody := B2CreateBucketRequest{
		AccountID:  b.accountID,
		BucketName: bucketName,
		BucketType: "allPrivate", // Can be "allPublic" or "allPrivate"
	}

	jsonBody, err := json.Marshal(requestBody)
	if err != nil {
		return "", fmt.Errorf("error marshaling bucket creation request: %w", err)
	}

	resp, err := b.makeAuthorizedRequest("POST", "/b2_create_bucket", strings.NewReader(string(jsonBody)))
	if err != nil {
		return "", b.transformError(err)
	}
	defer resp.Body.Close()

	var createResponse B2CreateBucketResponse
	if err := json.NewDecoder(resp.Body).Decode(&createResponse); err != nil {
		return "", fmt.Errorf("error parsing bucket creation response: %w", err)
	}

	// Return the bucket ID
	return createResponse.BucketID, nil
}

func (b *BackblazeServiceImpl) ContainerExists(bucketName string) bool {
	buckets, _, err := b.ListContainers(ListContainersInput{})
	if err != nil {
		return false
	}

	for _, bucket := range buckets {
		if bucket == bucketName {
			return true
		}
	}

	return false
}

func (b *BackblazeServiceImpl) ListContainers(input ListContainersInput) ([]string, string, error) {
	requestBody := B2ListBucketsRequest{
		AccountID: b.accountID,
	}

	jsonBody, err := json.Marshal(requestBody)
	if err != nil {
		return nil, "", fmt.Errorf("error marshaling list buckets request: %w", err)
	}

	resp, err := b.makeAuthorizedRequest("POST", "/b2_list_buckets", strings.NewReader(string(jsonBody)))
	if err != nil {
		return nil, "", fmt.Errorf("error listing buckets: %w", err)
	}
	defer resp.Body.Close()

	var listResponse B2ListBucketsResponse
	if err := json.NewDecoder(resp.Body).Decode(&listResponse); err != nil {
		return nil, "", fmt.Errorf("error parsing bucket list response: %w", err)
	}

	var bucketNames []string
	for _, bucket := range listResponse.Buckets {
		bucketNames = append(bucketNames, bucket.BucketName)
	}

	// B2 doesn't support pagination for bucket listing the same way S3 does
	// So we return an empty string for nextToken
	return bucketNames, "", nil
}

func (b *BackblazeServiceImpl) ShowContainer(bucketName string) (*ContainerMetadata, error) {
	requestBody := B2ListBucketsRequest{
		AccountID:  b.accountID,
		BucketName: bucketName,
	}

	jsonBody, err := json.Marshal(requestBody)
	if err != nil {
		return nil, fmt.Errorf("error marshaling show bucket request: %w", err)
	}

	resp, err := b.makeAuthorizedRequest("POST", "/b2_list_buckets", strings.NewReader(string(jsonBody)))
	if err != nil {
		return nil, b.transformError(err)
	}
	defer resp.Body.Close()

	var listResponse B2ListBucketsResponse
	if err := json.NewDecoder(resp.Body).Decode(&listResponse); err != nil {
		return nil, fmt.Errorf("error parsing bucket show response: %w", err)
	}

	if len(listResponse.Buckets) == 0 {
		return nil, errors.NewNotFoundError(fmt.Sprintf("bucket %q not found", bucketName))
	}

	bucket := listResponse.Buckets[0]
	return &ContainerMetadata{
		Identifier: bucket.BucketID,
		Region:     null.StringFrom(""), // B2 doesn't use regions in the same way as S3
	}, nil
}

func (b *BackblazeServiceImpl) DeleteContainer(bucketName string) error {
	// First, get the bucket ID
	containerMetadata, err := b.ShowContainer(bucketName)
	if err != nil {
		return err
	}

	bucketID := containerMetadata.Identifier

	// Delete the bucket
	requestBody := map[string]string{
		"accountId": b.accountID,
		"bucketId":  bucketID,
	}

	jsonBody, err := json.Marshal(requestBody)
	if err != nil {
		return fmt.Errorf("error marshaling delete bucket request: %w", err)
	}

	resp, err := b.makeAuthorizedRequest("POST", "/b2_delete_bucket", strings.NewReader(string(jsonBody)))
	if err != nil {
		return b.transformError(err)
	}
	resp.Body.Close()

	return nil
}

func (b *BackblazeServiceImpl) getUploadURL(bucketID string) (*B2GetUploadURLResponse, error) {
	requestBody := B2GetUploadURLRequest{
		BucketID: bucketID,
	}

	jsonBody, err := json.Marshal(requestBody)
	if err != nil {
		return nil, fmt.Errorf("error marshaling get upload URL request: %w", err)
	}

	resp, err := b.makeAuthorizedRequest("POST", "/b2_get_upload_url", strings.NewReader(string(jsonBody)))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var uploadURLResponse B2GetUploadURLResponse
	if err := json.NewDecoder(resp.Body).Decode(&uploadURLResponse); err != nil {
		return nil, fmt.Errorf("error parsing upload URL response: %w", err)
	}

	return &uploadURLResponse, nil
}

func (b *BackblazeServiceImpl) UploadFile(input UploadFileInput) error {
	// First, get the container metadata to get the bucket ID
	containerMetadata, err := b.ShowContainer(input.ContainerName)
	if err != nil {
		return err
	}

	bucketID := containerMetadata.Identifier

	// Get an upload URL
	uploadURLResponse, err := b.getUploadURL(bucketID)
	if err != nil {
		return fmt.Errorf("unable to get upload URL: %w", err)
	}

	// Create a request to upload the file
	req, err := http.NewRequest("POST", uploadURLResponse.UploadURL, strings.NewReader(string(input.FileBytes)))
	if err != nil {
		return fmt.Errorf("error creating upload request: %w", err)
	}

	// Add required headers
	req.Header.Add("Authorization", uploadURLResponse.AuthorizationToken)
	req.Header.Add("X-Bz-File-Name", url.QueryEscape(input.FileName))
	req.Header.Add("Content-Type", "application/octet-stream")
	req.Header.Add("X-Bz-Content-Sha1", "do_not_verify") // In production, calculate and verify SHA1

	// Make the request
	resp, err := b.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("error uploading file: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("upload failed with status %d: %s", resp.StatusCode, string(bodyBytes))
	}

	return nil
}

func (b *BackblazeServiceImpl) RenameFile(input RenameFileInput) error {
	// In B2, renaming a file requires copying it to a new name and then deleting the original
	// First, download the file
	fileInput := FileInput{
		ContainerName: input.ContainerName,
		FileName:      input.FileName,
	}

	fileData, err := b.DownloadFile(fileInput)
	if err != nil {
		return fmt.Errorf("unable to download file for rename: %w", err)
	}

	// Upload with the new name
	uploadInput := UploadFileInput{
		ContainerName: input.ContainerName,
		FileName:      input.NewFileName,
		FileBytes:     fileData,
	}

	if err := b.UploadFile(uploadInput); err != nil {
		return fmt.Errorf("unable to upload file with new name: %w", err)
	}

	// Delete the original file
	if err := b.DeleteFile(fileInput); err != nil {
		return fmt.Errorf("unable to delete original file after rename: %w", err)
	}

	return nil
}

func (b *BackblazeServiceImpl) DownloadFile(input FileInput) ([]byte, error) {
	// Get the file info
	containerMetadata, err := b.ShowContainer(input.ContainerName)
	if err != nil {
		return nil, err
	}

	bucketID := containerMetadata.Identifier

	// First, get the file ID
	_, err = b.getFileID(bucketID, input.FileName)
	if err != nil {
		return nil, fmt.Errorf("unable to get file ID: %w", err)
	}

	// Construct the download URL
	downloadURL := fmt.Sprintf("%s/file/%s/%s", b.downloadURL, input.ContainerName, url.QueryEscape(input.FileName))

	// Create request
	req, err := http.NewRequest("GET", downloadURL, nil)
	if err != nil {
		return nil, fmt.Errorf("error creating download request: %w", err)
	}

	// Add authorization
	req.Header.Add("Authorization", b.authorizationToken)

	// Make the request
	resp, err := b.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error downloading file: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("download failed with status %d: %s", resp.StatusCode, string(bodyBytes))
	}

	// Read the file content
	fileBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("unable to read file content: %w", err)
	}

	return fileBytes, nil
}

func (b *BackblazeServiceImpl) getFileID(bucketID, fileName string) (string, error) {
	requestBody := B2ListFilesRequest{
		BucketID:      bucketID,
		StartFileName: fileName,
		MaxFileCount:  1,
	}

	jsonBody, err := json.Marshal(requestBody)
	if err != nil {
		return "", fmt.Errorf("error marshaling list files request: %w", err)
	}

	resp, err := b.makeAuthorizedRequest("POST", "/b2_list_file_names", strings.NewReader(string(jsonBody)))
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	var listResponse B2ListFilesResponse
	if err := json.NewDecoder(resp.Body).Decode(&listResponse); err != nil {
		return "", fmt.Errorf("error parsing file list response: %w", err)
	}

	if len(listResponse.Files) == 0 || listResponse.Files[0].FileName != fileName {
		return "", errors.NewNotFoundError(fmt.Sprintf("file %q not found", fileName))
	}

	return listResponse.Files[0].FileID, nil
}

func (b *BackblazeServiceImpl) DeleteFile(input FileInput) error {
	// Get the container metadata to get the bucket ID
	containerMetadata, err := b.ShowContainer(input.ContainerName)
	if err != nil {
		return err
	}

	bucketID := containerMetadata.Identifier

	// Get the file ID
	fileID, err := b.getFileID(bucketID, input.FileName)
	if err != nil {
		return fmt.Errorf("unable to get file ID: %w", err)
	}

	// Delete the file
	requestBody := B2DeleteFileVersionRequest{
		FileID:   fileID,
		FileName: input.FileName,
	}

	jsonBody, err := json.Marshal(requestBody)
	if err != nil {
		return fmt.Errorf("error marshaling delete file request: %w", err)
	}

	resp, err := b.makeAuthorizedRequest("POST", "/b2_delete_file_version", strings.NewReader(string(jsonBody)))
	if err != nil {
		return b.transformError(err)
	}
	resp.Body.Close()

	return nil
}

func (b *BackblazeServiceImpl) transformError(err error) error {
	if err == nil {
		return nil
	}

	errorString := err.Error()

	if strings.Contains(errorString, "duplicate_bucket_name") {
		return errors.NewBadRequestError("backblaze.error.bucketAlreadyExists")
	}

	if strings.Contains(errorString, "not_found") && strings.Contains(errorString, "bucket") {
		return errors.NewNotFoundError("backblaze.error.bucketNotFound")
	}

	if strings.Contains(errorString, "file_not_present") {
		return errors.NewNotFoundError("backblaze.error.fileNotFound")
	}

	return err
}
