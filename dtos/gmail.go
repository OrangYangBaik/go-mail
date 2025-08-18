package dtos

import "go-mail/models"

type EmailJob struct {
	Preferences models.UserPreference `json:"preferences"`
	ThreadID    string                `json:"threadId"`
	Subject     string                `json:"subject"`
	Body        string                `json:"body"`
	Sender      string                `json:"sender"`
	Recipients  []string              `json:"recipients"`
}

type GmailListResponse struct {
	Messages           []GmailMessageMeta `json:"messages"`
	NextPageToken      string             `json:"nextPageToken"`
	ResultSizeEstimate int                `json:"resultSizeEstimate"`
}

type GmailMessageMeta struct {
	ID       string `json:"id"`
	ThreadID string `json:"threadId"`
}

type GmailMessageDetail struct {
	ID       string       `json:"id"`
	ThreadID string       `json:"threadId"`
	Payload  GmailPayload `json:"payload"`
	Size     int          `json:"sizeEstimate"`
}

type GmailPayload struct {
	PartID   string        `json:"partId"`
	MimeType string        `json:"mimeType"`
	Filename string        `json:"filename"`
	Headers  []GmailHeader `json:"headers"`
	Body     GmailBody     `json:"body"`
}

type GmailHeader struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}

type GmailBody struct {
	Size int    `json:"size"`
	Data string `json:"data"` // base64url
}

type GmailAttachment struct {
	AttachmentID string `json:"attachmentId"`
	Size         int    `json:"size"`
	Data         string `json:"data"` // base64url
}
