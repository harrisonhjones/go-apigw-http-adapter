package restadapter

import (
	"context"
	"io/ioutil"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTransformRequest_HappyPath(t *testing.T) {
	req := Request{
		Path:       "/my/path",
		HTTPMethod: "POST",
		MultiValueHeaders: map[string][]string{
			"Header1":      {"value1"},
			"Header2":      {"value1", "value2"},
			"header-three": {"value1"}, // Non-canonical key.
			"Cookie":       {"cookie1=val1; cookie2=val2"},
		},
		MultiValueQueryStringParameters: map[string][]string{
			"parameter1": {"value1", "value2"},
			"parameter2": {"value"},
		},
		RequestContext: requestContext{
			DomainName: "example.com",
		},
	}

	tstCtx := context.Background()

	t.Run("NotEncoded", func(t *testing.T) {
		req.Body = "Hello World!"
		req.IsBase64Encoded = false

		httpReq, err := TransformRequest(tstCtx, req)

		if !assert.NoError(t, err, "failed to transform request") {
			return
		}

		assert.Equal(t, tstCtx, httpReq.Context())

		assert.Equal(t,
			http.Header{
				"Cookie":       []string{"cookie1=val1; cookie2=val2"},
				"Header-Three": []string{"value1"},
				"Header1":      []string{"value1"},
				"Header2":      []string{"value1", "value2"},
			},
			httpReq.Header)

		assert.Equal(
			t,
			"https://example.com/my/path?parameter1=value1&parameter1=value2&parameter2=value",
			httpReq.URL.String(),
		)

		ck1, err := httpReq.Cookie("cookie1")
		if !assert.NoError(t, err, "failed to get cookie1") {
			return // ck1 is nil if err != nil so return early to prevent panics
		}
		assert.Equal(t, "val1", ck1.Value)

		ck2, err := httpReq.Cookie("cookie2")
		if !assert.NoError(t, err, "failed to get cookie2") {
			return // ck2 is nil if err != nil so return early to prevent panics
		}
		assert.Equal(t, "val2", ck2.Value)

		b, err := ioutil.ReadAll(httpReq.Body)
		if !assert.NoError(t, err, "failed to read body") {
			return
		}

		assert.Equal(t, []byte("Hello World!"), b)
	})

	t.Run("Encoded", func(t *testing.T) {
		req.Body = "SGVsbG8gRW5jb2RlZCBXb3JsZCE=" // FYI: base64.StdEncoding.EncodeToString([]byte("Hello Encoded World!"))
		req.IsBase64Encoded = true

		httpReq, err := TransformRequest(tstCtx, req)

		if !assert.NoError(t, err, "failed to transform request") {
			return
		}

		b, err := ioutil.ReadAll(httpReq.Body)
		if !assert.NoError(t, err, "failed to read body") {
			return
		}

		assert.Equal(t, []byte("Hello Encoded World!"), b)
	})

	t.Run("IncorrectlyEncoded", func(t *testing.T) {
		req.Body = "blarg"
		req.IsBase64Encoded = true

		_, err := TransformRequest(tstCtx, req)

		assert.EqualError(t, err, "failed to decode body: illegal base64 data at input byte 4")
	})
}
