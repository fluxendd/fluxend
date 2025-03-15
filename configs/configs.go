package configs

var AboutFluxton = map[string]string{
	"author":      "Fluxton",
	"email":       "chief@fluxton.io",
	"website":     "https://fluxton.io",
	"license":     "MIT",
	"version":     "0.0.1",
	"release":     "alpha",
	"releaseDate": "YYYY-MM-DD",
}

const (
	BackupBucketName              = "fluxton-backups"
	MaxTableNameLength            = 60
	MinTableNameLength            = 3
	MaxColumnNameLength           = 60
	MinColumnNameLength           = 2
	MaxIndexNameLength            = 60
	MinIndexNameLength            = 3
	MaxOrganizationNameLength     = 100
	MinOrganizationNameLength     = 3
	MaxProjectNameLength          = 100
	MinProjectNameLength          = 3
	MinFormFieldLabelLength       = 3
	MaxFormFieldLabelLength       = 100
	MinFormFieldDescriptionLength = 0
	MaxFormFieldDescriptionLength = 255
	MinBucketNameLength           = 3
	MaxBucketNameLength           = 63
	MinBucketDescriptionLength    = 0
	MaxBucketDescriptionLength    = 255
	MinFileNameLength             = 3
	MaxFileNameLength             = 63
)
