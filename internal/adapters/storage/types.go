package storage

import (
	"github.com/guregu/null/v6"
	"net/http"
)

type ListContainersInput struct {
	Path  string
	Limit int
	Token string
}

type ContainerMetadata struct {
	Identifier string
	Name       string
	Path       string
	Region     null.String
}

type FileInput struct {
	ContainerName string
	FileName      string
}

type UploadFileInput struct {
	ContainerName string
	FileName      string
	FileBytes     []byte
}

type RenameFileInput struct {
	ContainerName string
	FileName      string
	NewFileName   string
}

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
