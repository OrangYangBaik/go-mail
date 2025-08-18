package services

import (
	"encoding/base64"
	"fmt"
	"go-mail/adapters"
	"go-mail/dtos"
	"go-mail/models"
	"go-mail/repositories"
	"go-mail/utils"
	"log"
	"os"
	"sync"
	"time"
)

type PriorimailService interface {
	ProcessEmails() []string
}

type priorimailService struct {
	repositoryPreference repositories.PreferenceRepository
	repositoryUser       repositories.UserRepository
	adaptersTelegram     adapters.TelegramAdapter
	adapterGmail         adapters.GmailAdapter
	adapterAI            adapters.AIAdapter
	adapterGoAuth        adapters.GoAuthAdapter
}

func NewPriorimailService(
	repositoryPreference repositories.PreferenceRepository,
	repositoryUser repositories.UserRepository,
	adaptersTelegram adapters.TelegramAdapter,
	adapterGmail adapters.GmailAdapter,
	adapterAI adapters.AIAdapter,
	adapterGoAuth adapters.GoAuthAdapter,
) PriorimailService {
	return &priorimailService{
		repositoryPreference: repositoryPreference,
		repositoryUser:       repositoryUser,
		adaptersTelegram:     adaptersTelegram,
		adapterGmail:         adapterGmail,
		adapterAI:            adapterAI,
		adapterGoAuth:        adapterGoAuth,
	}
}

func (s *priorimailService) ProcessEmails() []string {
	var errorList []string
	const batchSize = 20

	emailJobs := make([]dtos.EmailJob, 0)
	var wg sync.WaitGroup

	offset := 0
	for {
		users := s.FetchUserPreferences(batchSize, offset)
		if len(users) == 0 {
			break
		}
		offset += batchSize

		wg.Add(len(users))
		for _, user := range users {
			go func(u models.UserPreference) {
				defer wg.Done()

				var sinceUnix int64
				if u.LastRun == 0 {
					now := time.Now()
					today := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
					sinceUnix = today.Unix()
				} else {
					sinceUnix = u.LastRun
				}

				emails, err := s.FetchUnreadEmails(u.GoogleID, sinceUnix)
				if err != nil {
					errorList = append(errorList, err.Error())
					return
				}

				u.LastRun = time.Now().Unix()
				err = s.repositoryPreference.UpdateUserPreferences(u)
				if err != nil {
					errorList = append(errorList, err.Error())
				}

				for _, email := range emails {
					threadId := email.ThreadID
					subject := utils.GetHeader(email.Payload.Headers, "Subject")
					from := utils.GetHeader(email.Payload.Headers, "From")
					to := utils.ParseRecipients(utils.GetHeader(email.Payload.Headers, "To"))
					//body := utils.ExtractEmailBody(email)

					emailJobs = append(emailJobs, dtos.EmailJob{
						Preferences: u,
						Subject:     subject,
						//Body:        body,
						ThreadID:   threadId,
						Sender:     from,
						Recipients: to,
					})
				}
			}(user)
		}
		wg.Wait()
	}

	for _, job := range emailJobs {
		important, err := s.adapterAI.AIVerdictEmail(
			job.Preferences.FilterCriteria,
			job.Sender,
			job.Subject,
			//job.Body,
		)
		if err != nil {
			errorList = append(errorList, fmt.Sprintf(
				"AI error for user %s: %v", job.Preferences.GoogleID, err))
			continue
		}

		if important {
			err := s.adaptersTelegram.SendEmailToTelegram(
				job.Preferences.TelegramBotToken,
				job.Preferences.TelegramChatID,
				job.Sender,
				job.Subject,
				job.ThreadID,
				// job.Body,
			)
			if err != nil {
				errorList = append(errorList, fmt.Sprintf(
					"Telegram error for user %s: %v", job.Preferences.GoogleID, err))
			}
		}
	}

	return errorList
}

func (s *priorimailService) FetchUserPreferences(batchSize, offset int) []models.UserPreference {
	preferences, err := s.repositoryPreference.FetchUserPreferences(batchSize, offset)
	if err != nil {
		log.Printf("Error fetching user preferences: %v", err)
		return nil
	}
	return preferences
}

func (s *priorimailService) FetchUnreadEmails(googleID string, lastRunUnix int64) ([]dtos.GmailMessageDetail, error) {
	var user *models.User
	var encryptKey64 = os.Getenv("ENCRYPTION_SECRET_KEY")

	user, err := s.repositoryUser.GetByGoogleID(googleID)
	if err != nil {
		return nil, fmt.Errorf("Error fetching user by Google ID %s: %v", googleID, err)
	}

	if user == nil {
		return nil, fmt.Errorf("User not found for Google ID %s", googleID)
	}

	if time.Until(user.Expiry) <= 5*time.Minute {
		err = s.adapterGoAuth.RefreshToken(user.RefreshToken, user.GoogleID)
		if err != nil {
			return nil, fmt.Errorf("Error refreshing token for user %s: %v", googleID, err)
		}

		user, err = s.repositoryUser.GetByGoogleID(googleID)
		if err != nil {
			return nil, fmt.Errorf("Error fetching user by Google ID %s: %v", googleID, err)
		}
	}

	key, err := base64.StdEncoding.DecodeString(encryptKey64)
	if err != nil {
		return nil, fmt.Errorf("Error decoding encryption key: %v", err)
	}

	accesstokenDecrypted, err := utils.DecryptAccessToken(user.AccessToken, key)
	if err != nil {
		return nil, fmt.Errorf("Error decrypting access token for user %d: %v", googleID, err)
	}

	emails, err := s.adapterGmail.FetchUnreadEmails(accesstokenDecrypted, lastRunUnix)
	if err != nil {
		return nil, fmt.Errorf("Error fetching unread emails for user %d: %v", googleID, err)
	}

	return emails, nil
}
