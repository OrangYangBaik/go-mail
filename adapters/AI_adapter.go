package adapters

import (
	"bytes"
	"encoding/json"
	"fmt"
	"go-mail/dtos"
	"io"
	"net/http"
	"regexp"
	"strings"
)

type AIAdapter interface {
	//AIVerdictEmail(preference string, sender string, subject string, body string) (bool, error)
	AIVerdictEmail(preference string, sender string, subject string) (bool, error)
}

type aiAdapter struct {
	baseLLMUrl string
	model      string
}

func NewAIAdapter(baseLLMUrl, model string) AIAdapter {
	return &aiAdapter{
		baseLLMUrl: baseLLMUrl,
		model:      model,
	}
}

// func (a *aiAdapter) AIVerdictEmail(preference, sender, subject, body string) (bool, error) {
// 	prompt := fmt.Sprintf(
// 		`You are an email filtering assistant.
// User preference: %s
// Email details:
// From: %s
// Subject: %s
// Body: %s

// Based on the user preference, respond only with a JSON object in the format:
// {"forward": true} or {"forward": false}`,
// 		preference, sender, subject, body,
// 	)

// 	requestBody, _ := json.Marshal(map[string]interface{}{
// 		"model":  a.model,
// 		"prompt": prompt,
// 		"stream": false,
// 	})

// 	resp, err := http.Post(a.baseLLMUrl+"/v1/completions", "application/json", bytes.NewBuffer(requestBody))
// 	if err != nil {
// 		return false, err
// 	}
// 	defer resp.Body.Close()

// 	respBody, _ := io.ReadAll(resp.Body)
// 	if resp.StatusCode != http.StatusOK {
// 		return false, fmt.Errorf("LM Studio API error: %s", string(respBody))
// 	}

// 	var lmResp struct {
// 		Choices []struct {
// 			Text string `json:"text"`
// 		} `json:"choices"`
// 	}
// 	if err := json.Unmarshal(respBody, &lmResp); err != nil {
// 		return false, err
// 	}

// 	if len(lmResp.Choices) == 0 {
// 		return false, fmt.Errorf("no response from AI")
// 	}

// 	var aiRes dtos.AIResponse
// 	if err := json.Unmarshal([]byte(lmResp.Choices[0].Text), &aiRes); err != nil {
// 		return false, fmt.Errorf("invalid AI response: %s", lmResp.Choices[0].Text)
// 	}

// 	return aiRes.Forward, nil
// }

func (a *aiAdapter) AIVerdictEmail(preference, sender, subject string) (bool, error) {
	systemMessage := "You are an email filtering assistant. Your task is to classify emails based on user preference. Respond with exactly one of these two strings: {\"forward\": true} or {\"forward\": false}. Do not include any extra text, formatting, or code fences."

	userMessage := fmt.Sprintf(
		`User preference: %s

Email details:
From: %s
Subject: %s

Decide if this email should be forwarded.`,
		preference, sender, subject,
	)

	requestBody, _ := json.Marshal(map[string]interface{}{
		"model": a.model,
		"messages": []map[string]interface{}{
			{
				"role":    "system",
				"content": systemMessage,
			},
			{
				"role":    "user",
				"content": userMessage,
			},
		},
		"temperature": 0.3,
		"stream":      false,
	})

	resp, err := http.Post(a.baseLLMUrl+"/v1/chat/completions", "application/json", bytes.NewBuffer(requestBody))
	if err != nil {
		return false, err
	}
	defer resp.Body.Close()
	respBody, _ := io.ReadAll(resp.Body)
	if resp.StatusCode != http.StatusOK {
		return false, fmt.Errorf("LLM API error: %s", string(respBody))
	}

	var lmResp struct {
		Choices []struct {
			Message struct {
				Content string `json:"content"`
			} `json:"message"`
		} `json:"choices"`
	}

	if err := json.Unmarshal(respBody, &lmResp); err != nil {
		return false, err
	}
	if len(lmResp.Choices) == 0 {
		return false, fmt.Errorf("no response from AI")
	}

	content := lmResp.Choices[0].Message.Content

	re := regexp.MustCompile(`(?s)<think>.*?</think>`)
	content = re.ReplaceAllString(content, "")

	content = strings.ReplaceAll(content, "\n", "")
	content = strings.ReplaceAll(content, "\r", "")
	content = strings.TrimSpace(content)

	var aiRes dtos.AIResponse
	if err := json.Unmarshal([]byte(content), &aiRes); err != nil {
		return false, fmt.Errorf("invalid AI response: %s", content)
	}
	return aiRes.Forward, nil
}
