package handlers

import (
	"Fisherman/models"
	"Fisherman/services"
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
)

func ReadHtmlApi(c *gin.Context) {
	urls := c.Query("urls")
	urlArr := strings.Split(urls, ",")
	resultsAsync := make(chan models.AsyncUrlCheckResult, len(urlArr))
	results := make([]models.UrlCheckResult, len(urlArr))

	for i, url := range urlArr {
		go func(i int, url string) {
			res := services.CheckUrl(url)
			resultsAsync <- models.AsyncUrlCheckResult{Index: i, Result: res}
		}(i, url)
	}

	for i := 0; i < len(urlArr); i++ {
		result := <-resultsAsync
		results[result.Index] = result.Result
	}
	c.JSON(http.StatusOK, results)
}
