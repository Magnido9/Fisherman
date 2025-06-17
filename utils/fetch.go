package utils

import (
	"io"
	"net/http"
)

func ReadHtml(link string) (string, error) {
	res, err := http.Get(link)
	if err != nil {
		return "", err
	}
	defer res.Body.Close()
	body, err := io.ReadAll(res.Body)
	return string(body), err
}
