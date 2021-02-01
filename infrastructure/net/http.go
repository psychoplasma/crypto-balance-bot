package net

import (
	"encoding/json"
	"errors"
	"net/http"
)

// GetJSON makes a HTTP GET call and expects a JSON response
func GetJSON(url string, v interface{}) error {
	resp, err := http.Get(url)
	if err != nil {
		return err
	}

	if resp.StatusCode >= 400 {
		return errors.New(resp.Status)
	}

	defer resp.Body.Close()
	return json.NewDecoder(resp.Body).Decode(v)
}
