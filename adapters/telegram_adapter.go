package adapters

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
)

type TelegramAdapter interface {
	IsValidChatBot(chatID string, botToken string) (bool, error)
	SendEmailToTelegram(botToken, chatID, sender, subject, threadId string) error
}

type telegramAdapter struct{}

func NewTelegramAdapter() TelegramAdapter {
	return &telegramAdapter{}
}

func (t *telegramAdapter) IsValidChatBot(chatID string, botToken string) (bool, error) {
	url := fmt.Sprintf("https://api.telegram.org/bot%s/getChat?chat_id=%s", botToken, chatID)
	resp, err := http.Get(url)
	if err != nil {
		return false, err
	}
	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)
	if resp.StatusCode == 200 {
		return true, nil
	}
	return false, fmt.Errorf("invalid chatID or botToken: %s", string(body))
}

func (t *telegramAdapter) SendEmailToTelegram(botToken, chatID, sender, subject, threadId string) error {
	apiURL := fmt.Sprintf("https://api.telegram.org/bot%s/sendMessage", botToken)

	message := fmt.Sprintf(
		"*New Email Received*\n\n*From:* %s\n*Subject:* %s\n*Email Link:* %s\n\nYou can view the email in your Gmail account, and it is only accessible on a PC or laptop.",
		sender,
		subject,
		fmt.Sprintf("https://mail.google.com/mail/u/0/#inbox/%s", threadId),
	)

	formData := url.Values{}
	formData.Set("chat_id", chatID)
	formData.Set("text", message)
	formData.Set("parse_mode", "Markdown")

	resp, err := http.PostForm(apiURL, formData)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	respBody, _ := io.ReadAll(resp.Body)
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to send message: %s", string(respBody))
	}

	return nil
}
