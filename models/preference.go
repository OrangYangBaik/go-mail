package models

type UserPreference struct {
	GoogleID         string `json:"google_id,omitempty" gorm:"primaryKey;not null"`
	ServiceEnabled   bool   `json:"service_enabled"`
	TelegramBotToken string `json:"telegram_bot_token"`
	TelegramChatID   string `json:"telegram_chat_id"`
	FilterCriteria   string `json:"filter_criteria"`
	LastRun          int64  `json:"last_run"`
}

func (UserPreference) TableName() string {
	return "user_preferences"
}
