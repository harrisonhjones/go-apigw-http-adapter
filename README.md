# go-apigw-http-adapter

Transform AWS API Gateway Lambda requests and responses to Go HTTP requests and
responses.

## Badges

- Build:
  ![Build](https://github.com/harrisonhjones/go-apigw-http-adapter/workflows/Go/badge.svg)
- Report Card:
  [![Go Report Card](https://goreportcard.com/badge/harrisonhjones.com/go-apigw-http-adapter)](https://goreportcard.com/report/harrisonhjones.com/go-apigw-http-adapter)

## Releases

- 0.x - Current - Under review / development. Candidate for 1.x.

## Links

- [pkg.go.dev (documentation)](https://pkg.go.dev/harrisonhjones.com/go-apigw-http-adapter)
- [Example Lambda Function](https://github.com/harrisonhjones/go-apigw-http-adapter-lambda-example)

## Limitations

1. Only handles HTTP APIs. REST APIs are not yet supported (PRs welcome!).

## Goals

1. _Once version 1 has released_: maintain the
   [Go compatibility promise](https://golang.org/doc/go1compat) as much as
   possible.
1. As few dependencies as possible.
   1. Looking at you
      [awslabs/aws-lambda-go-api-proxy](https://github.com/awslabs/aws-lambda-go-api-proxy).
   1. Currently only a single direct dependency on `github.com/stretchr/testify`
      for testing.

## HTTP Adapter Lambda Example

Example Lambda function that transforms the incoming request, routes it to a
http.ServerMux, and then returns the transformed result. See
[Example Lambda Function](https://github.com/harrisonhjones/go-apigw-http-adapter-lambda-example)
for a working example.

```go
package main

import (
	"context"
	"net/http"
	"net/http/httptest"

	"github.com/aws/aws-lambda-go/lambda"
	"harrisonhjones.com/go-apigw-http-adapter/httpadapter"
)

func HandleRequest(ctx context.Context, req httpadapter.Request) (*httpadapter.Response, error) {
	// FYI: Request transformation.
	httpReq, err := httpadapter.TransformRequest(ctx, req)
	if err != nil {
		return nil, err
	}

	// FYI: Handle Request.
	httpRec := httptest.NewRecorder()

	mux := http.NewServeMux()
	// TODO: Add your own handlers here.
	mux.ServeHTTP(httpRec, httpReq)

	// FYI: Response transformation.
	httpRes, err := httpadapter.TransformResponse(httpRec.Result(), func(response *http.Response) bool {
		// FYI: Here you might inspect the response Content-Type to determine if the response should be encoded or not.
		return false // FYI: Don't encode the response.
	})
	if err != nil {
		return nil, err
	}

	return httpRes, nil
}

func main() {
	lambda.Start(HandleRequest)
}
```
