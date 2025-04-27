package configs

var messages = map[string]string{
	// General

	// Authentication
	"auth.error.tokenRequired": "Token is required",
	"auth.error.tokenInvalid":  "Invalid token provided",
	"auth.error.tokenExpired":  "Token has expired",

	// User
	"user.error.notFound":              "User not found",
	"user.error.invalidCredentials":    "Invalid credentials provided",
	"user.error.updateForbidden":       "You don't have permission to update user",
	"user.error.unauthenticated":       "Unauthenticated",
	"user.error.invalidPayload":        "Invalid payload provided",
	"user.error.emailAlreadyExists":    "User with this email already exists",
	"user.error.usernameAlreadyExists": "User with this username already exists",

	// Organizations
	"organization.error.userNotFound":        "User not found in organization",
	"organization.error.notFound":            "Organization not found",
	"organization.error.viewForbidden":       "You don't have permission to view this organization",
	"organization.error.createForbidden":     "You don't have permission to create an organization",
	"organization.error.updateForbidden":     "You don't have permission to update this organization",
	"organization.error.createUserForbidden": "You don't have permission to create a user in this organization",
	"organization.error.userAlreadyExists":   "User already exists in this organization",
	"organization.error.deleteUserForbidden": "You don't have permission to delete this user from the organization",

	// Storage
	"bucket.error.notFound":        "Bucket not found",
	"bucket.error.listForbidden":   "You don't have permission to view buckets",
	"bucket.error.viewForbidden":   "You don't have permission to view this bucket",
	"bucket.error.createForbidden": "You don't have permission to create a bucket",
	"bucket.error.updateForbidden": "You don't have permission to update this bucket",
	"bucket.error.deleteWithFiles": "You can't delete this bucket because it contains files",
	"bucket.error.deleteForbidden": "You don't have permission to delete this bucket",
	"bucket.error.duplicateName":   "Bucket name already exists",

	// S3
	"s3.error.bucketAlreadyOwned":  "Bucket already owned by you",
	"s3.error.bucketAlreadyExists": "Bucket already exists",
	"s3.error.bucketNotFound":      "Bucket not found",

	// Files
	"file.error.notFound":        "File not found",
	"file.error.listForbidden":   "You don't have permission to view files",
	"file.error.viewForbidden":   "You don't have permission to view this file",
	"file.error.createForbidden": "You don't have permission to create a file",
	"file.error.updateForbidden": "You don't have permission to update this file",
	"file.error.deleteForbidden": "You don't have permission to delete this file",
	"file.error.invalidMimeType": "Invalid file type",
	"file.error.sizeExceeded":    "File size exceeds the maximum limit",
	"file.error.duplicateName":   "File name already exists",

	// Projects
	"project.error.notFound":        "Project not found",
	"project.error.viewForbidden":   "You don't have permission to view this project",
	"project.error.updateForbidden": "You don't have permission to update this project",
	"project.error.listForbidden":   "You don't have permission to view projects",
	"project.error.createForbidden": "You don't have permission to create a project",
	"project.error.duplicateName":   "Project name already exists",

	// Tables
	"table.error.createForbidden": "You don't have permission to create tables",
	"table.error.alreadyExists":   "Table already exists",

	// Columns
	"column.error.createForbidden":  "You don't have permission to create columns",
	"column.error.someAlreadyExist": "Some columns already exist",
	"column.error.someNotFound":     "Some columns not found",
	"column.error.notFound":         "Column not found",

	// Indexes
	"index.error.alreadyExists": "Index already exists",
	"index.error.notFound":      "Index not found",

	// Forms
	"form.error.notFound":        "Form not found",
	"form.error.listForbidden":   "You don't have permission to view forms",
	"form.error.viewForbidden":   "You don't have permission to view this form",
	"form.error.createForbidden": "You don't have permission to create a form",
	"form.error.updateForbidden": "You don't have permission to update this form",
	"form.error.deleteForbidden": "You don't have permission to delete this form",
	"form.error.duplicateName":   "Form name already exists",

	// Form Responses
	"formResponse.error.notFound":             "Form response not found",
	"formResponse.error.fieldRequired":        "This field is required",
	"formResponse.error.invalidNumber":        "Invalid number provided",
	"formResponse.error.numberTooLow":         "Number is too low",
	"formResponse.error.numberTooHigh":        "Number is too high",
	"formResponse.error.stringTooShort":       "String is too short",
	"formResponse.error.stringTooLong":        "String is too long",
	"formResponse.error.invalidPattern":       "Invalid pattern",
	"formResponse.error.invalidEmail":         "Invalid email address",
	"formResponse.error.invalidSelectOptions": "Invalid select options provided",
	"formResponse.error.invalidSelectOption":  "Invalid select option provided",

	// Form Fields
	"formField.error.listForbidden":       "You don't have permission to view form fields",
	"formField.error.viewForbidden":       "You don't have permission to view this form field",
	"formField.error.createForbidden":     "You don't have permission to create form fields",
	"formField.error.updateForbidden":     "You don't have permission to update this form field",
	"formField.error.deleteForbidden":     "You don't have permission to delete this form field",
	"formField.error.someDuplicateLabels": "Some labels already exist",
	"formField.error.duplicateLabel":      "Label already exists",

	// Backups
	"backup.error.notFound":         "Backup not found",
	"backup.error.listForbidden":    "You don't have permission to view backups",
	"backup.error.viewForbidden":    "You don't have permission to view this backup",
	"backup.error.createForbidden":  "You don't have permission to create a backup",
	"backup.error.deleteForbidden":  "You don't have permission to delete this backup",
	"backup.error.deleteInProgress": "Backup deletion is already in progress",

	// Settings
	"setting.error.listForbidden":   "You don't have permission to view settings",
	"setting.error.updateForbidden": "You don't have permission to update settings",
	"setting.error.resetForbidden":  "You don't have permission to reset settings",

	// Others
	"database_stats.error.forbidden": "You don't have permission to view database stats",
	"function.error.listForbidden":   "You don't have permission to view functions",
}

func Message(key string) string {
	if msg, ok := messages[key]; ok {
		return msg
	}

	return key
}
