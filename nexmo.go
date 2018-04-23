package sms

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/url"
)

// Nexmo represents a Nexmo request payload (https://developer.nexmo.com)
type Nexmo struct {
	APIKey    string
	APISecret string
	From      string
	to        string
	text      string
	callback  string
}

// NexmoResponse represents a Nexmo response payload in JSON
type NexmoResponse struct {
	MessageCount string                 `json:"message-count,omitempty"`
	Messages     []NexmoResponseMessage `json:"messages,omitempty"`
}

// NexmoResponseMessage represents a Nexmo message within the response payload
type NexmoResponseMessage struct {
	To               string `json:"to,omitempty"`
	MessageID        string `json:"message-id,omitempty"`
	Status           string `json:"status,omitempty"`
	ErrorText        string `json:"error-text,omitempty"`
	RemainingBalance string `json:"remaining-balance,omitempty"`
	MessagePrice     string `json:"message-price,omitempty"`
	Network          string `json:"network,omitempty"`
}

const (
	// NexmoResponseMessageStatusOK defines the Nexmo message status when OK
	NexmoResponseMessageStatusOK = "0"
)

var (
	// NexmoEndpoint defines the Nexmo ReST endpoint
	NexmoEndpoint = "https://rest.nexmo.com/sms/json"
)

// getResponse creates a http server request for Nexmo
func (n Nexmo) getResponse() (NexmoResponse, error) {
	var (
		b      []byte
		err    error
		req    *http.Request
		resp   *http.Response
		nResp  = NexmoResponse{}
		client = &http.Client{}
	)

	b = []byte(n.getQuery())

	if req, err = http.NewRequest(http.MethodPost, NexmoEndpoint, bytes.NewReader(b)); err != nil {
		return nResp, err
	}

	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	if resp, err = client.Do(req); err != nil {
		return nResp, err
	}

	defer resp.Body.Close()

	if b, err = ioutil.ReadAll(resp.Body); err != nil {
		return nResp, err
	}

	err = json.Unmarshal(b, &nResp)
	return nResp, err
}

// send sends the API request to Nexmo server with the `to` and `message` parameters
func (n Nexmo) send(to string, message Message, callback string) (interface{}, error) {

	n.to = to
	n.text = message.String()
	n.callback = callback

	return n.getResponse()
}

// MarshalJSON generates the JSON payload of a Nexmo object
func (n Nexmo) getQuery() string {
	u := url.Values{}

	if "" != n.APIKey {
		u.Set("api_key", n.APIKey)
	}

	if "" != n.APISecret {
		u.Set("api_secret", n.APISecret)
	}

	if "" != n.From {
		u.Set("from", n.From)
	}

	if "" != n.to {
		u.Set("to", n.to)
	}

	if "" != n.text {
		u.Set("text", n.text)
	}

	if "" != n.callback {
		u.Set("callback", n.callback)
	}

	return u.Encode()
}
