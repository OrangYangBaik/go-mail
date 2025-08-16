package adapters

import (
	"encoding/json"
	"fmt"
	"go-mail/dtos"
	"io"
	"net/http"
	"net/url"
)

type GmailAdapter interface {
	FetchUnreadEmails(accessToken string, lastRunUnix int64) ([]dtos.GmailMessageDetail, error)
}

type gmailAdapter struct{}

func NewGmailAdapter() GmailAdapter {
	return &gmailAdapter{}
}

func (g *gmailAdapter) FetchUnreadEmails(accessToken string, lastRunUnix int64) ([]dtos.GmailMessageDetail, error) {

	q := fmt.Sprintf("is:unread AND after:%d)", lastRunUnix)

	listURL := fmt.Sprintf("https://gmail.googleapis.com/gmail/v1/users/me/messages?q=%s", url.QueryEscape(q))

	req, err := http.NewRequest("GET", listURL, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", "Bearer "+accessToken)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("failed to fetch unread email list: %s", string(body))
	}

	var list dtos.GmailListResponse
	if err := json.Unmarshal(body, &list); err != nil {
		return nil, err
	}

	var emails []dtos.GmailMessageDetail
	for _, msg := range list.Messages {
		detail, err := g.fetchEmailDetail(accessToken, msg.ID)
		if err != nil {
			fmt.Printf("failed to fetch email detail for %s: %v\n", msg.ID, err)
			continue
		}
		emails = append(emails, detail)
	}

	return emails, nil
}

func (g *gmailAdapter) fetchEmailDetail(accessToken, messageID string) (dtos.GmailMessageDetail, error) {
	url := fmt.Sprintf("https://gmail.googleapis.com/gmail/v1/users/me/messages/%s?format=full", messageID)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return dtos.GmailMessageDetail{}, err
	}
	req.Header.Set("Authorization", "Bearer "+accessToken)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return dtos.GmailMessageDetail{}, err
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	if resp.StatusCode != 200 {
		return dtos.GmailMessageDetail{}, fmt.Errorf("failed to fetch email detail: %s", string(body))
	}

	var detail dtos.GmailMessageDetail
	if err := json.Unmarshal(body, &detail); err != nil {
		return dtos.GmailMessageDetail{}, err
	}

	return detail, nil
}
