# go-apigw-http-adapter

Deserialize/serialize AWS API Gateway Lambda requests/responses to Go HTTP
requests/responses

## Badges

- Build:
  ![Build](https://github.com/harrisonhjones/go-apigw-http-adapter/workflows/Go/badge.svg)
- Report Card:
  [![Go Report Card](https://goreportcard.com/badge/harrisonhjones.com/go-apigw-http-adapter)](https://goreportcard.com/report/harrisonhjones.com/go-apigw-http-adapter)

## Links

- [pkg.go.dev (documentation)](https://pkg.go.dev/harrisonhjones.com/go-apigw-http-adapter)
- [Example Lambda Function](https://github.com/harrisonhjones/go-apigw-http-adapter-lambda-example)

## HTTP Adapter Lambda Example

Example Lambda function that transforms the incoming request, routes it to a
http.ServerMux, and then returns the transformed result.

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
