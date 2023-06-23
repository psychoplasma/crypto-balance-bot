package net

import (
	"encoding/json"
	"errors"
	"net/http"
)

// GetJSON makes a HTTP GET call and expects a JSON response
func GetJSON(url string, v interface{}) error {
	client := &http.Client{}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return err
	}

	req.Header.Add("Accept", `text/html,application/xhtml+xml,application/xml;q=0.9,*/*;q=0.8`)
	req.Header.Add("User-Agent", `Mozilla/5.0 (Macintosh; Intel Mac OS X 10.15; rv:95.0) Gecko/20100101 Firefox/95.0`)
	resp, err := client.Do(req)
	if err != nil {
		return err
	}

	if resp.StatusCode >= 400 {
		return errors.New(resp.Status)
	}

	defer resp.Body.Close()
	return json.NewDecoder(resp.Body).Decode(v)
}
