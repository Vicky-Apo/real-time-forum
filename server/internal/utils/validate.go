package utils

import (
	"slices"
	"errors"
	"net/mail"
	"regexp"
	"strings"

	"platform.zone01.gr/git/gpapadopoulos/forum/config"
)

// more secure way to validate email, password, and username inputs
// Used net/mail package for email validation
func ValidateEmail(email string) error {
	email = strings.TrimSpace(email)

	_, err := mail.ParseAddress(email)
	if err != nil {
		return errors.New("invalid email format")
	}

	emailRegex := regexp.MustCompile(`^[a-z0-9._%+\-]+@[a-z0-9.\-]+\.[a-z]{2,}$`)
	if !emailRegex.MatchString(email) {
		return errors.New("invalid email format")
	}

	return nil
}

func ValidatePassword(password string) error {
	// Password validation using configuration
	if len(password) < config.Config.MinPasswordLen || len(password) > config.Config.MaxPasswordLen {
		return errors.New("password must be between " +
			string(rune(config.Config.MinPasswordLen)) + " and " +
			string(rune(config.Config.MaxPasswordLen)) + " characters")
	}

	// Check password complexity
	hasUppercase := regexp.MustCompile(`[A-Z]`).MatchString(password)
	hasLowercase := regexp.MustCompile(`[a-z]`).MatchString(password)
	hasNumber := regexp.MustCompile(`[0-9]`).MatchString(password)
	hasSpecial := regexp.MustCompile(`[^a-zA-Z0-9]`).MatchString(password)

	if !(hasUppercase && hasLowercase && hasNumber && hasSpecial) {
		return errors.New("password must contain uppercase, lowercase, number, and special character")
	}

	return nil
}

func ValidateUsername(username string) error {
	// Username validation using configuration
	if len(username) < config.Config.MinUsernameLen || len(username) > config.Config.MaxUsernameLen {
		return errors.New("username must be between " +
			string(rune(config.Config.MinUsernameLen)) + " and " +
			string(rune(config.Config.MaxUsernameLen)) + " characters")
	}

	// Username character validation (alphanumeric and underscore only)
	usernameRegex := regexp.MustCompile("^[a-zA-Z0-9_]+$")
	if !usernameRegex.MatchString(username) {
		return errors.New("username can only contain alphanumeric characters and underscores")
	}

	return nil
}

func ValidateName(name, fieldName string) error {
	// Trim whitespace
	name = strings.TrimSpace(name)

	// Check length
	if len(name) < 2 {
		return errors.New(fieldName + " must be at least 2 characters")
	}
	if len(name) > 50 {
		return errors.New(fieldName + " must be 50 characters or less")
	}

	// Check if name is only spaces or special characters
	if strings.TrimSpace(name) == "" {
		return errors.New(fieldName + " cannot be empty or only spaces")
	}

	// Allow letters (including unicode/accented), spaces, hyphens, and apostrophes
	// This supports names like: Jean-Pierre, O'Brien, María, François
	nameRegex := regexp.MustCompile(`^[\p{L}\s'-]+$`)
	if !nameRegex.MatchString(name) {
		return errors.New(fieldName + " can only contain letters, spaces, hyphens, and apostrophes")
	}

	// Prevent starting or ending with special characters (space, hyphen, apostrophe)
	if matched, _ := regexp.MatchString(`^[\s'-]|[\s'-]$`, name); matched {
		return errors.New(fieldName + " cannot start or end with spaces, hyphens, or apostrophes")
	}

	// Prevent multiple consecutive special characters
	if matched, _ := regexp.MatchString(`[\s'-]{2,}`, name); matched {
		return errors.New(fieldName + " cannot contain consecutive special characters")
	}

	return nil
}

func ValidateAge(age int) error {
	// Validate age - must be between 13 and 120
	if age < 13 {
		return errors.New("age must be at least 13 years old")
	}
	if age > 120 {
		return errors.New("age must be 120 or less")
	}
	return nil
}

func ValidateGender(gender string) error {
	// Validate gender - must be one of the allowed values
	allowedGenders := []string{"Male", "Female", "Other"}

	gender = strings.TrimSpace(gender)

	if slices.Contains(allowedGenders, gender) {
			return nil
		}

	return errors.New("gender must be one of: Male, Female, or Other")
}

func ValidateUserInput(username, email, password, gender, firstName, lastName string, age int) error {
	if err := ValidateUsername(username); err != nil {
		return err
	}
	if err := ValidateEmail(email); err != nil {
		return err
	}
	if err := ValidatePassword(password); err != nil {
		return err
	}
	if err := ValidateGender(gender); err != nil {
		return err
	}
	if err := ValidateAge(age); err != nil {
		return err
	}
	if err := ValidateName(firstName, "first name"); err != nil {
		return err
	}
	if err := ValidateName(lastName, "last name"); err != nil {
		return err
	}
	return nil
}

func ValidatePostContent(content string) error {
	// Content validation using configuration
	if len(content) < config.Config.MinPostContentLength || len(content) > config.Config.MaxPostContentLength {
		return errors.New("post content must be between " +
			string(rune(config.Config.MinPostContentLength)) + " and " +
			string(rune(config.Config.MaxPostContentLength)) + " characters")
	}

	// Check for prohibited words (example)
	prohibitedWords := []string{"fuck", "bitch", "asshole"}
	for _, word := range prohibitedWords {
		if strings.Contains(strings.ToLower(content), word) {
			return errors.New("post content contains prohibited words")
		}
	}

	return nil
}

func ValidateCommentContent(content string) error {
	// Content validation using configuration
	if len(content) < config.Config.MinCommentLength || len(content) > config.Config.MaxCommentLength {
		return errors.New("comment content must be between " +
			string(rune(config.Config.MinCommentLength)) + " and " +
			string(rune(config.Config.MaxCommentLength)) + " characters")
	}

	// Check for prohibited words (example)
	prohibitedWords := []string{"fuck", "bitch", "asshole"}
	for _, word := range prohibitedWords {
		if strings.Contains(strings.ToLower(content), word) {
			return errors.New("comment content contains prohibited words")
		}
	}

	return nil
}
