package httpadapter

import (
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
)

// Response configures the response to be returned by the API Gateway HTTP API for the request.
type Response struct {
	StatusCode      int               `json:"statusCode"`
	Headers         map[string]string `json:"headers"`
	Body            string            `json:"body"`
	IsBase64Encoded bool              `json:"isBase64Encoded,omitempty"`
	Cookies         []string          `json:"cookies"`
}

// TransformResponse transforms an http.Response to a Response.
func TransformResponse(res *http.Response, encRes func(*http.Response) bool) (*Response, error) {
	apigwRes := &Response{
		StatusCode: res.StatusCode,
		Headers:    map[string]string{},
	}

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %v", err)
	}

	if encRes != nil && encRes(res) {
		apigwRes.Body = base64.StdEncoding.EncodeToString(body)
		apigwRes.IsBase64Encoded = true
	} else {
		apigwRes.Body = string(body)
		apigwRes.IsBase64Encoded = false
	}

	// FYI: MultiValueHeaders aren't actually supported by HTTP APIs so don't use them.
	for k, v := range res.Header {
		// Cookies are handled further down.
		// Header names are case-insensitive.
		if strings.ToLower(k) == "set-cookie" {
			continue
		}
		apigwRes.Headers[k] = v[0]
	}

	for _, ck := range res.Cookies() {
		apigwRes.Cookies = append(apigwRes.Cookies, ck.Raw)
	}

	return apigwRes, nil
}
