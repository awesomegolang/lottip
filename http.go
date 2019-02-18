package main

import (
	"bytes"
	"fmt"
	"net/http"
)

func apiUrl(apiHost string, path string) string {
	return fmt.Sprintf("%s%s", apiHost, path)
}

func postData(apiUrl string, data []byte) {
	_, _ = http.Post(apiUrl, "application/json", bytes.NewBuffer(data))
}
