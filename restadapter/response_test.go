package restadapter

import (
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTransformResponse_HappyPathNotEncoded(t *testing.T) {
	recorder := httptest.NewRecorder()
	http.SetCookie(recorder, &http.Cookie{Name: "cookie1-name", Value: "cookie1-value"})
	http.SetCookie(recorder, &http.Cookie{Name: "cookie2-name", Value: "cookie2-value"})
	recorder.Header().Add("key1", "val1")
	recorder.Header().Add("key1", "val2")
	recorder.Header().Add("key2", "val2")
	recorder.WriteHeader(201)
	_, err := recorder.WriteString("Hello World!")
	assert.NoError(t, err, "failed to write string to test recorder")

	response, err := TransformResponse(recorder.Result(), nil)
	if !assert.NoError(t, err, "failed to transform request") {
		return
	}

	assert.Equal(t,
		&Response{
			StatusCode: 201,
			MultiValueHeaders: map[string][]string{
				"Key1":       {"val1", "val2"},
				"Key2":       {"val2"},
				"Set-Cookie": {"cookie1-name=cookie1-value", "cookie2-name=cookie2-value"},
			},
			Body:            "Hello World!",
			IsBase64Encoded: false,
		},
		response)
}

func TestTransformResponse_HappyPathNotEncodedWithEncRes(t *testing.T) {
	recorder := httptest.NewRecorder()
	http.SetCookie(recorder, &http.Cookie{Name: "cookie1-name", Value: "cookie1-value"})
	http.SetCookie(recorder, &http.Cookie{Name: "cookie2-name", Value: "cookie2-value"})
	recorder.Header().Add("key1", "val1")
	recorder.Header().Add("key1", "val2")
	recorder.Header().Add("key2", "val2")
	recorder.WriteHeader(201)
	_, err := recorder.WriteString("Hello World!")
	assert.NoError(t, err, "failed to write string to test recorder")

	response, err := TransformResponse(recorder.Result(), func(response *http.Response) bool {
		// FYI: Shallow check that it is the same http.Response.
		assert.Equal(t, "val1", response.Header.Get("key1"))
		return false
	})
	if !assert.NoError(t, err, "failed to transform request") {
		return
	}

	assert.Equal(t,
		&Response{
			StatusCode: 201,
			MultiValueHeaders: map[string][]string{
				"Key1":       {"val1", "val2"},
				"Key2":       {"val2"},
				"Set-Cookie": {"cookie1-name=cookie1-value", "cookie2-name=cookie2-value"},
			},
			Body:            "Hello World!",
			IsBase64Encoded: false,
		},
		response)
}

func TestTransformResponse_HappyPathEncoded(t *testing.T) {
	recorder := httptest.NewRecorder()
	http.SetCookie(recorder, &http.Cookie{Name: "cookie1-name", Value: "cookie1-value"})
	http.SetCookie(recorder, &http.Cookie{Name: "cookie2-name", Value: "cookie2-value"})
	recorder.Header().Add("key1", "val1")
	recorder.Header().Add("key1", "val2")
	recorder.Header().Add("key2", "val2")
	recorder.WriteHeader(201)
	_, err := recorder.WriteString("Hello Encoded World!")
	assert.NoError(t, err, "failed to write string to test recorder")

	response, err := TransformResponse(recorder.Result(), func(response *http.Response) bool {
		// FYI: Shallow check that it is the same http.Response.
		assert.Equal(t, "val1", response.Header.Get("key1"))
		return true
	})
	if !assert.NoError(t, err, "failed to transform request") {
		return
	}

	assert.Equal(t,
		&Response{
			StatusCode: 201,
			MultiValueHeaders: map[string][]string{
				"Key1":       {"val1", "val2"},
				"Key2":       {"val2"},
				"Set-Cookie": {"cookie1-name=cookie1-value", "cookie2-name=cookie2-value"},
			},
			Body:            "SGVsbG8gRW5jb2RlZCBXb3JsZCE=",
			IsBase64Encoded: true,
		},
		response)
}

func TestTransformResponse_BadBody(t *testing.T) {
	_, err := TransformResponse(&http.Response{
		Body: ioutil.NopCloser(&FailingReader{}),
	}, nil)
	assert.EqualError(t, err, "failed to read response body: boom")
}

type FailingReader struct{}

func (f FailingReader) Read([]byte) (n int, err error) {
	return 0, fmt.Errorf("boom")
}

var _ io.Reader = &FailingReader{}
