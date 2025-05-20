package pkg

// AlphaValidationPattern returns the regex pattern for alpha validation (only letters, no spaces)
// Example: "Hello" will match, but "Hello123" or "Hello World" will not match.
func AlphaValidationPattern() string {
	return "^[A-Za-z]+$"
}

// AlphaWithSpaceValidationPattern returns the regex pattern for alpha with space validation (letters and spaces)
// Example: "Hello World" will match, but "Hello123" or "Hello@World" will not match.
func AlphaWithSpaceValidationPattern() string {
	return "^[A-Za-z ]+$"
}

// AlphaWithSpaceAndSpecialCharsValidationPattern returns the regex pattern for alpha with space and special characters validation
// Example: "Hello World!" will match, but "Hello123" will not match.
func AlphaWithSpaceAndSpecialCharsValidationPattern() string {
	return "^[A-Za-z\\s!@#$%^&*()_+={}|:;,.<>?`~\\-]+$"
}

// NumericValidationPattern returns the regex pattern for numeric validation (only digits)
// Example: "12345" will match, but "123a5" or "12 345" will not match.
func NumericValidationPattern() string {
	return "^[0-9]+$"
}

// AlphanumericValidationPattern returns the regex pattern for alphanumeric validation (letters and digits)
// Example: "" will match, but "Hello!" or "123 456" will not match.
func AlphanumericValidationPattern() string {
	return "^[A-Za-z0-9]+$"
}

// AlphanumericWithSpaceValidationPattern returns the regex pattern for alphanumeric with space validation (letters, digits, and spaces)
// Example: "Hello 123" will match, but "Hello@123" or "Hello 123!" will not match.
func AlphanumericWithSpaceValidationPattern() string {
	return "^[A-Za-z0-9 ]+$"
}

// AlphanumericWithUnderscoreAndDashPattern returns the regex pattern for alphanumeric with underscores and dashes validation
// Example: "Hello_123" or "Hello-123" will match, but "Hello 123" or "Hello@123" will not match.
func AlphanumericWithUnderscoreAndDashPattern() string {
	return "^[A-Za-z0-9_-]+$"
}

// AlphanumericWithUnderscorePattern returns the regex pattern for alphanumeric with underscores validation
// Example: "Hello_123" will match, but "Hello-123" or "Hello 123" will not match.
func AlphanumericWithUnderscorePattern() string {
	return "^[A-Za-z0-9_]+$"
}

// AlphanumericWithSpaceUnderScoreAndDashPattern returns the regex pattern for alphanumeric with space, underscores, and dashes validation
// Example: "Hello_123" or "Hello-123" or "Hello 123" will match, but "Hello@123" will not match.
func AlphanumericWithSpaceUnderScoreAndDashPattern() string {
	return "^[A-Za-z0-9 _-]+$"
}

// EmailValidationPattern returns the regex pattern for email validation
// Example: "example@domain.com" will match, but "example@domain" or "example.com" will not match.
func EmailValidationPattern() string {
	return "^[A-Za-z0-9._%+-]+@[A-Za-z0-9.-]+\\.[A-Za-z]{2,}$"
}

// AlphanumericPattern returns the regex pattern for alphanumeric validation (letters and digits)
// Example: "Hello123" will match, but "Hello!" or "123 456" will not match.
func AlphanumericPattern() string {
	return "^[A-Za-z0-9]+$"
}

// AlphanumericWithSpacePattern returns the regex pattern for alphanumeric with space validation (letters, digits, and spaces)
// Example: "Hello 123" will match, but "Hello@123" will not match.
func AlphanumericWithSpacePattern() string {
	return "^[A-Za-z0-9 ]+$"
}

// AlphanumericWithSpaceAndSpecialCharsPattern returns the regex pattern for alphanumeric with space and special characters validation
// Example: "Hello123 !" will match, but "Hello123" or "Hello!" will not match.
func AlphanumericWithSpaceAndSpecialCharsPattern() string {
	return "^[A-Za-z0-9\\s!@#$%^&*()_+={}|:;,.<>?`~\\-]+$"
}
