package http_call

import (
	"io"
	"net/http"
	"strings"
	"time"
)

// 发起 GET 请求
func HttpGet(urlStr string, headers map[string]string) (string, error) {
	client := &http.Client{Timeout: 10 * time.Second}

	req, err := http.NewRequest("GET", urlStr, nil)
	if err != nil {
		return "", err
	}

	for k, v := range headers {
		req.Header.Set(k, v)
	}

	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	return string(body), err
}

func BuildCurl(urlStr string, headers map[string]string) string {
	var b strings.Builder
	b.WriteString("curl -X GET")
	b.WriteString(" '")
	b.WriteString(urlStr)
	b.WriteString("'")
	for k, v := range headers {
		b.WriteString(" -H '")
		b.WriteString(k)
		b.WriteString(": ")
		b.WriteString(v)
		b.WriteString("'")
	}
	return b.String()
}
