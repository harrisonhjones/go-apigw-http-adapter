package restadapter

import (
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"net/http"
)

// Response configures the response to be returned by the API Gateway REST API for the request.
// MultiValueHeaders are supported by REST APIs and are merged into Headers so Headers can be safely ignored.
// https://docs.aws.amazon.com/apigateway/latest/developerguide/set-up-lambda-proxy-integrations.html
type Response struct {
	StatusCode        int                 `json:"statusCode"`
	MultiValueHeaders map[string][]string `json:"multiValueHeaders"`
	Body              string              `json:"body"`
	IsBase64Encoded   bool                `json:"isBase64Encoded,omitempty"`
}

// TransformResponse transforms an http.Response to a Response.
func TransformResponse(res *http.Response, encRes func(*http.Response) bool) (*Response, error) {
	apigwRes := &Response{
		StatusCode: res.StatusCode,
		MultiValueHeaders:    map[string][]string{},
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

	apigwRes.MultiValueHeaders = res.Header

	return apigwRes, nil
}
