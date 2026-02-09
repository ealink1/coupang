package core

import (
	"context"
	"coupang/http_call"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"net/url"
	"strings"
	"time"
)

func (c *CoupangClient) generateHMAC(method, path, query string) string {
	dateTime := time.Now().UTC().Format("060102T150405Z")
	//m := strings.ToUpper(method)
	//q := query
	//if strings.HasPrefix(q, "?") {
	//	q = strings.TrimPrefix(q, "?")
	//}
	//message := dateTime + m + path + q

	message := dateTime + method + path + query

	h := hmac.New(sha256.New, []byte(c.SecretKey))
	h.Write([]byte(message))
	signature := hex.EncodeToString(h.Sum(nil))

	return fmt.Sprintf("CEA algorithm=HmacSHA256, access-key=%s, signed-date=%s, signature=%s",
		c.AccessKey, dateTime, signature)
}

func (c *CoupangClient) doRequest(ctx context.Context, method, path string, queryParams url.Values) (string, error) {
	query := queryParams.Encode()
	if query != "" {
		query = strings.ReplaceAll(query, "%3A00", ":00")
	}
	fullURL := fmt.Sprintf("%s://%s%s", Schema, Host, path)
	if query != "" {
		fullURL += "?" + query
	}

	fmt.Println("fullURL:", fullURL)

	if method != "GET" {
		return "", fmt.Errorf("only GET method is supported in this implementation")
	}

	authHeader := c.generateHMAC(method, path, query)

	headers := map[string]string{
		"Authorization":      authHeader,
		"Content-Type":       "application/json",
		"X-Requested-By":     c.VendorID,
		"X-MARKET":           "TW",
		"X-EXTENDED-TIMEOUT": "90000",
	}

	curlCmd := http_call.BuildCurl(fullURL, headers)
	fmt.Println("curlCmd:", curlCmd)

	respStr, err := http_call.HttpGet(fullURL, headers)
	if err != nil {
		return "", err
	}

	return respStr, nil
}

func (c *CoupangClient) doPostJSON(ctx context.Context, path string, body interface{}) (string, error) {
	query := ""
	fullURL := fmt.Sprintf("%s://%s%s", Schema, Host, path)
	fmt.Println("fullURL:", fullURL)

	authHeader := c.generateHMAC("POST", path, query)

	headers := map[string]string{
		"Authorization":      authHeader,
		"Content-Type":       http_call.ContentTypeJSON,
		"X-Requested-By":     c.VendorID,
		"X-MARKET":           "TW",
		"X-EXTENDED-TIMEOUT": "90000",
	}

	respBytes, err := http_call.HttpPost(fullURL, body, http_call.ContentTypeJSON, headers)
	if err != nil {
		return "", err
	}

	return string(respBytes), nil
}
