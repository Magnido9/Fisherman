package services

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

func CallModel(prompt string) (string, error) {
	log.Println("Calling model with prompt: ", prompt)
	apiKey := "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" // 🔒 Replace with your real API key
	referer := "https://fisherman.com"                                                    // Optional
	siteTitle := "fisherman"                                                              // Optional
	url := "https://openrouter.ai/api/v1/chat/completions"

	// Prepare request body
	requestBody := map[string]interface{}{
		"model": "deepseek/deepseek-r1-0528-qwen3-8b:free",
		"messages": []map[string]string{
			{
				"role":    "user",
				"content": prompt,
			},
		},
	}

	// Encode JSON
	jsonData, err := json.Marshal(requestBody)
	if err != nil {
		return "", err
	}

	// Create HTTP request
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return "", err
	}

	// Set headers
	req.Header.Set("Authorization", "Bearer "+apiKey)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("HTTP-Referer", referer) // Optional
	req.Header.Set("X-Title", siteTitle)    // Optional

	// Execute request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	// Read and parse response
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	// Decode JSON response
	var responseObject map[string]interface{}
	if err := json.Unmarshal(body, &responseObject); err != nil {
		return "", err
	}

	// Extract the assistant's reply from response
	if choices, ok := responseObject["choices"].([]interface{}); ok && len(choices) > 0 {
		if choice, ok := choices[0].(map[string]interface{}); ok {
			if message, ok := choice["message"].(map[string]interface{}); ok {
				if content, ok := message["content"].(string); ok {
					return content, nil
				}
			}
		}
	}

	// If the structure is unexpected
	return "", fmt.Errorf("unexpected response format: %s", string(body))
}
