package security

import "golang.org/x/crypto/bcrypt"

// HashPassword hashes a password using bcrypt with the default cost factor.
// It returns the bcrypt hash as a string, or an empty string if hashing fails.
//
// Parameters:
//   - password: The plain text password to hash
//
// Returns:
//   - string: The bcrypt hash (60 characters), or empty string on error
//
// Example:
//
//	hash := security.HashPassword("myPassword123")
//	if hash == "" {
//	    log.Fatal("Failed to hash password")
//	}
func HashPassword(password string) string {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return ""
	}
	return string(hash)
}

// CheckPasswordHash compares a plain text password with a bcrypt hash.
// It returns true if the password matches the hash, false otherwise.
//
// Parameters:
//   - password: The plain text password to verify
//   - hash: The bcrypt hash to compare against
//
// Returns:
//   - bool: true if password matches, false otherwise
//   - error: nil if comparison succeeded (even if password doesn't match),
//     error only for unexpected failures (e.g., corrupted hash)
//
// Example:
//
//	valid, err := security.CheckPasswordHash("myPassword123", hashedPassword)
//	if err != nil {
//	    log.Fatal(err)
//	}
//	if !valid {
//	    fmt.Println("Invalid password")
//	}
func CheckPasswordHash(password, hash string) (bool, error) {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	if err != nil {
		if err == bcrypt.ErrMismatchedHashAndPassword {
			return false, nil
		}
		return false, err
	}
	return true, nil
}
