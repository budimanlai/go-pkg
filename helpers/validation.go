package helpers

// IsValidPhoneNumber checks if the provided phone number is valid.
//
// Parameters:
//   - phone: The phone number string to be validated
//
// Returns:
//   - bool: True if the phone number is valid, false otherwise
func IsValidPhoneNumber(phone string) bool {
	// implement your phone number validation logic here
	// for simplicity, let's assume a valid phone number is 10-15 digits long
	if len(phone) < 10 || len(phone) > 15 {
		return false
	}
	for _, ch := range phone {
		if ch < '0' || ch > '9' {
			return false
		}
	}
	return true
}

// IsValidEmail checks if the provided email address is valid.
//
// Parameters:
//   - email: The email address string to be validated
//
// Returns:
//   - bool: True if the email address is valid, false otherwise
func IsValidEmail(email string) bool {
	// implement your email validation logic here
	// for simplicity, let's use a basic check
	at := false
	dot := false
	for i, ch := range email {
		if ch == '@' {
			at = true
		}
		if at && ch == '.' && i > 0 && i < len(email)-1 {
			dot = true
		}
	}
	return at && dot
}
