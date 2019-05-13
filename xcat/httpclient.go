package xcat

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
)

type HttpClient struct {
	Client  *http.Client
	Headers http.Header
}

func (s *HttpClient) request(method, url string, headers *http.Header, body io.Reader) (req *http.Request, err error) {
	req, err = http.NewRequest(method, url, body)
	if err != nil {
		return nil, err
	}
	if headers == nil && body != nil {
		headers = &http.Header{}
		headers.Add("Content-Type", "application/json")
		req.Header = *headers
	}
	return
}

func (s *HttpClient) Request(method string, url string, params *url.Values, headers *http.Header, body *[]byte) (data interface{}, err error) {
	// add params to url here
	if params != nil {
		url = url + "?" + params.Encode()
	}

	// Get the body if one is present
	var buf io.Reader
	if body != nil {
		buf = bytes.NewReader(*body)
	}
	req, err := s.request(method, url, headers, buf)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Accept", "application/json")
	var resp *http.Response
	resp, err = s.Do(req)
	if err != nil {
		return nil, err
	}
	err = CheckHTTPResponseStatusCode(resp)
	if err != nil {
		rbody, _ := ioutil.ReadAll(resp.Body)
		if rbody != nil {
			data := make(map[string]interface{})
			errjson := json.Unmarshal([]byte(rbody), &data)
                        if errjson != nil {
                            return "No message from response", err
                        }
			return data["message"], err
		}
		return nil, err
	}
	rbody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, errors.New("Can not read the message form response")
	}
	return string(rbody), nil
}

func (s *HttpClient) Do(req *http.Request) (*http.Response, error) {
	for k := range s.Headers {
		req.Header.Set(k, s.Headers.Get(k))
	}
	resp, err := s.Client.Do(req)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (s *HttpClient) Get(url string, params *url.Values, body interface{}) (interface{}, error) {
	var resp interface{}
	var err error
	if body != nil {
		bodyJson, _ := json.Marshal(body)
		resp, err = s.Request("GET", url, params, nil, &bodyJson)
	} else {
		resp, err = s.Request("GET", url, params, nil, nil)
	}
	return resp, err
}

func (s *HttpClient) Post(url string, params *url.Values, body interface{}) (interface{}, error) {
	bodyJson, err := json.Marshal(body)
	resp, err := s.Request("POST", url, params, nil, &bodyJson)
	return resp, err
}

func (s *HttpClient) Put(url string, params *url.Values, body interface{}) (interface{}, error) {
	var resp interface{}
	var err error
	if body != nil {
		bodyJson, _ := json.Marshal(body)
		resp, err = s.Request("PUT", url, params, nil, &bodyJson)
	} else {
		resp, err = s.Request("PUT", url, params, nil, nil)
	}
	return resp, err
}

func (s *HttpClient) Delete(url string, params *url.Values, body interface{}) (interface{}, error) {
	var resp interface{}
	var err error
	if body != nil {
		bodyJson, _ := json.Marshal(body)
		resp, err = s.Request("DELETE", url, params, nil, &bodyJson)
	} else {
		resp, err = s.Request("DELETE", url, params, nil, nil)
	}
	return resp, err
}

func (s *HttpClient) Patch(url string, params *url.Values, body interface{}) (interface{}, error) {
	bodyJson, err := json.Marshal(body)
	resp, err := s.Request("PATCH", url, params, nil, &bodyJson)
	return resp, err
}

func CheckHTTPResponseStatusCode(resp *http.Response) error {
	switch resp.StatusCode {
	case 200, 201, 202, 204, 206:
		return nil
	case 400:
		return errors.New("400 bad request")
	case 401:
		return errors.New("401 unauthorised")
	case 403:
		return errors.New("403 forbidden")
	case 404:
		return errors.New("404 not found")
	case 405:
		return errors.New("405 method not allowed")
	case 409:
		return errors.New("409 conflict")
	case 413:
		return errors.New("413 over limit")
	case 415:
		return errors.New("415 bad media type")
	case 422:
		return errors.New("422 unprocessable")
	case 429:
		return errors.New("429 too many request")
	case 500:
		return errors.New("500 instance fault / server err")
	case 501:
		return errors.New("501 not implemented")
	case 503:
		return errors.New("503 service unavailable")
	}
	return errors.New("Unexpected response status code" + string(resp.StatusCode))
}
