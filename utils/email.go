package utils

import (
	"encoding/base64"
	"go-mail/dtos"
	"strings"
)

func GetHeader(headers []dtos.GmailHeader, name string) string {
	for _, h := range headers {
		if strings.EqualFold(h.Name, name) {
			return h.Value
		}
	}
	return ""
}

func ParseRecipients(headerValue string) []string {
	if headerValue == "" {
		return nil
	}
	parts := strings.Split(headerValue, ",")
	for i := range parts {
		parts[i] = strings.TrimSpace(parts[i])
	}
	return parts
}

func ExtractEmailBody(email dtos.GmailMessageDetail) string {
	if email.Payload.Body.Data != "" {
		decoded, _ := base64.URLEncoding.DecodeString(email.Payload.Body.Data)
		return string(decoded)
	}
	return ""
}
