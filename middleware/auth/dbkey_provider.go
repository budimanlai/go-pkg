package auth

import "gorm.io/gorm"

// ApiKey represents the API key model in the database.
type ApiKey struct {
	gorm.Model
	ApiKey  string `gorm:"uniqueIndex;not null"`
	AuthKey string `gorm:"not null"`
	Status  string `gorm:"not null;default:'active'"`
}

// TableName sets the table name for the ApiKey model.
func (ApiKey) TableName() string {
	return "api_key"
}

type DbKeyProvider struct {
	db *gorm.DB
}

// NewDbKeyProvider creates a new instance of DbApiKey with the provided configuration.
func NewDbKeyProvider(db *gorm.DB) *DbKeyProvider {
	return &DbKeyProvider{
		db: db,
	}
}

// IsExists checks if the given key exists in the database.
func (dk *DbKeyProvider) IsExists(key string) bool {
	var count int64

	dk.db.Model(&ApiKey{}).Where("api_key = ? AND status = 'active'", key).Count(&count)
	return count > 0
}

// GetValue retrieves the value associated with the given key from the database.
func (dk *DbKeyProvider) GetValue(key string) (string, error) {
	var apiKey ApiKey
	result := dk.db.Where("api_key = ? and status = 'active'", key).First(&apiKey)
	if result.Error != nil {
		return "", result.Error
	}
	return apiKey.AuthKey, nil
}

// Add adds a new key with the same value to the database.
func (dk *DbKeyProvider) Add(key string) error {
	apiKey := ApiKey{
		ApiKey:  key,
		AuthKey: key,
	}
	result := dk.db.Create(&apiKey)
	return result.Error
}

// AddKeyValue adds a new key-value pair to the database.
func (dk *DbKeyProvider) AddKeyValue(key string, value string) error {
	apiKey := ApiKey{
		ApiKey:  key,
		AuthKey: value,
	}
	result := dk.db.Create(&apiKey)
	return result.Error
}

// Replace replaces all existing keys in the database with the provided key-value pairs.
func (dk *DbKeyProvider) Replace(newKeys map[string]string) error {
	// Start a transaction
	tx := dk.db.Begin()

	// Delete all existing keys
	if err := tx.Exec("DELETE FROM api_key").Error; err != nil {
		tx.Rollback()
		return err
	}

	// Insert new keys
	for key, value := range newKeys {
		apiKey := ApiKey{
			ApiKey:  key,
			AuthKey: value,
		}
		if err := tx.Create(&apiKey).Error; err != nil {
			tx.Rollback()
			return err
		}
	}

	// Commit the transaction
	return tx.Commit().Error
}
