package httpadapter

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

// Request contains all relevant API Gateway HTTP API request data needed to transform it into a http.Request.
type Request struct {
	Version         string            `json:"version"`
	RawQueryString  string            `json:"rawQueryString"`
	Cookies         []string          `json:"cookies,omitempty"`
	Headers         map[string]string `json:"headers"`
	RequestContext  RequestContext    `json:"requestContext"`
	Body            string            `json:"body,omitempty"`
	IsBase64Encoded bool              `json:"isBase64Encoded"`
}

// RequestContext contains all relevant data needed for Request transformation.
type RequestContext struct {
	DomainName string             `json:"domainName"`
	HTTP       RequestContextHTTP `json:"http"`
}

// RequestContextHTTP contains all relevant data needed for Request transformation.
type RequestContextHTTP struct {
	Method string `json:"method"`
	Path   string `json:"path"`
}

// TransformRequest transforms a Request to an http.Request.
func TransformRequest(ctx context.Context, req Request) (*http.Request, error) {
	if req.Version != "2.0" {
		return nil, fmt.Errorf("unsupported version %q", req.Version)
	}

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

	rawUrl := "https://" + req.RequestContext.DomainName + req.RequestContext.HTTP.Path
	if req.RawQueryString != "" {
		rawUrl = rawUrl + "?" + req.RawQueryString
	}
	u, err := url.Parse(rawUrl)
	if err != nil {
		return nil, fmt.Errorf("failed to parse url: %v", err)
	}

	hReq, err := http.NewRequest(req.RequestContext.HTTP.Method, u.String(), body)
	if err != nil {
		return nil, fmt.Errorf("failed to create new http request: %v", err)
	}

	hReq = hReq.WithContext(ctx)

	for k, v := range req.Headers {
		parts := strings.Split(v, ",")
		for _, part := range parts {
			hReq.Header.Add(k, part)
		}
	}

	if len(req.Cookies) > 0 {
		hReq.Header.Set("Cookie", strings.Join(req.Cookies, "; "))
	}

	return hReq, nil
}
