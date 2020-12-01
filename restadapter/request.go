package restadapter

import (
	"bytes"
	"context"
	"encoding/base64"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
)

// Request contains all relevant API Gateway REST API request data needed to transform it into a http.Request.
// MultiValueHeaders are supported by REST APIs and always contain all Headers so Headers can be safely ignored.
// MultiValueQueryStringParameters are supported by REST APIs and always contain all QueryStringParameters so
// QueryStringParameters can be safely ignored.
// https://docs.aws.amazon.com/apigateway/latest/developerguide/set-up-lambda-proxy-integrations.html
type Request struct {
	Path                            string              `json:"path"` // The url path for the caller
	HTTPMethod                      string              `json:"httpMethod"`
	MultiValueHeaders               map[string][]string `json:"multiValueHeaders"`
	MultiValueQueryStringParameters map[string][]string `json:"multiValueQueryStringParameters"`
	RequestContext                  requestContext      `json:"requestContext"`
	Body                            string              `json:"body"`
	IsBase64Encoded                 bool                `json:"isBase64Encoded,omitempty"`
}

type requestContext struct {
	DomainName string `json:"domainName"`
}

// TransformRequest transforms a Request to an http.Request.
func TransformRequest(ctx context.Context, req Request) (*http.Request, error) {
	// Mirror how http.Request bodies normally behave.
	// From the docs:
	// For server requests, the Request Body is always non-nil
	// but will return EOF immediately when no body is present.
	var body io.Reader
	if req.IsBase64Encoded {
		b, err := base64.StdEncoding.DecodeString(req.Body)
		if err != nil {
			return nil, fmt.Errorf("failed to decode body: %v", err)
		}
		body = bytes.NewBuffer(b)
	} else {
		body = strings.NewReader(req.Body)
	}

	u, err := url.Parse("https://" + req.RequestContext.DomainName + req.Path)
	if err != nil {
		return nil, fmt.Errorf("failed to parse url: %v", err)
	}

	qValues := u.Query()
	for k, parts := range req.MultiValueQueryStringParameters {
		for _, part := range parts {
			qValues.Add(k, part)
		}
	}
	u.RawQuery = qValues.Encode()

	hReq, err := http.NewRequest(req.HTTPMethod, u.String(), body)
	if err != nil {
		return nil, fmt.Errorf("failed to create new http request: %v", err)
	}

	hReq = hReq.WithContext(ctx)

	hReq.Header = req.MultiValueHeaders

	return hReq, nil
}
