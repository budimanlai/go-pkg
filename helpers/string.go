package helpers

import (
	"fmt"
	"math/rand"
	"strings"
	"time"

	"github.com/google/uuid"
)

// GenerateTrxIDWithPrefix generates a transaction ID with a specified prefix.
func GenerateTrxIDWithPrefix(prefix string) string {
	return prefix + GenerateTrxID()
}

// GenerateTrxIDWithSuffix generates a transaction ID with a specified suffix.
func GenerateTrxIDWithSuffix(suffix string) string {
	return GenerateTrxID() + suffix
}

// GenerateTrxID generates a transaction ID based on the current timestamp and a random 4-digit number.
func GenerateTrxID() string {
	// generate string dgn format YYMMDDHHiiss + 4 digit random
	now := time.Now()
	t := now.Format("060102150405") // YYMMDDHHMMSS
	rng := rand.New(rand.NewSource(now.UnixNano()))
	r := fmt.Sprintf("%04d", rng.Intn(10000))
	return t + r
}

// GenerateMessageID generates a new UUID string to be used as a message ID.
func GenerateMessageID() string {
	return uuid.New().String()
}

// GenerateUniqueID generates a new unique ID string. It returns the first 8 characters of a UUID.
func GenerateUniqueID() string {
	uuid := uuid.New().String()
	if len(uuid) >= 8 {
		return uuid[:8]
	}
	return uuid
}

// NormalizePhoneNumber ensures phone number starts with country code without +
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
