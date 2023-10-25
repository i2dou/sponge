// Package gohttp is http request client, which only supports returning json format.
package gohttp

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"
)

const defaultTimeout = 10 * time.Second

// Request HTTP request
type Request struct {
	customRequest func(req *http.Request, data *bytes.Buffer) // used to define HEADER, e.g. to add sign, etc.
	url           string
	params        map[string]interface{} // parameters after URL
	body          string                 // Body data
	bodyJSON      interface{}            // JSON marshal body data
	timeout       time.Duration          // Client timeout
	headers       map[string]string

	request  *http.Request
	response *Response
	method   string
	err      error
}

// Response HTTP response
type Response struct {
	*http.Response
	err error
}

// -----------------------------------  Request  -----------------------------------

// Reset set all fields to default value, use at pool
func (req *Request) Reset() {
	req.params = nil
	req.body = ""
	req.bodyJSON = nil
	req.timeout = 0
	req.headers = nil

	req.request = nil
	req.response = nil
	req.method = ""
	req.err = nil
}

// SetURL set URL
func (req *Request) SetURL(path string) *Request {
	req.url = path
	return req
}

// SetParams parameters after setting the URL
func (req *Request) SetParams(params map[string]interface{}) *Request {
	if req.params == nil {
		req.params = params
	} else {
		for k, v := range params {
			req.params[k] = v
		}
	}
	return req
}

// SetParam parameters after setting the URL
func (req *Request) SetParam(k string, v interface{}) *Request {
	if req.params == nil {
		req.params = make(map[string]interface{})
	}
	req.params[k] = v
	return req
}

// SetBody set body data
func (req *Request) SetBody(body string) *Request {
	req.body = body
	return req
}

// SetJSONBody set Body data, JSON format
func (req *Request) SetJSONBody(body interface{}) *Request {
	req.bodyJSON = body
	return req
}

// SetTimeout set timeout
func (req *Request) SetTimeout(t time.Duration) *Request {
	req.timeout = t
	return req
}

// SetContentType set ContentType
func (req *Request) SetContentType(a string) *Request {
	req.SetHeader("Content-Type", a)
	return req
}

// SetHeader set the value of the request header
func (req *Request) SetHeader(k, v string) *Request {
	if req.headers == nil {
		req.headers = make(map[string]string)
	}
	req.headers[k] = v
	return req
}

// SetHeaders set the value of Request Headers
func (req *Request) SetHeaders(headers map[string]string) *Request {
	if req.headers == nil {
		req.headers = make(map[string]string)
	}
	for k, v := range headers {
		req.headers[k] = v
	}
	return req
}

// CustomRequest customize request, e.g. add sign, set header, etc.
func (req *Request) CustomRequest(f func(req *http.Request, data *bytes.Buffer)) *Request {
	req.customRequest = f
	return req
}

// GET send a GET request
func (req *Request) GET() (*Response, error) {
	req.method = http.MethodGet
	return req.pull()
}

// DELETE send a DELETE request
func (req *Request) DELETE() (*Response, error) {
	req.method = http.MethodDelete
	return req.pull()
}

// POST send a POST request
func (req *Request) POST() (*Response, error) {
	req.method = http.MethodPost
	return req.push()
}

// PUT send a PUT request
func (req *Request) PUT() (*Response, error) {
	req.method = http.MethodPut
	return req.push()
}

// PATCH send PATCH requests
func (req *Request) PATCH() (*Response, error) {
	req.method = http.MethodPatch
	return req.push()
}

// Do a request
func (req *Request) Do(method string, data interface{}) (*Response, error) {
	req.method = method

	switch method {
	case http.MethodGet, http.MethodDelete:
		if data != nil {
			if params, ok := data.(map[string]interface{}); ok { //nolint
				req.SetParams(params)
			} else {
				req.err = errors.New("params is not a map[string]interface{}")
				return nil, req.err
			}
		}

		return req.pull()

	case http.MethodPost, http.MethodPut, http.MethodPatch:
		if data != nil {
			req.SetJSONBody(data)
		}

		return req.push()
	}

	req.err = errors.New("unknow method " + method)
	return nil, req.err
}

func (req *Request) pull() (*Response, error) {
	val := ""
	if len(req.params) > 0 {
		values := url.Values{}
		for k, v := range req.params {
			values.Add(k, fmt.Sprintf("%v", v))
		}
		val += values.Encode()
	}

	if val != "" {
		if strings.Contains(req.url, "?") {
			req.url += "&" + val
		} else {
			req.url += "?" + val
		}
	}

	var buf *bytes.Buffer
	if req.customRequest != nil {
		buf = bytes.NewBufferString(val)
	}

	return req.send(nil, buf)
}

func (req *Request) push() (*Response, error) {
	var buf *bytes.Buffer

	if req.bodyJSON != nil {
		body, err := json.Marshal(req.bodyJSON)
		if err != nil {
			req.err = err
			return nil, req.err
		}
		buf = bytes.NewBuffer(body)
	} else {
		buf = bytes.NewBufferString(req.body)
	}

	return req.send(buf, buf)
}

func (req *Request) send(body io.Reader, buf *bytes.Buffer) (*Response, error) {
	req.request, req.err = http.NewRequest(req.method, req.url, body)
	if req.err != nil {
		return nil, req.err
	}

	if req.customRequest != nil {
		req.customRequest(req.request, buf)
	}

	if req.headers != nil {
		for k, v := range req.headers {
			req.request.Header.Add(k, v)
		}
	}

	if req.timeout < 1 {
		req.timeout = defaultTimeout
	}

	client := http.Client{Timeout: req.timeout}
	resp := new(Response)
	resp.Response, resp.err = client.Do(req.request)

	req.response = resp
	req.err = resp.err

	return resp, resp.err
}

// Response return response
func (req *Request) Response() (*Response, error) {
	if req.err != nil {
		return nil, req.err
	}
	return req.response, req.response.Error()
}

// -----------------------------------  Response -----------------------------------

// Error return err
func (resp *Response) Error() error {
	return resp.err
}

// BodyString returns the body data of the HttpResponse
func (resp *Response) BodyString() (string, error) {
	if resp.err != nil {
		return "", resp.err
	}
	body, err := resp.ReadBody()
	return string(body), err
}

// ReadBody returns the body data of the HttpResponse
func (resp *Response) ReadBody() ([]byte, error) {
	if resp.err != nil {
		return []byte{}, resp.err
	}

	if resp.Response == nil {
		return []byte{}, errors.New("nil")
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return []byte{}, err
	}

	resp.Body = io.NopCloser(bytes.NewBuffer(body))
	return body, nil
}

// BindJSON parses the response's body as JSON
func (resp *Response) BindJSON(v interface{}) error {
	if resp.err != nil {
		return resp.err
	}
	body, err := resp.ReadBody()
	if err != nil {
		return err
	}
	return json.Unmarshal(body, v)
}

// -------------------------------------------------------------------------------------------------

// Simple crud function, no support for setting header, timeout, etc.

// Get request, return custom json format
func Get(result interface{}, urlStr string, params ...KV) error {
	var pms KV
	if len(params) > 0 {
		pms = params[0]
	}
	return gDo("GET", result, urlStr, pms)
}

// Delete request, return custom json format
func Delete(result interface{}, urlStr string, params ...KV) error {
	var pms KV
	if len(params) > 0 {
		pms = params[0]
	}
	return gDo("DELETE", result, urlStr, pms)
}

// Post request, return custom json format
func Post(result interface{}, urlStr string, body interface{}) error {
	return do("POST", result, urlStr, body)
}

// Put request, return custom json format
func Put(result interface{}, urlStr string, body interface{}) error {
	return do("PUT", result, urlStr, body)
}

// Patch request, return custom json format
func Patch(result interface{}, urlStr string, body interface{}) error {
	return do("PATCH", result, urlStr, body)
}

var requestErr = func(err error) error { return fmt.Errorf("request error, err=%v", err) }
var jsonParseErr = func(err error) error { return fmt.Errorf("json parsing error, err=%v", err) }
var notOKErr = func(resp *Response) error {
	body, err := resp.ReadBody()
	if err != nil {
		return err
	}
	if len(body) > 500 {
		body = append(body[:500], []byte(" ......")...)
	}
	return fmt.Errorf("statusCode=%d, body=%s", resp.StatusCode, body)
}

func do(method string, result interface{}, urlStr string, body interface{}, params ...KV) error {
	if result == nil {
		return fmt.Errorf("params 'result' is nil")
	}

	req := &Request{}
	req.SetURL(urlStr)
	req.SetContentType("application/json")
	if len(params) > 0 {
		req.SetParams(params[0])
	}
	req.SetJSONBody(body)

	var resp *Response
	var err error
	switch method {
	case "POST":
		resp, err = req.POST()
	case "PUT":
		resp, err = req.PUT()
	case "PATCH":
		resp, err = req.PATCH()
	}
	if err != nil {
		return requestErr(err)
	}
	defer resp.Body.Close() //nolint

	if resp.StatusCode != 200 {
		return notOKErr(resp)
	}

	err = resp.BindJSON(result)
	if err != nil {
		return jsonParseErr(err)
	}

	return nil
}

func gDo(method string, result interface{}, urlStr string, params KV) error {
	req := &Request{}
	req.SetURL(urlStr)
	req.SetParams(params)

	var resp *Response
	var err error
	switch method {
	case "GET":
		resp, err = req.GET()
	case "DELETE":
		resp, err = req.DELETE()
	}
	if err != nil {
		return requestErr(err)
	}
	defer resp.Body.Close() //nolint

	if resp.StatusCode != 200 {
		return notOKErr(resp)
	}

	err = resp.BindJSON(result)
	if err != nil {
		return jsonParseErr(err)
	}

	return nil
}

// StdResult standard return data
type StdResult struct {
	Code int         `json:"code"`
	Msg  string      `json:"msg"`
	Data interface{} `json:"data,omitempty"`
}

// KV string:interface{}
type KV map[string]interface{}
