package net

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
)

// MakePOSTRequest makes an HTTP POST request with params as payload and returns its body,
// status code and nil if successful. Returns nil, 0 and the error if one occurs.
func MakePOSTRequest(targetURL string, params map[string]interface{}) ([]byte, int, error) {
	postBody, err := json.Marshal(params)
	if err != nil {
		return nil, 0, err
	}

	resp, err := http.Post(targetURL, "application/json", bytes.NewBuffer(postBody))
	if err != nil {
		return nil, 0, err
	}
	defer resp.Body.Close()

	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, 0, err
	}

	return respBody, resp.StatusCode, nil
}

// MakeGETRequest makes an HTTP GET request with the parameters specified in req
// and returns its body and nil if successful.
func MakeGETRequest(req *http.Request) ([]byte, int, error) {
	client := &http.Client{}

	resp, err := client.Do(req)
	if err != nil {
		return nil, 0, err
	}
	defer resp.Body.Close()

	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, 0, err
	}

	return respBody, resp.StatusCode, nil
}
