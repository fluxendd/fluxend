package constants

const (
	DefaultOrganizationName = "Default"
	JWTSecretMinLength      = 32
	XProjectHeaderName      = "X-Project"
	BackupContainerName     = "fluxend-client-database-backups"
	StorageDriverFilesystem = "FILESYSTEM"
	StorageDriverS3         = "S3"
	StorageDriverDropbox    = "DROPBOX"
	StorageDriverBackBlaze  = "BACKBLAZE"
	EmailDriverSendGrid     = "SENDGRID"
	EmailDriverSMTP         = "SMTP"
	EmailDriverSES          = "SES"
	EmailDriverMailgun      = "MAILGUN"

	AlphanumericWithUnderscorePattern             = "^[A-Za-z0-9_]+$"
	AlphanumericWithUnderscoreAndDashPattern      = "^[A-Za-z0-9_-]+$"
	AlphanumericWithSpaceUnderScoreAndDashPattern = "^[A-Za-z0-9 _-]+$"
)
