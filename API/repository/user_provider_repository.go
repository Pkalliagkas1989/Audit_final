package repository

import "database/sql"

// UserProviderRepository manages OAuth provider links
type UserProviderRepository struct {
	DB *sql.DB
}

func NewUserProviderRepository(db *sql.DB) *UserProviderRepository {
	return &UserProviderRepository{DB: db}
}

// GetUserID returns the user_id linked to the provider credentials
func (r *UserProviderRepository) GetUserID(provider, providerID string) (string, error) {
	var userID string
	err := r.DB.QueryRow("SELECT user_id FROM user_providers WHERE provider=? AND provider_id=?", provider, providerID).Scan(&userID)
	if err != nil {
		if err == sql.ErrNoRows {
			return "", ErrUserNotFound
		}
		return "", err
	}
	return userID, nil
}

// LinkProvider associates a provider with a user
func (r *UserProviderRepository) LinkProvider(userID, provider, providerID string) error {
	_, err := r.DB.Exec("INSERT OR IGNORE INTO user_providers (user_id, provider, provider_id) VALUES (?, ?, ?)", userID, provider, providerID)
	return err
}
