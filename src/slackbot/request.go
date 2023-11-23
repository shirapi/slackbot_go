package slackbot

import "encoding/json"

type Request struct {
	Token     string `json:"token"`
	Challenge string `json:"challenge"`
	Type      string `json:"type"`
}

func ParseRequest(s string) (Request, error) {
	var req Request
	err := json.Unmarshal([]byte(s), &req)
	if err != nil {
		return Request{}, err
	}
	return req, err
}
