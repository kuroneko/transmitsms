package transmitsms

import (
	"errors"
	"io"
	"net/http"
	"net/url"
	"strings"
)

var (
	ErrMalformedRequest = errors.New("Malformed Request")
)

const (
	SmsTimestampFormat = "2006-01-02 15:04:05"
)

// SMSApi is the configuration object to access the SMS API.
//
// It contains a number of fields which control the various paremters relating
// to sending SMSes through the TransmitSMS API
type SMSApi struct {
	BaseURL   string
	APIKey    string
	APISecret string
}

type SendableRequest interface {
	RequestPath() string
	ToValues() (v url.Values, err error)
	DecodeResponse(hresp *http.Response) (resp interface{}, err error)
}

// newRequest creates a empty http.Request object with appropriate
// authorisation and TLS settings to communicate with the TransmitSMS API.
func (sms *SMSApi) newRequest(method, subPath string, body io.Reader) (req *http.Request, err error) {
	baseUrl, err := url.Parse(sms.BaseURL)
	if err != nil {
		return nil, err
	}
	reqUrl, err := baseUrl.Parse(subPath)
	if err != nil {
		return nil, err
	}
	req, err = http.NewRequest(method, reqUrl.String(), body)
	if err != nil {
		return nil, err
	}
	req.SetBasicAuth(sms.APIKey, sms.APISecret)

	return req, nil
}

func (sms *SMSApi) Send(r SendableRequest) (resp interface{}, err error) {
	v, err := r.ToValues()
	if err != nil {
		return nil, err
	}
	req, err := sms.newRequest("POST", r.RequestPath(), strings.NewReader(v.Encode()))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	hresp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	return r.DecodeResponse(hresp)
}
