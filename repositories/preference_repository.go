package repositories

import (
	"go-mail/models"

	"gorm.io/gorm"
)

type PreferenceRepository interface {
	GetUserPreferences(googleID string) (models.UserPreference, error)
	UpdateUserPreferences(preferences models.UserPreference) error
	FetchUserPreferences(batchSize, offset int) ([]models.UserPreference, error)
}

type preferenceRepository struct {
	db *gorm.DB
}

func NewPreferenceRepository(db *gorm.DB) PreferenceRepository {
	return &preferenceRepository{db: db}
}

func (r *preferenceRepository) GetUserPreferences(googleID string) (models.UserPreference, error) {
	var preferences models.UserPreference
	err := r.db.Where("google_id = ?", googleID).First(&preferences).Error
	if err != nil {
		return models.UserPreference{}, err
	}
	return preferences, nil
}

func (r *preferenceRepository) UpdateUserPreferences(preferences models.UserPreference) error {
	err := r.db.Save(&preferences).Error
	if err != nil {
		return err
	}

	return nil
}

func (r *preferenceRepository) FetchUserPreferences(batchSize, offset int) ([]models.UserPreference, error) {
	var preferences []models.UserPreference
	err := r.db.
		Where("service_enabled = ?", true).
		Offset(offset).
		Limit(batchSize).
		Find(&preferences).Error
	if err != nil {
		return nil, err
	}
	return preferences, nil
}
