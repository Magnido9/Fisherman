package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"time"
)

// TIP <p>To run your code, right-click the code and select <b>Run</b>.</p> <p>Alternatively, click
// the <icon src="AllIcons.Actions.Execute"/> icon in the gutter and select the <b>Run</b> menu item from here.</p>
var (
	ctx         = context.Background()
	redisClient *redis.Client
)

func main() {
	// Initialize Redis client
	redisClient = redis.NewClient(&redis.Options{
		Addr:     "localhost:6377", // Redis on Docker or local
		Password: "",               // No password set
		DB:       0,                // Use default DB
	})

	// Optional: Ping Redis
	if err := redisClient.Ping(ctx).Err(); err != nil {
		log.Fatalf("Could not connect to Redis: %v", err)
	}

	router := gin.Default()
	router.GET("/readUrl", readHtmlApi)
	router.Run("localhost:8080")
}
func readHtmlApi(c *gin.Context) {
	url := c.Query("url")
	if url == "" {
		c.String(http.StatusBadRequest, "Missing 'url' query parameter")
		return
	}

	// Auto-prepend https:// if missing
	if !strings.HasPrefix(url, "http://") && !strings.HasPrefix(url, "https://") {
		url = "https://" + url
	}

	// Check Redis cache
	cached, err := redisClient.Get(ctx, url).Result()
	if err == nil {
		c.String(http.StatusOK, cached)
		return
	}

	html := readHtml(url)
	prompt := "This is html of a site:\n" + string(html) + "\nThe url of the site: " + url + "\n please determine if the site is a phishing site\n" + "Answer in boolean:true or false ONLY THE BOOL KEYWORD"

	result, err := callModel(prompt)
	if err != nil {
		log.Printf("Failed chatgpt query: %v", err)
		c.String(http.StatusInternalServerError, err.Error())
		return
	}

	// Cache result for 24 hours
	err = redisClient.Set(ctx, url, result, 24*time.Hour).Err()
	if err != nil {
		log.Printf("Failed to set Redis cache: %v", err)
		c.String(http.StatusInternalServerError, err.Error())
		return
	}
	c.String(http.StatusOK, result)
}
func readHtml(link string) string {
	res, err := http.Get(link)
	if err != nil {
		log.Fatal(err)
	}
	content, err := io.ReadAll(res.Body)
	res.Body.Close()
	if err != nil {
		log.Fatal(err)
	}
	return string(content)
}

func callModel(prompt string) (string, error) {
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
