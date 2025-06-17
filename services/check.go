package services

import (
	"Fisherman/cache"
	"Fisherman/models"
	"Fisherman/utils"
	"errors"
	"io/ioutil"
	"log"
	"strconv"
	"strings"
	"time"
)

func CheckUrl(url string) models.UrlCheckResult {
	log.Println("Checking Url:", url)
	if url == "" {
		return models.UrlCheckResult{false, "missing url"}
	}
	if !strings.HasPrefix(url, "http://") && !strings.HasPrefix(url, "https://") {
		url = "https://" + url
	}

	// Redis cache
	cached, err := cache.Get(url)
	if err == nil {
		log.Println("Using cached url:", url, cached)
		val, _ := strconv.ParseBool(cached)
		return models.UrlCheckResult{val, ""}
	} else if !errors.Is(err, cache.ErrCacheMiss) {
		log.Println("Error checking url:", err)
		return models.UrlCheckResult{false, err.Error()}
	}

	html, err := utils.ReadHtml(url)
	if err != nil {
		log.Println("Error reading url:", err)
		return models.UrlCheckResult{false, err.Error()}
	}
	data, _ := ioutil.ReadFile("services/prompt.txt")
	prompt := string(data)
	prompt = strings.ReplaceAll(prompt, "<html>", html)
	prompt = strings.ReplaceAll(prompt, "<url>", url)

	log.Println("Calling model for url:", url)

	resp, err := CallModel(prompt)
	if err != nil {
		log.Println("Error calling model for url", url, err)
		return models.UrlCheckResult{false, err.Error()}
	}

	val, err := strconv.ParseBool(resp)
	if err != nil {
		log.Println("Error parsing model answer for url ", url, ", with response ", resp, " retrying...")
		return CheckUrl(url) // Retry
	}

	// Cache result
	cache.Set(url, resp, 24*time.Hour)
	log.Println("Result for url:", url, ":", resp)
	return models.UrlCheckResult{val, ""}
}
