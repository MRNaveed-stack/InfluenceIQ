package services

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
)

var geminiKey = os.Getenv("GEMINI_API_KEY")

const geminiEndpoint = "https://generativelanguage.googleapis.com/v1/models/gemini-2.5-flash:generateContent"

type GeminiRequest struct {
	Input string `json:"input"` // single string, not a list
}

type GeminiResponse struct {
	Candidates []struct {
		Content struct {
			Parts []struct {
				Text string `json:"text"`
			} `json:"parts"`
		} `json:"content"`
	} `json:"candidates"`
}

func CallGemini(prompt string) (string, error) {
	if geminiKey == "" {
		return "", errors.New("GEMINI_API_KEY not set")
	}

	url := geminiEndpoint + "?key=" + geminiKey

	payload := map[string]interface{}{
		"contents": []map[string]interface{}{
			{
				"parts": []map[string]interface{}{
					{"text": prompt},
				},
			},
		},
	}

	body, _ := json.Marshal(payload)

	req, _ := http.NewRequest("POST", url, bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	respBody, _ := ioutil.ReadAll(resp.Body)

	if resp.StatusCode != 200 {
		return "", fmt.Errorf("gemini api error: %s", string(respBody))
	}

	var result struct {
		Candidates []struct {
			Content struct {
				Parts []struct {
					Text string `json:"text"`
				} `json:"parts"`
			} `json:"content"`
		} `json:"candidates"`
	}

	if err := json.Unmarshal(respBody, &result); err != nil {
		return "", err
	}

	if len(result.Candidates) == 0 || len(result.Candidates[0].Content.Parts) == 0 {
		return "", errors.New("no response text returned")
	}

	return result.Candidates[0].Content.Parts[0].Text, nil
}
