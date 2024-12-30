package main

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"
)

var (
	DEFAULT_TIMEOUT = 30 // seconds
)

type WebTaskResponse struct {
	TaskId        string
	StatusCode    int
	ContentLength int64
	Body          []byte
	Headers       map[string]string
}

type WebTask struct {
	Url        string
	Method     string
	Headers    map[string]string
	Parameters map[string]string
	Body       []byte
}

type RequestManager struct {
	request  WebTask
	response WebTaskResponse
}

func newRequestManager(task WebTask) *RequestManager {
	r := &RequestManager{
		request:  task,
		response: WebTaskResponse{},
	}
	r.createTaskId()
	return r
}

func (r *RequestManager) createTaskId() {
	taskId := CreateTask()
	r.response.TaskId = taskId
}

func (r *RequestManager) MakeRequest() error {
	client := &http.Client{
		Transport: http.DefaultTransport,
		Timeout:   time.Duration(DEFAULT_TIMEOUT),
	}
	request, err := buildRequest(r.request.Url, r.request.Method, r.request.Headers, r.request.Body, r.request.Parameters)
	if err != nil {
		return err
	}

	response, err := client.Do(request)
	defer response.Body.Close() // Ensure the response body is closed

	if err != nil {
		return err
	}

	return r.ResponseDetails(response)
}

// ResponseDetails extracts headers and the body from the HTTP response
func (r *RequestManager) ResponseDetails(resp *http.Response) error {

	r.response.StatusCode = resp.StatusCode
	r.response.ContentLength = resp.ContentLength

	// Extract headers
	headers := make(map[string]string)
	for key, values := range resp.Header {
		headers[key] = values[0] // Take the first value (most headers have a single value)
	}
	r.response.Headers = headers

	// Read the body
	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	r.response.Body = bodyBytes

	return nil
}

func buildRequest(URL string, method string, headers map[string]string, body []byte, params map[string]string) (*http.Request, error) {
	// Build query parameters for GET Requests
	if len(params) > 0 {
		query := url.Values{}
		for key, value := range params {
			query.Add(key, value)
		}
		URL = fmt.Sprintf("%s?%s", URL, query.Encode())
	}
	// Build Body for POST Requests
	var requestBody *bytes.Reader
	if len(body) == 0 {
		requestBody = bytes.NewReader(nil)
	} else {
		requestBody = bytes.NewReader(body)
	}

	request, err := http.NewRequest(method, URL, requestBody)
	if err != nil {
		return request, err
	}
	// Add the headers
	for headerName, headerValue := range headers {
		request.Header.Set(headerName, headerValue) // using Set to ensure this header is overwritten
	}

	return request, err
}
