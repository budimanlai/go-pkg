package helpers

import (
	"fmt"
	"math/rand"
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

// NormalizePhoneNumber normalizes phone numbers to international format without the + prefix.
// It handles Indonesian phone numbers (62), US/Canada (1), and Singapore (65) country codes.
// If a phone number starts with 0, it's assumed to be Indonesian and converted to 62 format.
// Numbers without recognized country codes are defaulted to Indonesian format (62).
//
// Parameters:
//   - phone: Phone number string in various formats (with/without +, with/without country code)
//
// Returns:
//   - string: Normalized phone number with country code but without + prefix
//
// Examples:
//
//	NormalizePhoneNumber("+628123456789")   // Returns: 628123456789
//	NormalizePhoneNumber("08123456789")     // Returns: 628123456789
//	NormalizePhoneNumber("8123456789")      // Returns: 628123456789
//	NormalizePhoneNumber("+6591234567")     // Returns: 6591234567
//	NormalizePhoneNumber("+12025551234")    // Returns: 12025551234
func NormalizePhoneNumber(phone string) string {
	// Remove any + prefix
	phone = strings.TrimPrefix(phone, "+")

	// If starts with 0, replace with 62 (Indonesia country code)
	if strings.HasPrefix(phone, "0") {
		phone = "62" + phone[1:]
	}

	// Ensure it starts with country code
	if !strings.HasPrefix(phone, "62") && !strings.HasPrefix(phone, "1") && !strings.HasPrefix(phone, "65") {
		// Default to Indonesia if no country code
		phone = "62" + phone
	}

	return phone
}
