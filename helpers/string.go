package helpers

import (
	"fmt"
	"math/rand"
	"regexp"
	"strings"
	"time"

	"github.com/google/uuid"
)

// GenerateTrxIDWithPrefix generates a transaction ID with a specified prefix.
// The transaction ID follows the format: prefix + YYMMDDHHMMSS + 4-digit random number.
//
// Parameters:
//   - prefix: String to prepend to the transaction ID
//
// Returns:
//   - string: Transaction ID with the specified prefix
//
// Example:
//
//	trxID := GenerateTrxIDWithPrefix("TRX-")
//	// Output: TRX-2411131534251234
func GenerateTrxIDWithPrefix(prefix string) string {
	return prefix + GenerateTrxID()
}

// GenerateTrxIDWithSuffix generates a transaction ID with a specified suffix.
// The transaction ID follows the format: YYMMDDHHMMSS + 4-digit random number + suffix.
//
// Parameters:
//   - suffix: String to append to the transaction ID
//
// Returns:
//   - string: Transaction ID with the specified suffix
//
// Example:
//
//	trxID := GenerateTrxIDWithSuffix("-END")
//	// Output: 2411131534251234-END
func GenerateTrxIDWithSuffix(suffix string) string {
	return GenerateTrxID() + suffix
}

// GenerateTrxID generates a unique transaction ID based on the current timestamp and a random 4-digit number.
// The ID format is YYMMDDHHMMSS followed by a 4-digit random number (0000-9999).
//
// Returns:
//   - string: A 16-character transaction ID
//
// Example:
//
//	trxID := GenerateTrxID()
//	// Output: 2411131534251234 (YY=24, MM=11, DD=13, HH=15, MM=34, SS=25, Random=1234)
func GenerateTrxID() string {
	// generate string dgn format YYMMDDHHiiss + 4 digit random
	now := time.Now()
	t := now.Format("060102150405") // YYMMDDHHMMSS
	rng := rand.New(rand.NewSource(now.UnixNano()))
	r := fmt.Sprintf("%04d", rng.Intn(10000))
	return t + r
}

// GenerateMessageID generates a new UUID v4 string to be used as a message ID.
// Uses the standard UUID format with hyphens.
//
// Returns:
//   - string: A UUID v4 string in the format xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx
//
// Example:
//
//	messageID := GenerateMessageID()
//	// Output: 550e8400-e29b-41d4-a716-446655440000
func GenerateMessageID() string {
	return uuid.New().String()
}

// GenerateUniqueID generates a short unique ID string by extracting the first 8 characters of a UUID v4.
// This provides a shorter identifier while maintaining reasonable uniqueness for most use cases.
//
// Returns:
//   - string: The first 8 characters of a UUID, or the full UUID if it's shorter than 8 characters
//
// Example:
//
//	uniqueID := GenerateUniqueID()
//	// Output: 550e8400
func GenerateUniqueID() string {
	uuid := uuid.New().String()
	if len(uuid) >= 8 {
		return uuid[:8]
	}
	return uuid
}

// NormalizePhoneNumber normalizes phone numbers by removing all non-numeric characters
// and adding the Indonesian country code (62) only if the phone starts with "08".
//
// Parameters:
//   - phone: Phone number string in various formats
//
// Returns:
//   - string: Normalized phone number containing only digits
//
// Examples:
//
//	NormalizePhoneNumber("+62-812-3456-789")  // Returns: 6281234567890
//	NormalizePhoneNumber("0812-3456-789")     // Returns: 6281234567890
//	NormalizePhoneNumber("812-3456-789")      // Returns: 8123456789
//	NormalizePhoneNumber("+1-202-555-1234")   // Returns: 12025551234
func NormalizePhoneNumber(phone string) string {
	// Remove all non-numeric characters
	re := regexp.MustCompile(`[^0-9]`)
	phone = re.ReplaceAllString(phone, "")

	// If starts with 08, replace with 628
	if strings.HasPrefix(phone, "08") {
		phone = "628" + phone[2:]
	}

	return phone
}

// GenerateRandomString generates a random alphanumeric string of the specified length.
// The string consists of uppercase letters, lowercase letters, and digits.
//
// Parameters:
//   - length: Desired length of the random string
//
// Returns:
//   - string: Randomly generated alphanumeric string
//
// Example:
//
//	randomStr := GenerateRandomString(10)
//	// Output: a random string like "aZ3bC9dE1F"
func GenerateRandomString(length int) string {
	const charset = "abcdefghijklmnopqrstuvwxyz" +
		"ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

	seededRand := rand.New(rand.NewSource(time.Now().UnixNano()))
	b := make([]byte, length)
	for i := range b {
		b[i] = charset[seededRand.Intn(len(charset))]
	}
	return string(b)
}
